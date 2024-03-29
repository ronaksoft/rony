// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.16.24
// Source: imsg.proto

package msg

import (
	bytes "bytes"
	sync "sync"

	rony "github.com/ronaksoft/rony"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _ = pools.Imported

const C_TunnelMessage uint64 = 12904103981552626014

type poolTunnelMessage struct {
	pool sync.Pool
}

func (p *poolTunnelMessage) Get() *TunnelMessage {
	x, ok := p.pool.Get().(*TunnelMessage)
	if !ok {
		x = &TunnelMessage{}
	}

	x.Envelope = rony.PoolMessageEnvelope.Get()

	return x
}

func (p *poolTunnelMessage) Put(x *TunnelMessage) {
	if x == nil {
		return
	}

	x.SenderID = x.SenderID[:0]
	x.SenderReplicaSet = 0
	for _, z := range x.Store {
		rony.PoolKeyValue.Put(z)
	}
	x.Store = x.Store[:0]
	rony.PoolMessageEnvelope.Put(x.Envelope)

	p.pool.Put(x)
}

var PoolTunnelMessage = poolTunnelMessage{}

func (x *TunnelMessage) DeepCopy(z *TunnelMessage) {
	z.SenderID = append(z.SenderID[:0], x.SenderID...)
	z.SenderReplicaSet = x.SenderReplicaSet
	for idx := range x.Store {
		if x.Store[idx] == nil {
			continue
		}
		xx := rony.PoolKeyValue.Get()
		x.Store[idx].DeepCopy(xx)
		z.Store = append(z.Store, xx)
	}
	if x.Envelope != nil {
		if z.Envelope == nil {
			z.Envelope = rony.PoolMessageEnvelope.Get()
		}
		x.Envelope.DeepCopy(z.Envelope)
	} else {
		rony.PoolMessageEnvelope.Put(z.Envelope)
		z.Envelope = nil
	}
}

func (x *TunnelMessage) Clone() *TunnelMessage {
	z := &TunnelMessage{}
	x.DeepCopy(z)
	return z
}

func (x *TunnelMessage) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *TunnelMessage) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *TunnelMessage) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *TunnelMessage) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryTunnelMessage() registry.Message {
	return &TunnelMessage{}
}

const C_EdgeNode uint64 = 12741646567456486228

type poolEdgeNode struct {
	pool sync.Pool
}

func (p *poolEdgeNode) Get() *EdgeNode {
	x, ok := p.pool.Get().(*EdgeNode)
	if !ok {
		x = &EdgeNode{}
	}

	return x
}

func (p *poolEdgeNode) Put(x *EdgeNode) {
	if x == nil {
		return
	}

	x.ServerID = x.ServerID[:0]
	x.ReplicaSet = 0
	x.Hash = 0
	x.GatewayAddr = x.GatewayAddr[:0]
	x.TunnelAddr = x.TunnelAddr[:0]

	p.pool.Put(x)
}

var PoolEdgeNode = poolEdgeNode{}

func (x *EdgeNode) DeepCopy(z *EdgeNode) {
	z.ServerID = append(z.ServerID[:0], x.ServerID...)
	z.ReplicaSet = x.ReplicaSet
	z.Hash = x.Hash
	z.GatewayAddr = append(z.GatewayAddr[:0], x.GatewayAddr...)
	z.TunnelAddr = append(z.TunnelAddr[:0], x.TunnelAddr...)
}

func (x *EdgeNode) Clone() *EdgeNode {
	z := &EdgeNode{}
	x.DeepCopy(z)
	return z
}

func (x *EdgeNode) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *EdgeNode) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EdgeNode) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *EdgeNode) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryEdgeNode() registry.Message {
	return &EdgeNode{}
}

// register constructors of the messages to the registry package
func init() {
	registry.Register(12904103981552626014, "TunnelMessage", factoryTunnelMessage)
	registry.Register(12741646567456486228, "EdgeNode", factoryEdgeNode)

}

var _ = bytes.MinRead
