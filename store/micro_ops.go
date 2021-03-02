package store

import (
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2021 - Mar - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func Delete(txn *Txn, alloc *Allocator, keyParts ...interface{}) error {
	key := alloc.Gen(keyParts...)
	return txn.Delete(key)
}

func Move(txn *Txn, oldKey, newKey []byte) error {
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

func Set(txn *Txn, alloc *Allocator, val []byte, keyParts ...interface{}) error {
	key := alloc.Gen(keyParts...)
	return txn.Set(key, val)
}

func Get(txn *Txn, alloc *Allocator, keyParts ...interface{}) ([]byte, error) {
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

func GetByKey(txn *Txn, alloc *Allocator, key []byte) ([]byte, error) {
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

func Exists(txn *Txn, alloc *Allocator, keyParts ...interface{}) bool {
	_, err := Get(txn, alloc, keyParts...)
	if err != nil && err == ErrKeyNotFound {
		return false
	}
	return true
}

func Marshal(txn *Txn, alloc *Allocator, m proto.Message, keyParts ...interface{}) error {
	val := alloc.Marshal(m)
	return Set(txn, alloc, val, keyParts...)
}

func Unmarshal(txn *Txn, alloc *Allocator, m proto.Message, keyParts ...interface{}) error {
	val, err := Get(txn, alloc, keyParts...)
	if err != nil {
		return err
	}
	return proto.Unmarshal(val, m)
}

func UnmarshalMerge(txn *Txn, alloc *Allocator, m proto.Message, keyParts ...interface{}) error {
	val, err := Get(txn, alloc, keyParts...)
	if err != nil {
		return err
	}
	umo := proto.UnmarshalOptions{Merge: true}
	return umo.Unmarshal(val, m)
}
