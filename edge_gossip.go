package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/errors"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/msg"
	"github.com/gobwas/pool/pbytes"
	"github.com/hashicorp/memberlist"
	"go.uber.org/zap"
	"net"
)

/*
   Creation Time: 2020 - Feb - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	tagBundleID   = "BID"
	tagInstanceID = "IID"
	tagRaftPort   = "RP"
	tagRaftNodeID = "RID"
	tagRaftState  = "RS"
)

func (edge *EdgeServer) ClusterSend(serverID string, envelope *msg.MessageEnvelope) error {
	m := edge.cluster.GetByID(serverID)
	if m == nil {
		return errors.ErrEmpty
	}
	b := pbytes.GetLen(envelope.Size())
	_, err := envelope.MarshalTo(b)
	if err != nil {
		return err
	}
	err = edge.gossip.SendBestEffort(m.node, b)
	pbytes.Put(b)
	return err
}

func (edge *EdgeServer) updateCluster() error {
	return edge.gossip.UpdateNode(0)
}

// func (edge *EdgeServer) eventHandler(e serf.Event) {
// 	switch x := e.(type) {
// 	case serf.MemberEvent:
// 		switch x.Type {
// 		case serf.EventMemberJoin, serf.EventMemberUpdate:
// 			for idx := range x.Members {
// 				edge.cluster.AddMember(convertMember(&x.Members[idx]))
// 			}
// 		case serf.EventMemberLeave, serf.EventMemberFailed, serf.EventMemberReap:
// 			for idx := range x.Members {
// 				edge.cluster.RemoveMember(convertMember(&x.Members[idx]))
// 			}
// 		}
// 	case serf.UserEvent:
// 		if ce := log.Check(log.InfoLevel, "Cluster Message Received"); ce != nil {
// 			ce.Write(
// 				zap.String("ID", edge.GetServerID()),
// 				zap.String("Name", x.Name),
// 				zap.Int("Len", len(x.Payload)),
// 			)
// 		}
// 		envelope := pools.AcquireMessageEnvelope()
// 		err := envelope.Unmarshal(x.Payload)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		edge.dispatcher.DispatchClusterMessage(envelope)
// 	case *serf.Query:
//
// 	}
//
// }

// ClusterMember
type ClusterMember struct {
	ServerID   string
	BundleID   string
	InstanceID string
	Addr       net.IP
	Port       uint16
	node       *memberlist.Node
}

func convertMember(sm *memberlist.Node) *ClusterMember {
	edgeNode := EdgeNode{}
	err := edgeNode.Unmarshal(sm.Meta)
	if err != nil {
		log.Warn("Error On ConvertMember",
			zap.Error(err),
			zap.Int("Len", len(sm.Meta)),
		)
		return nil
	}

	return &ClusterMember{
		ServerID:   fmt.Sprintf("%s.%s", edgeNode.BundleID, edgeNode.InstanceID),
		BundleID:   edgeNode.BundleID,
		InstanceID: edgeNode.InstanceID,
		Addr:       sm.Addr,
		Port:       sm.Port,
		node:       sm,
	}
}

type Cluster struct {
	byServerID map[string]*ClusterMember
	byBundleID map[string][]*ClusterMember
}

func (c *Cluster) GetByID(id string) *ClusterMember {
	return c.byServerID[id]
}

func (c *Cluster) AddMember(m *ClusterMember) {
	if m == nil {
		return
	}
	if len(m.BundleID) == 0 {
		return
	}

	if c.byBundleID == nil {
		c.byBundleID = make(map[string][]*ClusterMember)
		c.byServerID = make(map[string]*ClusterMember)
		c.byServerID[m.ServerID] = m
		c.byBundleID[m.BundleID] = append(c.byBundleID[m.BundleID], m)
		return
	}

	c.byServerID[m.ServerID] = m
	if c.byBundleID[m.BundleID] == nil {
		c.byBundleID[m.BundleID] = append(c.byBundleID[m.BundleID], m)
		return
	}

	for idx := range c.byBundleID[m.BundleID] {
		if c.byBundleID[m.BundleID][idx].InstanceID == m.InstanceID {
			c.byBundleID[m.InstanceID][idx] = m
			return
		}
	}
	c.byBundleID[m.BundleID] = append(c.byBundleID[m.BundleID], m)
}

func (c *Cluster) RemoveMember(m *ClusterMember) {
	if m == nil {
		return
	}
	if len(m.BundleID) == 0 {
		return
	}

	if c.byBundleID == nil {
		c.byBundleID = make(map[string][]*ClusterMember)
		c.byServerID = make(map[string]*ClusterMember)
		return
	}
	delete(c.byServerID, m.ServerID)
	for idx := range c.byBundleID[m.BundleID] {
		if c.byBundleID[m.BundleID][idx].InstanceID == m.InstanceID {
			c.byBundleID[m.BundleID][idx] = c.byBundleID[m.BundleID][len(c.byBundleID[m.BundleID])-1]
			c.byBundleID[m.BundleID] = c.byBundleID[m.BundleID][:len(c.byBundleID[m.BundleID])-1]
			return
		}
	}
}

type delegateEvents struct {
	edge *EdgeServer
}

func (d delegateEvents) NotifyJoin(n *memberlist.Node) {
	d.edge.cluster.AddMember(convertMember(n))
}

func (d delegateEvents) NotifyLeave(n *memberlist.Node) {
	d.edge.cluster.RemoveMember(convertMember(n))
}

func (d delegateEvents) NotifyUpdate(n *memberlist.Node) {
	d.edge.cluster.AddMember(convertMember(n))
}

type delegateNode struct {
	edge *EdgeServer
}

func (d delegateNode) NodeMeta(limit int) []byte {
	n := EdgeNode{
		BundleID:   d.edge.bundleID,
		InstanceID: d.edge.instanceID,
		RaftPort:   uint32(d.edge.raftPort),
	}
	b, _ := n.Marshal()
	if len(b) > limit {
		log.Warn("Too Large Meta")
		return nil
	}
	return b
}

func (d delegateNode) NotifyMsg(data []byte) {
	log.Info("Message Received",
		zap.String("ServerID", d.edge.GetServerID()),
		zap.Int("Data", len(data)),
	)
}

func (d delegateNode) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d delegateNode) LocalState(join bool) []byte {
	return nil
}

func (d delegateNode) MergeRemoteState(buf []byte, join bool) {
	return
}
