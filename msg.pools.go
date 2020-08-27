package rony

import (
	"git.ronaksoft.com/ronak/rony/registry"
	"sync"
)

const C_MessageEnvelope int64 = 535232465

type poolMessageEnvelope struct {
	pool sync.Pool
}

func (p *poolMessageEnvelope) Get() *MessageEnvelope {
	x, ok := p.pool.Get().(*MessageEnvelope)
	if !ok {
		return &MessageEnvelope{}
	}
	return x
}

func (p *poolMessageEnvelope) Put(x *MessageEnvelope) {
	x.Constructor = 0
	x.RequestID = 0
	x.Message = x.Message[:0]
	p.pool.Put(x)
}

var PoolMessageEnvelope = poolMessageEnvelope{}

const C_MessageContainer int64 = 1972016308

type poolMessageContainer struct {
	pool sync.Pool
}

func (p *poolMessageContainer) Get() *MessageContainer {
	x, ok := p.pool.Get().(*MessageContainer)
	if !ok {
		return &MessageContainer{}
	}
	return x
}

func (p *poolMessageContainer) Put(x *MessageContainer) {
	x.Length = 0
	x.Envelopes = x.Envelopes[:0]
	p.pool.Put(x)
}

var PoolMessageContainer = poolMessageContainer{}

const C_Error int64 = 2619118453

type poolError struct {
	pool sync.Pool
}

func (p *poolError) Get() *Error {
	x, ok := p.pool.Get().(*Error)
	if !ok {
		return &Error{}
	}
	return x
}

func (p *poolError) Put(x *Error) {
	x.Code = ""
	x.Items = ""
	x.EnglishTemplate = ""
	x.EnglishItems = x.EnglishItems[:0]
	x.LocalTemplate = ""
	x.LocalItems = x.LocalItems[:0]
	p.pool.Put(x)
}

var PoolError = poolError{}

const C_Redirect int64 = 981138557

type poolRedirect struct {
	pool sync.Pool
}

func (p *poolRedirect) Get() *Redirect {
	x, ok := p.pool.Get().(*Redirect)
	if !ok {
		return &Redirect{}
	}
	return x
}

func (p *poolRedirect) Put(x *Redirect) {
	x.LeaderHostPort = x.LeaderHostPort[:0]
	x.ServerID = ""
	p.pool.Put(x)
}

var PoolRedirect = poolRedirect{}

const C_ClusterMessage int64 = 1078766375

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
	x.AuthID = 0
	if x.Envelope != nil {
		*x.Envelope = MessageEnvelope{}
	}
	p.pool.Put(x)
}

var PoolClusterMessage = poolClusterMessage{}

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
	x.AuthID = 0
	if x.Envelope != nil {
		*x.Envelope = MessageEnvelope{}
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
	x.ShardSet = 0
	x.RaftPort = 0
	x.RaftState = 0
	x.GatewayAddr = x.GatewayAddr[:0]
	p.pool.Put(x)
}

var PoolEdgeNode = poolEdgeNode{}

func init() {
	registry.RegisterConstructor(535232465, "MessageEnvelope")
	registry.RegisterConstructor(1972016308, "MessageContainer")
	registry.RegisterConstructor(2619118453, "Error")
	registry.RegisterConstructor(981138557, "Redirect")
	registry.RegisterConstructor(1078766375, "ClusterMessage")
	registry.RegisterConstructor(2919813429, "RaftCommand")
	registry.RegisterConstructor(999040174, "EdgeNode")
}
