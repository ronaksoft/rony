// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.15.8
// Rony ver. v0.11.4
// Source: model.proto

package model

import (
	bytes "bytes"
	rony "github.com/ronaksoft/rony"
	edge "github.com/ronaksoft/rony/edge"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	store "github.com/ronaksoft/rony/store"
	tools "github.com/ronaksoft/rony/tools"
	gocqlx "github.com/scylladb/gocqlx/v2"
	table "github.com/scylladb/gocqlx/v2/table"
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

func (x *Model1) Clone() *Model1 {
	z := &Model1{}
	x.DeepCopy(z)
	return z
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

func (x *Model2) Clone() *Model2 {
	z := &Model2{}
	x.DeepCopy(z)
	return z
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

func (x *Model3) Clone() *Model3 {
	z := &Model3{}
	x.DeepCopy(z)
	return z
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

var _ = bytes.MinRead

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
	if store.Exists(txn, alloc, 'M', C_Model1, 1248998560, m.ID, m.ShardKey, m.Enum) {
		return store.ErrAlreadyExists
	}

	// save table entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, 'M', C_Model1, 1248998560, m.ID, m.ShardKey, m.Enum)
	if err != nil {
		return
	}
	// save view entry
	err = store.Set(txn, alloc, val, 'M', C_Model1, 2535881670, m.Enum, m.ShardKey, m.ID)
	if err != nil {
		return err
	}

	key := alloc.Gen('M', C_Model1, 1248998560, m.ID, m.ShardKey, m.Enum)
	// update field index by saving new value: P1
	err = store.Set(txn, alloc, key, 'I', C_Model1, uint64(4843779728911368192), m.P1, m.ID, m.ShardKey, m.Enum)
	if err != nil {
		return
	}
	// update field index by saving new value: P2
	for idx := range m.P2 {
		err = store.Set(txn, alloc, key, 'I', C_Model1, uint64(4749204136736587776), m.P2[idx], m.ID, m.ShardKey, m.Enum)
		if err != nil {
			return
		}
	}

	return
}

func UpdateModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model1) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := DeleteModel1WithTxn(txn, alloc, m.ID, m.ShardKey, m.Enum)
	if err != nil {
		return err
	}

	return CreateModel1WithTxn(txn, alloc, m)
}

func UpdateModel1(id int32, shardKey int32, enum Enum, m *Model1) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		return store.ErrEmptyObject
	}

	err := store.Update(func(txn *rony.StoreTxn) (err error) {
		return UpdateModel1WithTxn(txn, alloc, m)
	})
	return err
}

func SaveModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Model1) (err error) {
	if store.Exists(txn, alloc, 'M', C_Model1, 1248998560, m.ID, m.ShardKey, m.Enum) {
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

func ReadModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int32, shardKey int32, enum Enum, m *Model1) (*Model1, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Model1, 1248998560, id, shardKey, enum)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func ReadModel1(id int32, shardKey int32, enum Enum, m *Model1) (*Model1, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Model1{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = ReadModel1WithTxn(txn, alloc, id, shardKey, enum, m)
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

func DeleteModel1WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int32, shardKey int32, enum Enum) error {
	m := &Model1{}
	err := store.Unmarshal(txn, alloc, m, 'M', C_Model1, 1248998560, id, shardKey, enum)
	if err != nil {
		return err
	}
	err = store.Delete(txn, alloc, 'M', C_Model1, 1248998560, m.ID, m.ShardKey, m.Enum)
	if err != nil {
		return err
	}

	// delete field index
	err = store.Delete(txn, alloc, 'I', C_Model1, uint64(4843779728911368192), m.P1, m.ID, m.ShardKey, m.Enum)
	if err != nil {
		return err
	}

	// delete field index
	for idx := range m.P2 {
		err = store.Delete(txn, alloc, 'I', C_Model1, uint64(4749204136736587776), m.P2[idx], m.ID, m.ShardKey, m.Enum)
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

func DeleteModel1(id int32, shardKey int32, enum Enum) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *rony.StoreTxn) error {
		return DeleteModel1WithTxn(txn, alloc, id, shardKey, enum)
	})
}

type Model1Order string

const (
	Model1OrderByEnum Model1Order = "Enum"
)

func (x *Model1) HasP2(xx string) bool {
	for idx := range x.P2 {
		if x.P2[idx] == xx {
			return true
		}
	}
	return false
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

	// save table entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, 'M', C_Model2, 1609271041, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return
	}
	// save view entry
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

	err := store.Update(func(txn *rony.StoreTxn) (err error) {
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

type Model2Order string

const (
	Model2OrderByP1 Model2Order = "P1"
)

func (x *Model2) HasP2(xx string) bool {
	for idx := range x.P2 {
		if x.P2[idx] == xx {
			return true
		}
	}
	return false
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

	// save table entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, 'M', C_Model3, 1609271041, m.ID, m.ShardKey, m.P1)
	if err != nil {
		return
	}
	// save view entry
	err = store.Set(txn, alloc, val, 'M', C_Model3, 2344331025, m.P1, m.ShardKey, m.ID)
	if err != nil {
		return err
	}
	// save view entry
	err = store.Set(txn, alloc, val, 'M', C_Model3, 3623577939, m.P1, m.ID)
	if err != nil {
		return err
	}

	key := alloc.Gen('M', C_Model3, 1609271041, m.ID, m.ShardKey, m.P1)
	// update field index by saving new value: P5
	for idx := range m.P5 {
		err = store.Set(txn, alloc, key, 'I', C_Model3, uint64(5041938112515670016), m.P5[idx], m.ID, m.ShardKey, m.P1)
		if err != nil {
			return
		}
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

	err := store.Update(func(txn *rony.StoreTxn) (err error) {
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

type Model3Order string

const (
	Model3OrderByP1 Model3Order = "P1"
)

func (x *Model3) HasP2(xx string) bool {
	for idx := range x.P2 {
		if x.P2[idx] == xx {
			return true
		}
	}
	return false
}

func (x *Model3) HasP5(xx []byte) bool {
	for idx := range x.P5 {
		if bytes.Equal(x.P5[idx], xx) {
			return true
		}
	}
	return false
}

type Model1RemoteRepo struct {
	qp map[string]*pools.QueryPool
	t  *table.Table
	v  map[string]*table.Table
	s  gocqlx.Session
}

func NewModel1RemoteRepo(s gocqlx.Session) *Model1RemoteRepo {
	r := &Model1RemoteRepo{
		s: s,
		t: table.New(table.Metadata{
			Name:    "tab_model_1",
			Columns: []string{"id", "shard_key", "enum", "sdata"},
			PartKey: []string{"id"},
			SortKey: []string{"shard_key", "enum"},
		}),
		v: map[string]*table.Table{
			"CustomerSort": table.New(table.Metadata{
				Name:    "view_model_1_customer_sort",
				Columns: []string{"enum", "shard_key", "id", "sdata"},
				PartKey: []string{"enum"},
				SortKey: []string{"shard_key", "id"},
			}),
		},
	}

	r.qp = map[string]*pools.QueryPool{
		"insertIF": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Unique().Query(s)
		}),
		"insert": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Query(s)
		}),
		"update": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.UpdateBuilder().Set("sdata").Query(s)
		}),
		"delete": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.DeleteBuilder().Query(s)
		}),
		"get": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.GetQuery(s)
		}),
		"getByCustomerSort": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.v["CustomerSort"].GetQuery(s)
		}),
	}
	return r
}

func (r *Model1RemoteRepo) Table() *table.Table {
	return r.t
}

func (r *Model1RemoteRepo) T() *table.Table {
	return r.t
}

func (r *Model1RemoteRepo) CustomerSort() *table.Table {
	return r.v["CustomerSort"]
}

func (r *Model1RemoteRepo) Insert(m *Model1, replace bool) error {
	buf := pools.Buffer.FromProto(m)
	defer pools.Buffer.Put(buf)

	var q *gocqlx.Queryx
	if replace {
		q = r.qp["insertIF"].GetQuery()
		defer r.qp["insertIF"].Put(q)
	} else {
		q = r.qp["insert"].GetQuery()
		defer r.qp["insert"].Put(q)
	}

	q.Bind(m.ID, m.ShardKey, m.Enum, *buf.Bytes())
	return q.Exec()
}

func (r *Model1RemoteRepo) Update(m *Model1) error {
	buf := pools.Buffer.FromProto(m)
	defer pools.Buffer.Put(buf)

	q := r.qp["update"].GetQuery()
	defer r.qp["update"].Put(q)

	q.Bind(*buf.Bytes(), m.ID, m.ShardKey, m.Enum)
	return q.Exec()
}

func (r *Model1RemoteRepo) Delete(id int32, shardKey int32, enum Enum) error {
	q := r.qp["delete"].GetQuery()
	defer r.qp["delete"].Put(q)

	q.Bind(id, shardKey, enum)
	return q.Exec()
}

func (r *Model1RemoteRepo) Get(id int32, shardKey int32, enum Enum, m *Model1) (*Model1, error) {
	q := r.qp["get"].GetQuery()
	defer r.qp["get"].Put(q)

	if m == nil {
		m = &Model1{}
	}

	q.Bind(id, shardKey, enum)

	var b []byte
	err := q.Scan(&m.ID, &m.ShardKey, &m.Enum, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}

func (r *Model1RemoteRepo) GetByCustomerSort(enum Enum, shardKey int32, id int32, m *Model1) (*Model1, error) {
	q := r.qp["getByCustomerSort"].GetQuery()
	defer r.qp["getByCustomerSort"].Put(q)

	if m == nil {
		m = &Model1{}
	}

	q.Bind(enum, shardKey, id)

	var b []byte
	err := q.Scan(&m.Enum, &m.ShardKey, &m.ID, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}

type Model1MountKey interface {
	makeItPrivate()
}

type Model1PK struct {
	ID int32
}

func (Model1PK) makeItPrivate() {}

type Model1CustomerSortPK struct {
	Enum Enum
}

func (Model1CustomerSortPK) makeItPrivate() {}

func (r *Model1RemoteRepo) List(mk Model1MountKey, limit uint) ([]*Model1, error) {
	var (
		q   *gocqlx.Queryx
		res []*Model1
		err error
	)

	switch mk := mk.(type) {
	case Model1PK:
		q = r.t.SelectBuilder().Limit(limit).Query(r.s)
		q.Bind(mk.ID)

	case Model1CustomerSortPK:
		q = r.v["CustomerSort"].SelectBuilder().Limit(limit).Query(r.s)
		q.Bind(mk.Enum)

	default:
		panic("BUG!! incorrect mount key")
	}
	err = q.SelectRelease(&res)

	return res, err
}

type Model2RemoteRepo struct {
	qp map[string]*pools.QueryPool
	t  *table.Table
	v  map[string]*table.Table
	s  gocqlx.Session
}

func NewModel2RemoteRepo(s gocqlx.Session) *Model2RemoteRepo {
	r := &Model2RemoteRepo{
		s: s,
		t: table.New(table.Metadata{
			Name:    "tab_model_2",
			Columns: []string{"id", "shard_key", "p_1", "sdata"},
			PartKey: []string{"id", "shard_key"},
			SortKey: []string{"p_1"},
		}),
		v: map[string]*table.Table{
			"P1ShardKeyID": table.New(table.Metadata{
				Name:    "view_model_2_p_1_shard_key_id",
				Columns: []string{"p_1", "shard_key", "id", "sdata"},
				PartKey: []string{"p_1"},
				SortKey: []string{"shard_key", "id"},
			}),
		},
	}

	r.qp = map[string]*pools.QueryPool{
		"insertIF": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Unique().Query(s)
		}),
		"insert": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Query(s)
		}),
		"update": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.UpdateBuilder().Set("sdata").Query(s)
		}),
		"delete": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.DeleteBuilder().Query(s)
		}),
		"get": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.GetQuery(s)
		}),
		"getByP1ShardKeyID": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.v["P1ShardKeyID"].GetQuery(s)
		}),
	}
	return r
}

func (r *Model2RemoteRepo) Table() *table.Table {
	return r.t
}

func (r *Model2RemoteRepo) T() *table.Table {
	return r.t
}

func (r *Model2RemoteRepo) MVP1ShardKeyID() *table.Table {
	return r.v["P1ShardKeyID"]
}

func (r *Model2RemoteRepo) Insert(m *Model2, replace bool) error {
	buf := pools.Buffer.FromProto(m)
	defer pools.Buffer.Put(buf)

	var q *gocqlx.Queryx
	if replace {
		q = r.qp["insertIF"].GetQuery()
		defer r.qp["insertIF"].Put(q)
	} else {
		q = r.qp["insert"].GetQuery()
		defer r.qp["insert"].Put(q)
	}

	q.Bind(m.ID, m.ShardKey, m.P1, *buf.Bytes())
	return q.Exec()
}

func (r *Model2RemoteRepo) Update(m *Model2) error {
	buf := pools.Buffer.FromProto(m)
	defer pools.Buffer.Put(buf)

	q := r.qp["update"].GetQuery()
	defer r.qp["update"].Put(q)

	q.Bind(*buf.Bytes(), m.ID, m.ShardKey, m.P1)
	return q.Exec()
}

func (r *Model2RemoteRepo) Delete(id int64, shardKey int32, p1 string) error {
	q := r.qp["delete"].GetQuery()
	defer r.qp["delete"].Put(q)

	q.Bind(id, shardKey, p1)
	return q.Exec()
}

func (r *Model2RemoteRepo) Get(id int64, shardKey int32, p1 string, m *Model2) (*Model2, error) {
	q := r.qp["get"].GetQuery()
	defer r.qp["get"].Put(q)

	if m == nil {
		m = &Model2{}
	}

	q.Bind(id, shardKey, p1)

	var b []byte
	err := q.Scan(&m.ID, &m.ShardKey, &m.P1, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}

func (r *Model2RemoteRepo) GetByP1ShardKeyID(p1 string, shardKey int32, id int64, m *Model2) (*Model2, error) {
	q := r.qp["getByP1ShardKeyID"].GetQuery()
	defer r.qp["getByP1ShardKeyID"].Put(q)

	if m == nil {
		m = &Model2{}
	}

	q.Bind(p1, shardKey, id)

	var b []byte
	err := q.Scan(&m.P1, &m.ShardKey, &m.ID, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}

type Model2MountKey interface {
	makeItPrivate()
}

type Model2PK struct {
	ID       int64
	ShardKey int32
}

func (Model2PK) makeItPrivate() {}

type Model2P1ShardKeyIDPK struct {
	P1 string
}

func (Model2P1ShardKeyIDPK) makeItPrivate() {}

func (r *Model2RemoteRepo) List(mk Model2MountKey, limit uint) ([]*Model2, error) {
	var (
		q   *gocqlx.Queryx
		res []*Model2
		err error
	)

	switch mk := mk.(type) {
	case Model2PK:
		q = r.t.SelectBuilder().Limit(limit).Query(r.s)
		q.Bind(mk.ID, mk.ShardKey)

	case Model2P1ShardKeyIDPK:
		q = r.v["P1ShardKeyID"].SelectBuilder().Limit(limit).Query(r.s)
		q.Bind(mk.P1)

	default:
		panic("BUG!! incorrect mount key")
	}
	err = q.SelectRelease(&res)

	return res, err
}

type Model3RemoteRepo struct {
	qp map[string]*pools.QueryPool
	t  *table.Table
	v  map[string]*table.Table
	s  gocqlx.Session
}

func NewModel3RemoteRepo(s gocqlx.Session) *Model3RemoteRepo {
	r := &Model3RemoteRepo{
		s: s,
		t: table.New(table.Metadata{
			Name:    "tab_model_3",
			Columns: []string{"id", "shard_key", "p_1", "sdata"},
			PartKey: []string{"id", "shard_key"},
			SortKey: []string{"p_1"},
		}),
		v: map[string]*table.Table{
			"P1ShardKeyID": table.New(table.Metadata{
				Name:    "view_model_3_p_1_shard_key_id",
				Columns: []string{"p_1", "shard_key", "id", "sdata"},
				PartKey: []string{"p_1"},
				SortKey: []string{"shard_key", "id"},
			}),
			"P1ID": table.New(table.Metadata{
				Name:    "view_model_3_p_1_id",
				Columns: []string{"p_1", "id", "sdata"},
				PartKey: []string{"p_1"},
				SortKey: []string{"id"},
			}),
		},
	}

	r.qp = map[string]*pools.QueryPool{
		"insertIF": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Unique().Query(s)
		}),
		"insert": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Query(s)
		}),
		"update": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.UpdateBuilder().Set("sdata").Query(s)
		}),
		"delete": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.DeleteBuilder().Query(s)
		}),
		"get": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.GetQuery(s)
		}),
		"getByP1ShardKeyID": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.v["P1ShardKeyID"].GetQuery(s)
		}),
		"getByP1ID": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.v["P1ID"].GetQuery(s)
		}),
	}
	return r
}

func (r *Model3RemoteRepo) Table() *table.Table {
	return r.t
}

func (r *Model3RemoteRepo) T() *table.Table {
	return r.t
}

func (r *Model3RemoteRepo) MVP1ShardKeyID() *table.Table {
	return r.v["P1ShardKeyID"]
}

func (r *Model3RemoteRepo) MVP1ID() *table.Table {
	return r.v["P1ID"]
}

func (r *Model3RemoteRepo) Insert(m *Model3, replace bool) error {
	buf := pools.Buffer.FromProto(m)
	defer pools.Buffer.Put(buf)

	var q *gocqlx.Queryx
	if replace {
		q = r.qp["insertIF"].GetQuery()
		defer r.qp["insertIF"].Put(q)
	} else {
		q = r.qp["insert"].GetQuery()
		defer r.qp["insert"].Put(q)
	}

	q.Bind(m.ID, m.ShardKey, m.P1, *buf.Bytes())
	return q.Exec()
}

func (r *Model3RemoteRepo) Update(m *Model3) error {
	buf := pools.Buffer.FromProto(m)
	defer pools.Buffer.Put(buf)

	q := r.qp["update"].GetQuery()
	defer r.qp["update"].Put(q)

	q.Bind(*buf.Bytes(), m.ID, m.ShardKey, m.P1)
	return q.Exec()
}

func (r *Model3RemoteRepo) Delete(id int64, shardKey int32, p1 []byte) error {
	q := r.qp["delete"].GetQuery()
	defer r.qp["delete"].Put(q)

	q.Bind(id, shardKey, p1)
	return q.Exec()
}

func (r *Model3RemoteRepo) Get(id int64, shardKey int32, p1 []byte, m *Model3) (*Model3, error) {
	q := r.qp["get"].GetQuery()
	defer r.qp["get"].Put(q)

	if m == nil {
		m = &Model3{}
	}

	q.Bind(id, shardKey, p1)

	var b []byte
	err := q.Scan(&m.ID, &m.ShardKey, &m.P1, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}

func (r *Model3RemoteRepo) GetByP1ShardKeyID(p1 []byte, shardKey int32, id int64, m *Model3) (*Model3, error) {
	q := r.qp["getByP1ShardKeyID"].GetQuery()
	defer r.qp["getByP1ShardKeyID"].Put(q)

	if m == nil {
		m = &Model3{}
	}

	q.Bind(p1, shardKey, id)

	var b []byte
	err := q.Scan(&m.P1, &m.ShardKey, &m.ID, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}

func (r *Model3RemoteRepo) GetByP1ID(p1 []byte, id int64, m *Model3) (*Model3, error) {
	q := r.qp["getByP1ID"].GetQuery()
	defer r.qp["getByP1ID"].Put(q)

	if m == nil {
		m = &Model3{}
	}

	q.Bind(p1, id)

	var b []byte
	err := q.Scan(&m.P1, &m.ID, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}

type Model3MountKey interface {
	makeItPrivate()
}

type Model3PK struct {
	ID       int64
	ShardKey int32
}

func (Model3PK) makeItPrivate() {}

type Model3P1ShardKeyIDPK struct {
	P1 []byte
}

func (Model3P1ShardKeyIDPK) makeItPrivate() {}

type Model3P1IDPK struct {
	P1 []byte
}

func (Model3P1IDPK) makeItPrivate() {}

func (r *Model3RemoteRepo) List(mk Model3MountKey, limit uint) ([]*Model3, error) {
	var (
		q   *gocqlx.Queryx
		res []*Model3
		err error
	)

	switch mk := mk.(type) {
	case Model3PK:
		q = r.t.SelectBuilder().Limit(limit).Query(r.s)
		q.Bind(mk.ID, mk.ShardKey)

	case Model3P1ShardKeyIDPK:
		q = r.v["P1ShardKeyID"].SelectBuilder().Limit(limit).Query(r.s)
		q.Bind(mk.P1)

	case Model3P1IDPK:
		q = r.v["P1ID"].SelectBuilder().Limit(limit).Query(r.s)
		q.Bind(mk.P1)

	default:
		panic("BUG!! incorrect mount key")
	}
	err = q.SelectRelease(&res)

	return res, err
}
