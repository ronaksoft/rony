// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.15.8
// Rony ver. v0.12.5
// Source: singleton.proto

package singleton

import (
	bytes "bytes"
	rony "github.com/ronaksoft/rony"
	edge "github.com/ronaksoft/rony/edge"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	store "github.com/ronaksoft/rony/store"
	tools "github.com/ronaksoft/rony/tools"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

var _ = pools.Imported

const C_Single1 int64 = 683727308

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
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
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

func (x *Single1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Single1, x)
}

const C_Single2 int64 = 2982774902

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
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
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

func (x *Single2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Single2, x)
}

func init() {
	registry.RegisterConstructor(683727308, "Single1")
	registry.RegisterConstructor(2982774902, "Single2")
}

var _ = bytes.MinRead

type Single1LocalSingleton struct {
	s rony.Store
}

func NewSingle1LocalSingleton(s rony.Store) *Single1LocalSingleton {
	return &Single1LocalSingleton{
		s: s,
	}
}

func (r *Single1LocalSingleton) SaveWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Single1) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err = store.Marshal(txn, alloc, m, 'S', C_Single1)
	if err != nil {
		return
	}
	return nil
}

func (r *Single1LocalSingleton) Save(m *Single1) (err error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.SaveWithTxn(txn, alloc, m)
	})
}

func (r *Single1LocalSingleton) ReadWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Single1) (*Single1, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'S', C_Single1)
	if err != nil {
		return nil, err
	}
	return m, err
}

func (r *Single1LocalSingleton) Read(m *Single1) (*Single1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Single1{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadWithTxn(txn, alloc, m)
		return
	})
	return m, err
}

type Single2LocalSingleton struct {
	s rony.Store
}

func NewSingle2LocalSingleton(s rony.Store) *Single2LocalSingleton {
	return &Single2LocalSingleton{
		s: s,
	}
}

func (r *Single2LocalSingleton) SaveWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Single2) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err = store.Marshal(txn, alloc, m, 'S', C_Single2)
	if err != nil {
		return
	}
	return nil
}

func (r *Single2LocalSingleton) Save(m *Single2) (err error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.SaveWithTxn(txn, alloc, m)
	})
}

func (r *Single2LocalSingleton) ReadWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Single2) (*Single2, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'S', C_Single2)
	if err != nil {
		return nil, err
	}
	return m, err
}

func (r *Single2LocalSingleton) Read(m *Single2) (*Single2, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Single2{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadWithTxn(txn, alloc, m)
		return
	})
	return m, err
}
