package rony

import (
	"fmt"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"github.com/hashicorp/serf/serf"
	"go.uber.org/zap"
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

func (edge *EdgeServer) updateCluster() {
	err := edge.gossip.SetTags(map[string]string{
		tagBundleID:   edge.bundleID,
		tagInstanceID: edge.instanceID,
		tagRaftNodeID: edge.GetServerID(),
		tagRaftPort:   fmt.Sprintf("%d", edge.raftPort),
	})
	if err != nil {
		log.Warn("Error On Update Cluster", zap.Error(err))
	}
}

func (edge *EdgeServer) eventHandler(e serf.Event) {

}


type Cluster struct {}

func (c *Cluster) Add(serverID string) {

}

func (c *Cluster) Remove(serverID string) {

}