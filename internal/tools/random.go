package tools

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/gobwas/pool/pbytes"
	mathRand "math/rand"
	"sync"
	"time"
)

/*
   Creation Time: 2019 - Oct - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	bufSize       = 1 << 14 // 16KB
	DIGITS        = "0123456789"
	ALPHANUMERICS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

type randomGenerator struct {
	sync.Pool
}

func (rg *randomGenerator) GetRand() *mathRand.Rand {
	return rg.Get().(*mathRand.Rand)
}

func (rg *randomGenerator) PutRand(r *mathRand.Rand) {
	rg.Put(r)
}

var rndGen randomGenerator

func init() {
	rndGen.New = func() interface{} {
		x := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
		return x
	}
}

// RandomID generates a pseudo-random string with length 'n' which characters are alphanumerics.
func RandomID(n int) string {
	rnd := rndGen.GetRand()
	defer rndGen.PutRand(rnd)
	b := make([]byte, n)
	for i := range b {
		b[i] = ALPHANUMERICS[rnd.Intn(len(ALPHANUMERICS))]
	}
	return ByteToStr(b)
}

// RandomDigit generates a pseudo-random string with length 'n' which characters are only digits (0-9)
func RandomDigit(n int) string {
	rnd := rndGen.GetRand()
	defer rndGen.PutRand(rnd)
	b := make([]byte, n)
	for i := 0; i < len(b); i++ {
		b[i] = DIGITS[rnd.Intn(len(DIGITS))]
	}
	return ByteToStr(b)
}

// RandomInt64 produces a pseudo-random number, if n == 0 there will be no limit otherwise
// the output will be smaller than n
func RandomInt64(n int64) (x int64) {
	rnd := rndGen.GetRand()
	if n == 0 {
		x = rnd.Int63()
	} else {
		x = rnd.Int63n(n)
	}
	rndGen.PutRand(rnd)
	return
}

func SecureRandomInt63(n int64) (x int64) {
	b := pbytes.GetLen(8)
	_, _ = rand.Read(b)
	xx := binary.BigEndian.Uint64(b)
	if n > 0 {
		x = int64(xx) % n
	} else {
		x = int64(xx >> 1)
	}
	return
}

func RandomInt(n int) (x int) {
	rnd := rndGen.GetRand()

	if n == 0 {
		x = rnd.Int()
	} else {
		x = rnd.Intn(n)
	}
	rndGen.PutRand(rnd)
	return
}

// RandUint64 produces a pseudo-random unsigned number
func RandomUint64() (x uint64) {
	rnd := rndGen.GetRand()
	x = rnd.Uint64()
	rndGen.PutRand(rnd)
	return
}
