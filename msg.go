package rony

import (
	"git.ronaksoftware.com/ronak/rony/pools"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2020 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate protoc -I=. --go_out=. msg.proto
//go:generate protoc -I=. --gorony_out=. msg.proto
var (
	ConstructorNames = map[int64]string{}
)

// ProtoBufferMessage
type ProtoBufferMessage interface {
	Size() int
	MarshalToSizedBuffer([]byte) (int, error)
	Unmarshal([]byte) error
}

func ErrorMessage(out *MessageEnvelope, reqID uint64, errCode, errItem string) {
	errMessage := PoolError.Get()
	errMessage.Code = &errCode
	errMessage.Items = &errItem
	out.Fill(reqID, C_Error, errMessage)
	PoolError.Put(errMessage)
	return
}

func (x *MessageEnvelope) Clone() *MessageEnvelope {
	c := PoolMessageEnvelope.Get()
	if x.Constructor != nil {
		if c.Constructor == nil {
			c.Constructor = new(int64)
		}
		*c.Constructor = *x.Constructor
	}

	if x.RequestID != nil {
		if c.RequestID == nil {
			c.RequestID = new(uint64)
		}
		*c.RequestID = *x.RequestID
	}
	c.Message = append(c.Message[:0], x.Message...)
	return c
}

func (x *MessageEnvelope) CopyTo(e *MessageEnvelope) *MessageEnvelope {
	if x.RequestID != nil {
		if e.RequestID == nil {
			e.RequestID = new(uint64)
		}
		*e.RequestID = *x.RequestID
	}
	if x.Constructor != nil {
		if e.Constructor == nil {
			e.Constructor = new(int64)
		}
		*e.Constructor = *x.Constructor
	}
	e.Message = append(e.Message[:0], x.Message...)
	return e
}

func (x *MessageEnvelope) Fill(reqID uint64, constructor int64, p proto.Message) {
	if x.RequestID == nil {
		x.RequestID = new(uint64)
	}
	*x.RequestID = reqID
	if x.Constructor == nil {
		x.Constructor = new(int64)
	}
	*x.Constructor = constructor

	mo := proto.MarshalOptions{
		UseCachedSize: true,
	}
	b := pools.Bytes.GetCap(mo.Size(p))
	b, _ = mo.MarshalAppend(b, p)
	x.Message = append(x.Message[:0], b...)
	pools.Bytes.Put(b)
}

func (x *RaftCommand) Fill(senderID []byte, authID int64, e *MessageEnvelope) {
	x.Sender = append(x.Sender[:0], senderID...)
	if x.AuthID == nil {
		x.AuthID = new(int64)
	}
	*x.AuthID = authID
	x.Envelope = e.CopyTo(x.Envelope)
}

func (x *ClusterMessage) Fill(senderID []byte, authID int64, e *MessageEnvelope) {
	x.Sender = append(x.Sender[:0], senderID...)
	if x.AuthID == nil {
		x.AuthID = new(int64)
	}
	x.Envelope = e.CopyTo(x.Envelope)
}
