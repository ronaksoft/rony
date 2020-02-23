package rony

import (
	"fmt"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/hashicorp/memberlist"
	"go.uber.org/zap"
	"time"
)

/*
   Creation Time: 2020 - Feb - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type gossipDelegate struct {
	edge *EdgeServer
}

func (g gossipDelegate) NodeMeta(limit int) []byte {
	log.Info("NodeMeta")
	// panic("implement me")
	return nil
}

func (g gossipDelegate) NotifyMsg(b []byte) {
	fmt.Println(tools.ByteToStr(b))
}

func (g gossipDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	log.Info("Broadcast", zap.Int("Overhead", overhead), zap.Int("Limit", limit))
	return nil
}

func (g gossipDelegate) LocalState(join bool) []byte {
	log.Info("LocalState", zap.Bool("Join", join))
	return nil
}

func (g gossipDelegate) MergeRemoteState(buf []byte, join bool) {
	// panic("implement me")
}

type gossipEventDelegate struct {
	edge *EdgeServer
}

func (g gossipEventDelegate) NotifyJoin(node *memberlist.Node) {
	log.Info("Gossip Join", zap.String("Name", node.Name), zap.Uint16("Port", node.Port))
}

func (g gossipEventDelegate) NotifyLeave(node *memberlist.Node) {
	log.Info("Gossip Leave", zap.String("Name", node.Name), zap.Uint16("Port", node.Port))
}

func (g gossipEventDelegate) NotifyUpdate(node *memberlist.Node) {
	log.Info("Gossip Update", zap.String("Name", node.Name), zap.Uint16("Port", node.Port))
}

type gossipPingDelegate struct {
	edge *EdgeServer
}

func (g gossipPingDelegate) AckPayload() []byte {
	log.Info("AckPayload")
	return tools.StrToByte("Ack Payload")
}

func (g gossipPingDelegate) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) {
	log.Info("Ping Complete")
}

type gossipAliveDelegate struct {
	edge *EdgeServer
}

func (g gossipAliveDelegate) NotifyAlive(node *memberlist.Node) error {
	// log.Info("Gossip Alive", zap.String("Name", node.Name), zap.Uint16("Port", node.Port))
	return nil
}
