package kv

import (
	"encoding/binary"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
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
	keys []*pools.ByteBuffer
}

func NewBulkKey() *BulkKey {
	return &BulkKey{}
}

func (bk *BulkKey) GenKey(v ...interface{}) []byte {
	b := pools.BytesBuffer.GetLen(getSize(v...))
	var buf [8]byte
	idx := 0
	for _, x := range v {
		switch x := x.(type) {
		case int:
			binary.BigEndian.PutUint64(buf[:], uint64(x))
			b.Fill(buf[:], idx, idx+8)
			idx += 8
		case uint:
			binary.BigEndian.PutUint64(buf[:], uint64(x))
			b.Fill(buf[:], idx, idx+8)
			idx += 8
		case int64:
			binary.BigEndian.PutUint64(buf[:], uint64(x))
			b.Fill(buf[:], idx, idx+8)
			idx += 8
		case uint64:
			binary.BigEndian.PutUint64(buf[:], x)
			b.Fill(buf[:], idx, idx+8)
			idx += 8
		case int32:
			binary.BigEndian.PutUint32(buf[:4], uint32(x))
			b.Fill(buf[:4], idx, idx+4)
			idx += 4
		case uint32:
			binary.BigEndian.PutUint32(buf[:4], x)
			b.Fill(buf[:4], idx, idx+4)
			idx += 4
		case []byte:
			b.Fill(x, idx, idx+len(x))
			idx += len(x)
		case string:
			b.Fill(tools.StrToByte(x), idx, idx+len(x))
			idx += len(x)
		default:
			panic("unsupported type")
		}
	}

	bk.keys = append(bk.keys, b)
	return *b.Bytes()
}

func (bk *BulkKey) ReleaseAll() {
	for _, b := range bk.keys {
		pools.BytesBuffer.Put(b)
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
