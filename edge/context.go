package edge

import (
	"bufio"
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cluster"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/proto"
	"net"
	"reflect"
	"sync"
	"time"
)

/*
   Creation Time: 2019 - Jun - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type MessageKind byte

const (
	_ MessageKind = iota
	GatewayMessage
	ReplicaMessage
	TunnelMessage
)

var (
	kindNames = map[MessageKind]string{
		GatewayMessage: "GatewayMessage",
		ReplicaMessage: "ReplicatedMessage",
		TunnelMessage:  "TunnelMessage",
	}
)

func (c MessageKind) String() string {
	return kindNames[c]
}

// DispatchCtx
type DispatchCtx struct {
	streamID          int64
	serverID          []byte
	conn              rony.Conn
	req               *rony.MessageEnvelope
	cluster           cluster.Cluster
	gatewayDispatcher Dispatcher
	kind              MessageKind
	buf               *tools.LinkedList
	// KeyValue Store Parameters
	mtx sync.RWMutex
	kv  map[string]interface{}
}

func newDispatchCtx(edge *Server) *DispatchCtx {
	return &DispatchCtx{
		cluster:           edge.cluster,
		gatewayDispatcher: edge.gatewayDispatcher,
		req:               &rony.MessageEnvelope{},
		kv:                make(map[string]interface{}, 3),
		buf:               tools.NewLinkedList(),
	}
}

func (ctx *DispatchCtx) reset() {
	for k := range ctx.kv {
		delete(ctx.kv, k)
	}
	ctx.buf.Reset()
}

func (ctx *DispatchCtx) ServerID() string {
	return string(ctx.serverID)
}

func (ctx *DispatchCtx) Debug() {
	fmt.Println("###")
	t := reflect.Indirect(reflect.ValueOf(ctx))
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Type().Field(i).Name, t.Type().Field(i).Offset, t.Type().Field(i).Type.Size())
	}
}

func (ctx *DispatchCtx) Conn() rony.Conn {
	return ctx.conn
}

func (ctx *DispatchCtx) StreamID() int64 {
	return ctx.streamID
}

func (ctx *DispatchCtx) FillEnvelope(requestID uint64, constructor int64, payload []byte, auth []byte, kv ...*rony.KeyValue) {
	ctx.req.RequestID = requestID
	ctx.req.Constructor = constructor
	ctx.req.Message = append(ctx.req.Message[:0], payload...)
	ctx.req.Auth = append(ctx.req.Auth[:0], auth...)
	if cap(ctx.req.Header) >= len(kv) {
		ctx.req.Header = ctx.req.Header[:len(kv)]
	} else {
		ctx.req.Header = make([]*rony.KeyValue, len(kv))
	}
	for idx, kv := range kv {
		if ctx.req.Header[idx] == nil {
			ctx.req.Header[idx] = &rony.KeyValue{}
		}
		kv.DeepCopy(ctx.req.Header[idx])
	}
}

func (ctx *DispatchCtx) Set(key string, v interface{}) {
	ctx.mtx.Lock()
	ctx.kv[key] = v
	ctx.mtx.Unlock()
}

func (ctx *DispatchCtx) Get(key string) interface{} {
	ctx.mtx.RLock()
	v := ctx.kv[key]
	ctx.mtx.RUnlock()
	return v
}

func (ctx *DispatchCtx) GetBytes(key string, defaultValue []byte) []byte {
	v, ok := ctx.Get(key).([]byte)
	if ok {
		return v
	}
	return defaultValue
}

func (ctx *DispatchCtx) GetString(key string, defaultValue string) string {
	v := ctx.Get(key)
	switch x := v.(type) {
	case []byte:
		return tools.ByteToStr(x)
	case string:
		return x
	default:
		return defaultValue
	}
}

func (ctx *DispatchCtx) Kind() MessageKind {
	return ctx.kind
}

func (ctx *DispatchCtx) GetInt64(key string, defaultValue int64) int64 {
	v, ok := ctx.Get(key).(int64)
	if ok {
		return v
	}
	return defaultValue
}

func (ctx *DispatchCtx) GetBool(key string) bool {
	v, ok := ctx.Get(key).(bool)
	if ok {
		return v
	}
	return false
}

func (ctx *DispatchCtx) UnmarshalEnvelope(data []byte) error {
	return proto.Unmarshal(data, ctx.req)
}

func (ctx *DispatchCtx) BufferPush(m *rony.MessageEnvelope) {
	ctx.buf.Append(m)
}

func (ctx *DispatchCtx) BufferPop() *rony.MessageEnvelope {
	v := ctx.buf.PickHeadData()
	if v != nil {
		return v.(*rony.MessageEnvelope)
	}
	return nil
}

func (ctx *DispatchCtx) BufferSize() int32 {
	return ctx.buf.Size()
}

// RequestCtx
type RequestCtx struct {
	dispatchCtx *DispatchCtx
	reqID       uint64
	nextChan    chan struct{}
	quickReturn bool
	stop        bool
}

func newRequestCtx() *RequestCtx {
	return &RequestCtx{
		nextChan: make(chan struct{}, 1),
	}
}

func (ctx *RequestCtx) Return() {
	if ctx.quickReturn {
		ctx.nextChan <- struct{}{}
	}
}

func (ctx *RequestCtx) ConnID() uint64 {
	if ctx.dispatchCtx.Conn() != nil {
		return ctx.dispatchCtx.Conn().ConnID()
	}
	return 0
}

func (ctx *RequestCtx) Conn() rony.Conn {
	return ctx.dispatchCtx.Conn()
}

func (ctx *RequestCtx) ReqID() uint64 {
	return ctx.reqID
}

func (ctx *RequestCtx) ServerID() string {
	return string(ctx.dispatchCtx.serverID)
}

func (ctx *RequestCtx) Kind() MessageKind {
	return ctx.dispatchCtx.kind
}

func (ctx *RequestCtx) StopExecution() {
	ctx.stop = true
}

func (ctx *RequestCtx) Stopped() bool {
	return ctx.stop
}

func (ctx *RequestCtx) Set(key string, v interface{}) {
	ctx.dispatchCtx.mtx.Lock()
	ctx.dispatchCtx.kv[key] = v
	ctx.dispatchCtx.mtx.Unlock()
}

func (ctx *RequestCtx) Get(key string) interface{} {
	return ctx.dispatchCtx.Get(key)
}

func (ctx *RequestCtx) GetBytes(key string, defaultValue []byte) []byte {
	return ctx.dispatchCtx.GetBytes(key, defaultValue)
}

func (ctx *RequestCtx) GetString(key string, defaultValue string) string {
	return ctx.dispatchCtx.GetString(key, defaultValue)
}

func (ctx *RequestCtx) GetInt64(key string, defaultValue int64) int64 {
	return ctx.dispatchCtx.GetInt64(key, defaultValue)
}

func (ctx *RequestCtx) GetBool(key string) bool {
	return ctx.dispatchCtx.GetBool(key)
}

func (ctx *RequestCtx) PushRedirectLeader() {
	ctxCluster := ctx.dispatchCtx.cluster
	if leaderID := ctxCluster.RaftLeaderID(); leaderID == "" {
		ctx.PushError(rony.ErrCodeUnavailable, rony.ErrItemRaftLeader)
		return
	}
	r := rony.PoolRedirect.Get()
	r.Reason = rony.RedirectReason_ReplicaMaster
	r.WaitInSec = 0
	members := ctxCluster.RaftMembers(ctxCluster.ReplicaSet())
	nodeInfos := make([]*rony.Edge, 0, len(members))
	for _, m := range members {
		ni := m.Proto(rony.PoolEdge.Get())
		if ni.Leader {
			r.Leader = ni
		} else {
			r.Followers = append(r.Followers, ni)
		}
		nodeInfos = append(nodeInfos, ni)
	}
	ctx.PushMessage(rony.C_Redirect, r)
	for _, p := range nodeInfos {
		rony.PoolEdge.Put(p)
	}
	rony.PoolRedirect.Put(r)
}

func (ctx *RequestCtx) PushRedirectSession(replicaSet uint64, wait time.Duration) {
	// TODO:: implement it
}

func (ctx *RequestCtx) PushRedirectRequest(replicaSet uint64) {
	// TODO:: implement it
}

func (ctx *RequestCtx) PushMessage(constructor int64, proto proto.Message) {
	ctx.PushCustomMessage(ctx.ReqID(), constructor, proto)
}

func (ctx *RequestCtx) PushCustomMessage(requestID uint64, constructor int64, proto proto.Message, kvs ...*rony.KeyValue) {
	if ctx.dispatchCtx.kind == ReplicaMessage {
		return
	}

	envelope := rony.PoolMessageEnvelope.Get()
	envelope.Fill(requestID, constructor, proto)
	envelope.Header = append(envelope.Header[:0], kvs...)

	switch ctx.dispatchCtx.kind {
	case GatewayMessage:
		ctx.dispatchCtx.gatewayDispatcher.OnMessage(ctx.dispatchCtx, envelope)
	case TunnelMessage:
		ctx.dispatchCtx.BufferPush(envelope.Clone())
	}

	rony.PoolMessageEnvelope.Put(envelope)
}

func (ctx *RequestCtx) PushError(code, item string) {
	ctx.PushCustomError(code, item, "")
}

func (ctx *RequestCtx) PushCustomError(code, item string, desc string) {
	ctx.PushMessage(rony.C_Error, &rony.Error{
		Code:        code,
		Items:       item,
		Description: desc,
	})
	ctx.stop = true
}

func (ctx *RequestCtx) Cluster() cluster.Cluster {
	return ctx.dispatchCtx.cluster
}

func (ctx *RequestCtx) ExecuteRemote(replicaSet uint64, onlyLeader bool, req, res *rony.MessageEnvelope) error {
	target := getReplicaMember(ctx.dispatchCtx.cluster, replicaSet, onlyLeader)
	if target == nil {
		return ErrMemberNotFound
	}
	if len(target.TunnelAddr()) == 0 {
		return ErrNoTunnelAddrs
	}

	conn, err := net.Dial("udp", target.TunnelAddr()[0])
	if err != nil {
		return err
	}

	// Get a rony.TunnelMessage from pool and put it back into the pool when we are done
	tmOut := rony.PoolTunnelMessage.Get()
	defer rony.PoolTunnelMessage.Put(tmOut)
	tmIn := rony.PoolTunnelMessage.Get()
	defer rony.PoolTunnelMessage.Put(tmIn)
	tmOut.SenderID = ctx.dispatchCtx.serverID
	tmOut.SenderReplicaSet = ctx.dispatchCtx.cluster.ReplicaSet()
	tmOut.Envelope = rony.PoolMessageEnvelope.Get()
	req.DeepCopy(tmOut.Envelope)

	// Marshal and send over the wire
	mo := proto.MarshalOptions{UseCachedSize: true}
	buf := pools.Buffer.GetCap(mo.Size(tmOut))
	b, _ := mo.MarshalAppend(*buf.Bytes(), tmOut)
	buf.SetBytes(&b)
	n, err := conn.Write(*buf.Bytes())
	if err != nil {
		return err
	}
	pools.Buffer.Put(buf)

	// Wait for response and unmarshal it
	buf = pools.Buffer.GetLen(4096)
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	n, err = bufio.NewReader(conn).Read(*buf.Bytes())
	if err != nil || n == 0 {
		return err
	}
	err = tmIn.Unmarshal((*buf.Bytes())[:n])
	if err != nil {
		return err
	}
	pools.Buffer.Put(buf)

	// deep copy
	tmIn.Envelope.DeepCopy(res)

	return nil
}

func getReplicaMember(cluster cluster.Cluster, replicaSet uint64, onlyLeader bool) (target cluster.Member) {
	members := cluster.RaftMembers(replicaSet)
	if len(members) == 0 {
		return nil
	}
	for idx := range members {
		if onlyLeader {
			if members[idx].RaftState() == rony.RaftState_Leader {
				target = members[idx]
				break
			}
		} else {
			target = members[idx]
		}
	}
	return
}

func (ctx *RequestCtx) Log() log.Logger {
	return log.DefaultLogger
}

func (ctx *RequestCtx) FindReplicaSet(pageID uint32) (uint64, error) {
	p := &rony.Page{}
	_, err := rony.ReadPage(pageID, p)
	if err == nil {
		return p.GetReplicaSet(), nil
	}
	thisReplicaSet := ctx.Cluster().ReplicaSet()
	p.ID = pageID
	p.ReplicaSet = thisReplicaSet
	if thisReplicaSet != 1 {
		req := rony.PoolMessageEnvelope.Get()
		defer rony.PoolMessageEnvelope.Put(req)
		res := rony.PoolMessageEnvelope.Get()
		defer rony.PoolMessageEnvelope.Put(res)
		getPage := rony.PoolGetPage.Get()
		defer rony.PoolGetPage.Put(getPage)
		getPage.PageID = pageID
		getPage.ReplicaSet = ctx.Cluster().ReplicaSet()
		req.Fill(uint64(tools.FastRand()<<31|tools.FastRand()), rony.C_GetPage, getPage)
		err = ctx.ExecuteRemote(1, true, req, res)
		if err != nil {
			return 0, err
		}

		switch res.GetConstructor() {
		case rony.C_Page:
			_ = p.Unmarshal(res.GetMessage())
		case rony.C_Error:
			x := &rony.Error{}
			_ = x.Unmarshal(res.GetMessage())
			return 0, x
		default:
			panic("BUG!! invalid response")
		}

	}

	err = rony.SavePage(p)
	if err != nil {
		return 0, err
	}
	return p.ReplicaSet, nil
}

func (ctx *RequestCtx) LocalReplicaSet() uint64 {
	if ctx.dispatchCtx.cluster == nil {
		return 0
	}
	return ctx.dispatchCtx.cluster.ReplicaSet()
}
