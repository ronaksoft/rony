// Code generated by Rony's protoc plugin; DO NOT EDIT.

package model

import (
	edge "github.com/ronaksoft/rony/edge"
	registry "github.com/ronaksoft/rony/registry"
	store "github.com/ronaksoft/rony/store"
	tools "github.com/ronaksoft/rony/tools"
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
		x = &Model1{}
	}
	return x
}

func (p *poolModel1) Put(x *Model1) {
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
		x = &Model2{}
	}
	return x
}

func (p *poolModel2) Put(x *Model2) {
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

type Model1Order string

const Model1OrderByEnum Model1Order = "Enum"

func CreateModel1(m *Model1) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *store.LTxn) error {
		return CreateModel1WithTxn(txn, alloc, m)
	})
}

func CreateModel1WithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Model1) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	if store.Exists(txn, alloc, 'M', C_Model1, 4018441491, m.ID, m.ShardKey) {
		return store.ErrAlreadyExists
	}
	// save entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, 'M', C_Model1, 4018441491, m.ID, m.ShardKey)
	if err != nil {
		return
	}

	// save views
	// save entry for view: [Enum ShardKey ID]
	err = store.Set(txn, alloc, val, 'M', C_Model1, 2535881670, m.Enum, m.ShardKey, m.ID)
	if err != nil {
		return
	}

	// save indices
	key := alloc.Gen('M', C_Model1, 4018441491, m.ID, m.ShardKey)
	// update field index by saving new values
	err = store.Set(txn, alloc, key, 'I', C_Model1, uint64(4843779728911368192), m.P1, m.ID, m.ShardKey)
	if err != nil {
		return
	}

	// update field index by saving new values
	for idx := range m.P2 {
		err = store.Set(txn, alloc, key, 'I', C_Model1, uint64(4749204136736587776), m.P2[idx], m.ID, m.ShardKey)
		if err != nil {
			return
		}
	}

	return

}

func ReadModel1WithTxn(txn *store.LTxn, alloc *tools.Allocator, id int32, shardKey int32, m *Model1) (*Model1, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model1, 4018441491, id, shardKey)
	if err != nil {
		return nil, err
	}
	return m, err
}

func ReadModel1(id int32, shardKey int32, m *Model1) (*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model1{}
	}

	err := store.View(func(txn *store.LTxn) (err error) {
		m, err = ReadModel1WithTxn(txn, alloc, id, shardKey, m)
		return err
	})
	return m, err
}

func ReadModel1ByEnumAndShardKeyAndIDWithTxn(txn *store.LTxn, alloc *tools.Allocator, enum Enum, shardKey int32, id int32, m *Model1) (*Model1, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model1, 2535881670, enum, shardKey, id)
	if err != nil {
		return nil, err
	}
	return m, err
}

func ReadModel1ByEnumAndShardKeyAndID(enum Enum, shardKey int32, id int32, m *Model1) (*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Model1{}
	}
	err := store.View(func(txn *store.LTxn) (err error) {
		m, err = ReadModel1ByEnumAndShardKeyAndIDWithTxn(txn, alloc, enum, shardKey, id, m)
		return err
	})
	return m, err
}

func UpdateModel1WithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Model1) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := DeleteModel1WithTxn(txn, alloc, m.ID, m.ShardKey)
	if err != nil {
		return err
	}

	return CreateModel1WithTxn(txn, alloc, m)
}

func UpdateModel1(id int32, shardKey int32, m *Model1) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		return store.ErrEmptyObject
	}

	err := store.View(func(txn *store.LTxn) (err error) {
		return UpdateModel1WithTxn(txn, alloc, m)
	})
	return err
}

func DeleteModel1WithTxn(txn *store.LTxn, alloc *tools.Allocator, id int32, shardKey int32) error {
	m := &Model1{}
	err := store.Unmarshal(txn, alloc, m, 'M', C_Model1, 4018441491, id, shardKey)
	if err != nil {
		return err
	}
	err = store.Delete(txn, alloc, 'M', C_Model1, 4018441491, m.ID, m.ShardKey)
	if err != nil {
		return err
	}

	// delete field index
	err = store.Delete(txn, alloc, 'I', C_Model1, uint64(4843779728911368192), m.P1, m.ID, m.ShardKey)
	if err != nil {
		return err
	}

	// delete field index
	for idx := range m.P2 {
		err = store.Delete(txn, alloc, 'I', C_Model1, uint64(4749204136736587776), m.P2[idx], m.ID, m.ShardKey)
		if err != nil {
			return err
		}
	}

	err = store.Delete(txn, alloc, 'M', C_Model1, 2535881670, m.Enum, m.ShardKey, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteModel1(id int32, shardKey int32) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *store.LTxn) error {
		return DeleteModel1WithTxn(txn, alloc, id, shardKey)
	})
}

func (x *Model1) HasP2(xx string) bool {
	for idx := range x.P2 {
		if x.P2[idx] == xx {
			return true
		}
	}
	return false
}

func SaveModel1WithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Model1) (err error) {
	if store.Exists(txn, alloc, 'M', C_Model1, 4018441491, m.ID, m.ShardKey) {
		return UpdateModel1WithTxn(txn, alloc, m)
	} else {
		return CreateModel1WithTxn(txn, alloc, m)
	}
}

func SaveModel1(m *Model1) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *store.LTxn) error {
		return SaveModel1WithTxn(txn, alloc, m)
	})
}

func IterModel1(txn *store.LTxn, alloc *tools.Allocator, cb func(m *Model1) bool, orderBy ...Model1Order) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	iterOpt := store.DefaultIteratorOptions
	if len(orderBy) == 0 {
		iterOpt.Prefix = alloc.Gen('M', C_Model1, 4018441491)
	} else {
		switch orderBy[0] {
		case Model1OrderByEnum:
			iterOpt.Prefix = alloc.Gen('M', C_Model1, 2535881670)
		}
	}
	iter := txn.NewIterator(iterOpt)
	for iter.Rewind(); iter.ValidForPrefix(iterOpt.Prefix); iter.Next() {
		_ = iter.Item().Value(func(val []byte) error {
			m := &Model1{}
			err := m.Unmarshal(val)
			if err != nil {
				return err
			}
			if !cb(m) {
				exitLoop = true
			}
			return nil
		})
		if exitLoop {
			break
		}
	}
	iter.Close()
	return nil
}

func ListModel1(
	offsetID int32, offsetShardKey int32, lo *store.ListOption, cond func(m *Model1) bool,
) ([]*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('M', C_Model1, 4018441491)
		opt.Reverse = lo.Backward()
		osk := alloc.Gen('M', C_Model1, 4018441491, offsetID)
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
				if cond == nil || cond(m) {
					res = append(res, m)
				} else {
					limit++
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func IterModel1ByID(txn *store.LTxn, alloc *tools.Allocator, id int32, cb func(m *Model1) bool) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	opt := store.DefaultIteratorOptions
	opt.Prefix = alloc.Gen('M', C_Model1, 4018441491, id)
	iter := txn.NewIterator(opt)
	for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
		_ = iter.Item().Value(func(val []byte) error {
			m := &Model1{}
			err := m.Unmarshal(val)
			if err != nil {
				return err
			}
			if !cb(m) {
				exitLoop = true
			}
			return nil
		})
		if exitLoop {
			break
		}
	}
	iter.Close()
	return nil
}

func IterModel1ByEnum(txn *store.LTxn, alloc *tools.Allocator, enum Enum, cb func(m *Model1) bool) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	opt := store.DefaultIteratorOptions
	opt.Prefix = alloc.Gen('M', C_Model1, 2535881670, enum)
	iter := txn.NewIterator(opt)
	for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
		_ = iter.Item().Value(func(val []byte) error {
			m := &Model1{}
			err := m.Unmarshal(val)
			if err != nil {
				return err
			}
			if !cb(m) {
				exitLoop = true
			}
			return nil
		})
		if exitLoop {
			break
		}
	}
	iter.Close()
	return nil
}

func ListModel1ByID(
	id int32, offsetShardKey int32, lo *store.ListOption, cond func(m *Model1) bool,
) ([]*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('M', C_Model1, 4018441491, id)
		opt.Reverse = lo.Backward()
		osk := alloc.Gen('M', C_Model1, 4018441491, id, offsetShardKey)
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
				if cond == nil || cond(m) {
					res = append(res, m)
				} else {
					limit++
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func ListModel1ByEnum(
	enum Enum, offsetShardKey int32, offsetID int32, lo *store.ListOption, cond func(m *Model1) bool,
) ([]*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('M', C_Model1, 2535881670, enum)
		opt.Reverse = lo.Backward()
		osk := alloc.Gen('M', C_Model1, 2535881670, enum, offsetShardKey, offsetID)
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
				if cond == nil || cond(m) {
					res = append(res, m)
				} else {
					limit++
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func ListModel1ByP1(
	p1 string, lo *store.ListOption, cond func(m *Model1) bool,
) ([]*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('I', C_Model1, uint64(4843779728911368192), p1)
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
					if cond == nil || cond(m) {
						res = append(res, m)
					} else {
						limit++
					}
					return nil
				})
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func ListModel1ByP2(
	p2 string, lo *store.ListOption, cond func(m *Model1) bool,
) ([]*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model1, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('I', C_Model1, uint64(4749204136736587776), p2)
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
					if cond == nil || cond(m) {
						res = append(res, m)
					} else {
						limit++
					}
					return nil
				})
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

type Model2Order string

const Model2OrderByP1 Model2Order = "P1"

func CreateModel2(m *Model2) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *store.LTxn) error {
		return CreateModel2WithTxn(txn, alloc, m)
	})
}

func CreateModel2WithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Model2) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	if store.Exists(txn, alloc, 'M', C_Model2, 1609271041, m.ID, m.ShardKey, m.P1) {
		return store.ErrAlreadyExists
	}
	// save entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, 'M', C_Model2, 1609271041, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return
	}

	// save views
	// save entry for view: [P1 ShardKey ID]
	err = store.Set(txn, alloc, val, 'M', C_Model2, 2344331025, m.P1, m.ShardKey, m.ID)
	if err != nil {
		return
	}

	return

}

func ReadModel2WithTxn(txn *store.LTxn, alloc *tools.Allocator, id int64, shardKey int32, p1 string, m *Model2) (*Model2, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model2, 1609271041, id, shardKey, p1)
	if err != nil {
		return nil, err
	}
	return m, err
}

func ReadModel2(id int64, shardKey int32, p1 string, m *Model2) (*Model2, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model2{}
	}

	err := store.View(func(txn *store.LTxn) (err error) {
		m, err = ReadModel2WithTxn(txn, alloc, id, shardKey, p1, m)
		return err
	})
	return m, err
}

func ReadModel2ByP1AndShardKeyAndIDWithTxn(txn *store.LTxn, alloc *tools.Allocator, p1 string, shardKey int32, id int64, m *Model2) (*Model2, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model2, 2344331025, p1, shardKey, id)
	if err != nil {
		return nil, err
	}
	return m, err
}

func ReadModel2ByP1AndShardKeyAndID(p1 string, shardKey int32, id int64, m *Model2) (*Model2, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Model2{}
	}
	err := store.View(func(txn *store.LTxn) (err error) {
		m, err = ReadModel2ByP1AndShardKeyAndIDWithTxn(txn, alloc, p1, shardKey, id, m)
		return err
	})
	return m, err
}

func UpdateModel2WithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Model2) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := DeleteModel2WithTxn(txn, alloc, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return err
	}

	return CreateModel2WithTxn(txn, alloc, m)
}

func UpdateModel2(id int64, shardKey int32, p1 string, m *Model2) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		return store.ErrEmptyObject
	}

	err := store.View(func(txn *store.LTxn) (err error) {
		return UpdateModel2WithTxn(txn, alloc, m)
	})
	return err
}

func DeleteModel2WithTxn(txn *store.LTxn, alloc *tools.Allocator, id int64, shardKey int32, p1 string) error {
	m := &Model2{}
	err := store.Unmarshal(txn, alloc, m, 'M', C_Model2, 1609271041, id, shardKey, p1)
	if err != nil {
		return err
	}
	err = store.Delete(txn, alloc, 'M', C_Model2, 1609271041, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return err
	}

	err = store.Delete(txn, alloc, 'M', C_Model2, 2344331025, m.P1, m.ShardKey, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteModel2(id int64, shardKey int32, p1 string) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *store.LTxn) error {
		return DeleteModel2WithTxn(txn, alloc, id, shardKey, p1)
	})
}

func (x *Model2) HasP2(xx string) bool {
	for idx := range x.P2 {
		if x.P2[idx] == xx {
			return true
		}
	}
	return false
}

func SaveModel2WithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Model2) (err error) {
	if store.Exists(txn, alloc, 'M', C_Model2, 1609271041, m.ID, m.ShardKey, m.P1) {
		return UpdateModel2WithTxn(txn, alloc, m)
	} else {
		return CreateModel2WithTxn(txn, alloc, m)
	}
}

func SaveModel2(m *Model2) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *store.LTxn) error {
		return SaveModel2WithTxn(txn, alloc, m)
	})
}

func IterModel2(txn *store.LTxn, alloc *tools.Allocator, cb func(m *Model2) bool, orderBy ...Model2Order) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	iterOpt := store.DefaultIteratorOptions
	if len(orderBy) == 0 {
		iterOpt.Prefix = alloc.Gen('M', C_Model2, 1609271041)
	} else {
		switch orderBy[0] {
		case Model2OrderByP1:
			iterOpt.Prefix = alloc.Gen('M', C_Model2, 2344331025)
		}
	}
	iter := txn.NewIterator(iterOpt)
	for iter.Rewind(); iter.ValidForPrefix(iterOpt.Prefix); iter.Next() {
		_ = iter.Item().Value(func(val []byte) error {
			m := &Model2{}
			err := m.Unmarshal(val)
			if err != nil {
				return err
			}
			if !cb(m) {
				exitLoop = true
			}
			return nil
		})
		if exitLoop {
			break
		}
	}
	iter.Close()
	return nil
}

func ListModel2(
	offsetID int64, offsetShardKey int32, offsetP1 string, lo *store.ListOption, cond func(m *Model2) bool,
) ([]*Model2, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model2, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('M', C_Model2, 1609271041)
		opt.Reverse = lo.Backward()
		osk := alloc.Gen('M', C_Model2, 1609271041, offsetID, offsetShardKey)
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
				if cond == nil || cond(m) {
					res = append(res, m)
				} else {
					limit++
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func IterModel2ByIDAndShardKey(txn *store.LTxn, alloc *tools.Allocator, id int64, shardKey int32, cb func(m *Model2) bool) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	opt := store.DefaultIteratorOptions
	opt.Prefix = alloc.Gen('M', C_Model2, 1609271041, id, shardKey)
	iter := txn.NewIterator(opt)
	for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
		_ = iter.Item().Value(func(val []byte) error {
			m := &Model2{}
			err := m.Unmarshal(val)
			if err != nil {
				return err
			}
			if !cb(m) {
				exitLoop = true
			}
			return nil
		})
		if exitLoop {
			break
		}
	}
	iter.Close()
	return nil
}

func IterModel2ByP1(txn *store.LTxn, alloc *tools.Allocator, p1 string, cb func(m *Model2) bool) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	opt := store.DefaultIteratorOptions
	opt.Prefix = alloc.Gen('M', C_Model2, 2344331025, p1)
	iter := txn.NewIterator(opt)
	for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
		_ = iter.Item().Value(func(val []byte) error {
			m := &Model2{}
			err := m.Unmarshal(val)
			if err != nil {
				return err
			}
			if !cb(m) {
				exitLoop = true
			}
			return nil
		})
		if exitLoop {
			break
		}
	}
	iter.Close()
	return nil
}

func ListModel2ByIDAndShardKey(
	id int64, shardKey int32, offsetP1 string, lo *store.ListOption, cond func(m *Model2) bool,
) ([]*Model2, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model2, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('M', C_Model2, 1609271041, id, shardKey)
		opt.Reverse = lo.Backward()
		osk := alloc.Gen('M', C_Model2, 1609271041, id, shardKey, offsetP1)
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
				if cond == nil || cond(m) {
					res = append(res, m)
				} else {
					limit++
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func ListModel2ByP1(
	p1 string, offsetShardKey int32, offsetID int64, lo *store.ListOption, cond func(m *Model2) bool,
) ([]*Model2, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Model2, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('M', C_Model2, 2344331025, p1)
		opt.Reverse = lo.Backward()
		osk := alloc.Gen('M', C_Model2, 2344331025, p1, offsetShardKey, offsetID)
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
				if cond == nil || cond(m) {
					res = append(res, m)
				} else {
					limit++
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}
