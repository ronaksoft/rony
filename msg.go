package rony

import "git.ronaksoftware.com/ronak/rony/internal/pools"

/*
   Creation Time: 2020 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate protoc -I=. --gorony_out=. msg.proto
//go:generate protoc -I=./vendor -I=.  --gogofaster_out=plugins=grpc:. msg.proto
var (
	ConstructorNames = map[int64]string{}
)

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
	c := &MessageEnvelope{
		Constructor: m.Constructor,
		RequestID:   m.RequestID,
	}
	c.Message = make([]byte, len(m.Message))
	copy(c.Message, m.Message)
	return c
}

func (m *MessageEnvelope) CopyTo(e *MessageEnvelope) {
	e.Constructor = m.Constructor
	e.RequestID = m.RequestID
	if len(m.Message) > cap(e.Message) {
		pools.Bytes.Put(e.Message)
		e.Message = pools.Bytes.GetLen(len(m.Message))
	} else {
		e.Message = e.Message[:len(m.Message)]
	}
	copy(e.Message, m.Message)
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
