package badgerRaft

import (
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - May - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Txn struct {
	ID     int64
	store  *Store
	update bool
}

func (txn *Txn) Delete(alloc *tools.Allocator, keyParts ...interface{}) error {
	req := PoolDelete.Get()
	defer PoolDelete.Put(req)
	key := alloc.Gen(keyParts...)
	req.Key = append(req.Key, key...)
	req.TxnID = txn.ID
	b := pools.Buffer.FromProto(req)
	f := txn.store.c.RaftApply(*b.Bytes())
	pools.Buffer.Put(b)

	err := f.Error()
	if err != nil {
		return err
	}

	res := f.Response()
	if res != nil {
		switch x := res.(type) {
		case nil:
		case error:
			return x
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (txn *Txn) Set(alloc *tools.Allocator, val []byte, keyParts ...interface{}) error {
	req := PoolSet.Get()
	defer PoolSet.Put(req)
	key := alloc.Gen(keyParts...)
	req.KV.Key = append(req.KV.Key, key...)
	req.KV.Value = append(req.KV.Value, val...)
	req.TxnID = txn.ID
	b := pools.Buffer.FromProto(req)
	f := txn.store.c.RaftApply(*b.Bytes())
	pools.Buffer.Put(b)

	err := f.Error()
	if err != nil {
		return err
	}

	res := f.Response()
	if res != nil {
		switch x := res.(type) {
		case nil:
		case error:
			return x
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (txn *Txn) Get(alloc *tools.Allocator, keyParts ...interface{}) ([]byte, error) {
	req := PoolGet.Get()
	defer PoolGet.Put(req)
	key := alloc.Gen(keyParts...)
	req.Key = append(req.Key, key...)
	req.TxnID = txn.ID
	b := pools.Buffer.FromProto(req)
	f := txn.store.c.RaftApply(*b.Bytes())
	pools.Buffer.Put(b)

	err := f.Error()
	if err != nil {
		return nil, err
	}

	res := f.Response()
	if res == nil {
		return nil, ErrUnknown
	}

	switch x := res.(type) {
	case error:
		return nil, x
	case *KeyValue:
		return x.Value, nil
	default:
		return nil, ErrUnknown
	}

}

func (txn *Txn) Exists(alloc *tools.Allocator, keyParts ...interface{}) bool {
	_, err := txn.Get(alloc, keyParts...)
	return err == nil
}
