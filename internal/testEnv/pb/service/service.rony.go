package service

import (
	fmt "fmt"
	rony "github.com/ronaksoft/rony"
	config "github.com/ronaksoft/rony/config"
	edge "github.com/ronaksoft/rony/edge"
	edgec "github.com/ronaksoft/rony/edgec"
	registry "github.com/ronaksoft/rony/registry"
	cobra "github.com/spf13/cobra"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

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
	x.Timestamp = 0
	x.ReplicaSet = 0
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
	x.Responder = ""
	x.Timestamp = 0
	x.Delay = 0
	x.ServerID = ""
	p.pool.Put(x)
}

var PoolEchoResponse = poolEchoResponse{}

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
	x.Param1 = 0
	x.Param2 = ""
	if x.M2 != nil {
		PoolMessage2.Put(x.M2)
		x.M2 = nil
	}
	x.M2S = x.M2S[:0]
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
	x.Param1 = 0
	x.P2 = x.P2[:0]
	x.P3 = x.P3[:0]
	if x.M1 != nil {
		PoolMessage1.Put(x.M1)
		x.M1 = nil
	}
	p.pool.Put(x)
}

var PoolMessage2 = poolMessage2{}

func init() {
	registry.RegisterConstructor(1904100324, "EchoRequest")
	registry.RegisterConstructor(4192619139, "EchoResponse")
	registry.RegisterConstructor(3131464828, "Message1")
	registry.RegisterConstructor(598674886, "Message2")
	registry.RegisterConstructor(3073810188, "Echo")
	registry.RegisterConstructor(27569121, "EchoLeaderOnly")
	registry.RegisterConstructor(3809767204, "EchoTunnel")
	registry.RegisterConstructor(3639218737, "EchoDelay")
}

func (x *EchoRequest) DeepCopy(z *EchoRequest) {
	z.Int = x.Int
	z.Timestamp = x.Timestamp
	z.ReplicaSet = x.ReplicaSet
}

func (x *EchoResponse) DeepCopy(z *EchoResponse) {
	z.Int = x.Int
	z.Responder = x.Responder
	z.Timestamp = x.Timestamp
	z.Delay = x.Delay
	z.ServerID = x.ServerID
}

func (x *Message1) DeepCopy(z *Message1) {
	z.Param1 = x.Param1
	z.Param2 = x.Param2
	if x.M2 != nil {
		z.M2 = PoolMessage2.Get()
		x.M2.DeepCopy(z.M2)
	}
	for idx := range x.M2S {
		if x.M2S[idx] != nil {
			xx := PoolMessage2.Get()
			x.M2S[idx].DeepCopy(xx)
			z.M2S = append(z.M2S, xx)
		}
	}
}

func (x *Message2) DeepCopy(z *Message2) {
	z.Param1 = x.Param1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P3 = append(z.P3[:0], x.P3...)
	if x.M1 != nil {
		z.M1 = PoolMessage1.Get()
		x.M1.DeepCopy(z.M1)
	}
}

func (x *EchoRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_EchoRequest, x)
}

func (x *EchoResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_EchoResponse, x)
}

func (x *Message1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Message1, x)
}

func (x *Message2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Message2, x)
}

func (x *EchoRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EchoResponse) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Message1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Message2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EchoRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *EchoResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Message1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Message2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_Echo int64 = 3073810188
const C_EchoLeaderOnly int64 = 27569121
const C_EchoTunnel int64 = 3809767204
const C_EchoDelay int64 = 3639218737

type ISample interface {
	Echo(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
	EchoLeaderOnly(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
	EchoTunnel(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
	EchoDelay(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
}

type sampleWrapper struct {
	h ISample
}

func RegisterSample(h ISample, e *edge.Server) {
	w := sampleWrapper{
		h: h,
	}
	w.Register(e)
}

func (sw *sampleWrapper) Register(e *edge.Server) {
	e.SetHandlers(C_Echo, false, sw.EchoWrapper)
	e.SetHandlers(C_EchoLeaderOnly, true, sw.EchoLeaderOnlyWrapper)
	e.SetHandlers(C_EchoTunnel, true, sw.EchoTunnelWrapper)
	e.SetHandlers(C_EchoDelay, true, sw.EchoDelayWrapper)
}

func (sw *sampleWrapper) EchoWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.Echo(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func (sw *sampleWrapper) EchoLeaderOnlyWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.EchoLeaderOnly(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func (sw *sampleWrapper) EchoTunnelWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.EchoTunnel(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func (sw *sampleWrapper) EchoDelayWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.EchoDelay(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func ExecuteRemoteEcho(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_EchoRequest, req)
	err := ctx.ExecuteRemote(replicaSet, false, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_EchoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return edge.ErrUnexpectedTunnelResponse
	}
}

func ExecuteRemoteEchoLeaderOnly(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_EchoRequest, req)
	err := ctx.ExecuteRemote(replicaSet, true, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_EchoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return edge.ErrUnexpectedTunnelResponse
	}
}

func ExecuteRemoteEchoTunnel(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_EchoRequest, req)
	err := ctx.ExecuteRemote(replicaSet, true, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_EchoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return edge.ErrUnexpectedTunnelResponse
	}
}

func ExecuteRemoteEchoDelay(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_EchoRequest, req)
	err := ctx.ExecuteRemote(replicaSet, true, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_EchoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return edge.ErrUnexpectedTunnelResponse
	}
}

type SampleClient struct {
	c edgec.Client
}

func NewSampleClient(ec edgec.Client) *SampleClient {
	return &SampleClient{
		c: ec,
	}
}

func (c *SampleClient) Echo(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Echo, req, kvs...)
	err := c.c.Send(out, in, false)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = proto.Unmarshal(in.Message, x)
		return nil, fmt.Errorf("%s:%s", x.GetCode(), x.GetItems())
	default:
		return nil, fmt.Errorf("unknown message: %d", in.GetConstructor())
	}
}

func (c *SampleClient) EchoLeaderOnly(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_EchoLeaderOnly, req, kvs...)
	err := c.c.Send(out, in, true)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = proto.Unmarshal(in.Message, x)
		return nil, fmt.Errorf("%s:%s", x.GetCode(), x.GetItems())
	default:
		return nil, fmt.Errorf("unknown message: %d", in.GetConstructor())
	}
}

func (c *SampleClient) EchoTunnel(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_EchoTunnel, req, kvs...)
	err := c.c.Send(out, in, true)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = proto.Unmarshal(in.Message, x)
		return nil, fmt.Errorf("%s:%s", x.GetCode(), x.GetItems())
	default:
		return nil, fmt.Errorf("unknown message: %d", in.GetConstructor())
	}
}

func (c *SampleClient) EchoDelay(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_EchoDelay, req, kvs...)
	err := c.c.Send(out, in, true)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = proto.Unmarshal(in.Message, x)
		return nil, fmt.Errorf("%s:%s", x.GetCode(), x.GetItems())
	default:
		return nil, fmt.Errorf("unknown message: %d", in.GetConstructor())
	}
}

func prepareSampleCommand(cmd *cobra.Command) (*SampleClient, error) {
	// Bind the current flags to registered flags in config package
	err := config.BindCmdFlags(cmd)
	if err != nil {
		return nil, err
	}

	httpC := edgec.NewHttp(edgec.HttpConfig{
		Name:         "",
		SeedHostPort: fmt.Sprintf("%s:%d", config.GetString("host"), config.GetInt("port")),
	})

	err = httpC.Start()
	if err != nil {
		return nil, err
	}
	return NewSampleClient(httpC), nil
}

var echoCmd = &cobra.Command{
	Use: "echo",
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := prepareSampleCommand(cmd)
		_ = cli
		if err != nil {
			return err
		}
		req := &EchoRequest{
			// int64
			// int64
			// uint64
		}
		_ = req
		return nil
	},
}
var echoLeaderOnlyCmd = &cobra.Command{
	Use: "echoLeaderOnly",
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := prepareSampleCommand(cmd)
		_ = cli
		if err != nil {
			return err
		}
		req := &EchoRequest{
			// int64
			// int64
			// uint64
		}
		_ = req
		return nil
	},
}
var echoTunnelCmd = &cobra.Command{
	Use: "echoTunnel",
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := prepareSampleCommand(cmd)
		_ = cli
		if err != nil {
			return err
		}
		req := &EchoRequest{
			// int64
			// int64
			// uint64
		}
		_ = req
		return nil
	},
}
var echoDelayCmd = &cobra.Command{
	Use: "echoDelay",
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := prepareSampleCommand(cmd)
		_ = cli
		if err != nil {
			return err
		}
		req := &EchoRequest{
			// int64
			// int64
			// uint64
		}
		_ = req
		return nil
	},
}

type ISampleCli interface {
	Echo(cli edgec.Client, cmd *cobra.Command, args []string) error
	EchoLeaderOnly(cli edgec.Client, cmd *cobra.Command, args []string) error
	EchoTunnel(cli edgec.Client, cmd *cobra.Command, args []string) error
	EchoDelay(cli edgec.Client, cmd *cobra.Command, args []string) error
}

func RegisterSampleCli(h ISampleCli, rootCmd *cobra.Command) {
	rootCmd.AddCommand(echoCmd)
	rootCmd.AddCommand(echoLeaderOnlyCmd)
	rootCmd.AddCommand(echoTunnelCmd)
	rootCmd.AddCommand(echoDelayCmd)
}
