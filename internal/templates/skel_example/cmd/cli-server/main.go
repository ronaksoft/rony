package main

import (
	"git.ronaksoft.com/ronak/rony/edge"
	tcpGateway "git.ronaksoft.com/ronak/rony/gateway/tcp"
	"os"
	"time"
)

const (
	serverID         = "YourServerID"
	gossipPort       = 2374
	replicaSet       = 1
	replicaPort      = 2000
	replicaFirstNode = true
)

func main() {
	s := NewServer(serverID,
		edge.WithTcpGateway(tcpGateway.Config{
			Concurrency:   100,
			ListenAddress: "",
			MaxBodySize:   0,
			MaxIdleTime:   time.Minute,
			Protocol:      tcpGateway.Auto,
		}),
		edge.WithGossipPort(gossipPort),
		edge.WithReplicaSet(replicaSet, replicaPort, replicaFirstNode),
	)
	err := s.Start()
	if err != nil {
		panic(err)
	}

	// Wait for Kill or Interrupt signal to shutdown the server gracefully
	s.Shutdown(os.Kill, os.Interrupt)
}
