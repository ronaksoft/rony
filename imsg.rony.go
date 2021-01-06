package rony

import (
	registry "github.com/ronaksoft/rony/registry"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_TunnelMessage int64 = 3271476222

type poolTunnelMessage struct {
	pool sync.Pool
}

func (p *poolTunnelMessage) Get() *TunnelMessage {
	x, ok := p.pool.Get().(*TunnelMessage)
	if !ok {
		return &TunnelMessage{}
	}
	return x
}

func (p *poolTunnelMessage) Put(x *TunnelMessage) {
	x.SenderID = x.SenderID[:0]
	x.SenderReplicaSet = 0
	x.Store = x.Store[:0]
	if x.Envelope != nil {
		PoolMessageEnvelope.Put(x.Envelope)
		x.Envelope = nil
	}
	p.pool.Put(x)
}

var PoolTunnelMessage = poolTunnelMessage{}

const C_RaftCommand int64 = 2919813429

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

const C_EdgeNode int64 = 999040174

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
	x.TunnelAddr = x.TunnelAddr[:0]
	p.pool.Put(x)
}

var PoolEdgeNode = poolEdgeNode{}

func init() {
	registry.RegisterConstructor(3271476222, "TunnelMessage")
	registry.RegisterConstructor(2919813429, "RaftCommand")
	registry.RegisterConstructor(999040174, "EdgeNode")
}

func (x *TunnelMessage) DeepCopy(z *TunnelMessage) {
	z.SenderID = append(z.SenderID[:0], x.SenderID...)
	z.SenderReplicaSet = x.SenderReplicaSet
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
	z.TunnelAddr = append(z.TunnelAddr[:0], x.TunnelAddr...)
}

func (x *TunnelMessage) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *RaftCommand) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EdgeNode) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *TunnelMessage) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *RaftCommand) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *EdgeNode) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}
