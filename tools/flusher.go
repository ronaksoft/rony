package tools

import (
	"sync/atomic"
)

/*
   Creation Time: 2020 - Dec - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type FlushEntry interface{}

type FlusherFunc func(targetID string, entries []FlushEntry)

type flusherPool struct {
	size        int32
	flusherFunc FlusherFunc
	poolMtx     SpinLock
	pool        map[string]*flusher
}

func NewFlusherPool(size int32, f FlusherFunc) *flusherPool {
	fp := &flusherPool{
		size:        size,
		flusherFunc: f,
		pool:        make(map[string]*flusher),
	}
	return fp
}

func (fp *flusherPool) Enter(targetID string, entry FlushEntry) {
	fp.poolMtx.Lock()
	f := fp.pool[targetID]
	fp.poolMtx.Unlock()
	if f == nil {
		f = &flusher{
			readyWorkers: fp.size,
			flusherFunc:  fp.flusherFunc,
			entryChan:    make(chan FlushEntry, fp.size),
			targetID:     targetID,
		}
		fp.poolMtx.Lock()
		fp.pool[targetID] = f
		fp.poolMtx.Unlock()
	}
	f.enter(entry)
}

type flusher struct {
	SpinLock
	readyWorkers int32
	flusherFunc  FlusherFunc
	entryChan    chan FlushEntry
	targetID     string
}

func (f *flusher) startWorker() {
	f.Lock()
	if atomic.AddInt32(&f.readyWorkers, -1) < 0 {
		atomic.AddInt32(&f.readyWorkers, 1)
		f.Unlock()
		return
	}
	f.Unlock()

	w := &worker{
		f: f,
	}
	go w.run()
}

func (f *flusher) enter(entry FlushEntry) {
	f.entryChan <- entry
	f.startWorker()
}

type worker struct {
	f *flusher
}

func (w *worker) run() {
	var (
		el = make([]FlushEntry, 0, 100)
	)
	for {
		el = el[:0]
		for {
			select {
			case e := <-w.f.entryChan:
				el = append(el, e)
				continue
			default:
			}
			break
		}

		w.f.Lock()
		if len(el) == 0 {
			// clean up and shutdown the worker
			atomic.AddInt32(&w.f.readyWorkers, 1)
			w.f.Unlock()
			break
		}
		w.f.Unlock()
		w.f.flusherFunc(w.f.targetID, el)
	}

}
