// Code generated by Rony's protoc plugin; DO NOT EDIT.

package badgerStore

import (
	edge "github.com/ronaksoft/rony/edge"
	registry "github.com/ronaksoft/rony/registry"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_StoreCommand int64 = 3294788005

type poolStoreCommand struct {
	pool sync.Pool
}

func (p *poolStoreCommand) Get() *StoreCommand {
	x, ok := p.pool.Get().(*StoreCommand)
	if !ok {
		x = &StoreCommand{}
	}
	return x
}

func (p *poolStoreCommand) Put(x *StoreCommand) {
	if x == nil {
		return
	}
	x.Type = 0
	for _, z := range x.KVs {
		PoolKeyValue.Put(z)
	}
	x.KVs = x.KVs[:0]
	p.pool.Put(x)
}

var PoolStoreCommand = poolStoreCommand{}

func (x *StoreCommand) DeepCopy(z *StoreCommand) {
	z.Type = x.Type
	for idx := range x.KVs {
		if x.KVs[idx] != nil {
			xx := PoolKeyValue.Get()
			x.KVs[idx].DeepCopy(xx)
			z.KVs = append(z.KVs, xx)
		}
	}
}

func (x *StoreCommand) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *StoreCommand) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *StoreCommand) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_StoreCommand, x)
}

const C_KeyValue int64 = 4276272820

type poolKeyValue struct {
	pool sync.Pool
}

func (p *poolKeyValue) Get() *KeyValue {
	x, ok := p.pool.Get().(*KeyValue)
	if !ok {
		x = &KeyValue{}
	}
	return x
}

func (p *poolKeyValue) Put(x *KeyValue) {
	if x == nil {
		return
	}
	x.Key = x.Key[:0]
	x.Value = x.Value[:0]
	p.pool.Put(x)
}

var PoolKeyValue = poolKeyValue{}

func (x *KeyValue) DeepCopy(z *KeyValue) {
	z.Key = append(z.Key[:0], x.Key...)
	z.Value = append(z.Value[:0], x.Value...)
}

func (x *KeyValue) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *KeyValue) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *KeyValue) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_KeyValue, x)
}

func init() {
	registry.RegisterConstructor(3294788005, "StoreCommand")
	registry.RegisterConstructor(4276272820, "KeyValue")
}
