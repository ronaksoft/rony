// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dev.proto

package pb

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	math "math"
	sync "sync"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

const C_ProtoMessage int64 = 2179260159

type poolProtoMessage struct {
	pool sync.Pool
}

func (p *poolProtoMessage) Get() *ProtoMessage {
	x, ok := p.pool.Get().(*ProtoMessage)
	if !ok {
		return &ProtoMessage{}
	}
	return x
}

func (p *poolProtoMessage) Put(x *ProtoMessage) {
	x.AuthID = 0
	x.MessageKey = x.MessageKey[:0]
	x.Payload = x.Payload[:0]
	p.pool.Put(x)
}

var PoolProtoMessage = poolProtoMessage{}

const C_EchoRequest int64 = 1904100324

type poolEchoRequest struct {
	pool sync.Pool
}

func (p *poolEchoRequest) Get() *EchoRequest {
	x, ok := p.pool.Get().(*EchoRequest)
	if !ok {
		return &EchoRequest{}
	}
	return x
}

func (p *poolEchoRequest) Put(x *EchoRequest) {
	p.pool.Put(x)
}

var PoolEchoRequest = poolEchoRequest{}

const C_EchoResponse int64 = 4192619139

type poolEchoResponse struct {
	pool sync.Pool
}

func (p *poolEchoResponse) Get() *EchoResponse {
	x, ok := p.pool.Get().(*EchoResponse)
	if !ok {
		return &EchoResponse{}
	}
	return x
}

func (p *poolEchoResponse) Put(x *EchoResponse) {
	p.pool.Put(x)
}

var PoolEchoResponse = poolEchoResponse{}

const C_AskRequest int64 = 3206229608

type poolAskRequest struct {
	pool sync.Pool
}

func (p *poolAskRequest) Get() *AskRequest {
	x, ok := p.pool.Get().(*AskRequest)
	if !ok {
		return &AskRequest{}
	}
	return x
}

func (p *poolAskRequest) Put(x *AskRequest) {
	p.pool.Put(x)
}

var PoolAskRequest = poolAskRequest{}

const C_AskResponse int64 = 489087205

type poolAskResponse struct {
	pool sync.Pool
}

func (p *poolAskResponse) Get() *AskResponse {
	x, ok := p.pool.Get().(*AskResponse)
	if !ok {
		return &AskResponse{}
	}
	return x
}

func (p *poolAskResponse) Put(x *AskResponse) {
	p.pool.Put(x)
}

var PoolAskResponse = poolAskResponse{}

const C_ReqSimple1 int64 = 396851362

type poolReqSimple1 struct {
	pool sync.Pool
}

func (p *poolReqSimple1) Get() *ReqSimple1 {
	x, ok := p.pool.Get().(*ReqSimple1)
	if !ok {
		return &ReqSimple1{}
	}
	return x
}

func (p *poolReqSimple1) Put(x *ReqSimple1) {
	x.P1 = x.P1[:0]
	p.pool.Put(x)
}

var PoolReqSimple1 = poolReqSimple1{}

const C_UpdateSimple1 int64 = 3946178210

type poolUpdateSimple1 struct {
	pool sync.Pool
}

func (p *poolUpdateSimple1) Get() *UpdateSimple1 {
	x, ok := p.pool.Get().(*UpdateSimple1)
	if !ok {
		return &UpdateSimple1{}
	}
	return x
}

func (p *poolUpdateSimple1) Put(x *UpdateSimple1) {
	x.P1 = x.P1[:0]
	p.pool.Put(x)
}

var PoolUpdateSimple1 = poolUpdateSimple1{}

const C_UpdateSimple2 int64 = 1916581656

type poolUpdateSimple2 struct {
	pool sync.Pool
}

func (p *poolUpdateSimple2) Get() *UpdateSimple2 {
	x, ok := p.pool.Get().(*UpdateSimple2)
	if !ok {
		return &UpdateSimple2{}
	}
	return x
}

func (p *poolUpdateSimple2) Put(x *UpdateSimple2) {
	p.pool.Put(x)
}

var PoolUpdateSimple2 = poolUpdateSimple2{}

const C_ResSimple1 int64 = 1434615775

type poolResSimple1 struct {
	pool sync.Pool
}

func (p *poolResSimple1) Get() *ResSimple1 {
	x, ok := p.pool.Get().(*ResSimple1)
	if !ok {
		return &ResSimple1{}
	}
	return x
}

func (p *poolResSimple1) Put(x *ResSimple1) {
	x.P1 = x.P1[:0]
	p.pool.Put(x)
}

var PoolResSimple1 = poolResSimple1{}

func init() {
	ConstructorNames[2179260159] = "ProtoMessage"
	ConstructorNames[1904100324] = "EchoRequest"
	ConstructorNames[4192619139] = "EchoResponse"
	ConstructorNames[3206229608] = "AskRequest"
	ConstructorNames[489087205] = "AskResponse"
	ConstructorNames[396851362] = "ReqSimple1"
	ConstructorNames[3946178210] = "UpdateSimple1"
	ConstructorNames[1916581656] = "UpdateSimple2"
	ConstructorNames[1434615775] = "ResSimple1"
}
