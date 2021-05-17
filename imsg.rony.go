// Code generated by Rony's protoc plugin; DO NOT EDIT.

package rony

import (
	registry "github.com/ronaksoft/rony/registry"
	store "github.com/ronaksoft/rony/store"
	tools "github.com/ronaksoft/rony/tools"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_GetPage int64 = 3721890413

type poolGetPage struct {
	pool sync.Pool
}

func (p *poolGetPage) Get() *GetPage {
	x, ok := p.pool.Get().(*GetPage)
	if !ok {
		x = &GetPage{}
	}
	return x
}

func (p *poolGetPage) Put(x *GetPage) {
	if x == nil {
		return
	}
	x.PageID = 0
	x.ReplicaSet = 0
	p.pool.Put(x)
}

var PoolGetPage = poolGetPage{}

func (x *GetPage) DeepCopy(z *GetPage) {
	z.PageID = x.PageID
	z.ReplicaSet = x.ReplicaSet
}

func (x *GetPage) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *GetPage) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_StoreMessage int64 = 4229276430

type poolStoreMessage struct {
	pool sync.Pool
}

func (p *poolStoreMessage) Get() *StoreMessage {
	x, ok := p.pool.Get().(*StoreMessage)
	if !ok {
		x = &StoreMessage{}
	}
	return x
}

func (p *poolStoreMessage) Put(x *StoreMessage) {
	if x == nil {
		return
	}
	x.Constructor = 0
	x.Payload = x.Payload[:0]
	x.Request = false
	p.pool.Put(x)
}

var PoolStoreMessage = poolStoreMessage{}

func (x *StoreMessage) DeepCopy(z *StoreMessage) {
	z.Constructor = x.Constructor
	z.Payload = append(z.Payload[:0], x.Payload...)
	z.Request = x.Request
}

func (x *StoreMessage) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *StoreMessage) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_TunnelMessage int64 = 3271476222

type poolTunnelMessage struct {
	pool sync.Pool
}

func (p *poolTunnelMessage) Get() *TunnelMessage {
	x, ok := p.pool.Get().(*TunnelMessage)
	if !ok {
		x = &TunnelMessage{}
	}
	return x
}

func (p *poolTunnelMessage) Put(x *TunnelMessage) {
	if x == nil {
		return
	}
	x.SenderID = x.SenderID[:0]
	x.SenderReplicaSet = 0
	for _, z := range x.Store {
		PoolKeyValue.Put(z)
	}
	x.Store = x.Store[:0]
	PoolMessageEnvelope.Put(x.Envelope)
	x.Envelope = nil
	p.pool.Put(x)
}

var PoolTunnelMessage = poolTunnelMessage{}

func (x *TunnelMessage) DeepCopy(z *TunnelMessage) {
	z.SenderID = append(z.SenderID[:0], x.SenderID...)
	z.SenderReplicaSet = x.SenderReplicaSet
	for idx := range x.Store {
		if x.Store[idx] != nil {
			xx := PoolKeyValue.Get()
			x.Store[idx].DeepCopy(xx)
			z.Store = append(z.Store, xx)
		}
	}
	if x.Envelope != nil {
		if z.Envelope == nil {
			z.Envelope = PoolMessageEnvelope.Get()
		}
		x.Envelope.DeepCopy(z.Envelope)
	} else {
		z.Envelope = nil
	}
}

func (x *TunnelMessage) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *TunnelMessage) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_EdgeNode int64 = 999040174

type poolEdgeNode struct {
	pool sync.Pool
}

func (p *poolEdgeNode) Get() *EdgeNode {
	x, ok := p.pool.Get().(*EdgeNode)
	if !ok {
		x = &EdgeNode{}
	}
	return x
}

func (p *poolEdgeNode) Put(x *EdgeNode) {
	if x == nil {
		return
	}
	x.ServerID = x.ServerID[:0]
	x.ReplicaSet = 0
	x.GatewayAddr = x.GatewayAddr[:0]
	x.TunnelAddr = x.TunnelAddr[:0]
	p.pool.Put(x)
}

var PoolEdgeNode = poolEdgeNode{}

func (x *EdgeNode) DeepCopy(z *EdgeNode) {
	z.ServerID = append(z.ServerID[:0], x.ServerID...)
	z.ReplicaSet = x.ReplicaSet
	z.GatewayAddr = append(z.GatewayAddr[:0], x.GatewayAddr...)
	z.TunnelAddr = append(z.TunnelAddr[:0], x.TunnelAddr...)
}

func (x *EdgeNode) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EdgeNode) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_Page int64 = 3023575326

type poolPage struct {
	pool sync.Pool
}

func (p *poolPage) Get() *Page {
	x, ok := p.pool.Get().(*Page)
	if !ok {
		x = &Page{}
	}
	return x
}

func (p *poolPage) Put(x *Page) {
	if x == nil {
		return
	}
	x.ID = 0
	x.ReplicaSet = 0
	p.pool.Put(x)
}

var PoolPage = poolPage{}

func (x *Page) DeepCopy(z *Page) {
	z.ID = x.ID
	z.ReplicaSet = x.ReplicaSet
}

func (x *Page) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Page) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func init() {
	registry.RegisterConstructor(3721890413, "GetPage")
	registry.RegisterConstructor(4229276430, "StoreMessage")
	registry.RegisterConstructor(3271476222, "TunnelMessage")
	registry.RegisterConstructor(999040174, "EdgeNode")
	registry.RegisterConstructor(3023575326, "Page")
}

func CreatePage(m *Page) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *store.LTxn) error {
		return CreatePageWithTxn(txn, alloc, m)
	})
}

func CreatePageWithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Page) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	if store.Exists(txn, alloc, 'M', C_Page, 299066170, m.ID) {
		return store.ErrAlreadyExists
	}
	// save entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, 'M', C_Page, 299066170, m.ID)
	if err != nil {
		return
	}

	// save views
	// save entry for view: [ReplicaSet ID]
	err = store.Set(txn, alloc, val, 'M', C_Page, 1040696757, m.ReplicaSet, m.ID)
	if err != nil {
		return
	}

	return

}

func ReadPageWithTxn(txn *store.LTxn, alloc *tools.Allocator, id uint32, m *Page) (*Page, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Page, 299066170, id)
	if err != nil {
		return nil, err
	}
	return m, err
}

func ReadPage(id uint32, m *Page) (*Page, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Page{}
	}

	err := store.View(func(txn *store.LTxn) (err error) {
		m, err = ReadPageWithTxn(txn, alloc, id, m)
		return err
	})
	return m, err
}

func ReadPageByReplicaSetAndIDWithTxn(txn *store.LTxn, alloc *tools.Allocator, replicaSet uint64, id uint32, m *Page) (*Page, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, 'M', C_Page, 1040696757, replicaSet, id)
	if err != nil {
		return nil, err
	}
	return m, err
}

func ReadPageByReplicaSetAndID(replicaSet uint64, id uint32, m *Page) (*Page, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Page{}
	}
	err := store.View(func(txn *store.LTxn) (err error) {
		m, err = ReadPageByReplicaSetAndIDWithTxn(txn, alloc, replicaSet, id, m)
		return err
	})
	return m, err
}

func UpdatePageWithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Page) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := DeletePageWithTxn(txn, alloc, m.ID)
	if err != nil {
		return err
	}

	return CreatePageWithTxn(txn, alloc, m)
}

func UpdatePage(id uint32, m *Page) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		return store.ErrEmptyObject
	}

	err := store.View(func(txn *store.LTxn) (err error) {
		return UpdatePageWithTxn(txn, alloc, m)
	})
	return err
}

func DeletePageWithTxn(txn *store.LTxn, alloc *tools.Allocator, id uint32) error {
	m := &Page{}
	err := store.Unmarshal(txn, alloc, m, 'M', C_Page, 299066170, id)
	if err != nil {
		return err
	}
	err = store.Delete(txn, alloc, 'M', C_Page, 299066170, m.ID)
	if err != nil {
		return err
	}

	err = store.Delete(txn, alloc, 'M', C_Page, 1040696757, m.ReplicaSet, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeletePage(id uint32) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *store.LTxn) error {
		return DeletePageWithTxn(txn, alloc, id)
	})
}

func SavePageWithTxn(txn *store.LTxn, alloc *tools.Allocator, m *Page) (err error) {
	if store.Exists(txn, alloc, 'M', C_Page, 299066170, m.ID) {
		return UpdatePageWithTxn(txn, alloc, m)
	} else {
		return CreatePageWithTxn(txn, alloc, m)
	}
}

func SavePage(m *Page) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *store.LTxn) error {
		return SavePageWithTxn(txn, alloc, m)
	})
}

func IterPages(txn *store.LTxn, alloc *tools.Allocator, cb func(m *Page) bool) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	iterOpt := store.DefaultIteratorOptions
	iterOpt.Prefix = alloc.Gen('M', C_Page, 299066170)
	iter := txn.NewIterator(iterOpt)
	for iter.Rewind(); iter.ValidForPrefix(iterOpt.Prefix); iter.Next() {
		_ = iter.Item().Value(func(val []byte) error {
			m := &Page{}
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

func ListPage(
	offsetID uint32, lo *store.ListOption, cond func(m *Page) bool,
) ([]*Page, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Page, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('M', C_Page, 299066170)
		opt.Reverse = lo.Backward()
		osk := alloc.Gen('M', C_Page, 299066170, offsetID)
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
				m := &Page{}
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

func IterPageByReplicaSet(txn *store.LTxn, alloc *tools.Allocator, replicaSet uint64, cb func(m *Page) bool) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	opt := store.DefaultIteratorOptions
	opt.Prefix = alloc.Gen('M', C_Page, 1040696757, replicaSet)
	iter := txn.NewIterator(opt)
	for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {
		_ = iter.Item().Value(func(val []byte) error {
			m := &Page{}
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

func ListPageByReplicaSet(
	replicaSet uint64, offsetID uint32, lo *store.ListOption, cond func(m *Page) bool,
) ([]*Page, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Page, 0, lo.Limit())
	err := store.View(func(txn *store.LTxn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.Gen('M', C_Page, 1040696757, replicaSet)
		opt.Reverse = lo.Backward()
		osk := alloc.Gen('M', C_Page, 1040696757, replicaSet, offsetID)
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
				m := &Page{}
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
