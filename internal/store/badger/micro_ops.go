package badgerStore

import (
	"github.com/ronaksoft/rony/internal/store"
)

/*
   Creation Time: 2021 - Mar - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func (fsm *Store) Delete(txn *Txn, alloc *store.Allocator, keyParts ...interface{}) error {
	key := alloc.Gen(keyParts...)
	return txn.Delete(key)
}

func (fsm *Store) Move(txn *Txn, oldKey, newKey []byte) error {
	item, err := txn.Get(oldKey)
	if err != nil {
		return err
	}
	err = item.Value(func(val []byte) error {
		return txn.Set(newKey, val)
	})
	if err != nil {
		return err
	}
	return txn.Delete(oldKey)
}

func (fsm *Store) Set(txn *Txn, alloc *store.Allocator, val []byte, keyParts ...interface{}) error {
	key := alloc.Gen(keyParts...)
	return txn.Set(key, val)
}

func (fsm *Store) Get(txn *Txn, alloc *store.Allocator, keyParts ...interface{}) ([]byte, error) {
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

func (fsm *Store) GetByKey(txn *Txn, alloc *store.Allocator, key []byte) ([]byte, error) {
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

func (fsm *Store) Exists(txn *Txn, alloc *store.Allocator, keyParts ...interface{}) bool {
	_, err := fsm.Get(txn, alloc, keyParts...)
	if err != nil && err == ErrKeyNotFound {
		return false
	}
	return true
}
