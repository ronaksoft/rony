package packageName

import (
	"sync"
)

const C_Message1 int64 = 3131464828

type poolMessage1 struct {
	pool sync.Pool
}

func (p *poolMessage1) Get() *Message1 {
	x, ok := p.pool.Get().(*Message1)
	if !ok {
		return &Message1{}
	}
	return x
}

func (p *poolMessage1) Put(x *Message1) {
	p.pool.Put(x)
}

var PoolMessage1 = poolMessage1{}

const C_Message2 int64 = 598674886

type poolMessage2 struct {
	pool sync.Pool
}

func (p *poolMessage2) Get() *Message2 {
	x, ok := p.pool.Get().(*Message2)
	if !ok {
		return &Message2{}
	}
	return x
}

func (p *poolMessage2) Put(x *Message2) {
	p.pool.Put(x)
}

var PoolMessage2 = poolMessage2{}

const C_Model1 int64 = 2074613123

type poolModel1 struct {
	pool sync.Pool
}

func (p *poolModel1) Get() *Model1 {
	x, ok := p.pool.Get().(*Model1)
	if !ok {
		return &Model1{}
	}
	return x
}

func (p *poolModel1) Put(x *Model1) {
	p.pool.Put(x)
}

var PoolModel1 = poolModel1{}

const C_Model2 int64 = 3802219577

type poolModel2 struct {
	pool sync.Pool
}

func (p *poolModel2) Get() *Model2 {
	x, ok := p.pool.Get().(*Model2)
	if !ok {
		return &Model2{}
	}
	return x
}

func (p *poolModel2) Put(x *Model2) {
	p.pool.Put(x)
}

var PoolModel2 = poolModel2{}

func init() {
	ConstructorNames[3131464828] = "Message1"
	ConstructorNames[598674886] = "Message2"
	ConstructorNames[2074613123] = "Model1"
	ConstructorNames[3802219577] = "Model2"
}
