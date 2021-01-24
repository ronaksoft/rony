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

func (x *Hook) DeepCopy(z *Hook) {
	z.ClientID = x.ClientID
	z.ID = x.ID
	z.Timestamp = x.Timestamp
	z.HookUrl = x.HookUrl
	z.Fired = x.Fired
	z.Success = x.Success
}

func (x *Hook) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Hook) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Hook) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Hook, x)
}

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

func (x *Model1) DeepCopy(z *Model1) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
}

func (x *Model1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model1, x)
}

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

func (x *Model2) DeepCopy(z *Model2) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
}

func (x *Model2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model2, x)
}

func init() {
	registry.RegisterConstructor(74116203, "Hook")
	registry.RegisterConstructor(2074613123, "Model1")
	registry.RegisterConstructor(3802219577, "Model2")
}

func SaveHook(m *Hook) error {
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

func ReadHook(clientID string, id string, m *Hook) (*Hook, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Hook{}
	}
	err := kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get(alloc.GenKey(C_Hook, clientID, id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
	})
	return m, err
}

func DeleteHook(clientID string, id string) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		err := txn.Delete(alloc.GenKey(C_Hook, clientID, id))
		if err != nil {
			return err
		}

		return nil
	})
}

func ListHook(
	offsetClientID string, offset int32, limit int32, cond func(m *Hook) bool,
) ([]*Hook, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Hook, 0, limit)
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey(C_Hook)
		osk := alloc.GenKey(C_Hook, offsetClientID)
		iter := txn.NewIterator(opt)
		for iter.Seek(osk); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			_ = iter.Item().Value(func(val []byte) error {
				m := &Hook{}
				err := m.Unmarshal(val)
				if err != nil {
					return err
				}
				if cond(m) {
					res = append(res, m)
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}
func SaveModel1(m *Model1) error {
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

func ReadModel1(id int32, shardKey int32, m *Model1) (*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Model1{}
	}
	err := kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get(alloc.GenKey(C_Model1, id, shardKey))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
	})
	return m, err
}

func ReadModel1ByShardKeyAndID(shardKey int32, id int32, m *Model1) (*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Model1{}
	}
	err := kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get(alloc.GenKey(C_Model1, 3512206964, m.ShardKey, m.ID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
	})
	return m, err
}

func DeleteModel1(id int32, shardKey int32) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		m := &Model1{}
		item, err := txn.Get(alloc.GenKey(C_Model1, id, shardKey))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
		if err != nil {
			return err
		}
		err = txn.Delete(alloc.GenKey(C_Model1, id, shardKey))
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

func ListModel1(
	offsetID int32, offset int32, limit int32, cond func(m *Model1) bool,
) ([]*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, limit)
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey(C_Model1)
		osk := alloc.GenKey(C_Model1, offsetID)
		iter := txn.NewIterator(opt)
		for iter.Seek(osk); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			_ = iter.Item().Value(func(val []byte) error {
				m := &Model1{}
				err := m.Unmarshal(val)
				if err != nil {
					return err
				}
				if cond(m) {
					res = append(res, m)
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}
func SaveModel2(m *Model2) error {
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

func ReadModel2(id int64, shardKey int32, p1 string, m *Model2) (*Model2, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Model2{}
	}
	err := kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get(alloc.GenKey(C_Model2, id, shardKey, p1))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
	})
	return m, err
}

func ReadModel2ByP1AndShardKeyAndID(p1 string, shardKey int32, id int64, m *Model2) (*Model2, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Model2{}
	}
	err := kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get(alloc.GenKey(C_Model2, 2344331025, m.P1, m.ShardKey, m.ID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
	})
	return m, err
}

func DeleteModel2(id int64, shardKey int32, p1 string) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		m := &Model2{}
		item, err := txn.Get(alloc.GenKey(C_Model2, id, shardKey, p1))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return m.Unmarshal(val)
		})
		if err != nil {
			return err
		}
		err = txn.Delete(alloc.GenKey(C_Model2, id, shardKey, p1))
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

func ListModel2(
	offsetID int64, offsetShardKey int32, offset int32, limit int32, cond func(m *Model2) bool,
) ([]*Model2, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model2, 0, limit)
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey(C_Model2)
		osk := alloc.GenKey(C_Model2, offsetID, offsetShardKey)
		iter := txn.NewIterator(opt)
		for iter.Seek(osk); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			_ = iter.Item().Value(func(val []byte) error {
				m := &Model2{}
				err := m.Unmarshal(val)
				if err != nil {
					return err
				}
				if cond(m) {
					res = append(res, m)
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}
