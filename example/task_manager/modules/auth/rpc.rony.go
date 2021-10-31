// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.14.31
// Source: rpc.proto

package auth

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
	sync "sync"
)

var _ = pools.Imported

const C_RegisterRequest uint64 = 11503803976169819412

type poolRegisterRequest struct {
	pool sync.Pool
}

func (p *poolRegisterRequest) Get() *RegisterRequest {
	x, ok := p.pool.Get().(*RegisterRequest)
	if !ok {
		x = &RegisterRequest{}
	}

	return x
}

func (p *poolRegisterRequest) Put(x *RegisterRequest) {
	if x == nil {
		return
	}

	x.Username = ""
	x.FirstName = ""
	x.LastName = ""
	x.Password = ""

	p.pool.Put(x)
}

var PoolRegisterRequest = poolRegisterRequest{}

func (x *RegisterRequest) DeepCopy(z *RegisterRequest) {
	z.Username = x.Username
	z.FirstName = x.FirstName
	z.LastName = x.LastName
	z.Password = x.Password
}

func (x *RegisterRequest) Clone() *RegisterRequest {
	z := &RegisterRequest{}
	x.DeepCopy(z)
	return z
}

func (x *RegisterRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *RegisterRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *RegisterRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *RegisterRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapRegisterRequest(e registry.Envelope) (proto.Message, error) {
	x := &RegisterRequest{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *RegisterRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_RegisterRequest, x)
}

const C_LoginRequest uint64 = 18198549685643443149

type poolLoginRequest struct {
	pool sync.Pool
}

func (p *poolLoginRequest) Get() *LoginRequest {
	x, ok := p.pool.Get().(*LoginRequest)
	if !ok {
		x = &LoginRequest{}
	}

	return x
}

func (p *poolLoginRequest) Put(x *LoginRequest) {
	if x == nil {
		return
	}

	x.Username = ""
	x.Password = ""

	p.pool.Put(x)
}

var PoolLoginRequest = poolLoginRequest{}

func (x *LoginRequest) DeepCopy(z *LoginRequest) {
	z.Username = x.Username
	z.Password = x.Password
}

func (x *LoginRequest) Clone() *LoginRequest {
	z := &LoginRequest{}
	x.DeepCopy(z)
	return z
}

func (x *LoginRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *LoginRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *LoginRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *LoginRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapLoginRequest(e registry.Envelope) (proto.Message, error) {
	x := &LoginRequest{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *LoginRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_LoginRequest, x)
}

const C_Authorization uint64 = 9434083844645341581

type poolAuthorization struct {
	pool sync.Pool
}

func (p *poolAuthorization) Get() *Authorization {
	x, ok := p.pool.Get().(*Authorization)
	if !ok {
		x = &Authorization{}
	}

	return x
}

func (p *poolAuthorization) Put(x *Authorization) {
	if x == nil {
		return
	}

	x.SessionID = ""

	p.pool.Put(x)
}

var PoolAuthorization = poolAuthorization{}

func (x *Authorization) DeepCopy(z *Authorization) {
	z.SessionID = x.SessionID
}

func (x *Authorization) Clone() *Authorization {
	z := &Authorization{}
	x.DeepCopy(z)
	return z
}

func (x *Authorization) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Authorization) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Authorization) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Authorization) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func unwrapAuthorization(e registry.Envelope) (proto.Message, error) {
	x := &Authorization{}
	err := x.Unmarshal(e.GetMessage())
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Authorization) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Authorization, x)
}

const C_AuthRegister uint64 = 10692114844155491067
const C_AuthLogin uint64 = 13609756835726243961

// register constructors of the messages to the registry package
func init() {
	registry.Register(11503803976169819412, "RegisterRequest", unwrapRegisterRequest)
	registry.Register(18198549685643443149, "LoginRequest", unwrapLoginRequest)
	registry.Register(9434083844645341581, "Authorization", unwrapAuthorization)
	registry.Register(10692114844155491067, "AuthRegister", unwrapRegisterRequest)
	registry.Register(13609756835726243961, "AuthLogin", unwrapLoginRequest)

}

var _ = tools.TimeUnix()

type IAuth interface {
	Register(ctx *edge.RequestCtx, req *RegisterRequest, res *Authorization) *rony.Error
	Login(ctx *edge.RequestCtx, req *LoginRequest, res *Authorization) *rony.Error
}

func RegisterAuth(h IAuth, e *edge.Server, preHandlers ...edge.Handler) {
	w := authWrapper{
		h: h,
	}
	w.Register(e, func(c uint64) []edge.Handler {
		return preHandlers
	})
}

func RegisterAuthWithFunc(h IAuth, e *edge.Server, handlerFunc func(c uint64) []edge.Handler) {
	w := authWrapper{
		h: h,
	}
	w.Register(e, handlerFunc)
}

type authWrapper struct {
	h IAuth
}

func (sw *authWrapper) registerWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &RegisterRequest{}
	res := &Authorization{}

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.Register(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_Authorization, res)
	}
}
func (sw *authWrapper) loginWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &LoginRequest{}
	res := &Authorization{}

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.Login(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_Authorization, res)
	}
}

func (sw *authWrapper) Register(e *edge.Server, handlerFunc func(c uint64) []edge.Handler) {
	if handlerFunc == nil {
		handlerFunc = func(c uint64) []edge.Handler {
			return nil
		}
	}
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_AuthRegister).
			SetHandler(handlerFunc(C_AuthRegister)...).
			Append(sw.registerWrapper),
	)
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_AuthLogin).
			SetHandler(handlerFunc(C_AuthLogin)...).
			Append(sw.loginWrapper),
	)
}

func TunnelRequestAuthRegister(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *RegisterRequest, res *Authorization,
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_AuthRegister, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_Authorization:
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
func TunnelRequestAuthLogin(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *LoginRequest, res *Authorization,
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_AuthLogin, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_Authorization:
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

type AuthClient struct {
	c edgec.Client
}

func NewAuthClient(ec edgec.Client) *AuthClient {
	return &AuthClient{
		c: ec,
	}
}
func (c *AuthClient) Register(
	req *RegisterRequest, kvs ...*rony.KeyValue,
) (*Authorization, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_AuthRegister, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_Authorization:
		x := &Authorization{}
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
func (c *AuthClient) Login(
	req *LoginRequest, kvs ...*rony.KeyValue,
) (*Authorization, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_AuthLogin, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_Authorization:
		x := &Authorization{}
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

func prepareAuthCommand(cmd *cobra.Command, c edgec.Client) (*AuthClient, error) {
	// Bind current flags to the registered flags in config package
	err := config.BindCmdFlags(cmd)
	if err != nil {
		return nil, err
	}

	return NewAuthClient(c), nil
}

var genAuthRegisterCmd = func(h IAuthCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "register",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareAuthCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Register(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.StringFlag("username", "", ""),
		config.StringFlag("firstName", "", ""),
		config.StringFlag("lastName", "", ""),
		config.StringFlag("password", "", ""),
	)
	return cmd
}
var genAuthLoginCmd = func(h IAuthCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareAuthCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Login(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.StringFlag("username", "", ""),
		config.StringFlag("password", "", ""),
	)
	return cmd
}

type IAuthCli interface {
	Register(cli *AuthClient, cmd *cobra.Command, args []string) error
	Login(cli *AuthClient, cmd *cobra.Command, args []string) error
}

func RegisterAuthCli(h IAuthCli, c edgec.Client, rootCmd *cobra.Command) {
	subCommand := &cobra.Command{
		Use: "Auth",
	}
	subCommand.AddCommand(
		genAuthRegisterCmd(h, c),
		genAuthLoginCmd(h, c),
	)

	rootCmd.AddCommand(subCommand)
}

var _ = bytes.MinRead
