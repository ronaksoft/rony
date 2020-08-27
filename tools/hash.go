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

// Sha256 returns a 32bytes array which is sha256(in)
func Sha256(in []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(in); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

var poolSha512 = sync.Pool{
	New: func() interface{} {
		return sha512.New()
	},
}

// Sha512 returns a 64bytes array which is sha512(in)
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
	return
}

func MustSha256(in []byte) []byte {
	h, err := Sha256(in)
	if err != nil {
		panic(err)
	}
	return h
}
