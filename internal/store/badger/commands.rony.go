// Code generated by Rony's protoc plugin; DO NOT EDIT.

package badgerStore

import (
	edge "github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/pools"
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
	x.Payload = x.Payload[:0]
	p.pool.Put(x)
}

var PoolStoreCommand = poolStoreCommand{}

func (x *StoreCommand) DeepCopy(z *StoreCommand) {
	z.Type = x.Type
	z.Payload = append(z.Payload[:0], x.Payload...)
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

const C_StartTxn int64 = 605208098

type poolStartTxn struct {
	pool sync.Pool
}

func (p *poolStartTxn) Get() *StartTxn {
	x, ok := p.pool.Get().(*StartTxn)
	if !ok {
		x = &StartTxn{}
	}
	return x
}

func (p *poolStartTxn) Put(x *StartTxn) {
	if x == nil {
		return
	}
	x.ID = 0
	x.Update = false
	p.pool.Put(x)
}

var PoolStartTxn = poolStartTxn{}

func (x *StartTxn) DeepCopy(z *StartTxn) {
	z.ID = x.ID
	z.Update = x.Update
}

func (x *StartTxn) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *StartTxn) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *StartTxn) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_StartTxn, x)
}

const C_StopTxn int64 = 1239816782

type poolStopTxn struct {
	pool sync.Pool
}

func (p *poolStopTxn) Get() *StopTxn {
	x, ok := p.pool.Get().(*StopTxn)
	if !ok {
		x = &StopTxn{}
	}
	return x
}

func (p *poolStopTxn) Put(x *StopTxn) {
	if x == nil {
		return
	}
	x.ID = 0
	x.Commit = false
	p.pool.Put(x)
}

var PoolStopTxn = poolStopTxn{}

func (x *StopTxn) DeepCopy(z *StopTxn) {
	z.ID = x.ID
	z.Commit = x.Commit
}

func (x *StopTxn) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *StopTxn) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *StopTxn) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_StopTxn, x)
}

const C_CommitTxn int64 = 15774688

type poolCommitTxn struct {
	pool sync.Pool
}

func (p *poolCommitTxn) Get() *CommitTxn {
	x, ok := p.pool.Get().(*CommitTxn)
	if !ok {
		x = &CommitTxn{}
	}
	return x
}

func (p *poolCommitTxn) Put(x *CommitTxn) {
	if x == nil {
		return
	}
	x.ID = 0
	p.pool.Put(x)
}

var PoolCommitTxn = poolCommitTxn{}

func (x *CommitTxn) DeepCopy(z *CommitTxn) {
	z.ID = x.ID
}

func (x *CommitTxn) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *CommitTxn) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *CommitTxn) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_CommitTxn, x)
}

const C_Set int64 = 3730400060

type poolSet struct {
	pool sync.Pool
}

func (p *poolSet) Get() *Set {
	x, ok := p.pool.Get().(*Set)
	if !ok {
		x = &Set{}
	}
	return x
}

func (p *poolSet) Put(x *Set) {
	if x == nil {
		return
	}
	x.TxnID = 0
	for _, z := range x.KVs {
		PoolKeyValue.Put(z)
	}
	x.KVs = x.KVs[:0]
	p.pool.Put(x)
}

var PoolSet = poolSet{}

func (x *Set) DeepCopy(z *Set) {
	z.TxnID = x.TxnID
	for idx := range x.KVs {
		if x.KVs[idx] != nil {
			xx := PoolKeyValue.Get()
			x.KVs[idx].DeepCopy(xx)
			z.KVs = append(z.KVs, xx)
		}
	}
}

func (x *Set) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Set) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Set) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Set, x)
}

const C_Delete int64 = 1035893169

type poolDelete struct {
	pool sync.Pool
}

func (p *poolDelete) Get() *Delete {
	x, ok := p.pool.Get().(*Delete)
	if !ok {
		x = &Delete{}
	}
	return x
}

func (p *poolDelete) Put(x *Delete) {
	if x == nil {
		return
	}
	x.TxnID = 0
	for _, z := range x.Keys {
		pools.Bytes.Put(z)
	}
	x.Keys = x.Keys[:0]
	p.pool.Put(x)
}

var PoolDelete = poolDelete{}

func (x *Delete) DeepCopy(z *Delete) {
	z.TxnID = x.TxnID
	z.Keys = append(z.Keys[:0], x.Keys...)
}

func (x *Delete) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Delete) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Delete) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Delete, x)
}

const C_Get int64 = 3312871568

type poolGet struct {
	pool sync.Pool
}

func (p *poolGet) Get() *Get {
	x, ok := p.pool.Get().(*Get)
	if !ok {
		x = &Get{}
	}
	return x
}

func (p *poolGet) Put(x *Get) {
	if x == nil {
		return
	}
	x.TxnID = 0
	for _, z := range x.Keys {
		pools.Bytes.Put(z)
	}
	x.Keys = x.Keys[:0]
	p.pool.Put(x)
}

var PoolGet = poolGet{}

func (x *Get) DeepCopy(z *Get) {
	z.TxnID = x.TxnID
	z.Keys = append(z.Keys[:0], x.Keys...)
}

func (x *Get) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Get) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Get) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Get, x)
}

func init() {
	registry.RegisterConstructor(3294788005, "StoreCommand")
	registry.RegisterConstructor(4276272820, "KeyValue")
	registry.RegisterConstructor(605208098, "StartTxn")
	registry.RegisterConstructor(1239816782, "StopTxn")
	registry.RegisterConstructor(15774688, "CommitTxn")
	registry.RegisterConstructor(3730400060, "Set")
	registry.RegisterConstructor(1035893169, "Delete")
	registry.RegisterConstructor(3312871568, "Get")
}
