package store

import (
	"github.com/ronaksoft/rony/tools"
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

func Delete(txn *LTxn, alloc *tools.Allocator, keyParts ...interface{}) error {
	key := alloc.Gen(keyParts...)
	return txn.Delete(key)
}

func Move(txn *LTxn, oldKey, newKey []byte) error {
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

func Set(txn *LTxn, alloc *tools.Allocator, val []byte, keyParts ...interface{}) error {
	key := alloc.Gen(keyParts...)
	return txn.Set(key, val)
}

func Get(txn *LTxn, alloc *tools.Allocator, keyParts ...interface{}) ([]byte, error) {
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

func GetByKey(txn *LTxn, alloc *tools.Allocator, key []byte) ([]byte, error) {
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

func Exists(txn *LTxn, alloc *tools.Allocator, keyParts ...interface{}) bool {
	_, err := Get(txn, alloc, keyParts...)
	if err != nil && err == ErrKeyNotFound {
		return false
	}
	return true
}

func Marshal(txn *LTxn, alloc *tools.Allocator, m proto.Message, keyParts ...interface{}) error {
	val := alloc.Marshal(m)
	return Set(txn, alloc, val, keyParts...)
}

func Unmarshal(txn *LTxn, alloc *tools.Allocator, m proto.Message, keyParts ...interface{}) error {
	val, err := Get(txn, alloc, keyParts...)
	if err != nil {
		return err
	}
	return proto.Unmarshal(val, m)
}

func UnmarshalMerge(txn *LTxn, alloc *tools.Allocator, m proto.Message, keyParts ...interface{}) error {
	val, err := Get(txn, alloc, keyParts...)
	if err != nil {
		return err
	}
	umo := proto.UnmarshalOptions{Merge: true}
	return umo.Unmarshal(val, m)
}
