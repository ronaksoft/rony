package gossipCluster

import (
	"github.com/hashicorp/memberlist"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
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

func (d *clusterDelegate) addMember(n *memberlist.Node) {
	cm, err := newMember(n)
	if err != nil {
		log.Warn("Error On Cluster Node Add", zap.Error(err))
		return
	}

	d.c.addMember(cm)
	if d.c.subscriber != nil {
		d.c.subscriber.OnJoin(cm.hash)
	}
}

func (d *clusterDelegate) nodeData() []byte {
	n := rony.EdgeNode{
		ServerID:    d.c.localServerID,
		Hash:        crc64.Checksum(d.c.localServerID, crcTable),
		ReplicaSet:  d.c.cfg.ReplicaSet,
		GatewayAddr: d.c.localGatewayAddr,
		TunnelAddr:  d.c.localTunnelAddr,
	}

	b, _ := proto.Marshal(&n)

	return b
}

func (d *clusterDelegate) NotifyJoin(n *memberlist.Node) {
	d.addMember(n)

}

func (d *clusterDelegate) NotifyUpdate(n *memberlist.Node) {
	d.addMember(n)
}

func (d *clusterDelegate) NotifyAlive(n *memberlist.Node) error {
	d.addMember(n)
	return nil
}

func (d *clusterDelegate) NotifyLeave(n *memberlist.Node) {
	en := rony.PoolEdgeNode.Get()
	defer rony.PoolEdgeNode.Put(en)
	err := extractNode(n, en)
	if err != nil {
		log.Warn("Error On Cluster Node Update", zap.Error(err))
		return
	}

	d.c.removeMember(en)
	if d.c.subscriber != nil && n.State == memberlist.StateLeft {
		d.c.subscriber.OnLeave(en.Hash)
	}
}

func (d *clusterDelegate) NodeMeta(limit int) []byte {
	b := d.nodeData()
	if len(b) > limit {
		log.Warn("Too Large Meta", zap.ByteString("ServerID", d.c.localServerID))
		return nil
	}
	return b
}

func (d *clusterDelegate) NotifyMsg(data []byte) {}

func (d *clusterDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d *clusterDelegate) LocalState(join bool) []byte {
	return d.nodeData()
}

func (d *clusterDelegate) MergeRemoteState(buf []byte, join bool) {
}
