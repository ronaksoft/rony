// Code generated by Rony's protoc plugin; DO NOT EDIT.

package model

import (
	badger "github.com/dgraph-io/badger/v2"
	edge "github.com/ronaksoft/rony/edge"
	registry "github.com/ronaksoft/rony/registry"
	kv "github.com/ronaksoft/rony/repo/kv"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

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
	x.Enum = 0
	p.pool.Put(x)
}

var PoolModel1 = poolModel1{}

func (x *Model1) DeepCopy(z *Model1) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
	z.Enum = x.Enum
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
	registry.RegisterConstructor(2074613123, "Model1")
	registry.RegisterConstructor(3802219577, "Model2")
}

func SaveModel1WithTxn(txn *badger.Txn, alloc *kv.Allocator, m *Model1) (err error) {
	if alloc == nil {
		alloc = kv.NewAllocator()
		defer alloc.ReleaseAll()
	}

	// Try to read old value
	om := &Model1{}
	om, err = ReadModel1WithTxn(txn, alloc, m.ID, m.ShardKey, om)
	if err != nil && err != badger.ErrKeyNotFound {
		return
	}

	if om != nil {
		// update field index by deleting old values
		if om.P1 != m.P1 {
			err = txn.Delete(alloc.GenKey("IDX", C_Model1, 2864467857, om.P1, om.ID, om.ShardKey))
			if err != nil {
				return
			}
		}

		// update field index by deleting old values
		for i := 0; i < len(om.P2); i++ {
			found := false
			for j := 0; j < len(m.P2); j++ {
				if om.P2[i] == m.P2[j] {
					found = true
					break
				}
			}
			if !found {
				err = txn.Delete(alloc.GenKey("IDX", C_Model1, 867507755, om.P2[i], om.ID, om.ShardKey))
				if err != nil {
					return
				}
			}
		}

		// update field index by deleting old values
		if om.Enum != m.Enum {
			err = txn.Delete(alloc.GenKey("IDX", C_Model1, 2928410991, om.Enum, om.ID, om.ShardKey))
			if err != nil {
				return
			}
		}

	}

	// save entry
	b := alloc.GenValue(m)
	key := alloc.GenKey(C_Model1, 4018441491, m.ID, m.ShardKey)
	err = txn.Set(key, b)
	if err != nil {
		return
	}

	// save entry for view[Enum ShardKey ID]
	err = txn.Set(alloc.GenKey(C_Model1, 2535881670, m.Enum, m.ShardKey, m.ID), b)
	if err != nil {
		return
	}

	// update field index by saving new values
	err = txn.Set(alloc.GenKey("IDX", C_Model1, 2864467857, m.P1, m.ID, m.ShardKey), key)
	if err != nil {
		return
	}

	// update field index by saving new values
	for idx := range m.P2 {
		err = txn.Set(alloc.GenKey("IDX", C_Model1, 867507755, m.P2[idx], m.ID, m.ShardKey), key)
		if err != nil {
			return
		}
	}

	// update field index by saving new values
	err = txn.Set(alloc.GenKey("IDX", C_Model1, 2928410991, m.Enum, m.ID, m.ShardKey), key)
	if err != nil {
		return
	}

	return

}

func SaveModel1(m *Model1) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		return SaveModel1WithTxn(txn, alloc, m)
	})
}

func ReadModel1WithTxn(txn *badger.Txn, alloc *kv.Allocator, id int32, shardKey int32, m *Model1) (*Model1, error) {
	if alloc == nil {
		alloc = kv.NewAllocator()
		defer alloc.ReleaseAll()
	}

	item, err := txn.Get(alloc.GenKey(C_Model1, 4018441491, id, shardKey))
	if err != nil {
		return nil, err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	return m, err
}

func ReadModel1(id int32, shardKey int32, m *Model1) (*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model1{}
	}

	err := kv.View(func(txn *badger.Txn) (err error) {
		m, err = ReadModel1WithTxn(txn, alloc, id, shardKey, m)
		return err
	})
	return m, err
}

func ReadModel1ByEnumAndShardKeyAndIDWithTxn(txn *badger.Txn, alloc *kv.Allocator, enum Enum, shardKey int32, id int32, m *Model1) (*Model1, error) {
	if alloc == nil {
		alloc = kv.NewAllocator()
		defer alloc.ReleaseAll()
	}

	item, err := txn.Get(alloc.GenKey(C_Model1, 2535881670, enum, shardKey, id))
	if err != nil {
		return nil, err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	return m, err
}

func ReadModel1ByEnumAndShardKeyAndID(enum Enum, shardKey int32, id int32, m *Model1) (*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Model1{}
	}
	err := kv.View(func(txn *badger.Txn) (err error) {
		m, err = ReadModel1ByEnumAndShardKeyAndIDWithTxn(txn, alloc, enum, shardKey, id, m)
		return err
	})
	return m, err
}

func DeleteModel1WithTxn(txn *badger.Txn, alloc *kv.Allocator, id int32, shardKey int32) error {
	m := &Model1{}
	item, err := txn.Get(alloc.GenKey(C_Model1, 4018441491, id, shardKey))
	if err != nil {
		return err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	if err != nil {
		return err
	}
	err = txn.Delete(alloc.GenKey(C_Model1, 4018441491, m.ID, m.ShardKey))
	if err != nil {
		return err
	}

	err = txn.Delete(alloc.GenKey(C_Model1, 2535881670, m.Enum, m.ShardKey, m.ID))
	if err != nil {
		return err
	}

	return nil
}

func DeleteModel1(id int32, shardKey int32) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	return kv.Update(func(txn *badger.Txn) error {
		return DeleteModel1WithTxn(txn, alloc, id, shardKey)
	})
}

func ListModel1(
	offsetID int32, offsetShardKey int32, lo *kv.ListOption, cond func(m *Model1) bool,
) ([]*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey(C_Model1, 4018441491)
		opt.Reverse = lo.Backward()
		osk := alloc.GenKey(C_Model1, 4018441491, offsetID)
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
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

func (x *Model1) HasP2(xx string) bool {
	for idx := range x.P2 {
		if x.P2[idx] == xx {
			return true
		}
	}
	return false
}

func ListModel1ByP1(p1 string, lo *kv.ListOption) ([]*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey("IDX", C_Model1, 2864467857, p1)
		opt.Reverse = lo.Backward()
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
		for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			_ = iter.Item().Value(func(val []byte) error {
				item, err := txn.Get(val)
				if err != nil {
					return err
				}
				return item.Value(func(val []byte) error {
					m := &Model1{}
					err := m.Unmarshal(val)
					if err != nil {
						return err
					}
					res = append(res, m)
					return nil
				})
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func ListModel1ByP2(p2 string, lo *kv.ListOption) ([]*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey("IDX", C_Model1, 867507755, p2)
		opt.Reverse = lo.Backward()
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
		for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			_ = iter.Item().Value(func(val []byte) error {
				item, err := txn.Get(val)
				if err != nil {
					return err
				}
				return item.Value(func(val []byte) error {
					m := &Model1{}
					err := m.Unmarshal(val)
					if err != nil {
						return err
					}
					res = append(res, m)
					return nil
				})
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func ListModel1ByEnum(enum Enum, lo *kv.ListOption) ([]*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey("IDX", C_Model1, 2928410991, enum)
		opt.Reverse = lo.Backward()
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
		for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			_ = iter.Item().Value(func(val []byte) error {
				item, err := txn.Get(val)
				if err != nil {
					return err
				}
				return item.Value(func(val []byte) error {
					m := &Model1{}
					err := m.Unmarshal(val)
					if err != nil {
						return err
					}
					res = append(res, m)
					return nil
				})
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func ListModel1ByID(id int32, offsetShardKey int32, lo *kv.ListOption) ([]*Model1, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey(C_Model1, 4018441491, id)
		opt.Reverse = lo.Backward()
		osk := alloc.GenKey(C_Model1, 4018441491, id, offsetShardKey)
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
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
				res = append(res, m)
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func SaveModel2WithTxn(txn *badger.Txn, alloc *kv.Allocator, m *Model2) (err error) {
	if alloc == nil {
		alloc = kv.NewAllocator()
		defer alloc.ReleaseAll()
	}

	// Try to read old value
	om := &Model2{}
	om, err = ReadModel2WithTxn(txn, alloc, m.ID, m.ShardKey, m.P1, om)
	if err != nil && err != badger.ErrKeyNotFound {
		return
	}

	if om != nil {
		// update field index by deleting old values
		if om.P1 != m.P1 {
			err = txn.Delete(alloc.GenKey("IDX", C_Model2, 2864467857, om.P1, om.ID, om.ShardKey, om.P1))
			if err != nil {
				return
			}
		}

	}

	// save entry
	b := alloc.GenValue(m)
	key := alloc.GenKey(C_Model2, 1609271041, m.ID, m.ShardKey, m.P1)
	err = txn.Set(key, b)
	if err != nil {
		return
	}

	// save entry for view[P1 ShardKey ID]
	err = txn.Set(alloc.GenKey(C_Model2, 2344331025, m.P1, m.ShardKey, m.ID), b)
	if err != nil {
		return
	}

	// update field index by saving new values
	err = txn.Set(alloc.GenKey("IDX", C_Model2, 2864467857, m.P1, m.ID, m.ShardKey, m.P1), key)
	if err != nil {
		return
	}

	return

}

func SaveModel2(m *Model2) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	return kv.Update(func(txn *badger.Txn) error {
		return SaveModel2WithTxn(txn, alloc, m)
	})
}

func ReadModel2WithTxn(txn *badger.Txn, alloc *kv.Allocator, id int64, shardKey int32, p1 string, m *Model2) (*Model2, error) {
	if alloc == nil {
		alloc = kv.NewAllocator()
		defer alloc.ReleaseAll()
	}

	item, err := txn.Get(alloc.GenKey(C_Model2, 1609271041, id, shardKey, p1))
	if err != nil {
		return nil, err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	return m, err
}

func ReadModel2(id int64, shardKey int32, p1 string, m *Model2) (*Model2, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model2{}
	}

	err := kv.View(func(txn *badger.Txn) (err error) {
		m, err = ReadModel2WithTxn(txn, alloc, id, shardKey, p1, m)
		return err
	})
	return m, err
}

func ReadModel2ByP1AndShardKeyAndIDWithTxn(txn *badger.Txn, alloc *kv.Allocator, p1 string, shardKey int32, id int64, m *Model2) (*Model2, error) {
	if alloc == nil {
		alloc = kv.NewAllocator()
		defer alloc.ReleaseAll()
	}

	item, err := txn.Get(alloc.GenKey(C_Model2, 2344331025, p1, shardKey, id))
	if err != nil {
		return nil, err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	return m, err
}

func ReadModel2ByP1AndShardKeyAndID(p1 string, shardKey int32, id int64, m *Model2) (*Model2, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Model2{}
	}
	err := kv.View(func(txn *badger.Txn) (err error) {
		m, err = ReadModel2ByP1AndShardKeyAndIDWithTxn(txn, alloc, p1, shardKey, id, m)
		return err
	})
	return m, err
}

func DeleteModel2WithTxn(txn *badger.Txn, alloc *kv.Allocator, id int64, shardKey int32, p1 string) error {
	m := &Model2{}
	item, err := txn.Get(alloc.GenKey(C_Model2, 1609271041, id, shardKey, p1))
	if err != nil {
		return err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	if err != nil {
		return err
	}
	err = txn.Delete(alloc.GenKey(C_Model2, 1609271041, m.ID, m.ShardKey, m.P1))
	if err != nil {
		return err
	}

	err = txn.Delete(alloc.GenKey(C_Model2, 2344331025, m.P1, m.ShardKey, m.ID))
	if err != nil {
		return err
	}

	return nil
}

func DeleteModel2(id int64, shardKey int32, p1 string) error {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	return kv.Update(func(txn *badger.Txn) error {
		return DeleteModel2WithTxn(txn, alloc, id, shardKey, p1)
	})
}

func ListModel2(
	offsetID int64, offsetShardKey int32, offsetP1 string, lo *kv.ListOption, cond func(m *Model2) bool,
) ([]*Model2, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model2, 0, lo.Limit())
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey(C_Model2, 1609271041)
		opt.Reverse = lo.Backward()
		osk := alloc.GenKey(C_Model2, 1609271041, offsetID, offsetShardKey)
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
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

func (x *Model2) HasP2(xx string) bool {
	for idx := range x.P2 {
		if x.P2[idx] == xx {
			return true
		}
	}
	return false
}

func ListModel2ByP1(p1 string, lo *kv.ListOption) ([]*Model2, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model2, 0, lo.Limit())
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey("IDX", C_Model2, 2864467857, p1)
		opt.Reverse = lo.Backward()
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
		for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			_ = iter.Item().Value(func(val []byte) error {
				item, err := txn.Get(val)
				if err != nil {
					return err
				}
				return item.Value(func(val []byte) error {
					m := &Model2{}
					err := m.Unmarshal(val)
					if err != nil {
						return err
					}
					res = append(res, m)
					return nil
				})
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func ListModel2ByIDAndShardKey(id int64, shardKey int32, offsetP1 string, lo *kv.ListOption) ([]*Model2, error) {
	alloc := kv.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model2, 0, lo.Limit())
	err := kv.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey(C_Model2, 1609271041, id, shardKey)
		opt.Reverse = lo.Backward()
		osk := alloc.GenKey(C_Model2, 1609271041, id, shardKey, offsetP1)
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
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
				res = append(res, m)
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}
