package rony

import (
	registry "github.com/ronaksoft/rony/registry"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_ClusterMessage int64 = 56353988

type poolClusterMessage struct {
	pool sync.Pool
}

func (p *poolClusterMessage) Get() *ClusterMessage {
	x, ok := p.pool.Get().(*ClusterMessage)
	if !ok {
		return &ClusterMessage{}
	}
	return x
}

func (p *poolClusterMessage) Put(x *ClusterMessage) {
	x.Sender = x.Sender[:0]
	x.Store = x.Store[:0]
	if x.Envelope != nil {
		PoolMessageEnvelope.Put(x.Envelope)
		x.Envelope = nil
	}
	p.pool.Put(x)
}

var PoolClusterMessage = poolClusterMessage{}

const C_RaftCommand int64 = 3965951915

type poolRaftCommand struct {
	pool sync.Pool
}

func (p *poolRaftCommand) Get() *RaftCommand {
	x, ok := p.pool.Get().(*RaftCommand)
	if !ok {
		return &RaftCommand{}
	}
	return x
}

func (p *poolRaftCommand) Put(x *RaftCommand) {
	x.Sender = x.Sender[:0]
	x.Store = x.Store[:0]
	if x.Envelope != nil {
		PoolMessageEnvelope.Put(x.Envelope)
		x.Envelope = nil
	}
	p.pool.Put(x)
}

var PoolRaftCommand = poolRaftCommand{}

const C_EdgeNode int64 = 2474233262

type poolEdgeNode struct {
	pool sync.Pool
}

func (p *poolEdgeNode) Get() *EdgeNode {
	x, ok := p.pool.Get().(*EdgeNode)
	if !ok {
		return &EdgeNode{}
	}
	return x
}

func (p *poolEdgeNode) Put(x *EdgeNode) {
	x.ServerID = x.ServerID[:0]
	x.ReplicaSet = 0
	x.ShardRangeMin = 0
	x.ShardRangeMax = 0
	x.RaftPort = 0
	x.RaftState = 0
	x.GatewayAddr = x.GatewayAddr[:0]
	p.pool.Put(x)
}

var PoolEdgeNode = poolEdgeNode{}

func init() {
	registry.RegisterConstructor(56353988, "rony.ClusterMessage")
	registry.RegisterConstructor(3965951915, "rony.RaftCommand")
	registry.RegisterConstructor(2474233262, "rony.EdgeNode")
}

func (x *ClusterMessage) DeepCopy(z *ClusterMessage) {
	z.Sender = append(z.Sender[:0], x.Sender...)
	for idx := range x.Store {
		if x.Store[idx] != nil {
			xx := PoolKeyValue.Get()
			x.Store[idx].DeepCopy(xx)
			z.Store = append(z.Store, xx)
		}
	}
	if x.Envelope != nil {
		z.Envelope = PoolMessageEnvelope.Get()
		x.Envelope.DeepCopy(z.Envelope)
	}
}

func (x *RaftCommand) DeepCopy(z *RaftCommand) {
	z.Sender = append(z.Sender[:0], x.Sender...)
	for idx := range x.Store {
		if x.Store[idx] != nil {
			xx := PoolKeyValue.Get()
			x.Store[idx].DeepCopy(xx)
			z.Store = append(z.Store, xx)
		}
	}
	if x.Envelope != nil {
		z.Envelope = PoolMessageEnvelope.Get()
		x.Envelope.DeepCopy(z.Envelope)
	}
}

func (x *EdgeNode) DeepCopy(z *EdgeNode) {
	z.ServerID = append(z.ServerID[:0], x.ServerID...)
	z.ReplicaSet = x.ReplicaSet
	z.ShardRangeMin = x.ShardRangeMin
	z.ShardRangeMax = x.ShardRangeMax
	z.RaftPort = x.RaftPort
	z.RaftState = x.RaftState
	z.GatewayAddr = append(z.GatewayAddr[:0], x.GatewayAddr...)
}

func (x *ClusterMessage) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *RaftCommand) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EdgeNode) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *ClusterMessage) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *RaftCommand) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *EdgeNode) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}
