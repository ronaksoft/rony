package natsBridge

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/internal/flusher"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/pool/pbytes"
	"github.com/nats-io/nats.go"
	"github.com/scylladb/go-set/u64set"
	"go.uber.org/zap"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2019 - Feb - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type ErrHandler func(err error)
type NotifyHandler func(connIDs []uint64)
type MessageHandler func(c *Container) bool
type Config struct {
	BundleID   string
	InstanceID string
	Timeout    time.Duration
	Retries    int
	Options    nats.Options
	MaxWorkers int
}

type Bridge struct {
	ErrHandler
	NotifyHandler
	MessageHandler
	messageID       uint64
	prefixMessageID uint64

	bundleID         string
	instanceID       string
	bridgeID         string
	deliveryEndpoint string
	notifyEndpoint   string
	subNotifier      *nats.Subscription
	subDelivery      *nats.Subscription
	subMessage       *nats.Subscription

	deliveryDispatcher *flusher.LifoFlusher
	notifyDispatcher   *flusher.LifoFlusher
	messageDispatcher  *flusher.LifoFlusher

	counters       Counters
	unconfirmedMtx sync.Mutex
	unconfirmed    map[uint64]chan bool
	conn           *nats.Conn
	waitGroup      sync.WaitGroup
	maxWorkers     int
}

type Counters struct {
	MessageReceived              int32
	MessageSend                  int32
	MessageSendTimeOut           int32
	MessageSendDelivered         int32
	NotifierSendTimeOut          int32
	NotifierSendDelivered        int32
	NotifierSendConnections      int32
	NotifierReceived             int32
	NotifierReceivedConnections  int32
	DeliveryReceived             int32
	DeliveryReceivedAcknowledges int32
	DeliverySend                 int32
	DeliverySendAcknowledges     int32
}

func NewBridge(conf Config) (*Bridge, error) {
	b := &Bridge{
		bundleID:    conf.BundleID,
		instanceID:  conf.InstanceID,
		bridgeID:    fmt.Sprintf("%s.%s", conf.BundleID, conf.InstanceID),
		unconfirmed: make(map[uint64]chan bool),
	}

	if conf.MaxWorkers <= 0 {
		b.maxWorkers = defaultMaxWorkers
	} else {
		b.maxWorkers = conf.MaxWorkers
	}
	b.prefixMessageID = uint64(tools.RandomInt(1<<31)) << 32

	b.ErrHandler = func(err error) {}
	b.NotifyHandler = func(connIDs []uint64) {}
	b.MessageHandler = func(c *Container) bool { return false }

	conf.Options.Name = fmt.Sprintf("%s.%s", b.bundleID, b.instanceID)
	conf.Options.MaxReconnect = -1
	if c, err := conf.Options.Connect(); err != nil {
		return nil, err
	} else {
		b.conn = c
	}

	// Setup Subscriptions
	err := b.setupSubscriptions()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bridge) nextMessageID() uint64 {
	return b.prefixMessageID | atomic.AddUint64(&b.messageID, 1)
}

func (b *Bridge) setupSubscriptions() (err error) {
	// Initialize Delivery Subscription
	chDelivery := make(chan *nats.Msg, flushMaxSize)
	b.deliveryEndpoint = fmt.Sprintf("%s.%s.%s", b.bundleID, b.instanceID, deliveryReportSubject)
	if b.subDelivery, err = b.conn.ChanSubscribe(
		b.deliveryEndpoint,
		chDelivery,
	); err != nil {
		return err
	}
	for i := 0; i < b.maxWorkers; i++ {
		go b.deliveryReceiver(chDelivery)
	}
	b.deliveryDispatcher = flusher.NewLifo(flushMaxSize, flushMaxWorkers, flushPeriod, b.deliverySender)

	// Initialize Notifier Subscription
	chNotifier := make(chan *nats.Msg, flushMaxSize)
	b.notifyEndpoint = fmt.Sprintf("%s.%s.%s", b.bundleID, b.instanceID, notifierSubject)
	if b.subNotifier, err = b.conn.ChanSubscribe(
		b.notifyEndpoint,
		chNotifier,
	); err != nil {
		return err
	}
	for i := 0; i < b.maxWorkers; i++ {
		go b.notifierReceiver(chNotifier)
	}
	b.notifyDispatcher = flusher.NewLifo(flushMaxSize, flushMaxWorkers, flushPeriod, b.notifySender)

	chMessage := make(chan *nats.Msg, flushMaxSize)
	if b.subMessage, err = b.conn.ChanQueueSubscribe(
		fmt.Sprintf("%s.%s", b.bundleID, messageSubject),
		fmt.Sprintf("%s.%s", b.bundleID, messageSubject),
		chMessage,
	); err != nil {
		return err
	}
	for i := 0; i < b.maxWorkers; i++ {
		go b.messageReceiver(chMessage)
	}
	b.messageDispatcher = flusher.NewLifo(flushMaxSize, flushMaxWorkers, flushPeriod, b.messageSender)
	return nil
}

func (b *Bridge) deliverySender(items []flusher.Entry) {
	m := make(map[string]*u64set.Set)
	for idx := range items {
		key := items[idx].Key.(string)
		val := items[idx].Value.(uint64)
		if _, ok := m[key]; !ok {
			m[key] = u64set.New(val)
		} else {
			m[key].Add(val)
		}
	}

	for subject := range m {
		if ce := log.Check(log.DebugLevel, "Bridge Delivery Sent"); ce != nil {
			ce.Write(
				zap.Int("Count", m[subject].Size()),
				zap.String("BundleID", b.bundleID),
				zap.String("InstanceID", b.instanceID),
				zap.String("Receiver", subject),
			)
		}
		deliverySendFunc(b, subject, m[subject].List())
	}
}
func deliverySendFunc(b *Bridge, dstSubject string, messageIDs []uint64) {
	deliveryReport := new(DeliveryReport)
	deliveryReport.MessageIDs = messageIDs

	bytes := pbytes.GetLen(deliveryReport.Size())
	defer pbytes.Put(bytes)

	n, err := deliveryReport.MarshalTo(bytes)
	if err != nil {
		b.ErrHandler(err)
		return
	}

	err = b.conn.Publish(dstSubject, bytes[:n])
	if err != nil {
		b.ErrHandler(err)
		return
	}

	// for debugging purpose
	atomic.AddInt32(&b.counters.DeliverySend, 1)
	atomic.AddInt32(&b.counters.DeliverySendAcknowledges, int32(len(messageIDs)))
}

func (b *Bridge) notifySender(items []flusher.Entry) {
	m := make(map[string]*u64set.Set)
	for idx := range items {
		key := items[idx].Key.(string)
		val := items[idx].Value.(uint64)
		if _, ok := m[key]; !ok {
			m[key] = u64set.New(val)
		} else {
			m[key].Add(val)
		}
	}

	for subject := range m {
		if ce := log.Check(log.DebugLevel, "Bridge Notify Sent"); ce != nil {
			ce.Write(
				zap.Int("Count", m[subject].Size()),
				zap.String("BundleID", b.bundleID),
				zap.String("InstanceID", b.instanceID),
				zap.String("Receiver", subject),
			)
		}
		notifySendFunc(b, subject, m[subject].List())
	}
}
func notifySendFunc(b *Bridge, dstSubject string, connIDs []uint64) {
	if dstSubject == b.notifyEndpoint {
		b.NotifyHandler(connIDs)
		return
	}

	notifier := new(Notifier)
	notifier.MessageID = b.nextMessageID()
	notifier.ConnectionIDs = connIDs

	bytes := pbytes.GetLen(notifier.Size())
	defer pbytes.Put(bytes)

	_, err := notifier.MarshalTo(bytes)
	if err != nil {
		b.ErrHandler(err)
		return
	}

	chDelivered := make(chan bool, 1)
	b.unconfirmedMtx.Lock()
	b.unconfirmed[notifier.MessageID] = chDelivered
	b.unconfirmedMtx.Unlock()

	timer := pools.AcquireTimer(sendTimeout)
	defer pools.ReleaseTimer(timer)

	err = b.conn.PublishRequest(dstSubject, b.deliveryEndpoint, bytes)
	if err != nil {
		b.ErrHandler(err)
		return
	}

	atomic.AddInt32(&b.counters.NotifierSendConnections, int32(len(connIDs)))
	select {
	case <-timer.C:
		atomic.AddInt32(&b.counters.NotifierSendTimeOut, 1)
		b.unconfirmedMtx.Lock()
		delete(b.unconfirmed, notifier.MessageID)
		b.unconfirmedMtx.Unlock()
		for _, connID := range notifier.ConnectionIDs {
			b.notifyDispatcher.Enter(dstSubject, connID)
		}
	case <-chDelivered:
		atomic.AddInt32(&b.counters.NotifierSendDelivered, 1)
	}
}

func (b *Bridge) messageSender(items []flusher.Entry) {
	m := make(map[string][]messageRequest)
	md := make(map[string]bool) // message delivered
	for idx := range items {
		key := items[idx].Key.(string)
		val := items[idx].Value.(messageRequest)
		if _, ok := m[key]; !ok {
			m[key] = make([]messageRequest, 0, len(items))
		}
		m[key] = append(m[key], val)
	}

	waitGroup := pools.AcquireWaitGroup()
	lock := sync.Mutex{}
	waitGroup.Add(len(m))
	for subject := range m {
		go func(subject string) {
			if ce := log.Check(log.DebugLevel, "Bridge Message Sent"); ce != nil {
				ce.Write(
					zap.Int("Count", len(m[subject])),
					zap.String("BundleID", b.bundleID),
					zap.String("InstanceID", b.instanceID),
					zap.String("Receiver", subject),
				)
			}
			delivered := messageSendFunc(b, subject, m[subject])
			lock.Lock()
			md[subject] = delivered
			lock.Unlock()

			waitGroup.Done()
		}(subject)
	}
	waitGroup.Wait()
	pools.ReleaseWaitGroup(waitGroup)

	for idx := range items {
		items[idx].Callback(md[items[idx].Key.(string)])
	}
}
func messageSendFunc(b *Bridge, dstSubject string, reqs []messageRequest) bool {
	mc := new(Container)
	mc.MessageID = b.nextMessageID()

	for idx := range reqs {
		mc.Messages = append(mc.Messages, &Message{
			AuthID:  reqs[idx].AuthID,
			UserID:  reqs[idx].UserID,
			Payload: reqs[idx].Data,
		})
	}

	bytes := pbytes.GetLen(mc.Size())
	defer pbytes.Put(bytes)

	_, err := mc.MarshalTo(bytes)
	if err != nil {
		b.ErrHandler(err)
		return false
	}

	chDelivered := make(chan bool, 1)
	b.unconfirmedMtx.Lock()
	b.unconfirmed[mc.MessageID] = chDelivered
	b.unconfirmedMtx.Unlock()

	retry := 0
	delivered := false
	keepGoing := true

	t := pools.AcquireTimer(sendTimeout)
	defer pools.ReleaseTimer(t)

	for keepGoing {
		err := b.conn.PublishRequest(dstSubject, b.deliveryEndpoint, bytes)
		if err != nil {
			b.ErrHandler(err)
			return false

		}
		atomic.AddInt32(&b.counters.MessageSend, 1)
		select {
		case <-t.C:
			atomic.AddInt32(&b.counters.MessageSendTimeOut, 1)
			retry++
			if retry > sendRetries {
				b.unconfirmedMtx.Lock()
				delete(b.unconfirmed, mc.MessageID)
				b.unconfirmedMtx.Unlock()
				keepGoing = false
			} else {
				pools.ResetTimer(t, sendTimeout)
			}

		case <-chDelivered:
			atomic.AddInt32(&b.counters.MessageSendDelivered, 1)
			delivered = true
			keepGoing = false
		}
	}
	return delivered
}

func (b *Bridge) deliveryReceiver(ch <-chan *nats.Msg) {
	b.waitGroup.Add(1)
	defer b.waitGroup.Done()
	deliveryReport := new(DeliveryReport)
	for m := range ch {
		if ce := log.Check(log.DebugLevel, "Bridge Delivery Received"); ce != nil {
			ce.Write(
				zap.String("BundleID", b.bundleID),
				zap.String("InstanceID", b.instanceID),
			)
		}
		deliveryReport.MessageIDs = deliveryReport.MessageIDs[:0]
		atomic.AddInt32(&b.counters.DeliveryReceived, 1)

		if err := deliveryReport.Unmarshal(m.Data); err != nil {
			b.ErrHandler(err)
			continue
		}
		atomic.AddInt32(&b.counters.DeliveryReceivedAcknowledges, int32(len(deliveryReport.MessageIDs)))

		// For each of the MessageIDs send the delivery signal
		for _, messageID := range deliveryReport.MessageIDs {
			b.unconfirmedMtx.Lock()
			if v, ok := b.unconfirmed[messageID]; ok {
				v <- true
				delete(b.unconfirmed, messageID)
			}
			b.unconfirmedMtx.Unlock()
		}
	}

	return
}

func (b *Bridge) notifierReceiver(ch <-chan *nats.Msg) {
	b.waitGroup.Add(1)
	defer b.waitGroup.Done()
	notifier := new(Notifier)
	for m := range ch {
		if ce := log.Check(log.DebugLevel, "Bridge Notifier Received"); ce != nil {
			ce.Write(
				zap.String("BundleID", b.bundleID),
				zap.String("InstanceID", b.instanceID),
			)
		}
		notifier.ConnectionIDs = notifier.ConnectionIDs[:0]
		atomic.AddInt32(&b.counters.NotifierReceived, 1)

		if err := notifier.Unmarshal(m.Data); err != nil {
			b.ErrHandler(err)
			continue
		}

		// Send Delivery Report
		b.deliveryDispatcher.Enter(m.Reply, notifier.MessageID)

		// Call the NotifyHandler
		b.NotifyHandler(notifier.ConnectionIDs)

		// for debugging purpose
		atomic.AddInt32(&b.counters.NotifierReceivedConnections, int32(len(notifier.ConnectionIDs)))
	}
}

func (b *Bridge) messageReceiver(ch <-chan *nats.Msg) {
	b.waitGroup.Add(1)
	defer b.waitGroup.Done()
	for m := range ch {
		if ce := log.Check(log.DebugLevel, "Bridge Message Received"); ce != nil {
			ce.Write(
				zap.String("BundleID", b.bundleID),
				zap.String("InstanceID", b.instanceID),
			)
		}

		// for debugging purpose
		atomic.AddInt32(&b.counters.MessageReceived, 1)

		// Try to unmarshal the message
		container := &Container{}
		if err := container.Unmarshal(m.Data); err != nil {
			b.ErrHandler(err)
			continue
		}

		// Call the message handler
		if b.MessageHandler(container) {
			// Send Delivery Report
			b.deliveryDispatcher.Enter(m.Reply, container.MessageID)
		}
	}
}

func (b *Bridge) GetID() string {
	return fmt.Sprintf("%s.%s", b.bundleID, b.instanceID)
}

func (b *Bridge) Counters() Counters {
	return b.counters
}

func (b *Bridge) SendNotifier(bridgeID string, connID uint64) {
	dstSubject := strings.Builder{}
	dstSubject.WriteString(bridgeID)
	dstSubject.WriteString(".")
	dstSubject.WriteString(notifierSubject)
	b.notifyDispatcher.Enter(dstSubject.String(), connID)
}

func (b *Bridge) SendMessage(bridgeID string, authID int64, userID int64, data []byte) bool {
	dstSubject := strings.Builder{}
	dstSubject.WriteString(bridgeID)
	dstSubject.WriteString(".")
	dstSubject.WriteString(messageSubject)
	return b.messageDispatcher.EnterWithResult(
		dstSubject.String(),
		messageRequest{
			AuthID: authID,
			UserID: userID,
			Data:   data,
		},
	).(bool)
}

func (b *Bridge) Shutdown() error {
	if err := b.subDelivery.Unsubscribe(); err != nil {
		return err
	}
	if err := b.subNotifier.Unsubscribe(); err != nil {
		return err
	}
	if err := b.subMessage.Unsubscribe(); err != nil {
		return err
	}

	// Close the NATS connection
	b.conn.Close()

	// Wait for all routines to finish
	b.waitGroup.Wait()

	return nil
}

type messageRequest struct {
	AuthID int64
	UserID int64
	Data   []byte
}
