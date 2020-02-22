package rony

import (
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/gateway"
	grpcGateway "git.ronaksoftware.com/ronak/rony/gateway/grpc"
	httpGateway "git.ronaksoftware.com/ronak/rony/gateway/http"
	quicGateway "git.ronaksoftware.com/ronak/rony/gateway/quic"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
)

/*
   Creation Time: 2020 - Feb - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	gatewayProtocol  gateway.Protocol
	gatewayWebsocket *websocketGateway.Gateway
	gatewayHTTP      *httpGateway.Gateway
	gatewayQuic      *quicGateway.Gateway
	gatewayGrpc      *grpcGateway.Gateway
)

func InitWebsocketGateway(config websocketGateway.Config) (err error) {
	if gatewayProtocol != gateway.Undefined {
		return errors.ErrGatewayAlreadyInitialized
	}
	gatewayWebsocket, err = websocketGateway.New(config)
	if err != nil {
		return
	}
	gatewayProtocol = gateway.Websocket
	return
}

func InitHttpGateway(config httpGateway.Config) (err error) {
	if gatewayProtocol != gateway.Undefined {
		return errors.ErrGatewayAlreadyInitialized
	}
	gatewayHTTP = httpGateway.New(config)
	gatewayProtocol = gateway.HTTP
	return
}

func InitQuicGateway(config quicGateway.Config) (err error) {
	if gatewayProtocol != gateway.Undefined {
		return errors.ErrGatewayAlreadyInitialized
	}
	gatewayQuic, err = quicGateway.New(config)
	if err != nil {
		return
	}
	gatewayProtocol = gateway.QUIC
	return
}

func InitGrpcGateway(config grpcGateway.Config) (err error) {
	if gatewayProtocol != gateway.Undefined {
		return errors.ErrGatewayAlreadyInitialized
	}

	gatewayGrpc = grpcGateway.New(config)
	gatewayProtocol = gateway.GRPC
	return
}
