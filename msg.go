package rony

/*
   Creation Time: 2020 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate protoc -I=.  --gogofaster_out=. msg.proto
//go:generate protoc -I=. --gohelpers_out=. msg.proto
var (
	ConstructorNames = map[int64]string{}
)

func ErrorMessage(out *MessageEnvelope, errCode, errItem string) {
	errMessage := PoolError.Get()
	errMessage.Code = errCode
	errMessage.Items = errItem
	ResultError(out, errMessage)
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