package gossipCluster

import (
	"bytes"
	"encoding/binary"
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/raft"
	"github.com/pkg/errors"
	"github.com/ronaksoft/rony/store"
)

/*
   Creation Time: 2021 - Feb - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	// Prefix names to distinguish between logs and conf
	prefixLogs = []byte{0x0}
	prefixConf = []byte{0x1}

	// ErrKeyNotFound is an error indicating a given key does not exist
	ErrKeyNotFound = errors.New("not found")
)

// BadgerStore provides access to Badger for Raft to store and retrieve
// log entries. It also provides key/value storage, and can be used as
// a LogStore and StableStore.
type BadgerStore struct{}

// FirstIndex returns the first known index from the Raft log.
func (b *BadgerStore) FirstIndex() (uint64, error) {
	var value uint64
	err := store.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Reverse:        false,
		})
		defer it.Close()

		it.Seek(prefixLogs)
		if it.ValidForPrefix(prefixLogs) {
			value = bytesToUint64(it.Item().Key()[1:])
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return value, nil
}

// LastIndex returns the last known index from the Raft log.
func (b *BadgerStore) LastIndex() (uint64, error) {
	var value uint64
	err := store.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Reverse:        true,
		})
		defer it.Close()

		it.Seek(append(prefixLogs, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff))
		if it.ValidForPrefix(prefixLogs) {
			value = bytesToUint64(it.Item().Key()[1:])
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return value, nil
}

// GetLog gets a log entry from Badger at a given index.
func (b *BadgerStore) GetLog(index uint64, log *raft.Log) error {
	return store.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(prefixLogs, uint64ToBytes(index)...))
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				return raft.ErrLogNotFound
			default:
				return err
			}
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return decodeMsgPack(val, log)
	})
}

// StoreLog stores a single raft log.
func (b *BadgerStore) StoreLog(log *raft.Log) error {
	val, err := encodeMsgPack(log)
	if err != nil {
		return err
	}
	return store.Update(func(txn *badger.Txn) error {
		return txn.Set(append(prefixLogs, uint64ToBytes(log.Index)...), val.Bytes())
	})
}

// StoreLogs stores a set of raft logs.
func (b *BadgerStore) StoreLogs(logs []*raft.Log) error {
	return store.Update(func(txn *badger.Txn) error {
		for _, log := range logs {
			key := append(prefixLogs, uint64ToBytes(log.Index)...)
			val, err := encodeMsgPack(log)
			if err != nil {
				return err
			}
			if err := txn.Set(key, val.Bytes()); err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteRange deletes logs within a given range inclusively.
func (b *BadgerStore) DeleteRange(min, max uint64) error {
	// we manage the transaction manually in order to avoid ErrTxnTooBig errors
	txn := store.DB().NewTransaction(true)
	it := txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: false,
		Reverse:        false,
	})

	start := append(prefixLogs, uint64ToBytes(min)...)
	for it.Seek(start); it.Valid(); it.Next() {
		key := make([]byte, 16)
		it.Item().KeyCopy(key)
		// Handle out-of-range log index
		if bytesToUint64(key[8:]) > max {
			break
		}
		// Delete in-range log index
		if err := txn.Delete(key); err != nil {
			if err == badger.ErrTxnTooBig {
				it.Close()
				err = txn.Commit()
				if err != nil {
					return err
				}
				return b.DeleteRange(bytesToUint64(key[1:]), max)
			}
			return err
		}
	}
	it.Close()
	err := txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Set is used to set a key/value set outside of the raft log.
func (b *BadgerStore) Set(key []byte, val []byte) error {
	return store.Update(func(txn *badger.Txn) error {
		return txn.Set(append(prefixConf, key...), val)
	})
}

// Get is used to retrieve a value from the k/v store by key
func (b *BadgerStore) Get(key []byte) ([]byte, error) {
	var value []byte
	err := store.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(prefixConf, key...))
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				return ErrKeyNotFound
			default:
				return err
			}
		}
		value, err = item.ValueCopy(value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

// SetUint64 is like Set, but handles uint64 values
func (b *BadgerStore) SetUint64(key []byte, val uint64) error {
	return b.Set(key, uint64ToBytes(val))
}

// GetUint64 is like Get, but handles uint64 values
func (b *BadgerStore) GetUint64(key []byte) (uint64, error) {
	val, err := b.Get(key)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(val), nil
}

// Decode reverses the encode operation on a byte slice input
func decodeMsgPack(buf []byte, out interface{}) error {
	r := bytes.NewBuffer(buf)
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(r, &hd)
	return dec.Decode(out)
}

// Encode writes an encoded object to a new bytes buffer
func encodeMsgPack(in interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	hd := codec.MsgpackHandle{}
	enc := codec.NewEncoder(buf, &hd)
	err := enc.Encode(in)
	return buf, err
}

// Converts bytes to an integer
func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Converts a uint to a byte slice
func uint64ToBytes(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}
