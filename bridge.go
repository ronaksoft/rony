package rony

import (
	"git.ronaksoftware.com/ronak/rony/bridge"
	natsBridge "git.ronaksoftware.com/ronak/rony/bridge/nats"
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"go.uber.org/zap"
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
	bridgeMedium bridge.Medium
	bridgeNats   *natsBridge.Bridge
)

func bridgeNotifyHandler(connIDs []uint64) {
	if ce := log.Check(log.DebugLevel, "Notify Arrived"); ce != nil {
		ce.Write(
			zap.Uint64s("ConnIDs", connIDs),
		)
	}
	switch gatewayProtocol {
	case gateway.Websocket:
		for idx := range connIDs {
			wsConn := gatewayWebsocket.GetConnection(connIDs[idx])
			if wsConn == nil {
				if ce := log.Check(log.DebugLevel, "Flush Signal For Nil Websocket"); ce != nil {
					ce.Write(zap.Uint64("ConnID", connIDs[idx]))
				}
				continue
			}
			wsConn.Flush()
		}
	// case gateway.Grpc:
	// 	for idx := range connIDs {
	// 		conn := gatewayGRPC.GetConnection(connIDs[idx])
	// 		if conn == nil {
	// 			log.Warn("Flush Signal For Nil Grpc stream", zap.Uint64("ConnID", connIDs[idx]))
	// 			continue
	// 		}
	// 		conn.Flush()
	// 	}
	default:

	}
}

func InitNatsBridge(config natsBridge.Config) (err error) {
	if bridgeMedium != bridge.Undefined {
		return errors.ErrBridgeAlreadyInitialized
	}

	bridgeNats, err = natsBridge.NewBridge(config)
	if err != nil {
		return
	}
	bridgeNats.NotifyHandler = bridgeNotifyHandler
	bridgeMedium = bridge.Nats
	return
}
