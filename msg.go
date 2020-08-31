package rony

import (
	"git.ronaksoft.com/ronak/rony/pools"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2020 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. --go_out=. msg.proto
//go:generate protoc -I=. --gorony_out=. msg.proto
var (
	ConstructorNames = map[int64]string{}
)

func ErrorMessage(out *MessageEnvelope, reqID uint64, errCode, errItem string) {
	errMessage := PoolError.Get()
	errMessage.Code = errCode
	errMessage.Items = errItem
	out.Fill(reqID, C_Error, errMessage)
	PoolError.Put(errMessage)
	return
}

func (x *MessageEnvelope) Clone() *MessageEnvelope {
	c := PoolMessageEnvelope.Get()
	c.Constructor = x.Constructor
	c.RequestID = x.RequestID
	c.Message = append(c.Message[:0], x.Message...)
	return c
}

func (x *MessageEnvelope) Fill(reqID uint64, constructor int64, p proto.Message) {
	x.RequestID = reqID
	x.Constructor = constructor
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
	x.AuthID = authID
	e.DeepCopy(x.Envelope)
}

func (x *ClusterMessage) Fill(senderID []byte, authID int64, e *MessageEnvelope) {
	x.Sender = append(x.Sender[:0], senderID...)
	x.AuthID = authID
	e.DeepCopy(x.Envelope)
}
