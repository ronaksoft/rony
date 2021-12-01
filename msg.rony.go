// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.16.4
// Source: msg.proto

package rony

import (
	bytes "bytes"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

var _ = pools.Imported

const C_MessageEnvelope uint64 = 11385363697628734430

type poolMessageEnvelope struct {
	pool sync.Pool
}

func (p *poolMessageEnvelope) Get() *MessageEnvelope {
	x, ok := p.pool.Get().(*MessageEnvelope)
	if !ok {
		x = &MessageEnvelope{}
	}

	return x
}

func (p *poolMessageEnvelope) Put(x *MessageEnvelope) {
	if x == nil {
		return
	}

	x.Constructor = 0
	x.RequestID = 0
	x.JsonEncoded = false
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
	z.JsonEncoded = x.JsonEncoded
	z.Message = append(z.Message[:0], x.Message...)
	z.Auth = append(z.Auth[:0], x.Auth...)
	for idx := range x.Header {
		if x.Header[idx] == nil {
			continue
		}
		xx := PoolKeyValue.Get()
		x.Header[idx].DeepCopy(xx)
		z.Header = append(z.Header, xx)
	}
}

func (x *MessageEnvelope) Clone() *MessageEnvelope {
	z := &MessageEnvelope{}
	x.DeepCopy(z)
	return z
}

func (x *MessageEnvelope) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *MessageEnvelope) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func factoryMessageEnvelope() registry.Message {
	return &MessageEnvelope{}
}

const C_Ping uint64 = 3584172304239689728

type poolPing struct {
	pool sync.Pool
}

func (p *poolPing) Get() *Ping {
	x, ok := p.pool.Get().(*Ping)
	if !ok {
		x = &Ping{}
	}

	return x
}

func (p *poolPing) Put(x *Ping) {
	if x == nil {
		return
	}

	x.ID = 0

	p.pool.Put(x)
}

var PoolPing = poolPing{}

func (x *Ping) DeepCopy(z *Ping) {
	z.ID = x.ID
}

func (x *Ping) Clone() *Ping {
	z := &Ping{}
	x.DeepCopy(z)
	return z
}

func (x *Ping) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Ping) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Ping) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Ping) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryPing() registry.Message {
	return &Ping{}
}

const C_Pong uint64 = 3584174915579805696

type poolPong struct {
	pool sync.Pool
}

func (p *poolPong) Get() *Pong {
	x, ok := p.pool.Get().(*Pong)
	if !ok {
		x = &Pong{}
	}

	return x
}

func (p *poolPong) Put(x *Pong) {
	if x == nil {
		return
	}

	x.ID = 0

	p.pool.Put(x)
}

var PoolPong = poolPong{}

func (x *Pong) DeepCopy(z *Pong) {
	z.ID = x.ID
}

func (x *Pong) Clone() *Pong {
	z := &Pong{}
	x.DeepCopy(z)
	return z
}

func (x *Pong) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Pong) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Pong) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Pong) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryPong() registry.Message {
	return &Pong{}
}

const C_KeyValue uint64 = 9712539143259722476

type poolKeyValue struct {
	pool sync.Pool
}

func (p *poolKeyValue) Get() *KeyValue {
	x, ok := p.pool.Get().(*KeyValue)
	if !ok {
		x = &KeyValue{}
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

func (x *KeyValue) Clone() *KeyValue {
	z := &KeyValue{}
	x.DeepCopy(z)
	return z
}

func (x *KeyValue) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *KeyValue) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *KeyValue) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *KeyValue) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryKeyValue() registry.Message {
	return &KeyValue{}
}

const C_MessageContainer uint64 = 8705479839606992738

type poolMessageContainer struct {
	pool sync.Pool
}

func (p *poolMessageContainer) Get() *MessageContainer {
	x, ok := p.pool.Get().(*MessageContainer)
	if !ok {
		x = &MessageContainer{}
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
	x.SyncRun = false

	p.pool.Put(x)
}

var PoolMessageContainer = poolMessageContainer{}

func (x *MessageContainer) DeepCopy(z *MessageContainer) {
	z.Length = x.Length
	for idx := range x.Envelopes {
		if x.Envelopes[idx] == nil {
			continue
		}
		xx := PoolMessageEnvelope.Get()
		x.Envelopes[idx].DeepCopy(xx)
		z.Envelopes = append(z.Envelopes, xx)
	}
	z.SyncRun = x.SyncRun
}

func (x *MessageContainer) Clone() *MessageContainer {
	z := &MessageContainer{}
	x.DeepCopy(z)
	return z
}

func (x *MessageContainer) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *MessageContainer) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *MessageContainer) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *MessageContainer) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryMessageContainer() registry.Message {
	return &MessageContainer{}
}

const C_Error uint64 = 3300063058275074048

type poolError struct {
	pool sync.Pool
}

func (p *poolError) Get() *Error {
	x, ok := p.pool.Get().(*Error)
	if !ok {
		x = &Error{}
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

func (x *Error) Clone() *Error {
	z := &Error{}
	x.DeepCopy(z)
	return z
}

func (x *Error) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Error) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Error) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Error) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryError() registry.Message {
	return &Error{}
}

const C_Redirect uint64 = 12783312295490183931

type poolRedirect struct {
	pool sync.Pool
}

func (p *poolRedirect) Get() *Redirect {
	x, ok := p.pool.Get().(*Redirect)
	if !ok {
		x = &Redirect{}
	}

	return x
}

func (p *poolRedirect) Put(x *Redirect) {
	if x == nil {
		return
	}

	x.Reason = 0
	for _, z := range x.Edges {
		PoolEdge.Put(z)
	}
	x.Edges = x.Edges[:0]
	x.WaitInSec = 0

	p.pool.Put(x)
}

var PoolRedirect = poolRedirect{}

func (x *Redirect) DeepCopy(z *Redirect) {
	z.Reason = x.Reason
	for idx := range x.Edges {
		if x.Edges[idx] == nil {
			continue
		}
		xx := PoolEdge.Get()
		x.Edges[idx].DeepCopy(xx)
		z.Edges = append(z.Edges, xx)
	}
	z.WaitInSec = x.WaitInSec
}

func (x *Redirect) Clone() *Redirect {
	z := &Redirect{}
	x.DeepCopy(z)
	return z
}

func (x *Redirect) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Redirect) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Redirect) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Redirect) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryRedirect() registry.Message {
	return &Redirect{}
}

const C_Edge uint64 = 3661905149145645056

type poolEdge struct {
	pool sync.Pool
}

func (p *poolEdge) Get() *Edge {
	x, ok := p.pool.Get().(*Edge)
	if !ok {
		x = &Edge{}
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

	p.pool.Put(x)
}

var PoolEdge = poolEdge{}

func (x *Edge) DeepCopy(z *Edge) {
	z.ReplicaSet = x.ReplicaSet
	z.ServerID = x.ServerID
	z.HostPorts = append(z.HostPorts[:0], x.HostPorts...)
}

func (x *Edge) Clone() *Edge {
	z := &Edge{}
	x.DeepCopy(z)
	return z
}

func (x *Edge) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Edge) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Edge) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Edge) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryEdge() registry.Message {
	return &Edge{}
}

const C_Edges uint64 = 3202852853167161344

type poolEdges struct {
	pool sync.Pool
}

func (p *poolEdges) Get() *Edges {
	x, ok := p.pool.Get().(*Edges)
	if !ok {
		x = &Edges{}
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
		if x.Nodes[idx] == nil {
			continue
		}
		xx := PoolEdge.Get()
		x.Nodes[idx].DeepCopy(xx)
		z.Nodes = append(z.Nodes, xx)
	}
}

func (x *Edges) Clone() *Edges {
	z := &Edges{}
	x.DeepCopy(z)
	return z
}

func (x *Edges) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Edges) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Edges) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Edges) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryEdges() registry.Message {
	return &Edges{}
}

const C_GetNodes uint64 = 17614373553603896039

type poolGetNodes struct {
	pool sync.Pool
}

func (p *poolGetNodes) Get() *GetNodes {
	x, ok := p.pool.Get().(*GetNodes)
	if !ok {
		x = &GetNodes{}
	}

	return x
}

func (p *poolGetNodes) Put(x *GetNodes) {
	if x == nil {
		return
	}

	x.ReplicaSet = x.ReplicaSet[:0]

	p.pool.Put(x)
}

var PoolGetNodes = poolGetNodes{}

func (x *GetNodes) DeepCopy(z *GetNodes) {
	z.ReplicaSet = append(z.ReplicaSet[:0], x.ReplicaSet...)
}

func (x *GetNodes) Clone() *GetNodes {
	z := &GetNodes{}
	x.DeepCopy(z)
	return z
}

func (x *GetNodes) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *GetNodes) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *GetNodes) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *GetNodes) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryGetNodes() registry.Message {
	return &GetNodes{}
}

const C_GetAllNodes uint64 = 2426091802953969362

type poolGetAllNodes struct {
	pool sync.Pool
}

func (p *poolGetAllNodes) Get() *GetAllNodes {
	x, ok := p.pool.Get().(*GetAllNodes)
	if !ok {
		x = &GetAllNodes{}
	}

	return x
}

func (p *poolGetAllNodes) Put(x *GetAllNodes) {
	if x == nil {
		return
	}

	p.pool.Put(x)
}

var PoolGetAllNodes = poolGetAllNodes{}

func (x *GetAllNodes) DeepCopy(z *GetAllNodes) {
}

func (x *GetAllNodes) Clone() *GetAllNodes {
	z := &GetAllNodes{}
	x.DeepCopy(z)
	return z
}

func (x *GetAllNodes) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *GetAllNodes) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *GetAllNodes) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *GetAllNodes) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryGetAllNodes() registry.Message {
	return &GetAllNodes{}
}

const C_HttpBody uint64 = 4148778432730789982

type poolHttpBody struct {
	pool sync.Pool
}

func (p *poolHttpBody) Get() *HttpBody {
	x, ok := p.pool.Get().(*HttpBody)
	if !ok {
		x = &HttpBody{}
	}

	return x
}

func (p *poolHttpBody) Put(x *HttpBody) {
	if x == nil {
		return
	}

	x.ContentType = ""
	for _, z := range x.Header {
		PoolKeyValue.Put(z)
	}
	x.Header = x.Header[:0]
	x.Body = x.Body[:0]

	p.pool.Put(x)
}

var PoolHttpBody = poolHttpBody{}

func (x *HttpBody) DeepCopy(z *HttpBody) {
	z.ContentType = x.ContentType
	for idx := range x.Header {
		if x.Header[idx] == nil {
			continue
		}
		xx := PoolKeyValue.Get()
		x.Header[idx].DeepCopy(xx)
		z.Header = append(z.Header, xx)
	}
	z.Body = append(z.Body[:0], x.Body...)
}

func (x *HttpBody) Clone() *HttpBody {
	z := &HttpBody{}
	x.DeepCopy(z)
	return z
}

func (x *HttpBody) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *HttpBody) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *HttpBody) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *HttpBody) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryHttpBody() registry.Message {
	return &HttpBody{}
}

// register constructors of the messages to the registry package
func init() {
	registry.Register(11385363697628734430, "MessageEnvelope", factoryMessageEnvelope)
	registry.Register(3584172304239689728, "Ping", factoryPing)
	registry.Register(3584174915579805696, "Pong", factoryPong)
	registry.Register(9712539143259722476, "KeyValue", factoryKeyValue)
	registry.Register(8705479839606992738, "MessageContainer", factoryMessageContainer)
	registry.Register(3300063058275074048, "Error", factoryError)
	registry.Register(12783312295490183931, "Redirect", factoryRedirect)
	registry.Register(3661905149145645056, "Edge", factoryEdge)
	registry.Register(3202852853167161344, "Edges", factoryEdges)
	registry.Register(17614373553603896039, "GetNodes", factoryGetNodes)
	registry.Register(2426091802953969362, "GetAllNodes", factoryGetAllNodes)
	registry.Register(4148778432730789982, "HttpBody", factoryHttpBody)

}

var _ = bytes.MinRead
