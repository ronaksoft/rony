package registry

import (
	"fmt"

	"github.com/goccy/go-json"
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

// Message defines the behavior of messages in Rony. They must be able to
// encode/decode ProtocolBuffer and JSON formats.
type Message interface {
	proto.Message
	json.Marshaler
	json.Unmarshaler
}

type Envelope interface {
	GetConstructor() uint64
	GetMessage() []byte
}

type Factory func() Message

var (
	constructors     = map[uint64]string{}
	constructorNames = map[string]uint64{}
	factories        = map[uint64]Factory{}
)

func Register(c uint64, n string, f Factory) {
	if old, ok := constructors[c]; ok {
		panic(fmt.Sprintf("constructor already exists %s:%s", old, n))
	} else {
		constructors[c] = n
		constructorNames[n] = c
	}
	if f != nil {
		factories[c] = f
	}
}

func Get(c uint64) (Message, error) {
	f := factories[c]
	if f == nil {
		return nil, errFactoryNotExists
	}

	return f(), nil
}

func Unwrap(envelope Envelope) (Message, error) {
	m, err := Get(envelope.GetConstructor())
	if err != nil {
		return nil, err
	}

	err = proto.UnmarshalOptions{Merge: true}.Unmarshal(envelope.GetMessage(), m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func ConstructorName(c uint64) string {
	return constructors[c]
}

func C(c uint64) string {
	return constructors[c]
}

func ConstructorNumber(n string) uint64 {
	return constructorNames[n]
}

func N(n string) uint64 {
	return constructorNames[n]
}

type JSONEnvelope struct {
	RequestID   uint64            `json:"requestId,omitempty"`
	Header      map[string]string `json:"header,omitempty"`
	Constructor string            `json:"constructor"`
	Message     json.RawMessage   `json:"message,omitempty"`
}

var (
	errFactoryNotExists = fmt.Errorf("factory not exists")
)
