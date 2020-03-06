package flusher

import (
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"sync"
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2019 - Oct - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Entry struct {
	Key            interface{}
	Value          interface{}
	Ret            chan interface{}
	CallbackSynced bool
	Callback       func(ret interface{})
}

type Flusher struct {
	entries        chan Entry
	next           chan struct{}
	nextID         int64
	flushPeriod    time.Duration
	maxWorkers     int32
	maxBatchSize   int
	runningWorkers int32
	workerFunc     Func

	// Ticket Storage System
	ticket uint64
	done   uint64
	lock   int32
}

type Func func(items []Entry)

func New(maxBatchSize, maxConcurrency int32, flushPeriod time.Duration, ff Func) *Flusher {
	f := new(Flusher)
	f.flushPeriod = flushPeriod
	f.maxWorkers = maxConcurrency
	f.maxBatchSize = int(maxBatchSize)
	f.entries = make(chan Entry, f.maxBatchSize)
	f.next = make(chan struct{}, f.maxWorkers)
	f.workerFunc = ff

	// Run the 1st instance in the background
	for i := int32(0); i < f.maxWorkers; i++ {
		go f.worker()
	}

	// Set the first worker to listen on the channel
	f.next <- struct{}{}

	return f
}

// EnterWithChan
func (f *Flusher) EnterWithChan(key, value interface{}) chan interface{} {
	ch := make(chan interface{}, 1)
	f.entries <- Entry{
		Key:   key,
		Value: value,
		Ret:   ch,
	}
	return ch
}

// EnterWithResult
func (f *Flusher) EnterWithResult(key, value interface{}) (res interface{}) {
	waitGroup := acquireWaitGroup()
	waitGroup.Add(1)
	f.entries <- Entry{
		Key:            key,
		Value:          value,
		CallbackSynced: true,
		Callback: func(ret interface{}) {
			res = ret
			waitGroup.Done()
		},
	}
	waitGroup.Wait()
	releaseWaitGroup(waitGroup)
	return res
}

// EnterWithSyncCallback
// Use this function if the cbFunc is fast enough since, they will be called synchronously
func (f *Flusher) EnterWithSyncCallback(key, value interface{}, cbFunc func(interface{})) {
	f.entries <- Entry{
		Key:            key,
		Value:          value,
		CallbackSynced: true,
		Callback:       cbFunc,
	}
}

// EnterWithAsyncCallback
func (f *Flusher) EnterWithAsyncCallback(key, value interface{}, cbFunc func(interface{})) {
	f.entries <- Entry{
		Key:            key,
		Value:          value,
		CallbackSynced: false,
		Callback:       cbFunc,
	}
}

// Enter
// If you calling this function make sure that the Func registered with this Flusher
// can handle nil channels otherwise you may encounter undetermined results
func (f *Flusher) Enter(key, value interface{}) {
	f.entries <- Entry{
		Key:   key,
		Value: value,
		Ret:   nil,
	}
}

// PendingItems returns the number of items are still in the channel and are not picked by any job worker
func (f *Flusher) PendingItems() int {
	return len(f.entries)
}

// RunningJobs returns the number of job workers currently picking up items from the channel and/or running the flusher func
func (f *Flusher) RunningJobs() int {
	return int(f.runningWorkers)
}

func (f *Flusher) worker() {
	items := make([]Entry, 0, f.maxBatchSize)
	t := pools.AcquireTimer(f.flushPeriod)
	defer pools.ReleaseTimer(t)

	for {
		items = items[:0]

		// Wait for next signal to run the job
		<-f.next
		atomic.AddInt32(&f.runningWorkers, 1)

		// Wait for next entry
		select {
		case item := <-f.entries:
			items = append(items, item)
		}
		pools.ResetTimer(t, f.flushPeriod)
	InnerLoop:
		for {
			select {
			case item := <-f.entries:
				items = append(items, item)
				if len(items) >= f.maxBatchSize {
					// Send signal for the next worker to listen for entries
					f.next <- struct{}{}

					// Run the job synchronously
					f.workerFunc(items)

					atomic.AddInt32(&f.runningWorkers, -1)
					break InnerLoop
				}
			case <-t.C:
				// Send signal for the next worker to listen for entries
				f.next <- struct{}{}

				// Run the job synchronously
				f.workerFunc(items)

				atomic.AddInt32(&f.runningWorkers, -1)
				break InnerLoop
			}
		}
	}
}

const (
	maxWorkerIdleTime = 3 * time.Second
)

// LifoFlusher
type LifoFlusher struct {
	sync.Mutex
	waitingList    *tools.LinkedList
	flushPeriod    time.Duration
	maxWorkers     int32
	maxBatchSize   int
	runningWorkers int32
	currentWorker  *worker
	workerFunc     Func
	state          int32
	nextID         uint32
}

func NewLifo(maxBatchSize, maxConcurrency int32, flushPeriod time.Duration, ff Func) *LifoFlusher {
	f := new(LifoFlusher)
	f.flushPeriod = flushPeriod
	f.maxWorkers = maxConcurrency
	f.maxBatchSize = int(maxBatchSize)
	f.waitingList = tools.NewLinkedList()
	f.workerFunc = ff

	// Run the garbage collector in the background
	go f.gc()

	// Run the 1st worker
	f.nextWorker()

	return f
}

func (f *LifoFlusher) gc() {
	for {
		time.Sleep(maxWorkerIdleTime)
		currentTime := time.Now()
		n := f.waitingList.Head()
		for n != nil {
			if currentTime.Sub(n.GetData().(*worker).lastRun) >= maxWorkerIdleTime {
				w := f.waitingList.PickHeadData().(*worker)
				n = f.waitingList.Head()
				w.stop()
			} else {
				break
			}
		}
		// _Log.Debug("GC",
		// 	zap.Int32("Running", f.runningWorkers),
		// 	zap.Int32("Waiting", f.waitingList.Size()),
		// )
	}
}

func (f *LifoFlusher) nextWorker() bool {
	f.Lock()
	wn := f.waitingList.PickTailData()
	if wn == nil {
		if atomic.LoadInt32(&f.runningWorkers) < f.maxWorkers {
			atomic.AddInt32(&f.runningWorkers, 1)
			f.currentWorker = &worker{
				id:       atomic.AddUint32(&f.nextID, 1),
				flusher:  f,
				stopChan: make(chan struct{}, 1),
				entries:  make(chan Entry, f.maxBatchSize),
			}
			go f.currentWorker.run()
		} else {
			f.Unlock()
			return false
		}
	} else {
		f.currentWorker = wn.(*worker)
	}
	f.Unlock()
	return true
}

func (f *LifoFlusher) enter(entry Entry) {
	f.currentWorker.entries <- entry
}

type worker struct {
	id       uint32
	lastRun  time.Time
	flusher  *LifoFlusher
	entries  chan Entry
	stopChan chan struct{}
}

func (fw *worker) run() {
	items := make([]Entry, 0, fw.flusher.maxBatchSize)
	t := pools.AcquireTimer(fw.flusher.flushPeriod)
	defer pools.ReleaseTimer(t)
	for {
		fw.lastRun = time.Now()
		items = items[:0]
		// Wait until there is at least one item or worker has been stopped by GC.
		select {
		case item := <-fw.entries:
			items = append(items, item)
		case <-fw.stopChan:
			if len(fw.entries) > 0 {
				fw.stop()
			}
			atomic.AddInt32(&fw.flusher.runningWorkers, -1)
			return
		}
		pools.ResetTimer(t, fw.flusher.flushPeriod)
	InnerLoop:
		for {
			select {
			case item := <-fw.entries:
				items = append(items, item)
				if len(items) >= fw.flusher.maxBatchSize {
					fw.flusher.nextWorker()
					fw.flusher.workerFunc(items)
					break InnerLoop
				}
			case <-t.C:
				// Send signal for the next worker to listen for entries and run the job synchronously
				fw.flusher.nextWorker()
				fw.flusher.workerFunc(items)
				break InnerLoop
			}
		}
		// If there is no more entry in the list, then pass the worker to the waiting list,
		// then it either will be selected for future requests, or will be stopped by GC.
		if len(fw.entries) == 0 {
			fw.flusher.waitingList.Append(fw)
		}
	}
}

func (fw *worker) stop() bool {
	fw.stopChan <- struct{}{}
	// select {
	// case
	// 	return true
	// default:
	// 	// _Log.Warn("Cannot stop the Worker", zap.Uint32("ID", fw.id))
	// }
	// return false
	return true
}

// Enter
// If you calling this function make sure that the Func registered with this Flusher
// can handle nil channels otherwise you may encounter undetermined results
func (f *LifoFlusher) Enter(key, value interface{}) {
	f.enter(Entry{
		Key:   key,
		Value: value,
		Ret:   nil,
	})
}

// EnterWithChan
func (f *LifoFlusher) EnterWithChan(key, value interface{}) chan interface{} {
	ch := make(chan interface{}, 1)
	f.enter(Entry{
		Key:   key,
		Value: value,
		Ret:   ch,
	})
	return ch
}

// EnterWithResult
func (f *LifoFlusher) EnterWithResult(key, value interface{}) (res interface{}) {
	waitGroup := acquireWaitGroup()
	waitGroup.Add(1)
	f.enter(Entry{
		Key:            key,
		Value:          value,
		CallbackSynced: true,
		Callback: func(ret interface{}) {
			res = ret
			waitGroup.Done()
		},
	})
	waitGroup.Wait()
	releaseWaitGroup(waitGroup)
	return res
}

// RunningJobs returns the number of job workers currently picking up items from the channel and/or running the flusher func
func (f *LifoFlusher) RunningJobs() int {
	return int(f.runningWorkers)
}

/*
	Helper Functions
*/

type entryPool struct {
	pool sync.Pool
}

func (p *entryPool) Get() *Entry {
	fe, ok := p.pool.Get().(*Entry)
	if !ok {
		return &Entry{}
	}
	fe.Key = nil
	fe.Value = nil
	fe.CallbackSynced = false
	fe.Callback = nil
	return fe
}

func (p *entryPool) Put(entry *Entry) {
	p.pool.Put(entry)
}
