package buf

import (
	"io"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Bytes struct {
	ri int
	b  []byte
}

func (bb *Bytes) Write(p []byte) (n int, err error) {
	bb.b = append(bb.b, p...)

	return len(p), nil
}

func (bb *Bytes) Read(p []byte) (n int, err error) {
	if bb.ri >= len(bb.b)-1 {
		return 0, io.EOF
	}
	n = copy(p, bb.b[bb.ri:])
	bb.ri += n

	return n, nil
}

func NewBytes(n, c int) *Bytes {
	if n > c {
		panic("requested length is greater than capacity")
	}

	return &Bytes{b: make([]byte, n, c)}
}

func (bb *Bytes) Reset() {
	bb.ri = 0
	bb.b = bb.b[:bb.ri]
}

func (bb *Bytes) Bytes() *[]byte {
	return &bb.b
}

func (bb *Bytes) SetBytes(b *[]byte) {
	if b == nil {
		return
	}
	bb.b = *b
}

func (bb *Bytes) Fill(data []byte, start, end int) {
	copy(bb.b[start:end], data)
}

func (bb *Bytes) CopyFromWithOffset(data []byte, offset int) {
	copy(bb.b[offset:], data)
}

func (bb *Bytes) CopyFrom(data []byte) {
	copy(bb.b, data)
}

func (bb *Bytes) CopyTo(data []byte) []byte {
	copy(data, bb.b)

	return data
}

func (bb *Bytes) AppendFrom(data []byte) {
	bb.b = append(bb.b, data...)
}

func (bb *Bytes) AppendTo(data []byte) []byte {
	data = append(data, bb.b...)

	return data
}

func (bb *Bytes) Len() int {
	return len(bb.b)
}

func (bb Bytes) Cap() int {
	return cap(bb.b)
}

// bytesPool. contains logic of reusing objects distinguishable by size in generic
// way.
type bytesPool struct {
	pool map[int]*sync.Pool
}

// NewBytesPool creates new bytesPool that reuses objects which size is in logarithmic range
// [min, max].
//
// Note that it is a shortcut for Custom() constructor with Options provided by
// WithLogSizeMapping() and WithLogSizeRange(min, max) calls.
func NewBytesPool(min, max int) *bytesPool {
	p := &bytesPool{
		pool: make(map[int]*sync.Pool),
	}
	logarithmicRange(min, max, func(n int) {
		p.pool[n] = &sync.Pool{}
	})

	return p
}

// Get returns probably reused slice of bytes with at least capacity of c and
// exactly len of n.
func (p *bytesPool) Get(n, c int) *Bytes {
	if n > c {
		panic("requested length is greater than capacity")
	}

	size := ceilToPowerOfTwo(c)
	if pool := p.pool[size]; pool != nil {
		v := pool.Get()
		if v != nil {
			bb := v.(*Bytes)
			bb.b = bb.b[:n]

			return bb
		} else {
			return NewBytes(n, size)
		}
	}

	return NewBytes(n, c)
}

// Put returns given slice to reuse pool.
// It does not reuse bytes whose size is not power of two or is out of pool
// min/max range.
func (p *bytesPool) Put(bb *Bytes) {
	if pool := p.pool[cap(bb.b)]; pool != nil {
		pool.Put(bb)
	}
}

// GetCap returns probably reused slice of bytes with at least capacity of n.
func (p *bytesPool) GetCap(c int) *Bytes {
	return p.Get(0, c)
}

// GetLen returns probably reused slice of bytes with at least capacity of n
// and exactly len of n.
func (p *bytesPool) GetLen(n int) *Bytes {
	return p.Get(n, n)
}

func (p *bytesPool) FromProto(m proto.Message) *Bytes {
	mo := proto.MarshalOptions{UseCachedSize: true}
	buf := p.GetCap(mo.Size(m))
	bb, _ := mo.MarshalAppend(*buf.Bytes(), m)
	buf.SetBytes(&bb)

	return buf
}

func (p *bytesPool) FromBytes(b []byte) *Bytes {
	buf := p.GetCap(len(b))
	buf.AppendFrom(b)

	return buf
}

const (
	bitSize       = 32 << (^uint(0) >> 63)
	maxintHeadBit = 1 << (bitSize - 2)
)

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
