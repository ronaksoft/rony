package tools

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/ronaksoft/rony/pools/buf"

	"github.com/ronaksoft/rony/pools"
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

var (
	prefix = []byte{0xFF}
)

type Allocator struct {
	blocks []*buf.Bytes
}

func NewAllocator() *Allocator {
	return &Allocator{
		blocks: make([]*buf.Bytes, 0, 8),
	}
}

// Gen acquired a byte slice fitted to hold all the v variables.
func (bk *Allocator) Gen(v ...interface{}) []byte {
	b := pools.Buffer.GetLen(1 + getSize(v...))
	var buf [8]byte
	b.CopyFrom(prefix)
	idx := 1
	for _, x := range v {
		t := reflect.TypeOf(x)
		switch t.Kind() {
		case reflect.Int:
			binary.BigEndian.PutUint64(buf[:], uint64(reflect.ValueOf(x).Int()))
			b.Fill(buf[:], idx, idx+8)
			idx += 8
		case reflect.Uint:
			binary.BigEndian.PutUint64(buf[:], reflect.ValueOf(x).Uint())
			b.Fill(buf[:], idx, idx+8)
			idx += 8
		case reflect.Int64:
			binary.BigEndian.PutUint64(buf[:], uint64(reflect.ValueOf(x).Int()))
			b.Fill(buf[:], idx, idx+8)
			idx += 8
		case reflect.Uint64:
			binary.BigEndian.PutUint64(buf[:], reflect.ValueOf(x).Uint())
			b.Fill(buf[:], idx, idx+8)
			idx += 8
		case reflect.Int32:
			binary.BigEndian.PutUint32(buf[:4], uint32(reflect.ValueOf(x).Int()))
			b.Fill(buf[:4], idx, idx+4)
			idx += 4
		case reflect.Uint32:
			binary.BigEndian.PutUint32(buf[:4], uint32(reflect.ValueOf(x).Uint()))
			b.Fill(buf[:4], idx, idx+4)
			idx += 4
		case reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.Uint8:
				xb := reflect.ValueOf(x).Bytes()
				b.Fill(reflect.ValueOf(x).Bytes(), idx, idx+len(xb))
				idx += len(xb)
			default:
				panic(fmt.Sprintf("unsupported slice type: %s", t.Elem().Kind().String()))
			}

		case reflect.String:
			xb := StrToByte(reflect.ValueOf(x).String())
			b.Fill(xb, idx, idx+len(xb))
			idx += len(xb)
		default:
			panic("unsupported type")
		}
	}

	bk.blocks = append(bk.blocks, b)

	return *b.Bytes()
}

// Marshal acquires a byte slice fitted for message 'm'
func (bk *Allocator) Marshal(m proto.Message) []byte {
	buf := pools.Buffer.FromProto(m)
	bk.blocks = append(bk.blocks, buf)

	return *buf.Bytes()
}

// FillWith acquired a byte slice with the capacity of 'v' and append/copy v into it.
func (bk *Allocator) FillWith(v []byte) []byte {
	b := pools.Buffer.GetCap(len(v))
	b.AppendFrom(v)
	bk.blocks = append(bk.blocks, b)

	return *b.Bytes()
}

// ReleaseAll releases all the byte slices.
func (bk *Allocator) ReleaseAll() {
	for _, b := range bk.blocks {
		pools.Buffer.Put(b)
	}
	bk.blocks = bk.blocks[:0]
}

func getSize(v ...interface{}) int {
	s := 0
	for _, x := range v {
		t := reflect.TypeOf(x)
		switch t.Kind() {
		case reflect.Int64, reflect.Uint64, reflect.Int, reflect.Uint:
			s += 8
		case reflect.Int32, reflect.Uint32:
			s += 4
		case reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.Uint8:
				xb := reflect.ValueOf(x).Bytes()
				s += len(xb)
			default:
				panic(fmt.Sprintf("unsupported slice type: %s", t.Elem().Kind().String()))
			}
		case reflect.String:
			xb := StrToByte(reflect.ValueOf(x).String())
			s += len(xb)
		default:
			panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(x).Kind()))
		}
	}

	return s
}
