package badgerStore

import (
	"github.com/ronaksoft/rony/internal/store"
	"github.com/ronaksoft/rony/pools"
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
	ID    int64
	store *Store
}

func (txn *Txn) Delete(alloc *store.Allocator, keyParts ...interface{}) error {
	req := PoolDelete.Get()
	defer PoolDelete.Put(req)
	key := alloc.Gen(keyParts...)
	req.Keys = append(req.Keys, key)
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

func (txn *Txn) Set(alloc *store.Allocator, val []byte, keyParts ...interface{}) error {
	req := PoolSet.Get()
	defer PoolSet.Put(req)
	key := alloc.Gen(keyParts...)
	req.KVs = append(req.KVs, &KeyValue{
		Key:   key,
		Value: val,
	})
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

func (txn *Txn) Get(alloc *store.Allocator, keyParts ...interface{}) ([]byte, error) {
	item, err := txn.Get(alloc.Gen(keyParts...))
	if err != nil {
		return nil, err
	}

	var b []byte
	_ = item.Value(func(val []byte) error {
		b = alloc.FillWith(val)
		return nil
	})
	return b, nil
}

func (txn *Txn) GetByKey(alloc *store.Allocator, key []byte) ([]byte, error) {
	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	var b []byte
	_ = item.Value(func(val []byte) error {
		b = alloc.FillWith(val)
		return nil
	})
	return b, nil
}

func (txn *Txn) Exists(alloc *store.Allocator, keyParts ...interface{}) bool {
	_, err := txn.store.db.Get(txn, alloc, keyParts...)
	if err != nil && err == ErrKeyNotFound {
		return false
	}
	return true
}
