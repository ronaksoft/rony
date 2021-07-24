// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.15.8
// Rony ver. v0.12.14
// Source: sample.proto

package service

import (
	bytes "bytes"
	fmt "fmt"
	rony "github.com/ronaksoft/rony"
	config "github.com/ronaksoft/rony/config"
	edge "github.com/ronaksoft/rony/edge"
	edgec "github.com/ronaksoft/rony/edgec"
	errors "github.com/ronaksoft/rony/errors"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	tools "github.com/ronaksoft/rony/tools"
	cobra "github.com/spf13/cobra"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	http "net/http"
	sync "sync"
)

var _ = pools.Imported

const C_EchoRequest int64 = 1904100324

type poolEchoRequest struct {
	pool sync.Pool
}

func (p *poolEchoRequest) Get() *EchoRequest {
	x, ok := p.pool.Get().(*EchoRequest)
	if !ok {
		x = &EchoRequest{}
	}

	return x
}

func (p *poolEchoRequest) Put(x *EchoRequest) {
	if x == nil {
		return
	}

	x.ID = 0
	x.RandomText = ""

	p.pool.Put(x)
}

var PoolEchoRequest = poolEchoRequest{}

func (x *EchoRequest) DeepCopy(z *EchoRequest) {
	z.ID = x.ID
	z.RandomText = x.RandomText
}

func (x *EchoRequest) Clone() *EchoRequest {
	z := &EchoRequest{}
	x.DeepCopy(z)
	return z
}

func (x *EchoRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *EchoRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EchoRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *EchoRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

const C_EchoResponse int64 = 4192619139

type poolEchoResponse struct {
	pool sync.Pool
}

func (p *poolEchoResponse) Get() *EchoResponse {
	x, ok := p.pool.Get().(*EchoResponse)
	if !ok {
		x = &EchoResponse{}
	}

	return x
}

func (p *poolEchoResponse) Put(x *EchoResponse) {
	if x == nil {
		return
	}

	x.ReqID = 0
	x.RandomText = ""

	p.pool.Put(x)
}

var PoolEchoResponse = poolEchoResponse{}

func (x *EchoResponse) DeepCopy(z *EchoResponse) {
	z.ReqID = x.ReqID
	z.RandomText = x.RandomText
}

func (x *EchoResponse) Clone() *EchoResponse {
	z := &EchoResponse{}
	x.DeepCopy(z)
	return z
}

func (x *EchoResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *EchoResponse) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EchoResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *EchoResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

const C_SampleEcho int64 = 3852587671

func init() {
	registry.RegisterConstructor(1904100324, "EchoRequest")
	registry.RegisterConstructor(4192619139, "EchoResponse")
	registry.RegisterConstructor(3852587671, "SampleEcho")
}

var _ = tools.TimeUnix()

type ISample interface {
	Echo(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
}

func RegisterSample(h ISample, e *edge.Server, preHandlers ...edge.Handler) {
	w := sampleWrapper{
		h: h,
	}
	w.Register(e, func(c int64) []edge.Handler {
		return preHandlers
	})
}

func RegisterSampleWithFunc(h ISample, e *edge.Server, handlerFunc func(c int64) []edge.Handler) {
	w := sampleWrapper{
		h: h,
	}
	w.Register(e, handlerFunc)
}

type sampleWrapper struct {
	h ISample
}

func (sw *sampleWrapper) echoWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	sw.h.Echo(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func (sw *sampleWrapper) Register(e *edge.Server, handlerFunc func(c int64) []edge.Handler) {
	if handlerFunc == nil {
		handlerFunc = func(c int64) []edge.Handler {
			return nil
		}
	}
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_SampleEcho).
			SetHandler(handlerFunc(C_SampleEcho)...).
			Append(sw.echoWrapper),
	)
	e.SetRestProxy(
		"get", "/echo/:ID/:Random",
		edge.NewRestProxy(sw.echoRestClient, sw.echoRestServer),
	)
}

func TunnelRequestSampleEcho(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse, kvs ...*rony.KeyValue) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleEcho, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
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
		return errors.ErrUnexpectedTunnelResponse
	}
}

func (sw *sampleWrapper) echoRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	req.ID = tools.StrToInt64(tools.GetString(conn.Get("ID"), "0"))
	req.RandomText = tools.GetString(conn.Get("Random"), "")

	ctx.FillEnvelope(conn.ConnID(), C_SampleEcho, req)
	return nil
}
func (sw *sampleWrapper) echoRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_EchoResponse:
			x := &EchoResponse{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			b, err = x.MarshalJSON()
			if err != nil {
				return
			}
			err = conn.WriteBinary(ctx.StreamID(), b)
			return
		case rony.C_Error:
			x := &rony.Error{}
			_ = x.Unmarshal(envelope.Message)
			err = x
			return
		case rony.C_Redirect:
			x := &rony.Redirect{}
			_ = x.Unmarshal(envelope.Message)
			if len(x.Edges) == 0 || len(x.Edges[0].HostPorts) == 0 {
				break
			}
			switch x.Reason {
			case rony.RedirectReason_ReplicaSetSession:
				conn.Redirect(http.StatusPermanentRedirect, x.Edges[0].HostPorts[0])
			case rony.RedirectReason_ReplicaSetRequest:
				conn.Redirect(http.StatusTemporaryRedirect, x.Edges[0].HostPorts[0])
			}
			return
		}
		err = errors.ErrUnexpectedResponse
	}) {
		err = errors.ErrInternalServer
	}

	return
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
	out.Fill(c.c.GetRequestID(), C_SampleEcho, req, kvs...)
	err := c.c.Send(out, in)
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
		_ = x.Unmarshal(in.Message)
		return nil, x
	default:
		return nil, fmt.Errorf("unknown message :%d", in.GetConstructor())
	}
}

func prepareSampleCommand(cmd *cobra.Command, c edgec.Client) (*SampleClient, error) {
	// Bind current flags to the registered flags in config package
	err := config.BindCmdFlags(cmd)
	if err != nil {
		return nil, err
	}

	return NewSampleClient(c), nil
}

var genSampleEchoCmd = func(h ISampleCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "echo",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareSampleCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Echo(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.Int64Flag("id", 0, ""),
		config.StringFlag("randomText", "", ""),
	)
	return cmd
}

type ISampleCli interface {
	Echo(cli *SampleClient, cmd *cobra.Command, args []string) error
}

func RegisterSampleCli(h ISampleCli, c edgec.Client, rootCmd *cobra.Command) {
	subCommand := &cobra.Command{
		Use: "Sample",
	}
	subCommand.AddCommand(
		genSampleEchoCmd(h, c),
	)

	rootCmd.AddCommand(subCommand)
}

var _ = bytes.MinRead
