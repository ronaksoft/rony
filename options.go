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

func WithUpdateDispatcher(h func(authID int64, envelope *msg.UpdateEnvelope)) Option {
	return func(edge *EdgeServer) {
		edge.updateDispatcher = h
	}
}

func WithMessageDispatcher(h func(authID int64, envelope *msg.MessageEnvelope)) Option {
	return func(edge *EdgeServer) {
		edge.messageDispatcher = h
	}
}

func WithRaft(bindPort int, bootstrap bool) Option {
	return func(edge *EdgeServer) {
		edge.raftFSM = raftFSM{edge: edge}
		edge.raftEnabled = true
		edge.raftPort = bindPort
		edge.raftBootstrap = bootstrap
	}
}

func WithCustomConstructorName(h func(constructor int64) (name string)) Option {
	return func(edge *EdgeServer) {
		edge.getConstructorName = h
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

func WithGossipPort(gossipPort int) Option {
	return func(edge *EdgeServer) {
		edge.gossipPort = gossipPort
	}
}
