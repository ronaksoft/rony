package kv

import (
	"encoding/binary"
	"github.com/ronaksoft/rony/pools"
)

/*
   Creation Time: 2020 - Nov - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type BulkKey struct {
	keys [][]byte
}

func NewBulkKey() *BulkKey {
	return &BulkKey{}
}

func (bk *BulkKey) GenKey(v ...interface{}) []byte {
	b := pools.TinyBytes.GetLen(getSize(v...))
	idx := 0
	for _, x := range v {
		switch x := x.(type) {
		case int:
			binary.BigEndian.PutUint64(b[idx:idx+8], uint64(x))
			idx += 8
		case uint:
			binary.BigEndian.PutUint64(b[idx:idx+8], uint64(x))
			idx += 8
		case int64:
			binary.BigEndian.PutUint64(b[idx:idx+8], uint64(x))
			idx += 8
		case uint64:
			binary.BigEndian.PutUint64(b[idx:idx+8], x)
			idx += 8
		case int32:
			binary.BigEndian.PutUint32(b[idx:idx+4], uint32(x))
			idx += 4
		case uint32:
			binary.BigEndian.PutUint32(b[idx:idx+4], x)
			idx += 4
		case []byte:
			copy(b[idx:], x)
			idx += len(x)
		case string:
			copy(b[idx:], x)
			idx += len(x)
		default:
			panic("unsupported type")
		}
	}

	bk.keys = append(bk.keys, b)
	return b
}

func (bk *BulkKey) ReleaseAll() {
	for _, b := range bk.keys {
		pools.TinyBytes.Put(b)
	}
}

func getSize(v ...interface{}) int {
	s := 0
	for _, x := range v {
		switch x := x.(type) {
		case int64, uint64, int, uint:
			s += 8
		case int32, uint32:
			s += 4
		case []byte:
			s += len(x)
		case string:
			s += len(x)
		default:
			panic("unsupported type")
		}
	}
	return s
}
