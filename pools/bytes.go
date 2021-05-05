package pools

import (
	"google.golang.org/protobuf/proto"
	"io"
	"sync"
)

/*
   Creation Time: 2020 - Mar - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var TinyBytes = NewByteSlice(4, 128)
var Bytes = NewByteSlice(32, 64<<10)
var Buffer = NewByteBuffer(4, 64<<10)

const (
	bitSize       = 32 << (^uint(0) >> 63)
	maxintHeadBit = 1 << (bitSize - 2)
)

// byteSlicePool contains logic of reusing objects distinguishable by size in generic
// way.
type byteSlicePool struct {
	pool map[int]*sync.Pool
}

// NewByteSlice creates new byteSlicePool that reuses objects which size is in logarithmic range
// [min, max].
//
// Note that it is a shortcut for Custom() constructor with Options provided by
// WithLogSizeMapping() and WithLogSizeRange(min, max) calls.
func NewByteSlice(min, max int) *byteSlicePool {
	p := &byteSlicePool{
		pool: make(map[int]*sync.Pool),
	}
	logarithmicRange(min, max, func(n int) {
		p.pool[n] = &sync.Pool{}
	})
	return p
}

// Get returns probably reused slice of bytes with at least capacity of c and
// exactly len of n.
func (p *byteSlicePool) Get(n, c int) []byte {
	if n > c {
		panic("requested length is greater than capacity")
	}

	size := ceilToPowerOfTwo(c)
	if pool := p.pool[size]; pool != nil {
		v := pool.Get()
		if v != nil {
			bts := v.([]byte)
			bts = bts[:n]
			return bts
		} else {
			return make([]byte, n, size)
		}
	}

	return make([]byte, n, c)
}

// Put returns given slice to reuse pool.
// It does not reuse bytes whose size is not power of two or is out of pool
// min/max range.
func (p *byteSlicePool) Put(bts []byte) {
	if pool := p.pool[cap(bts)]; pool != nil {
		pool.Put(bts)
	}
}

// GetCap returns probably reused slice of bytes with at least capacity of n.
func (p *byteSlicePool) GetCap(c int) []byte {
	return p.Get(0, c)
}

// GetLen returns probably reused slice of bytes with at least capacity of n
// and exactly len of n.
func (p *byteSlicePool) GetLen(n int) []byte {
	return p.Get(n, n)
}

type ByteBuffer struct {
	ri int
	b  []byte
}

func (bb *ByteBuffer) Read(p []byte) (n int, err error) {
	if bb.ri == len(bb.b) - 1 {
		return 0, io.EOF
	}
	n = copy(p, bb.b[bb.ri:])
	return n, nil
}

func newByteBuffer(n, c int) *ByteBuffer {
	if n > c {
		panic("requested length is greater than capacity")
	}
	return &ByteBuffer{b: make([]byte, n, c)}
}

func (bb *ByteBuffer) Reset() {
	bb.ri = 0
	bb.b = bb.b[:bb.ri]
}

func (bb *ByteBuffer) Bytes() *[]byte {
	return &bb.b
}

func (bb *ByteBuffer) SetBytes(b *[]byte) {
	if b == nil {
		return
	}
	bb.b = *b
}

func (bb *ByteBuffer) Fill(data []byte, start, end int) {
	copy(bb.b[start:end], data)
}

func (bb *ByteBuffer) Copy(data []byte) {
	copy(bb.b, data)
}

func (bb *ByteBuffer) Append(data []byte) {
	bb.b = append(bb.b, data...)
}

func (bb *ByteBuffer) Len() int {
	return len(bb.b)
}

func (bb ByteBuffer) Cap() int {
	return cap(bb.b)
}

// byteBufferPool. contains logic of reusing objects distinguishable by size in generic
// way.
type byteBufferPool struct {
	pool map[int]*sync.Pool
}

// NewByteBuffer creates new byteBufferPool that reuses objects which size is in logarithmic range
// [min, max].
//
// Note that it is a shortcut for Custom() constructor with Options provided by
// WithLogSizeMapping() and WithLogSizeRange(min, max) calls.
func NewByteBuffer(min, max int) *byteBufferPool {
	p := &byteBufferPool{
		pool: make(map[int]*sync.Pool),
	}
	logarithmicRange(min, max, func(n int) {
		p.pool[n] = &sync.Pool{}
	})
	return p
}

// Get returns probably reused slice of bytes with at least capacity of c and
// exactly len of n.
func (p *byteBufferPool) Get(n, c int) *ByteBuffer {
	if n > c {
		panic("requested length is greater than capacity")
	}

	size := ceilToPowerOfTwo(c)
	if pool := p.pool[size]; pool != nil {
		v := pool.Get()
		if v != nil {
			bb := v.(*ByteBuffer)
			bb.b = bb.b[:n]
			return bb
		} else {
			return newByteBuffer(n, size)
		}
	}

	return newByteBuffer(n, c)
}

// Put returns given slice to reuse pool.
// It does not reuse bytes whose size is not power of two or is out of pool
// min/max range.
func (p *byteBufferPool) Put(bb *ByteBuffer) {
	if pool := p.pool[cap(bb.b)]; pool != nil {
		pool.Put(bb)
	}
}

// GetCap returns probably reused slice of bytes with at least capacity of n.
func (p *byteBufferPool) GetCap(c int) *ByteBuffer {
	return p.Get(0, c)
}

// GetLen returns probably reused slice of bytes with at least capacity of n
// and exactly len of n.
func (p *byteBufferPool) GetLen(n int) *ByteBuffer {
	return p.Get(n, n)
}

func (p *byteBufferPool) FromProto(m proto.Message) *ByteBuffer {
	mo := proto.MarshalOptions{UseCachedSize: true}
	buf := p.GetCap(mo.Size(m))
	bb, _ := mo.MarshalAppend(*buf.Bytes(), m)
	buf.SetBytes(&bb)
	return buf
}

func (p *byteBufferPool) FromBytes(b []byte) *ByteBuffer {
	buf := p.GetCap(len(b))
	buf.Append(b)
	return buf
}

// logarithmicRange iterates from ceil to power of two min to max,
// calling cb on each iteration.
func logarithmicRange(min, max int, cb func(int)) {
	if min == 0 {
		min = 1
	}
	for n := ceilToPowerOfTwo(min); n <= max; n <<= 1 {
		cb(n)
	}
}

// ceilToPowerOfTwo returns the least power of two integer value greater than
// or equal to n.
func ceilToPowerOfTwo(n int) int {
	if n&maxintHeadBit != 0 && n > maxintHeadBit {
		panic("argument is too large")
	}
	if n <= 2 {
		return n
	}
	n--
	n = fillBits(n)
	n++
	return n
}

func fillBits(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n
}
