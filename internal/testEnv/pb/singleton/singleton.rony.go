// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.14.33
// Source: singleton.proto

package singleton

import (
	bytes "bytes"
	edge "github.com/ronaksoft/rony/edge"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

var _ = pools.Imported

const C_Single1 uint64 = 4833170411250891328

type poolSingle1 struct {
	pool sync.Pool
}

func (p *poolSingle1) Get() *Single1 {
	x, ok := p.pool.Get().(*Single1)
	if !ok {
		x = &Single1{}
	}

	return x
}

func (p *poolSingle1) Put(x *Single1) {
	if x == nil {
		return
	}

	x.ID = 0
	x.ShardKey = 0
	x.P1 = ""
	x.P2 = x.P2[:0]
	x.P5 = 0
	x.Enum = 0

	p.pool.Put(x)
}

var PoolSingle1 = poolSingle1{}

func (x *Single1) DeepCopy(z *Single1) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
	z.Enum = x.Enum
}

func (x *Single1) Clone() *Single1 {
	z := &Single1{}
	x.DeepCopy(z)
	return z
}

func (x *Single1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Single1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Single1) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Single1) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapSingle1(e registry.Envelope) (proto.Message, error) {
	x := &Single1{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Single1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Single1, x)
}

const C_Single2 uint64 = 4738594819076110912

type poolSingle2 struct {
	pool sync.Pool
}

func (p *poolSingle2) Get() *Single2 {
	x, ok := p.pool.Get().(*Single2)
	if !ok {
		x = &Single2{}
	}

	return x
}

func (p *poolSingle2) Put(x *Single2) {
	if x == nil {
		return
	}

	x.ID = 0
	x.ShardKey = 0
	x.P1 = ""
	x.P2 = x.P2[:0]
	x.P5 = 0

	p.pool.Put(x)
}

var PoolSingle2 = poolSingle2{}

func (x *Single2) DeepCopy(z *Single2) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
}

func (x *Single2) Clone() *Single2 {
	z := &Single2{}
	x.DeepCopy(z)
	return z
}

func (x *Single2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Single2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Single2) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Single2) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapSingle2(e registry.Envelope) (proto.Message, error) {
	x := &Single2{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Single2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Single2, x)
}

// register constructors of the messages to the registry package
func init() {
	registry.Register(4833170411250891328, "Single1", unwrapSingle1)
	registry.Register(4738594819076110912, "Single2", unwrapSingle2)

}

var _ = bytes.MinRead
