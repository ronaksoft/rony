// Code generated by Rony's protoc plugin; DO NOT EDIT.

package rony

import (
	registry "github.com/ronaksoft/rony/registry"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
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
	if x == nil {
		return
	}
	x.Constructor = 0
	x.RequestID = 0
	x.Message = x.Message[:0]
	x.Auth = x.Auth[:0]
	for _, z := range x.Header {
		PoolKeyValue.Put(z)
	}
	x.Header = x.Header[:0]
	p.pool.Put(x)
}

var PoolMessageEnvelope = poolMessageEnvelope{}

func (x *MessageEnvelope) DeepCopy(z *MessageEnvelope) {
	z.Constructor = x.Constructor
	z.RequestID = x.RequestID
	z.Message = append(z.Message[:0], x.Message...)
	z.Auth = append(z.Auth[:0], x.Auth...)
	for idx := range x.Header {
		if x.Header[idx] != nil {
			xx := PoolKeyValue.Get()
			x.Header[idx].DeepCopy(xx)
			z.Header = append(z.Header, xx)
		}
	}
}

func (x *MessageEnvelope) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *MessageEnvelope) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_KeyValue int64 = 4276272820

type poolKeyValue struct {
	pool sync.Pool
}

func (p *poolKeyValue) Get() *KeyValue {
	x, ok := p.pool.Get().(*KeyValue)
	if !ok {
		return &KeyValue{}
	}
	return x
}

func (p *poolKeyValue) Put(x *KeyValue) {
	if x == nil {
		return
	}
	x.Key = ""
	x.Value = ""
	p.pool.Put(x)
}

var PoolKeyValue = poolKeyValue{}

func (x *KeyValue) DeepCopy(z *KeyValue) {
	z.Key = x.Key
	z.Value = x.Value
}

func (x *KeyValue) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *KeyValue) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

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
	if x == nil {
		return
	}
	x.Length = 0
	for _, z := range x.Envelopes {
		PoolMessageEnvelope.Put(z)
	}
	x.Envelopes = x.Envelopes[:0]
	p.pool.Put(x)
}

var PoolMessageContainer = poolMessageContainer{}

func (x *MessageContainer) DeepCopy(z *MessageContainer) {
	z.Length = x.Length
	for idx := range x.Envelopes {
		if x.Envelopes[idx] != nil {
			xx := PoolMessageEnvelope.Get()
			x.Envelopes[idx].DeepCopy(xx)
			z.Envelopes = append(z.Envelopes, xx)
		}
	}
}

func (x *MessageContainer) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *MessageContainer) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

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
	if x == nil {
		return
	}
	x.Code = ""
	x.Items = ""
	x.Description = ""
	p.pool.Put(x)
}

var PoolError = poolError{}

func (x *Error) DeepCopy(z *Error) {
	z.Code = x.Code
	z.Items = x.Items
	z.Description = x.Description
}

func (x *Error) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Error) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

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
	if x == nil {
		return
	}
	x.Reason = 0
	PoolEdge.Put(x.Leader)
	x.Leader = nil
	for _, z := range x.Followers {
		PoolEdge.Put(z)
	}
	x.Followers = x.Followers[:0]
	x.WaitInSec = 0
	p.pool.Put(x)
}

var PoolRedirect = poolRedirect{}

func (x *Redirect) DeepCopy(z *Redirect) {
	z.Reason = x.Reason
	if x.Leader != nil {
		z.Leader = PoolEdge.Get()
		x.Leader.DeepCopy(z.Leader)
	}
	for idx := range x.Followers {
		if x.Followers[idx] != nil {
			xx := PoolEdge.Get()
			x.Followers[idx].DeepCopy(xx)
			z.Followers = append(z.Followers, xx)
		}
	}
	z.WaitInSec = x.WaitInSec
}

func (x *Redirect) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Redirect) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_Edge int64 = 3576986712

type poolEdge struct {
	pool sync.Pool
}

func (p *poolEdge) Get() *Edge {
	x, ok := p.pool.Get().(*Edge)
	if !ok {
		return &Edge{}
	}
	return x
}

func (p *poolEdge) Put(x *Edge) {
	if x == nil {
		return
	}
	x.ReplicaSet = 0
	x.ServerID = ""
	x.HostPorts = x.HostPorts[:0]
	x.Leader = false
	p.pool.Put(x)
}

var PoolEdge = poolEdge{}

func (x *Edge) DeepCopy(z *Edge) {
	z.ReplicaSet = x.ReplicaSet
	z.ServerID = x.ServerID
	z.HostPorts = append(z.HostPorts[:0], x.HostPorts...)
	z.Leader = x.Leader
}

func (x *Edge) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Edge) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_Edges int64 = 2120950449

type poolEdges struct {
	pool sync.Pool
}

func (p *poolEdges) Get() *Edges {
	x, ok := p.pool.Get().(*Edges)
	if !ok {
		return &Edges{}
	}
	return x
}

func (p *poolEdges) Put(x *Edges) {
	if x == nil {
		return
	}
	for _, z := range x.Nodes {
		PoolEdge.Put(z)
	}
	x.Nodes = x.Nodes[:0]
	p.pool.Put(x)
}

var PoolEdges = poolEdges{}

func (x *Edges) DeepCopy(z *Edges) {
	for idx := range x.Nodes {
		if x.Nodes[idx] != nil {
			xx := PoolEdge.Get()
			x.Nodes[idx].DeepCopy(xx)
			z.Nodes = append(z.Nodes, xx)
		}
	}
}

func (x *Edges) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Edges) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func init() {
	registry.RegisterConstructor(535232465, "MessageEnvelope")
	registry.RegisterConstructor(4276272820, "KeyValue")
	registry.RegisterConstructor(1972016308, "MessageContainer")
	registry.RegisterConstructor(2619118453, "Error")
	registry.RegisterConstructor(981138557, "Redirect")
	registry.RegisterConstructor(3576986712, "Edge")
	registry.RegisterConstructor(2120950449, "Edges")
}
