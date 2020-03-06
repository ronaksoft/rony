package quicGateway

import (
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"github.com/lucas-clemente/quic-go"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

/*
   Creation Time: 2019 - Aug - 26
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Conn
type Conn struct {
	sync.RWMutex
	ConnID     uint64
	AuthID     int64
	UserID     int64
	StreamID   quic.StreamID
	AuthKey    []byte
	ClientIP   string
	Counters   ConnCounters
	ServerSeq  int64
	ClientSeq  int64
	ClientType string

	// Internals
	gateway     *Gateway
	conn        quic.Session
	writeErrors int32
	closed      bool
	flushChan   chan bool
	flushing    int32
	streams     map[quic.StreamID]quic.Stream
	noStreams   int32
}

func (qc *Conn) GetAuthID() int64 {
	return qc.AuthID
}

func (qc *Conn) SetAuthID(authID int64) {
	qc.AuthID = authID
}

func (qc *Conn) GetAuthKey() []byte {
	return qc.AuthKey
}

func (qc *Conn) SetAuthKey(authKey []byte) {
	qc.AuthKey = authKey
}

func (qc *Conn) IncServerSeq(x int64) int64 {
	return atomic.AddInt64(&qc.ServerSeq, x)
}

func (qc *Conn) GetConnID() uint64 {
	return qc.ConnID
}

func (qc *Conn) GetClientIP() string {
	return qc.ClientIP
}

func (qc *Conn) GetUserID() int64 {
	return qc.UserID
}

func (qc *Conn) SetUserID(userID int64) {
	qc.UserID = userID
}

func (qc *Conn) Persistent() bool {
	return true
}

type ConnCounters struct {
	IncomingCount uint64
	IncomingBytes uint64
	OutgoingCount uint64
	OutgoingBytes uint64
}

// SendBinary
// Make sure you don't use payload after calling this function, because its underlying
// array will be put back into the pool to be reused.
// You MUST NOT re-use the underlying array of payload, otherwise you might get unexpected results.
func (qc *Conn) SendBinary(streamID int64, payload []byte) error {
	if qc.closed {
		return ErrWriteToClosedConn
	}
	qc.RLock()
	stream, ok := qc.streams[quic.StreamID(streamID)]
	qc.RUnlock()
	if !ok {
		return nil
	}
	_, err := stream.Write(payload)
	_ = stream.Close()
	return err
}

func (qc *Conn) Disconnect() {
	qc.gateway.removeConnection(qc.ConnID)
}

// Flush
func (qc *Conn) Flush() {
	if atomic.CompareAndSwapInt32(&qc.flushing, 0, 1) {
		go qc.flushJob()
	}
	select {
	case qc.flushChan <- true:
	default:
	}
}

func (qc *Conn) flushJob() {
	timer := pools.AcquireTimer(flushDebounceTime)

	keepGoing := true
	cnt := 0
	for keepGoing {
		if cnt++; cnt > flushDebounceMax {
			break
		}
		select {
		case <-qc.flushChan:
			pools.ResetTimer(timer, flushDebounceTime)
		case <-timer.C:
			// Stop the loop
			keepGoing = false
		}
	}

	// Read the flush function
	bytesSlice := qc.gateway.FlushFunc(qc)

	// Reset the flushing flag, let the flusher run again
	atomic.StoreInt32(&qc.flushing, 0)

	qc.Lock()
	if len(qc.streams) == 0 {
		s, err := qc.conn.OpenStream()
		if err != nil {
			qc.Unlock()
			pools.ReleaseTimer(timer)
			return
		}
		qc.streams[s.StreamID()] = s
	}
	for idx := 0; idx < len(bytesSlice); idx++ {
		for _, s := range qc.streams {
			_, err := s.Write(bytesSlice[idx])
			_ = s.Close()
			if err == nil {
				break
			}
			if ce := log.Check(log.DebugLevel, "Error On Write To Quic Stream"); ce != nil {
				ce.Write(
					zap.Uint64("ConnID", qc.ConnID),
					zap.Int64("StreamID", int64(s.StreamID())),
					zap.Int64("AuthID", qc.AuthID),
					zap.Error(err),
				)
			}
		}
	}
	qc.Unlock()
	pools.ReleaseTimer(timer)
}
