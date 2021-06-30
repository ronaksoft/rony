package gateway

/*
   Creation Time: 2021 - Jun - 30
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Protocol int32

const (
	Undefined Protocol = 0
	Dummy     Protocol = 1 << iota
	Http
	Websocket
	Quic
	Grpc
	TCP = Http | Websocket // Http & Websocket
)

var protocolNames = map[Protocol]string{
	Undefined: "Undefined",
	Dummy:     "Dummy",
	Http:      "Http",
	Websocket: "Websocket",
	Quic:      "Quic",
	Grpc:      "Grpc",
	TCP:       "TCP",
}

func (p Protocol) String() string {
	return protocolNames[p]
}
