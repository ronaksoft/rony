package rony

import (
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
	p.pool.Put(x)
}

var PoolEdgeNode = poolEdgeNode{}

func init() {
	ConstructorNames[535232465] = "MessageEnvelope"
	ConstructorNames[1972016308] = "MessageContainer"
	ConstructorNames[2619118453] = "Error"
	ConstructorNames[981138557] = "Redirect"
	ConstructorNames[1078766375] = "ClusterMessage"
	ConstructorNames[2919813429] = "RaftCommand"
	ConstructorNames[999040174] = "EdgeNode"
}
