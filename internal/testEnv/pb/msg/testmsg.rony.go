// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.16.0
// Source: testmsg.proto

package testmsg

import (
	bytes "bytes"
	sync "sync"

	edge "github.com/ronaksoft/rony/edge"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _ = pools.Imported

const C_Envelope1 uint64 = 12704881787996916493

type poolEnvelope1 struct {
	pool sync.Pool
}

func (p *poolEnvelope1) Get() *Envelope1 {
	x, ok := p.pool.Get().(*Envelope1)
	if !ok {
		x = &Envelope1{}
	}

	x.Embed1 = PoolEmbed1.Get()

	return x
}

func (p *poolEnvelope1) Put(x *Envelope1) {
	if x == nil {
		return
	}

	x.E1 = ""
	PoolEmbed1.Put(x.Embed1)

	p.pool.Put(x)
}

var PoolEnvelope1 = poolEnvelope1{}

func (x *Envelope1) DeepCopy(z *Envelope1) {
	z.E1 = x.E1
	if x.Embed1 != nil {
		if z.Embed1 == nil {
			z.Embed1 = PoolEmbed1.Get()
		}
		x.Embed1.DeepCopy(z.Embed1)
	} else {
		PoolEmbed1.Put(z.Embed1)
		z.Embed1 = nil
	}
}

func (x *Envelope1) Clone() *Envelope1 {
	z := &Envelope1{}
	x.DeepCopy(z)
	return z
}

func (x *Envelope1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Envelope1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Envelope1) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Envelope1) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapEnvelope1(e registry.Envelope) (registry.Message, error) {
	x := &Envelope1{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Envelope1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Envelope1, x)
}

const C_Envelope2 uint64 = 12862507774954883853

type poolEnvelope2 struct {
	pool sync.Pool
}

func (p *poolEnvelope2) Get() *Envelope2 {
	x, ok := p.pool.Get().(*Envelope2)
	if !ok {
		x = &Envelope2{}
	}

	return x
}

func (p *poolEnvelope2) Put(x *Envelope2) {
	if x == nil {
		return
	}

	x.Constructor = 0
	x.Message = x.Message[:0]

	p.pool.Put(x)
}

var PoolEnvelope2 = poolEnvelope2{}

func (x *Envelope2) DeepCopy(z *Envelope2) {
	z.Constructor = x.Constructor
	z.Message = append(z.Message[:0], x.Message...)
}

func (x *Envelope2) Clone() *Envelope2 {
	z := &Envelope2{}
	x.DeepCopy(z)
	return z
}

func (x *Envelope2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Envelope2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Envelope2) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Envelope2) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapEnvelope2(e registry.Envelope) (registry.Message, error) {
	x := &Envelope2{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Envelope2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Envelope2, x)
}

const C_Embed1 uint64 = 4833315499840692224

type poolEmbed1 struct {
	pool sync.Pool
}

func (p *poolEmbed1) Get() *Embed1 {
	x, ok := p.pool.Get().(*Embed1)
	if !ok {
		x = &Embed1{}
	}

	return x
}

func (p *poolEmbed1) Put(x *Embed1) {
	if x == nil {
		return
	}

	x.F1 = ""

	p.pool.Put(x)
}

var PoolEmbed1 = poolEmbed1{}

func (x *Embed1) DeepCopy(z *Embed1) {
	z.F1 = x.F1
}

func (x *Embed1) Clone() *Embed1 {
	z := &Embed1{}
	x.DeepCopy(z)
	return z
}

func (x *Embed1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Embed1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Embed1) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Embed1) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapEmbed1(e registry.Envelope) (registry.Message, error) {
	x := &Embed1{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Embed1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Embed1, x)
}

const C_Embed2 uint64 = 4738739907665911808

type poolEmbed2 struct {
	pool sync.Pool
}

func (p *poolEmbed2) Get() *Embed2 {
	x, ok := p.pool.Get().(*Embed2)
	if !ok {
		x = &Embed2{}
	}

	return x
}

func (p *poolEmbed2) Put(x *Embed2) {
	if x == nil {
		return
	}

	x.F2 = ""

	p.pool.Put(x)
}

var PoolEmbed2 = poolEmbed2{}

func (x *Embed2) DeepCopy(z *Embed2) {
	z.F2 = x.F2
}

func (x *Embed2) Clone() *Embed2 {
	z := &Embed2{}
	x.DeepCopy(z)
	return z
}

func (x *Embed2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Embed2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Embed2) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Embed2) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapEmbed2(e registry.Envelope) (registry.Message, error) {
	x := &Embed2{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Embed2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Embed2, x)
}

// register constructors of the messages to the registry package
func init() {
	registry.Register(12704881787996916493, "Envelope1", unwrapEnvelope1)
	registry.Register(12862507774954883853, "Envelope2", unwrapEnvelope2)
	registry.Register(4833315499840692224, "Embed1", unwrapEmbed1)
	registry.Register(4738739907665911808, "Embed2", unwrapEmbed2)

}

var _ = bytes.MinRead
