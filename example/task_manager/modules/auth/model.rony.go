// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.12.44
// Source: model.proto

package auth

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

const C_User uint64 = 3297422589340680192

type poolUser struct {
	pool sync.Pool
}

func (p *poolUser) Get() *User {
	x, ok := p.pool.Get().(*User)
	if !ok {
		x = &User{}
	}

	return x
}

func (p *poolUser) Put(x *User) {
	if x == nil {
		return
	}

	x.Username = ""
	x.Password = ""
	x.FirstName = ""
	x.LastName = ""

	p.pool.Put(x)
}

var PoolUser = poolUser{}

func (x *User) DeepCopy(z *User) {
	z.Username = x.Username
	z.Password = x.Password
	z.FirstName = x.FirstName
	z.LastName = x.LastName
}

func (x *User) Clone() *User {
	z := &User{}
	x.DeepCopy(z)
	return z
}

func (x *User) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *User) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *User) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *User) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapUser(e registry.Envelope) (proto.Message, error) {
	x := &User{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *User) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_User, x)
}

const C_Session uint64 = 4434983008330775104

type poolSession struct {
	pool sync.Pool
}

func (p *poolSession) Get() *Session {
	x, ok := p.pool.Get().(*Session)
	if !ok {
		x = &Session{}
	}

	return x
}

func (p *poolSession) Put(x *Session) {
	if x == nil {
		return
	}

	x.ID = ""
	x.Username = ""

	p.pool.Put(x)
}

var PoolSession = poolSession{}

func (x *Session) DeepCopy(z *Session) {
	z.ID = x.ID
	z.Username = x.Username
}

func (x *Session) Clone() *Session {
	z := &Session{}
	x.DeepCopy(z)
	return z
}

func (x *Session) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Session) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Session) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Session) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapSession(e registry.Envelope) (proto.Message, error) {
	x := &Session{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Session) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Session, x)
}

// register constructors of the messages to the registry package
func init() {
	registry.Register(3297422589340680192, "User", unwrapUser)
	registry.Register(4434983008330775104, "Session", unwrapSession)

}

var _ = bytes.MinRead

type UserPrimaryKey interface {
	makeUserPrivate()
}

type UserPK struct {
	Username string
}

func (UserPK) makeUserPrivate() {}

type UserLocalRepo struct {
	s *store.Store
}

func NewUserLocalRepo(s *store.Store) *UserLocalRepo {
	return &UserLocalRepo{
		s: s,
	}
}

func (r *UserLocalRepo) CreateWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *User) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	key := alloc.Gen('M', C_User, uint64(12744997956659367007), m.Username)
	if store.ExistsByKey(txn, alloc, key) {
		return errors.ErrAlreadyExists
	}

	// save table entry
	val := alloc.Marshal(m)
	err = store.SetByKey(txn, val, key)
	if err != nil {
		return
	}

	return
}

func (r *UserLocalRepo) Create(m *User) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.CreateWithTxn(txn, alloc, m)
	})
}

func (r *UserLocalRepo) UpdateWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *User) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := r.DeleteWithTxn(txn, alloc, m.Username)
	if err != nil {
		return err
	}

	return r.CreateWithTxn(txn, alloc, m)
}

func (r *UserLocalRepo) Update(username string, m *User) error {
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

func (r *UserLocalRepo) SaveWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *User) (err error) {
	if store.Exists(txn, alloc, 'M', C_User, uint64(12744997956659367007), m.Username) {
		return r.UpdateWithTxn(txn, alloc, m)
	} else {
		return r.CreateWithTxn(txn, alloc, m)
	}
}

func (r *UserLocalRepo) Save(m *User) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.SaveWithTxn(txn, alloc, m)
	})
}

func (r *UserLocalRepo) ReadWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, username string, m *User) (*User, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_User, uint64(12744997956659367007), username)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *UserLocalRepo) Read(username string, m *User) (*User, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &User{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadWithTxn(txn, alloc, username, m)
		return err
	})
	return m, err
}

func (r *UserLocalRepo) DeleteWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, username string) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Delete(txn, alloc, 'M', C_User, uint64(12744997956659367007), username)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserLocalRepo) Delete(username string) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.DeleteWithTxn(txn, alloc, username)
	})
}

func (r *UserLocalRepo) ListWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator, offset UserPrimaryKey, lo *store.ListOption, cond func(m *User) bool,
) ([]*User, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	var seekKey []byte
	opt := store.DefaultIteratorOptions
	opt.Reverse = lo.Backward()
	res := make([]*User, 0, lo.Limit())

	switch offset := offset.(type) {
	case UserPK:
		opt.Prefix = alloc.Gen('M', C_User, uint64(12744997956659367007), offset.Username)
		seekKey = alloc.Gen('M', C_User, uint64(12744997956659367007), offset.Username)

	default:
		opt.Prefix = alloc.Gen('M', C_User, uint64(12744997956659367007))
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
				m := &User{}
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

func (r *UserLocalRepo) List(
	pk UserPrimaryKey, lo *store.ListOption, cond func(m *User) bool,
) ([]*User, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	var (
		res []*User
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

func (r *UserLocalRepo) IterWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator, offset UserPrimaryKey, ito *store.IterOption, cb func(m *User) bool,
) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	var seekKey []byte
	opt := store.DefaultIteratorOptions
	opt.Reverse = ito.Backward()

	switch offset := offset.(type) {
	case UserPK:
		opt.Prefix = alloc.Gen('M', C_User, uint64(12744997956659367007), offset.Username)
		seekKey = alloc.Gen('M', C_User, uint64(12744997956659367007), offset.Username)

	default:
		opt.Prefix = alloc.Gen('M', C_User, uint64(12744997956659367007))
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
				m := &User{}
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

func (r *UserLocalRepo) Iter(
	pk UserPrimaryKey, ito *store.IterOption, cb func(m *User) bool,
) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.View(func(txn *rony.StoreTxn) error {
		return r.IterWithTxn(txn, alloc, pk, ito, cb)
	})
}

type SessionPrimaryKey interface {
	makeSessionPrivate()
}

type SessionPK struct {
	ID string
}

func (SessionPK) makeSessionPrivate() {}

type SessionLocalRepo struct {
	s *store.Store
}

func NewSessionLocalRepo(s *store.Store) *SessionLocalRepo {
	return &SessionLocalRepo{
		s: s,
	}
}

func (r *SessionLocalRepo) CreateWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Session) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	key := alloc.Gen('M', C_Session, uint64(387204014837596160), m.ID)
	if store.ExistsByKey(txn, alloc, key) {
		return errors.ErrAlreadyExists
	}

	// save table entry
	val := alloc.Marshal(m)
	err = store.SetByKey(txn, val, key)
	if err != nil {
		return
	}

	// key := alloc.Gen('M', C_Session, uint64(387204014837596160), m.ID)
	// update field index by saving new value: Username
	err = store.Set(txn, alloc, key, 'I', C_Session, uint64(12744997956659367007), m.Username, m.ID)
	if err != nil {
		return
	}

	return
}

func (r *SessionLocalRepo) Create(m *Session) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.CreateWithTxn(txn, alloc, m)
	})
}

func (r *SessionLocalRepo) UpdateWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Session) error {
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

func (r *SessionLocalRepo) Update(id string, m *Session) error {
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

func (r *SessionLocalRepo) SaveWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *Session) (err error) {
	if store.Exists(txn, alloc, 'M', C_Session, uint64(387204014837596160), m.ID) {
		return r.UpdateWithTxn(txn, alloc, m)
	} else {
		return r.CreateWithTxn(txn, alloc, m)
	}
}

func (r *SessionLocalRepo) Save(m *Session) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.SaveWithTxn(txn, alloc, m)
	})
}

func (r *SessionLocalRepo) ReadWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id string, m *Session) (*Session, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Session, uint64(387204014837596160), id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *SessionLocalRepo) Read(id string, m *Session) (*Session, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Session{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadWithTxn(txn, alloc, id, m)
		return err
	})
	return m, err
}

func (r *SessionLocalRepo) DeleteWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, id string) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	m := &Session{}
	err := store.Unmarshal(txn, alloc, m, 'M', C_Session, uint64(387204014837596160), id)
	if err != nil {
		return err
	}
	err = store.Delete(txn, alloc, 'M', C_Session, uint64(387204014837596160), m.ID)
	if err != nil {
		return err
	}

	// delete field index
	err = store.Delete(txn, alloc, 'I', C_Session, uint64(12744997956659367007), m.Username, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionLocalRepo) Delete(id string) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.DeleteWithTxn(txn, alloc, id)
	})
}

func (r *SessionLocalRepo) ListWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator, offset SessionPrimaryKey, lo *store.ListOption, cond func(m *Session) bool,
) ([]*Session, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	var seekKey []byte
	opt := store.DefaultIteratorOptions
	opt.Reverse = lo.Backward()
	res := make([]*Session, 0, lo.Limit())

	switch offset := offset.(type) {
	case SessionPK:
		opt.Prefix = alloc.Gen('M', C_Session, uint64(387204014837596160), offset.ID)
		seekKey = alloc.Gen('M', C_Session, uint64(387204014837596160), offset.ID)

	default:
		opt.Prefix = alloc.Gen('M', C_Session, uint64(387204014837596160))
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
				m := &Session{}
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

func (r *SessionLocalRepo) List(
	pk SessionPrimaryKey, lo *store.ListOption, cond func(m *Session) bool,
) ([]*Session, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	var (
		res []*Session
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

func (r *SessionLocalRepo) IterWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator, offset SessionPrimaryKey, ito *store.IterOption, cb func(m *Session) bool,
) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	var seekKey []byte
	opt := store.DefaultIteratorOptions
	opt.Reverse = ito.Backward()

	switch offset := offset.(type) {
	case SessionPK:
		opt.Prefix = alloc.Gen('M', C_Session, uint64(387204014837596160), offset.ID)
		seekKey = alloc.Gen('M', C_Session, uint64(387204014837596160), offset.ID)

	default:
		opt.Prefix = alloc.Gen('M', C_Session, uint64(387204014837596160))
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
				m := &Session{}
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

func (r *SessionLocalRepo) Iter(
	pk SessionPrimaryKey, ito *store.IterOption, cb func(m *Session) bool,
) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.View(func(txn *rony.StoreTxn) error {
		return r.IterWithTxn(txn, alloc, pk, ito, cb)
	})
}

func (r *SessionLocalRepo) ListByUsername(username string, lo *store.ListOption, cond func(*Session) bool) ([]*Session, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	opt := store.DefaultIteratorOptions
	opt.Reverse = lo.Backward()
	opt.Prefix = alloc.Gen('I', C_Session, uint64(12744997956659367007), username)
	res := make([]*Session, 0, lo.Limit())
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
				m := &Session{}
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
	di.MustProvide(NewUserLocalRepo)

	di.MustProvide(NewSessionLocalRepo)

}