package rony

import (
	"git.ronaksoftware.com/ronak/rony/internal/pools"
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
//go:generate protoc -I=. --gogofaster_out=. msg.proto
var (
	ConstructorNames = map[int64]string{}
)

// ProtoBufferMessage
type ProtoBufferMessage interface {
	Size() int
	MarshalToSizedBuffer([]byte) (int, error)
}

func ErrorMessage(out *MessageEnvelope, reqID uint64, errCode, errItem string) {
	errMessage := PoolError.Get()
	errMessage.Code = errCode
	errMessage.Items = errItem
	out.Fill(reqID, C_Error, errMessage)
	PoolError.Put(errMessage)
	return
}

func (m *MessageEnvelope) Clone() *MessageEnvelope {
	// c := &MessageEnvelope{}
	c := PoolMessageEnvelope.Get()
	c.Constructor = m.Constructor
	c.RequestID = m.RequestID
	c.Message = append(c.Message[:0], m.Message...)
	return c
}

func (m *MessageEnvelope) CopyTo(e *MessageEnvelope) *MessageEnvelope {
	e.Constructor = m.Constructor
	e.RequestID = m.RequestID
	e.Message = append(e.Message[:0], m.Message...)
	return e
}

func (m *MessageEnvelope) Fill(reqID uint64, constructor int64, p ProtoBufferMessage) {
	m.RequestID = reqID
	m.Constructor = constructor
	protoSize := p.Size()
	b := pools.Bytes.GetLen(protoSize)
	_, _ = p.MarshalToSizedBuffer(b)
	m.Message = append(m.Message[:0], b...)
	pools.Bytes.Put(b)
}

func (m *RaftCommand) Fill(senderID []byte, authID int64, e *MessageEnvelope) {
	m.Sender = append(m.Sender[:0], senderID...)
	m.AuthID = authID
	m.Envelope = e.CopyTo(m.Envelope)
}

func (m *ClusterMessage) Fill(senderID []byte, authID int64, e *MessageEnvelope) {
	m.Sender = append(m.Sender[:0], senderID...)
	m.AuthID = authID
	m.Envelope = e.CopyTo(m.Envelope)
}
