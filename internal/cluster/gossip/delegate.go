package gossipCluster

import (
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/internal/msg"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"hash/crc64"
)

/*
   Creation Time: 2021 - Jan - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type clusterDelegate struct {
	c *Cluster
}

func (d *clusterDelegate) nodeData() []byte {
	n := &msg.EdgeNode{
		ServerID:    d.c.localServerID,
		Hash:        crc64.Checksum(d.c.localServerID, crcTable),
		ReplicaSet:  d.c.cfg.ReplicaSet,
		GatewayAddr: d.c.localGatewayAddr,
		TunnelAddr:  d.c.localTunnelAddr,
	}
	b, _ := proto.Marshal(n)
	return b
}

func (d *clusterDelegate) NodeMeta(limit int) []byte {
	b := d.nodeData()
	if len(b) > limit {
		log.Warn("Too Large Meta", zap.ByteString("ServerID", d.c.localServerID))
		return nil
	}
	return b
}

func (d *clusterDelegate) NotifyMsg(_ []byte) {}

func (d *clusterDelegate) GetBroadcasts(_, _ int) [][]byte {
	return nil
}

func (d *clusterDelegate) LocalState(_ bool) []byte {
	return nil
}

func (d *clusterDelegate) MergeRemoteState(_ []byte, _ bool) {
}
