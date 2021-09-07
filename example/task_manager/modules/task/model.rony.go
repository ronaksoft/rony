// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.12.49
// Source: model.proto

package task

import (
	bytes "bytes"
	rony "github.com/ronaksoft/rony"
	di "github.com/ronaksoft/rony/di"
	edge "github.com/ronaksoft/rony/edge"
	errors "github.com/ronaksoft/rony/errors"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	store "github.com/ronaksoft/rony/store"
	tools "github.com/ronaksoft/rony/tools"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

var _ = pools.Imported

const C_Task uint64 = 4245896359682506752

type poolTask struct {
	pool sync.Pool
}

func (p *poolTask) Get() *Task {
	x, ok := p.pool.Get().(*Task)
	if !ok {
		x = &Task{}
	}

	return x
}

func (p *poolTask) Put(x *Task) {
	if x == nil {
		return
	}

	x.ID = 0
	x.Title = ""
	x.TODOs = x.TODOs[:0]
	x.DueDate = 0
	x.Username = ""

	p.pool.Put(x)
}

var PoolTask = poolTask{}

func (x *Task) DeepCopy(z *Task) {
	z.ID = x.ID
	z.Title = x.Title
	z.TODOs = append(z.TODOs[:0], x.TODOs...)
	z.DueDate = x.DueDate
	z.Username = x.Username
}

func (x *Task) Clone() *Task {
	z := &Task{}
	x.DeepCopy(z)
	return z
}

func (x *Task) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Task) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Task) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Task) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapTask(e registry.Envelope) (proto.Message, error) {
	x := &Task{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Task) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Task, x)
}

// register constructors of the messages to the registry package
func init() {
	registry.Register(4245896359682506752, "Task", unwrapTask)

}

var _ = bytes.MinRead

func (x *Task) HasTODO(xx string) bool {
	for idx := range x.TODOs {
		if x.TODOs[idx] == xx {
			return true
		}
	}
	return false
}

type TaskPrimaryKey interface {
	makeTaskPrivate()
}

type TaskPK struct {
	ID int64
}

func (TaskPK) makeTaskPrivate() {}

type TaskLocalRepo struct {
	s *store.Store
}

func NewTaskLocalRepo(s *store.Store) *TaskLocalRepo {
	return &TaskLocalRepo{
		s: s,
	}
}

func (r *TaskLocalRepo) CreateWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Task) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	key := alloc.Gen('M', C_Task, uint64(387204014837596160), m.ID)
	if store.ExistsByKey(txn, alloc, key) {
		return errors.ErrAlreadyExists
	}

	// save table entry
	val := alloc.Marshal(m)
	err = store.SetByKey(txn, val, key)
	if err != nil {
		return
	}

	// key := alloc.Gen('M', C_Task, uint64(387204014837596160), m.ID)
	// update field index by saving new value: Username
	err = store.Set(txn, alloc, key, 'I', C_Task, uint64(12744997956659367007), m.Username, m.ID)
	if err != nil {
		return
	}

	return
}

func (r *TaskLocalRepo) Create(m *Task) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.CreateWithTxn(txn, alloc, m)
	})
}

func (r *TaskLocalRepo) UpdateWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Task) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := r.DeleteWithTxn(txn, alloc, m.ID)
	if err != nil {
		return err
	}

	return r.CreateWithTxn(txn, alloc, m)
}

func (r *TaskLocalRepo) Update(id int64, m *Task) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		return errors.ErrEmptyObject
	}

	err := r.s.Update(func(txn *rony.StoreTxn) (err error) {
		return r.UpdateWithTxn(txn, alloc, m)
	})

	return err
}

func (r *TaskLocalRepo) SaveWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Task) (err error) {
	if store.Exists(txn, alloc, 'M', C_Task, uint64(387204014837596160), m.ID) {
		return r.UpdateWithTxn(txn, alloc, m)
	} else {
		return r.CreateWithTxn(txn, alloc, m)
	}
}

func (r *TaskLocalRepo) Save(m *Task) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.SaveWithTxn(txn, alloc, m)
	})
}

func (r *TaskLocalRepo) ReadWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int64, m *Task) (*Task, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Task, uint64(387204014837596160), id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *TaskLocalRepo) Read(id int64, m *Task) (*Task, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Task{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadWithTxn(txn, alloc, id, m)
		return err
	})
	return m, err
}

func (r *TaskLocalRepo) DeleteWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id int64) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	m := &Task{}
	err := store.Unmarshal(txn, alloc, m, 'M', C_Task, uint64(387204014837596160), id)
	if err != nil {
		return err
	}
	err = store.Delete(txn, alloc, 'M', C_Task, uint64(387204014837596160), m.ID)
	if err != nil {
		return err
	}

	// delete field index
	err = store.Delete(txn, alloc, 'I', C_Task, uint64(12744997956659367007), m.Username, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskLocalRepo) Delete(id int64) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.DeleteWithTxn(txn, alloc, id)
	})
}

func (r *TaskLocalRepo) ListWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator, offset TaskPrimaryKey, lo *store.ListOption, cond func(m *Task) bool,
) ([]*Task, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	var seekKey []byte
	opt := store.DefaultIteratorOptions
	opt.Reverse = lo.Backward()
	res := make([]*Task, 0, lo.Limit())

	switch offset := offset.(type) {
	case TaskPK:
		opt.Prefix = alloc.Gen('M', C_Task, uint64(387204014837596160), offset.ID)
		seekKey = alloc.Gen('M', C_Task, uint64(387204014837596160), offset.ID)

	default:
		opt.Prefix = alloc.Gen('M', C_Task, uint64(387204014837596160))
		seekKey = opt.Prefix
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
		for iter.Seek(seekKey); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			err = iter.Item().Value(func(val []byte) error {
				m := &Task{}
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
			if err != nil {
				return err
			}
		}
		iter.Close()
		return
	})

	return res, err
}

func (r *TaskLocalRepo) List(
	pk TaskPrimaryKey, lo *store.ListOption, cond func(m *Task) bool,
) ([]*Task, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	var (
		res []*Task
		err error
	)
	err = r.s.View(func(txn *rony.StoreTxn) error {
		res, err = r.ListWithTxn(txn, alloc, pk, lo, cond)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *TaskLocalRepo) IterWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator, offset TaskPrimaryKey, ito *store.IterOption, cb func(m *Task) bool,
) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	var seekKey []byte
	opt := store.DefaultIteratorOptions
	opt.Reverse = ito.Backward()

	switch offset := offset.(type) {
	case TaskPK:
		opt.Prefix = alloc.Gen('M', C_Task, uint64(387204014837596160), offset.ID)
		seekKey = alloc.Gen('M', C_Task, uint64(387204014837596160), offset.ID)

	default:
		opt.Prefix = alloc.Gen('M', C_Task, uint64(387204014837596160))
		seekKey = opt.Prefix
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		iter := txn.NewIterator(opt)
		if ito.OffsetKey() == nil {
			iter.Seek(seekKey)
		} else {
			iter.Seek(ito.OffsetKey())
		}
		exitLoop := false
		for ; iter.ValidForPrefix(opt.Prefix); iter.Next() {
			err = iter.Item().Value(func(val []byte) error {
				m := &Task{}
				err := m.Unmarshal(val)
				if err != nil {
					return err
				}
				if !cb(m) {
					exitLoop = true
				}
				return nil
			})
			if err != nil || exitLoop {
				break
			}
		}
		iter.Close()

		return
	})

	return err
}

func (r *TaskLocalRepo) Iter(
	pk TaskPrimaryKey, ito *store.IterOption, cb func(m *Task) bool,
) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.View(func(txn *rony.StoreTxn) error {
		return r.IterWithTxn(txn, alloc, pk, ito, cb)
	})
}

func (r *TaskLocalRepo) ListByUsername(username string, lo *store.ListOption, cond func(*Task) bool) ([]*Task, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	opt := store.DefaultIteratorOptions
	opt.Reverse = lo.Backward()
	opt.Prefix = alloc.Gen('I', C_Task, uint64(12744997956659367007), username)
	res := make([]*Task, 0, lo.Limit())
	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
		for iter.Seek(opt.Prefix); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			err = iter.Item().Value(func(val []byte) error {
				b, err := store.GetByKey(txn, alloc, val)
				if err != nil {
					return err
				}
				m := &Task{}
				err = m.Unmarshal(b)
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
			if err != nil {
				break
			}
		}
		iter.Close()
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// register provider constructors for dependency injection
func init() {
	di.MustProvide(NewTaskLocalRepo)

}
