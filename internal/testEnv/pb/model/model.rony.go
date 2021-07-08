// Code generated by Rony's protoc plugin; DO NOT EDIT.

package model

import (
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

func (x *Model1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model1) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Model1) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
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

func (x *Model2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model2) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Model2) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *Model2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model2, x)
}

const C_Model3 int64 = 2510714031

type poolModel3 struct {
	pool sync.Pool
}

func (p *poolModel3) Get() *Model3 {
	x, ok := p.pool.Get().(*Model3)
	if !ok {
		x = &Model3{}
	}
	return x
}

func (p *poolModel3) Put(x *Model3) {
	if x == nil {
		return
	}

	x.ID = 0
	x.ShardKey = 0
	x.P1 = x.P1[:0]
	x.P2 = x.P2[:0]
	for _, z := range x.P5 {
		pools.Bytes.Put(z)
	}
	x.P5 = x.P5[:0]

	p.pool.Put(x)
}

var PoolModel3 = poolModel3{}

func (x *Model3) DeepCopy(z *Model3) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = append(z.P1[:0], x.P1...)
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = z.P5[:0]
	zl := len(z.P5)
	for idx := range x.P5 {
		if idx < zl {
			z.P5 = append(z.P5, append(z.P5[idx][:0], x.P5[idx]...))
		} else {
			zb := pools.Bytes.GetCap(len(x.P5[idx]))
			z.P5 = append(z.P5, append(zb, x.P5[idx]...))
		}
	}
}

func (x *Model3) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model3) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model3) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Model3) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *Model3) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model3, x)
}

func init() {
	registry.RegisterConstructor(2074613123, "Model1")
	registry.RegisterConstructor(3802219577, "Model2")
	registry.RegisterConstructor(2510714031, "Model3")
}

func CreateModel1(m *Model1) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *rony.StoreTxn) error {
		return CreateModel1WithTxn(txn, alloc, m)
	})
}

func CreateModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model1) (err error) {
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
	// save view [{Enum enum unsupported Enum } {ShardKey int32 int int32 ASC} {ID int32 int int32 ASC}]
	err = store.Set(txn, alloc, val, 'M', C_Model1, 2535881670, m.Enum, m.ShardKey, m.ID)
	if err != nil {
		return err
	}

	key := alloc.Gen('M', C_Model1, 4018441491, m.ID, m.ShardKey)
	// update field index by saving new value: P1
	err = store.Set(txn, alloc, key, 'I', C_Model1, uint64(4843779728911368192), m.P1, m.ID, m.ShardKey)
	if err != nil {
		return
	}
	// update field index by saving new value: P2
	err = store.Set(txn, alloc, key, 'I', C_Model1, uint64(4749204136736587776), m.P2, m.ID, m.ShardKey)
	if err != nil {
		return
	}

	return
}

func UpdateModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model1) error {
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

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		return UpdateModel1WithTxn(txn, alloc, m)
	})
	return err
}

func SaveModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model1) (err error) {
	if store.Exists(txn, alloc, 'M', C_Model1, 4018441491, m.ID, m.ShardKey) {
		return UpdateModel1WithTxn(txn, alloc, m)
	} else {
		return CreateModel1WithTxn(txn, alloc, m)
	}
}

func SaveModel1(m *Model1) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *rony.StoreTxn) error {
		return SaveModel1WithTxn(txn, alloc, m)
	})
}

func ReadModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int32, shardKey int32, m *Model1) (*Model1, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model1, 4018441491, id, shardKey)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func ReadModel1(id int32, shardKey int32, m *Model1) (*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model1{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = ReadModel1WithTxn(txn, alloc, id, shardKey, m)
		return err
	})
	return m, err
}

func ReadModel1ByEnumAndShardKeyAndIDWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator,
	enum Enum, shardKey int32, id int32, m *Model1,
) (*Model1, error) {
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

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = ReadModel1ByEnumAndShardKeyAndIDWithTxn(txn, alloc, enum, shardKey, id, m)
		return err
	})
	return m, err
}

func DeleteModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int32, shardKey int32) error {
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

	return store.Update(func(txn *rony.StoreTxn) error {
		return DeleteModel1WithTxn(txn, alloc, id, shardKey)
	})
}

func CreateModel2(m *Model2) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *rony.StoreTxn) error {
		return CreateModel2WithTxn(txn, alloc, m)
	})
}

func CreateModel2WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model2) (err error) {
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
	// save view [{P1 string blob string } {ShardKey int32 int int32 ASC} {ID int64 bigint int64 ASC}]
	err = store.Set(txn, alloc, val, 'M', C_Model2, 2344331025, m.P1, m.ShardKey, m.ID)
	if err != nil {
		return err
	}

	return
}

func UpdateModel2WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model2) error {
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

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		return UpdateModel2WithTxn(txn, alloc, m)
	})
	return err
}

func SaveModel2WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model2) (err error) {
	if store.Exists(txn, alloc, 'M', C_Model2, 1609271041, m.ID, m.ShardKey, m.P1) {
		return UpdateModel2WithTxn(txn, alloc, m)
	} else {
		return CreateModel2WithTxn(txn, alloc, m)
	}
}

func SaveModel2(m *Model2) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *rony.StoreTxn) error {
		return SaveModel2WithTxn(txn, alloc, m)
	})
}

func ReadModel2WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int64, shardKey int32, p1 string, m *Model2) (*Model2, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model2, 1609271041, id, shardKey, p1)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func ReadModel2(id int64, shardKey int32, p1 string, m *Model2) (*Model2, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model2{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = ReadModel2WithTxn(txn, alloc, id, shardKey, p1, m)
		return err
	})
	return m, err
}

func ReadModel2ByP1AndShardKeyAndIDWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator,
	p1 string, shardKey int32, id int64, m *Model2,
) (*Model2, error) {
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

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = ReadModel2ByP1AndShardKeyAndIDWithTxn(txn, alloc, p1, shardKey, id, m)
		return err
	})
	return m, err
}

func DeleteModel2WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int64, shardKey int32, p1 string) error {
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

	return store.Update(func(txn *rony.StoreTxn) error {
		return DeleteModel2WithTxn(txn, alloc, id, shardKey, p1)
	})
}

func CreateModel3(m *Model3) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *rony.StoreTxn) error {
		return CreateModel3WithTxn(txn, alloc, m)
	})
}

func CreateModel3WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model3) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	if store.Exists(txn, alloc, 'M', C_Model3, 1609271041, m.ID, m.ShardKey, m.P1) {
		return store.ErrAlreadyExists
	}

	// save entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, 'M', C_Model3, 1609271041, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return
	}
	// save view [{P1 bytes blob []byte } {ShardKey int32 int int32 ASC} {ID int64 bigint int64 ASC}]
	err = store.Set(txn, alloc, val, 'M', C_Model3, 2344331025, m.P1, m.ShardKey, m.ID)
	if err != nil {
		return err
	}
	// save view [{P1 bytes blob []byte } {ID int64 bigint int64 ASC}]
	err = store.Set(txn, alloc, val, 'M', C_Model3, 3623577939, m.P1, m.ID)
	if err != nil {
		return err
	}

	key := alloc.Gen('M', C_Model3, 1609271041, m.ID, m.ShardKey, m.P1)
	// update field index by saving new value: P5
	err = store.Set(txn, alloc, key, 'I', C_Model3, uint64(5041938112515670016), m.P5, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return
	}

	return
}

func UpdateModel3WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model3) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := DeleteModel3WithTxn(txn, alloc, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return err
	}

	return CreateModel3WithTxn(txn, alloc, m)
}

func UpdateModel3(id int64, shardKey int32, p1 []byte, m *Model3) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		return store.ErrEmptyObject
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		return UpdateModel3WithTxn(txn, alloc, m)
	})
	return err
}

func SaveModel3WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model3) (err error) {
	if store.Exists(txn, alloc, 'M', C_Model3, 1609271041, m.ID, m.ShardKey, m.P1) {
		return UpdateModel3WithTxn(txn, alloc, m)
	} else {
		return CreateModel3WithTxn(txn, alloc, m)
	}
}

func SaveModel3(m *Model3) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *rony.StoreTxn) error {
		return SaveModel3WithTxn(txn, alloc, m)
	})
}

func ReadModel3WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int64, shardKey int32, p1 []byte, m *Model3) (*Model3, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model3, 1609271041, id, shardKey, p1)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func ReadModel3(id int64, shardKey int32, p1 []byte, m *Model3) (*Model3, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model3{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = ReadModel3WithTxn(txn, alloc, id, shardKey, p1, m)
		return err
	})
	return m, err
}

func ReadModel3ByP1AndShardKeyAndIDWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator,
	p1 []byte, shardKey int32, id int64, m *Model3,
) (*Model3, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model3, 2344331025, p1, shardKey, id)
	if err != nil {
		return nil, err
	}
	return m, err
}

func ReadModel3ByP1AndShardKeyAndID(p1 []byte, shardKey int32, id int64, m *Model3) (*Model3, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model3{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = ReadModel3ByP1AndShardKeyAndIDWithTxn(txn, alloc, p1, shardKey, id, m)
		return err
	})
	return m, err
}

func ReadModel3ByP1AndIDWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator,
	p1 []byte, id int64, m *Model3,
) (*Model3, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model3, 3623577939, p1, id)
	if err != nil {
		return nil, err
	}
	return m, err
}

func ReadModel3ByP1AndID(p1 []byte, id int64, m *Model3) (*Model3, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model3{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = ReadModel3ByP1AndIDWithTxn(txn, alloc, p1, id, m)
		return err
	})
	return m, err
}

func DeleteModel3WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int64, shardKey int32, p1 []byte) error {
	m := &Model3{}
	err := store.Unmarshal(txn, alloc, m, 'M', C_Model3, 1609271041, id, shardKey, p1)
	if err != nil {
		return err
	}
	err = store.Delete(txn, alloc, 'M', C_Model3, 1609271041, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return err
	}

	// delete field index
	for idx := range m.P5 {
		err = store.Delete(txn, alloc, 'I', C_Model3, uint64(5041938112515670016), m.P5[idx], m.ID, m.ShardKey, m.P1)
		if err != nil {
			return err
		}
	}
	err = store.Delete(txn, alloc, 'M', C_Model3, 2344331025, m.P1, m.ShardKey, m.ID)
	if err != nil {
		return err
	}
	err = store.Delete(txn, alloc, 'M', C_Model3, 3623577939, m.P1, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteModel3(id int64, shardKey int32, p1 []byte) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *rony.StoreTxn) error {
		return DeleteModel3WithTxn(txn, alloc, id, shardKey, p1)
	})
}
