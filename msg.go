package rony

import (
	"git.ronaksoftware.com/ronak/rony/pools"
)

/*
   Creation Time: 2020 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate protoc -I=. --gopool_out=. msg.proto
//go:generate protoc -I=.  --gogofaster_out=. msg.proto
var (
	ConstructorNames = map[int64]string{}
)

// ProtoBufferMessage
type ProtoBufferMessage interface {
	Size() int
	MarshalToSizedBuffer([]byte) (int, error)
}

func ErrorMessage(out *MessageEnvelope, errCode, errItem string) {
	errMessage := PoolError.Get()
	errMessage.Code = errCode
	errMessage.Items = errItem
	out.Message, _ = errMessage.Marshal()
	out.Constructor = C_Error
	PoolError.Put(errMessage)
	return
}

func (m *MessageEnvelope) Clone() *MessageEnvelope {
	c := PoolMessageEnvelope.Get()
	c.Constructor = m.Constructor
	c.RequestID = m.RequestID
	c.Message = append(c.Message[:0], m.Message...)
	return c
}

func (m *MessageEnvelope) CopyTo(e *MessageEnvelope) {
	e.Constructor = m.Constructor
	e.RequestID = m.RequestID
	e.Message = append(e.Message[:0], m.Message...)
}

func (m *MessageEnvelope) Fill(reqID uint64, constructor int64, p ProtoBufferMessage) {
	m.RequestID = reqID
	m.Constructor = constructor
	b := pools.Bytes.GetLen(p.Size())
	_, _ = p.MarshalToSizedBuffer(b)
	m.Message = append(m.Message[:0], b...)
	pools.Bytes.Put(b)
}

func (m *UpdateEnvelope) Clone() *UpdateEnvelope {
	c := &UpdateEnvelope{
		Constructor: m.Constructor,
		UCount:      m.UCount,
		UpdateID:    m.UpdateID,
		Timestamp:   m.Timestamp,
	}
	c.Update = make([]byte, len(m.Update))
	copy(c.Update, m.Update)
	return c
}

func (m *UpdateEnvelope) CopyTo(u *UpdateEnvelope) {
	u.Constructor = m.Constructor
	u.UCount = m.UCount
	u.UpdateID = m.UpdateID
	u.Timestamp = m.Timestamp
	if len(m.Update) > cap(u.Update) {
		pools.Bytes.Put(u.Update)
		u.Update = pools.Bytes.GetLen(len(m.Update))
	} else {
		u.Update = u.Update[:len(m.Update)]
	}
	copy(u.Update, m.Update)
}
