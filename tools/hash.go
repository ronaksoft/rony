package tools

import (
	"crypto/sha256"
	"crypto/sha512"
)

/*
   Creation Time: 2019 - Oct - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Sha256 returns a 32bytes array which is sha256(in)
func Sha256(in []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(in); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Sha512 returns a 64bytes array which is sha512(in)
func Sha512(in []byte) ([]byte, error) {
	h := sha512.New()
	if _, err := h.Write(in); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func MustSha512(in []byte) []byte {
	h, err := Sha512(in)
	if err != nil {
		panic(err)
	}
	return h
}

func MustSha256(in []byte) []byte {
	h, err := Sha256(in)
	if err != nil {
		panic(err)
	}
	return h
}
