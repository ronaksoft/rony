package rony

import (
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/gateway"
	httpGateway "git.ronaksoftware.com/ronak/rony/gateway/http"
	quicGateway "git.ronaksoftware.com/ronak/rony/gateway/quic"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"git.ronaksoftware.com/ronak/rony/msg"
)

/*
   Creation Time: 2020 - Feb - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Option func(edge *EdgeServer)

func WithReplicaSet(replicaSet string, bindPort int, bootstrap bool) Option {
	return func(edge *EdgeServer) {
		edge.raftFSM = raftFSM{edge: edge}
		edge.replicaSet = replicaSet
		edge.raftEnabled = true
		edge.raftPort = bindPort
		edge.raftBootstrap = bootstrap
	}
}

func WithGossipPort(gossipPort int) Option {
	return func(edge *EdgeServer) {
		edge.gossipPort = gossipPort
	}
}

func WithCustomConstructorName(h func(constructor int64) (name string)) Option {
	return func(edge *EdgeServer) {
		edge.getConstructorName = func(constructor int64) string {
			// Lookup internal messages first
			n := msg.ConstructorNames[constructor]
			if len(n) > 0 {
				return n
			}
			return h(constructor)
		}
	}
}

func WithDataPath(path string) Option {
	return func(edge *EdgeServer) {
		edge.dataPath = path
	}
}

func WithWebsocketGateway(config websocketGateway.Config) Option {
	return func(edge *EdgeServer) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(errors.ErrGatewayAlreadyInitialized)
		}
		gatewayWebsocket, err := websocketGateway.New(config)
		if err != nil {
			panic(err)
		}
		gatewayWebsocket.MessageHandler = edge.onMessage
		gatewayWebsocket.ConnectHandler = edge.onConnect
		gatewayWebsocket.CloseHandler = edge.onClose
		gatewayWebsocket.FlushFunc = edge.onFlush
		edge.gatewayProtocol = gateway.Websocket
		edge.gateway = gatewayWebsocket
		return
	}
}

func WithHttpGateway(config httpGateway.Config) Option {
	return func(edge *EdgeServer) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(errors.ErrGatewayAlreadyInitialized)
		}
		gatewayHttp := httpGateway.New(config)
		gatewayHttp.MessageHandler = edge.onMessage
		gatewayHttp.FlushFunc = edge.onFlush

		edge.gatewayProtocol = gateway.HTTP
		edge.gateway = gatewayHttp
		return
	}
}

func WithQuicGateway(config quicGateway.Config) Option {
	return func(edge *EdgeServer) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(errors.ErrGatewayAlreadyInitialized)
		}
		gatewayQuic, err := quicGateway.New(config)
		if err != nil {
			panic(err)
		}
		gatewayQuic.MessageHandler = edge.onMessage
		gatewayQuic.ConnectHandler = edge.onConnect
		gatewayQuic.CloseHandler = edge.onClose
		gatewayQuic.FlushFunc = edge.onFlush

		edge.gatewayProtocol = gateway.QUIC
		edge.gateway = gatewayQuic
		return
	}
}
