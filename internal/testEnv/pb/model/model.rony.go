package model

import (
	badger "github.com/dgraph-io/badger/v2"
	edge "github.com/ronaksoft/rony/edge"
	registry "github.com/ronaksoft/rony/registry"
	kv "github.com/ronaksoft/rony/repo/kv"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_Hook int64 = 74116203

type poolHook struct {
	pool sync.Pool
}

func (p *poolHook) Get() *Hook {
	x, ok := p.pool.Get().(*Hook)
	if !ok {
		return &Hook{}
	}
	return x
}

func (p *poolHook) Put(x *Hook) {
	x.ClientID = ""
	x.ID = ""
	x.Timestamp = ""
	x.HookUrl = ""
	x.Fired = false
	x.Success = false
	p.pool.Put(x)
}

var PoolHook = poolHook{}

const C_Model1 int64 = 2074613123

type poolModel1 struct {
	pool sync.Pool
}

func (p *poolModel1) Get() *Model1 {
	x, ok := p.pool.Get().(*Model1)
	if !ok {
		return &Model1{}
	}
	return x
}

func (p *poolModel1) Put(x *Model1) {
	x.ID = 0
	x.ShardKey = 0
	x.P1 = ""
	x.P2 = x.P2[:0]
	x.P5 = 0
	p.pool.Put(x)
}

var PoolModel1 = poolModel1{}

const C_Model2 int64 = 3802219577

type poolModel2 struct {
	pool sync.Pool
}

func (p *poolModel2) Get() *Model2 {
	x, ok := p.pool.Get().(*Model2)
	if !ok {
		return &Model2{}
	}
	return x
}

func (p *poolModel2) Put(x *Model2) {
	x.ID = 0
	x.ShardKey = 0
	x.P1 = ""
	x.P2 = x.P2[:0]
	x.P5 = 0
	p.pool.Put(x)
}

var PoolModel2 = poolModel2{}

func init() {
	registry.RegisterConstructor(74116203, "Hook")
	registry.RegisterConstructor(2074613123, "Model1")
	registry.RegisterConstructor(3802219577, "Model2")
}

func (x *Hook) DeepCopy(z *Hook) {
	z.ClientID = x.ClientID
	z.ID = x.ID
	z.Timestamp = x.Timestamp
	z.HookUrl = x.HookUrl
	z.Fired = x.Fired
	z.Success = x.Success
}

func (x *Model1) DeepCopy(z *Model1) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
}

func (x *Model2) DeepCopy(z *Model2) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
}

func (x *Hook) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Hook, x)
}

func (x *Model1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model1, x)
}

func (x *Model2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model2, x)
}

func (x *Hook) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Hook) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func CreateHook(m *Hook) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		b := alloc.GenValue(m)
		err := txn.Set(alloc.GenKey(C_Hook, m.ClientID, m.ID), b)
		if err != nil {
			return err
		}

		return nil
	})
}

func ReadHook(m *Hook) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get(alloc.GenKey(C_Hook, m.ClientID, m.ID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
	})
}

func DeleteHook(m *Hook) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		err := txn.Delete(alloc.GenKey(C_Hook, m.ClientID, m.ID))
		if err != nil {
			return err
		}

		return nil
	})
}

func CreateModel1(m *Model1) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		b := alloc.GenValue(m)
		err := txn.Set(alloc.GenKey(C_Model1, m.ID, m.ShardKey), b)
		if err != nil {
			return err
		}

		err = txn.Set(alloc.GenKey(C_Model1, 3113976552, m.ShardKey, m.ID), b)
		if err != nil {
			return err
		}

		return nil
	})
}

func ReadModel1(m *Model1) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get(alloc.GenKey(C_Model1, m.ID, m.ShardKey))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
	})
}

func DeleteModel1(m *Model1) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		err := txn.Delete(alloc.GenKey(C_Model1, m.ID, m.ShardKey))
		if err != nil {
			return err
		}

		err = txn.Delete(alloc.GenKey(C_Model1, 3113976552, m.ShardKey, m.ID))
		if err != nil {
			return err
		}

		return nil
	})
}

func CreateModel2(m *Model2) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		b := alloc.GenValue(m)
		err := txn.Set(alloc.GenKey(C_Model2, m.ID, m.ShardKey, m.P1), b)
		if err != nil {
			return err
		}

		err = txn.Set(alloc.GenKey(C_Model2, 3495323833, m.P1, m.ShardKey, m.ID), b)
		if err != nil {
			return err
		}

		return nil
	})
}

func ReadModel2(m *Model2) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get(alloc.GenKey(C_Model2, m.ID, m.ShardKey, m.P1))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
	})
}

func DeleteModel2(m *Model2) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		err := txn.Delete(alloc.GenKey(C_Model2, m.ID, m.ShardKey, m.P1))
		if err != nil {
			return err
		}

		err = txn.Delete(alloc.GenKey(C_Model2, 3495323833, m.P1, m.ShardKey, m.ID))
		if err != nil {
			return err
		}

		return nil
	})
}
