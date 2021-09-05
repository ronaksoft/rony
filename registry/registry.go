package registry

import (
	"fmt"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2020 - Aug - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Envelope interface {
	GetConstructor() uint64
	GetMessage() []byte
}

type UnwrapFunc func(envelope Envelope) (proto.Message, error)

var constructors = map[uint64]string{}
var unwrapFunctions = map[uint64]UnwrapFunc{}

func Register(c uint64, n string, unwrapFunc UnwrapFunc) {
	if old, ok := constructors[c]; ok {
		panic(fmt.Sprintf("constructor already exists %s:%s", old, n))
	} else {
		constructors[c] = n
	}
	if unwrapFunc != nil {
		unwrapFunctions[c] = unwrapFunc
	}
}

func ConstructorName(c uint64) string {
	return constructors[c]
}

func Unwrap(envelope Envelope) (proto.Message, error) {
	unwrapFunc := unwrapFunctions[envelope.GetConstructor()]
	if unwrapFunc == nil {
		return nil, fmt.Errorf("not found")
	}
	return unwrapFunc(envelope)
}

func C(c uint64) string {
	return constructors[c]
}
