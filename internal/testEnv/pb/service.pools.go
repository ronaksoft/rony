package pb

import (
	"git.ronaksoft.com/ronak/rony/registry"
	"sync"
)

const C_Req1 int64 = 1772509555

type poolReq1 struct {
	pool sync.Pool
}

func (p *poolReq1) Get() *Req1 {
	x, ok := p.pool.Get().(*Req1)
	if !ok {
		return &Req1{}
	}
	return x
}

func (p *poolReq1) Put(x *Req1) {
	x.Item1 = 0
	p.pool.Put(x)
}

var PoolReq1 = poolReq1{}

const C_Req2 int64 = 4038002889

type poolReq2 struct {
	pool sync.Pool
}

func (p *poolReq2) Get() *Req2 {
	x, ok := p.pool.Get().(*Req2)
	if !ok {
		return &Req2{}
	}
	return x
}

func (p *poolReq2) Put(x *Req2) {
	x.Item1 = ""
	p.pool.Put(x)
}

var PoolReq2 = poolReq2{}

const C_Res1 int64 = 1536179185

type poolRes1 struct {
	pool sync.Pool
}

func (p *poolRes1) Get() *Res1 {
	x, ok := p.pool.Get().(*Res1)
	if !ok {
		return &Res1{}
	}
	return x
}

func (p *poolRes1) Put(x *Res1) {
	x.Item1 = 0
	p.pool.Put(x)
}

var PoolRes1 = poolRes1{}

const C_Res2 int64 = 3264834123

type poolRes2 struct {
	pool sync.Pool
}

func (p *poolRes2) Get() *Res2 {
	x, ok := p.pool.Get().(*Res2)
	if !ok {
		return &Res2{}
	}
	return x
}

func (p *poolRes2) Put(x *Res2) {
	x.Item1 = ""
	p.pool.Put(x)
}

var PoolRes2 = poolRes2{}

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
	x.Int = 0
	x.Bool = false
	x.Timestamp = 0
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
	x.Int = 0
	x.Bool = false
	x.Timestamp = 0
	x.Delay = 0
	x.ServerID = ""
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
	x.ServerID = ""
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
	x.Coordinator = ""
	x.Responder = ""
	p.pool.Put(x)
}

var PoolAskResponse = poolAskResponse{}

func init() {
	registry.RegisterConstructor(1772509555, "Req1")
	registry.RegisterConstructor(4038002889, "Req2")
	registry.RegisterConstructor(1536179185, "Res1")
	registry.RegisterConstructor(3264834123, "Res2")
	registry.RegisterConstructor(1904100324, "EchoRequest")
	registry.RegisterConstructor(4192619139, "EchoResponse")
	registry.RegisterConstructor(3206229608, "AskRequest")
	registry.RegisterConstructor(489087205, "AskResponse")
}
