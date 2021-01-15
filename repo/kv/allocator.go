package kv

import (
	"encoding/binary"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2020 - Nov - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Allocator struct {
	blocks []*pools.ByteBuffer
}

func NewAllocator() *Allocator {
	return &Allocator{
		blocks: make([]*pools.ByteBuffer, 0, 8),
	}
}

func (bk *Allocator) GenKey(v ...interface{}) []byte {
	b := pools.Buffer.GetLen(getSize(v...))
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

	bk.blocks = append(bk.blocks, b)
	return *b.Bytes()
}

func (bk *Allocator) GenValue(m proto.Message) []byte {
	mo := proto.MarshalOptions{UseCachedSize: true}
	bb := pools.Buffer.GetCap(mo.Size(m))
	b, _ := mo.MarshalAppend(*bb.Bytes(), m)
	bb.SetBytes(&b)
	bk.blocks = append(bk.blocks, bb)
	return *bb.Bytes()
}

func (bk *Allocator) ReleaseAll() {
	for _, b := range bk.blocks {
		pools.Buffer.Put(b)
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
