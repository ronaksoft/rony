package tools

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"sync"
)

/*
   Creation Time: 2019 - Oct - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var poolSha512 = sync.Pool{
	New: func() interface{} {
		return sha512.New()
	},
}

// Sha512 appends a 64bytes array which is sha512(in) to out
func Sha512(in, out []byte) error {
	h := poolSha512.Get().(hash.Hash)
	if _, err := h.Write(in); err != nil {
		h.Reset()
		poolSha512.Put(h)
		return err
	}
	h.Sum(out)
	h.Reset()
	poolSha512.Put(h)
	return nil
}

func MustSha512(in, out []byte) {
	err := Sha512(in, out)
	if err != nil {
		panic(err)
	}
}

var poolSha256 = sync.Pool{
	New: func() interface{} {
		return sha256.New()
	},
}

// Sha256 appends a 32bytes array which is sha256(in) to out
func Sha256(in, out []byte) error {
	h := poolSha256.Get().(hash.Hash)
	if _, err := h.Write(in); err != nil {
		h.Reset()
		poolSha256.Put(h)
		return err
	}
	h.Sum(out)
	h.Reset()
	poolSha256.Put(h)
	return nil
}

func MustSha256(in, out []byte) {
	err := Sha256(in, out)
	if err != nil {
		panic(err)
	}
}
