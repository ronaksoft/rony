package rony

import (
	registry "github.com/ronaksoft/rony/registry"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_MessageEnvelope int64 = 648436572

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
	x.Auth = x.Auth[:0]
	x.Header = x.Header[:0]
	p.pool.Put(x)
}

var PoolMessageEnvelope = poolMessageEnvelope{}

const C_MessageContainer int64 = 3870960525

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

const C_Error int64 = 1239766265

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
	x.Template = ""
	x.TemplateItems = x.TemplateItems[:0]
	x.LocalTemplate = ""
	x.LocalTemplateItems = x.LocalTemplateItems[:0]
	p.pool.Put(x)
}

var PoolError = poolError{}

const C_Redirect int64 = 2458850685

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
	x.HostPorts = x.HostPorts[:0]
	x.ServerID = ""
	x.WaitInSec = 0
	p.pool.Put(x)
}

var PoolRedirect = poolRedirect{}

const C_KeyValue int64 = 1444370356

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
	x.Key = ""
	x.Value = ""
	p.pool.Put(x)
}

var PoolKeyValue = poolKeyValue{}

func init() {
	registry.RegisterConstructor(648436572, "rony.MessageEnvelope")
	registry.RegisterConstructor(3870960525, "rony.MessageContainer")
	registry.RegisterConstructor(1239766265, "rony.Error")
	registry.RegisterConstructor(2458850685, "rony.Redirect")
	registry.RegisterConstructor(1444370356, "rony.KeyValue")
}

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

func (x *Error) DeepCopy(z *Error) {
	z.Code = x.Code
	z.Items = x.Items
	z.Template = x.Template
	z.TemplateItems = append(z.TemplateItems[:0], x.TemplateItems...)
	z.LocalTemplate = x.LocalTemplate
	z.LocalTemplateItems = append(z.LocalTemplateItems[:0], x.LocalTemplateItems...)
}

func (x *Redirect) DeepCopy(z *Redirect) {
	z.LeaderHostPort = append(z.LeaderHostPort[:0], x.LeaderHostPort...)
	z.HostPorts = append(z.HostPorts[:0], x.HostPorts...)
	z.ServerID = x.ServerID
	z.WaitInSec = x.WaitInSec
}

func (x *KeyValue) DeepCopy(z *KeyValue) {
	z.Key = x.Key
	z.Value = x.Value
}

func (x *MessageEnvelope) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *MessageContainer) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Error) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Redirect) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *KeyValue) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *MessageEnvelope) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *MessageContainer) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Error) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Redirect) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *KeyValue) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}
