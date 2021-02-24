// Code generated by Rony's protoc plugin; DO NOT EDIT.

package rony

import (
	registry "github.com/ronaksoft/rony/registry"
	store "github.com/ronaksoft/rony/store"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_GetNodes int64 = 362407405

type poolGetNodes struct {
	pool sync.Pool
}

func (p *poolGetNodes) Get() *GetNodes {
	x, ok := p.pool.Get().(*GetNodes)
	if !ok {
		return &GetNodes{}
	}
	return x
}

func (p *poolGetNodes) Put(x *GetNodes) {
	x.ReplicaSet = x.ReplicaSet[:0]
	p.pool.Put(x)
}

var PoolGetNodes = poolGetNodes{}

func (x *GetNodes) DeepCopy(z *GetNodes) {
	z.ReplicaSet = append(z.ReplicaSet[:0], x.ReplicaSet...)
}

func (x *GetNodes) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *GetNodes) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_GetPage int64 = 3721890413

type poolGetPage struct {
	pool sync.Pool
}

func (p *poolGetPage) Get() *GetPage {
	x, ok := p.pool.Get().(*GetPage)
	if !ok {
		return &GetPage{}
	}
	return x
}

func (p *poolGetPage) Put(x *GetPage) {
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

const C_TunnelMessage int64 = 3271476222

type poolTunnelMessage struct {
	pool sync.Pool
}

func (p *poolTunnelMessage) Get() *TunnelMessage {
	x, ok := p.pool.Get().(*TunnelMessage)
	if !ok {
		return &TunnelMessage{}
	}
	return x
}

func (p *poolTunnelMessage) Put(x *TunnelMessage) {
	x.SenderID = x.SenderID[:0]
	x.SenderReplicaSet = 0
	x.Store = x.Store[:0]
	if x.Envelope != nil {
		PoolMessageEnvelope.Put(x.Envelope)
		x.Envelope = nil
	}
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
		z.Envelope = PoolMessageEnvelope.Get()
		x.Envelope.DeepCopy(z.Envelope)
	}
}

func (x *TunnelMessage) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *TunnelMessage) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_RaftCommand int64 = 2919813429

type poolRaftCommand struct {
	pool sync.Pool
}

func (p *poolRaftCommand) Get() *RaftCommand {
	x, ok := p.pool.Get().(*RaftCommand)
	if !ok {
		return &RaftCommand{}
	}
	return x
}

func (p *poolRaftCommand) Put(x *RaftCommand) {
	x.Sender = x.Sender[:0]
	x.Store = x.Store[:0]
	if x.Envelope != nil {
		PoolMessageEnvelope.Put(x.Envelope)
		x.Envelope = nil
	}
	p.pool.Put(x)
}

var PoolRaftCommand = poolRaftCommand{}

func (x *RaftCommand) DeepCopy(z *RaftCommand) {
	z.Sender = append(z.Sender[:0], x.Sender...)
	for idx := range x.Store {
		if x.Store[idx] != nil {
			xx := PoolKeyValue.Get()
			x.Store[idx].DeepCopy(xx)
			z.Store = append(z.Store, xx)
		}
	}
	if x.Envelope != nil {
		z.Envelope = PoolMessageEnvelope.Get()
		x.Envelope.DeepCopy(z.Envelope)
	}
}

func (x *RaftCommand) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *RaftCommand) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_EdgeNode int64 = 999040174

type poolEdgeNode struct {
	pool sync.Pool
}

func (p *poolEdgeNode) Get() *EdgeNode {
	x, ok := p.pool.Get().(*EdgeNode)
	if !ok {
		return &EdgeNode{}
	}
	return x
}

func (p *poolEdgeNode) Put(x *EdgeNode) {
	x.ServerID = x.ServerID[:0]
	x.ReplicaSet = 0
	x.RaftPort = 0
	x.RaftState = 0
	x.GatewayAddr = x.GatewayAddr[:0]
	x.TunnelAddr = x.TunnelAddr[:0]
	p.pool.Put(x)
}

var PoolEdgeNode = poolEdgeNode{}

func (x *EdgeNode) DeepCopy(z *EdgeNode) {
	z.ServerID = append(z.ServerID[:0], x.ServerID...)
	z.ReplicaSet = x.ReplicaSet
	z.RaftPort = x.RaftPort
	z.RaftState = x.RaftState
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
		return &Page{}
	}
	return x
}

func (p *poolPage) Put(x *Page) {
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
	registry.RegisterConstructor(362407405, "GetNodes")
	registry.RegisterConstructor(3721890413, "GetPage")
	registry.RegisterConstructor(3271476222, "TunnelMessage")
	registry.RegisterConstructor(2919813429, "RaftCommand")
	registry.RegisterConstructor(999040174, "EdgeNode")
	registry.RegisterConstructor(3023575326, "Page")
}

func SavePageWithTxn(txn *store.Txn, alloc *store.Allocator, m *Page) (err error) {
	if alloc == nil {
		alloc = store.NewAllocator()
		defer alloc.ReleaseAll()
	}

	// save entry
	b := alloc.GenValue(m)
	key := alloc.GenKey('M', C_Page, 299066170, m.ID)
	err = txn.Set(key, b)
	if err != nil {
		return
	}

	// save entry for view[ReplicaSet ID]
	err = txn.Set(alloc.GenKey('M', C_Page, 1040696757, m.ReplicaSet, m.ID), b)
	if err != nil {
		return
	}

	return

}

func SavePage(m *Page) error {
	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *store.Txn) error {
		return SavePageWithTxn(txn, alloc, m)
	})
}

func ReadPageWithTxn(txn *store.Txn, alloc *store.Allocator, id uint32, m *Page) (*Page, error) {
	if alloc == nil {
		alloc = store.NewAllocator()
		defer alloc.ReleaseAll()
	}

	item, err := txn.Get(alloc.GenKey('M', C_Page, 299066170, id))
	if err != nil {
		return nil, err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	return m, err
}

func ReadPage(id uint32, m *Page) (*Page, error) {
	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &Page{}
	}

	err := store.View(func(txn *store.Txn) (err error) {
		m, err = ReadPageWithTxn(txn, alloc, id, m)
		return err
	})
	return m, err
}

func ReadPageByReplicaSetAndIDWithTxn(txn *store.Txn, alloc *store.Allocator, replicaSet uint64, id uint32, m *Page) (*Page, error) {
	if alloc == nil {
		alloc = store.NewAllocator()
		defer alloc.ReleaseAll()
	}

	item, err := txn.Get(alloc.GenKey('M', C_Page, 1040696757, replicaSet, id))
	if err != nil {
		return nil, err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	return m, err
}

func ReadPageByReplicaSetAndID(replicaSet uint64, id uint32, m *Page) (*Page, error) {
	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()
	if m == nil {
		m = &Page{}
	}
	err := store.View(func(txn *store.Txn) (err error) {
		m, err = ReadPageByReplicaSetAndIDWithTxn(txn, alloc, replicaSet, id, m)
		return err
	})
	return m, err
}

func DeletePageWithTxn(txn *store.Txn, alloc *store.Allocator, id uint32) error {
	m := &Page{}
	item, err := txn.Get(alloc.GenKey('M', C_Page, 299066170, id))
	if err != nil {
		return err
	}
	err = item.Value(func(val []byte) error {
		return m.Unmarshal(val)
	})
	if err != nil {
		return err
	}
	err = txn.Delete(alloc.GenKey('M', C_Page, 299066170, m.ID))
	if err != nil {
		return err
	}

	err = txn.Delete(alloc.GenKey('M', C_Page, 1040696757, m.ReplicaSet, m.ID))
	if err != nil {
		return err
	}

	return nil
}

func DeletePage(id uint32) error {
	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *store.Txn) error {
		return DeletePageWithTxn(txn, alloc, id)
	})
}

func ListPage(
	offsetID uint32, lo *store.ListOption, cond func(m *Page) bool,
) ([]*Page, error) {
	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Page, 0, lo.Limit())
	err := store.View(func(txn *store.Txn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey('M', C_Page, 299066170)
		opt.Reverse = lo.Backward()
		osk := alloc.GenKey('M', C_Page, 299066170, offsetID)
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
				}
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func IterPages(txn *store.Txn, alloc *store.Allocator, cb func(m *Page) bool) error {
	if alloc == nil {
		alloc = store.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	iterOpt := store.DefaultIteratorOptions
	iterOpt.Prefix = alloc.GenKey('M', C_Page, 299066170)
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

func ListPageByReplicaSet(replicaSet uint64, offsetID uint32, lo *store.ListOption) ([]*Page, error) {
	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()

	res := make([]*Page, 0, lo.Limit())
	err := store.View(func(txn *store.Txn) error {
		opt := store.DefaultIteratorOptions
		opt.Prefix = alloc.GenKey('M', C_Page, 1040696757, replicaSet)
		opt.Reverse = lo.Backward()
		osk := alloc.GenKey('M', C_Page, 1040696757, replicaSet, offsetID)
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
				res = append(res, m)
				return nil
			})
		}
		iter.Close()
		return nil
	})
	return res, err
}

func IterPageByReplicaSet(txn *store.Txn, alloc *store.Allocator, replicaSet uint64, cb func(m *Page) bool) error {
	if alloc == nil {
		alloc = store.NewAllocator()
		defer alloc.ReleaseAll()
	}

	exitLoop := false
	opt := store.DefaultIteratorOptions
	opt.Prefix = alloc.GenKey('M', C_Page, 1040696757, replicaSet)
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
