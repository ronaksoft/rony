// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.12.34
// Source: rpc.proto

package task

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

const C_CreateRequest int64 = 2229707971

type poolCreateRequest struct {
	pool sync.Pool
}

func (p *poolCreateRequest) Get() *CreateRequest {
	x, ok := p.pool.Get().(*CreateRequest)
	if !ok {
		x = &CreateRequest{}
	}

	return x
}

func (p *poolCreateRequest) Put(x *CreateRequest) {
	if x == nil {
		return
	}

	x.Title = ""
	x.TODOs = x.TODOs[:0]
	x.DueDate = 0

	p.pool.Put(x)
}

var PoolCreateRequest = poolCreateRequest{}

func (x *CreateRequest) DeepCopy(z *CreateRequest) {
	z.Title = x.Title
	z.TODOs = append(z.TODOs[:0], x.TODOs...)
	z.DueDate = x.DueDate
}

func (x *CreateRequest) Clone() *CreateRequest {
	z := &CreateRequest{}
	x.DeepCopy(z)
	return z
}

func (x *CreateRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *CreateRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *CreateRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *CreateRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *CreateRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_CreateRequest, x)
}

const C_GetRequest int64 = 3359917651

type poolGetRequest struct {
	pool sync.Pool
}

func (p *poolGetRequest) Get() *GetRequest {
	x, ok := p.pool.Get().(*GetRequest)
	if !ok {
		x = &GetRequest{}
	}

	return x
}

func (p *poolGetRequest) Put(x *GetRequest) {
	if x == nil {
		return
	}

	x.TaskID = 0

	p.pool.Put(x)
}

var PoolGetRequest = poolGetRequest{}

func (x *GetRequest) DeepCopy(z *GetRequest) {
	z.TaskID = x.TaskID
}

func (x *GetRequest) Clone() *GetRequest {
	z := &GetRequest{}
	x.DeepCopy(z)
	return z
}

func (x *GetRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *GetRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *GetRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *GetRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *GetRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_GetRequest, x)
}

const C_DeleteRequest int64 = 4088608791

type poolDeleteRequest struct {
	pool sync.Pool
}

func (p *poolDeleteRequest) Get() *DeleteRequest {
	x, ok := p.pool.Get().(*DeleteRequest)
	if !ok {
		x = &DeleteRequest{}
	}

	return x
}

func (p *poolDeleteRequest) Put(x *DeleteRequest) {
	if x == nil {
		return
	}

	x.TaskID = 0

	p.pool.Put(x)
}

var PoolDeleteRequest = poolDeleteRequest{}

func (x *DeleteRequest) DeepCopy(z *DeleteRequest) {
	z.TaskID = x.TaskID
}

func (x *DeleteRequest) Clone() *DeleteRequest {
	z := &DeleteRequest{}
	x.DeepCopy(z)
	return z
}

func (x *DeleteRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *DeleteRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *DeleteRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *DeleteRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *DeleteRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_DeleteRequest, x)
}

const C_ListRequest int64 = 307567737

type poolListRequest struct {
	pool sync.Pool
}

func (p *poolListRequest) Get() *ListRequest {
	x, ok := p.pool.Get().(*ListRequest)
	if !ok {
		x = &ListRequest{}
	}

	return x
}

func (p *poolListRequest) Put(x *ListRequest) {
	if x == nil {
		return
	}

	x.Offset = 0
	x.Limit = 0

	p.pool.Put(x)
}

var PoolListRequest = poolListRequest{}

func (x *ListRequest) DeepCopy(z *ListRequest) {
	z.Offset = x.Offset
	z.Limit = x.Limit
}

func (x *ListRequest) Clone() *ListRequest {
	z := &ListRequest{}
	x.DeepCopy(z)
	return z
}

func (x *ListRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *ListRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *ListRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *ListRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *ListRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_ListRequest, x)
}

const C_TaskView int64 = 2236647323

type poolTaskView struct {
	pool sync.Pool
}

func (p *poolTaskView) Get() *TaskView {
	x, ok := p.pool.Get().(*TaskView)
	if !ok {
		x = &TaskView{}
	}

	return x
}

func (p *poolTaskView) Put(x *TaskView) {
	if x == nil {
		return
	}

	x.ID = 0
	x.Title = ""
	x.TODOs = x.TODOs[:0]
	x.DueDate = 0

	p.pool.Put(x)
}

var PoolTaskView = poolTaskView{}

func (x *TaskView) DeepCopy(z *TaskView) {
	z.ID = x.ID
	z.Title = x.Title
	z.TODOs = append(z.TODOs[:0], x.TODOs...)
	z.DueDate = x.DueDate
}

func (x *TaskView) Clone() *TaskView {
	z := &TaskView{}
	x.DeepCopy(z)
	return z
}

func (x *TaskView) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *TaskView) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *TaskView) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *TaskView) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *TaskView) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_TaskView, x)
}

const C_TaskViewMany int64 = 2061119489

type poolTaskViewMany struct {
	pool sync.Pool
}

func (p *poolTaskViewMany) Get() *TaskViewMany {
	x, ok := p.pool.Get().(*TaskViewMany)
	if !ok {
		x = &TaskViewMany{}
	}

	return x
}

func (p *poolTaskViewMany) Put(x *TaskViewMany) {
	if x == nil {
		return
	}

	for _, z := range x.Tasks {
		PoolTaskView.Put(z)
	}
	x.Tasks = x.Tasks[:0]

	p.pool.Put(x)
}

var PoolTaskViewMany = poolTaskViewMany{}

func (x *TaskViewMany) DeepCopy(z *TaskViewMany) {
	for idx := range x.Tasks {
		if x.Tasks[idx] == nil {
			continue
		}
		xx := PoolTaskView.Get()
		x.Tasks[idx].DeepCopy(xx)
		z.Tasks = append(z.Tasks, xx)
	}
}

func (x *TaskViewMany) Clone() *TaskViewMany {
	z := &TaskViewMany{}
	x.DeepCopy(z)
	return z
}

func (x *TaskViewMany) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *TaskViewMany) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *TaskViewMany) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *TaskViewMany) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *TaskViewMany) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_TaskViewMany, x)
}

const C_Bool int64 = 4122188204

type poolBool struct {
	pool sync.Pool
}

func (p *poolBool) Get() *Bool {
	x, ok := p.pool.Get().(*Bool)
	if !ok {
		x = &Bool{}
	}

	return x
}

func (p *poolBool) Put(x *Bool) {
	if x == nil {
		return
	}

	x.Result = false

	p.pool.Put(x)
}

var PoolBool = poolBool{}

func (x *Bool) DeepCopy(z *Bool) {
	z.Result = x.Result
}

func (x *Bool) Clone() *Bool {
	z := &Bool{}
	x.DeepCopy(z)
	return z
}

func (x *Bool) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Bool) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Bool) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Bool) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *Bool) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Bool, x)
}

const C_TaskManagerCreate int64 = 613086944
const C_TaskManagerGet int64 = 1803980423
const C_TaskManagerDelete int64 = 2437835420
const C_TaskManagerList int64 = 1736909319

func init() {
	registry.RegisterConstructor(2229707971, "CreateRequest")
	registry.RegisterConstructor(3359917651, "GetRequest")
	registry.RegisterConstructor(4088608791, "DeleteRequest")
	registry.RegisterConstructor(307567737, "ListRequest")
	registry.RegisterConstructor(2236647323, "TaskView")
	registry.RegisterConstructor(2061119489, "TaskViewMany")
	registry.RegisterConstructor(4122188204, "Bool")
	registry.RegisterConstructor(613086944, "TaskManagerCreate")
	registry.RegisterConstructor(1803980423, "TaskManagerGet")
	registry.RegisterConstructor(2437835420, "TaskManagerDelete")
	registry.RegisterConstructor(1736909319, "TaskManagerList")
}

var _ = tools.TimeUnix()

type ITaskManager interface {
	Create(ctx *edge.RequestCtx, req *CreateRequest, res *TaskView) *rony.Error
	Get(ctx *edge.RequestCtx, req *GetRequest, res *TaskView) *rony.Error
	Delete(ctx *edge.RequestCtx, req *DeleteRequest, res *Bool) *rony.Error
	List(ctx *edge.RequestCtx, req *ListRequest, res *TaskViewMany) *rony.Error
}

func RegisterTaskManager(h ITaskManager, e *edge.Server, preHandlers ...edge.Handler) {
	w := taskManagerWrapper{
		h: h,
	}
	w.Register(e, func(c int64) []edge.Handler {
		return preHandlers
	})
}

func RegisterTaskManagerWithFunc(h ITaskManager, e *edge.Server, handlerFunc func(c int64) []edge.Handler) {
	w := taskManagerWrapper{
		h: h,
	}
	w.Register(e, handlerFunc)
}

type taskManagerWrapper struct {
	h ITaskManager
}

func (sw *taskManagerWrapper) createWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolCreateRequest.Get()
	defer PoolCreateRequest.Put(req)
	res := PoolTaskView.Get()
	defer PoolTaskView.Put(res)

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.Create(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_TaskView, res)
	}
}
func (sw *taskManagerWrapper) getWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolGetRequest.Get()
	defer PoolGetRequest.Put(req)
	res := PoolTaskView.Get()
	defer PoolTaskView.Put(res)

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.Get(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_TaskView, res)
	}
}
func (sw *taskManagerWrapper) deleteWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolDeleteRequest.Get()
	defer PoolDeleteRequest.Put(req)
	res := PoolBool.Get()
	defer PoolBool.Put(res)

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.Delete(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_Bool, res)
	}
}
func (sw *taskManagerWrapper) listWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolListRequest.Get()
	defer PoolListRequest.Put(req)
	res := PoolTaskViewMany.Get()
	defer PoolTaskViewMany.Put(res)

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.List(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_TaskViewMany, res)
	}
}

func (sw *taskManagerWrapper) Register(e *edge.Server, handlerFunc func(c int64) []edge.Handler) {
	if handlerFunc == nil {
		handlerFunc = func(c int64) []edge.Handler {
			return nil
		}
	}
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_TaskManagerCreate).
			SetHandler(handlerFunc(C_TaskManagerCreate)...).
			Append(sw.createWrapper),
	)
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_TaskManagerGet).
			SetHandler(handlerFunc(C_TaskManagerGet)...).
			Append(sw.getWrapper),
	)
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_TaskManagerDelete).
			SetHandler(handlerFunc(C_TaskManagerDelete)...).
			Append(sw.deleteWrapper),
	)
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_TaskManagerList).
			SetHandler(handlerFunc(C_TaskManagerList)...).
			Append(sw.listWrapper),
	)
}

func TunnelRequestTaskManagerCreate(ctx *edge.RequestCtx, replicaSet uint64, req *CreateRequest, res *TaskView, kvs ...*rony.KeyValue) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_TaskManagerCreate, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_TaskView:
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
func TunnelRequestTaskManagerGet(ctx *edge.RequestCtx, replicaSet uint64, req *GetRequest, res *TaskView, kvs ...*rony.KeyValue) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_TaskManagerGet, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_TaskView:
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
func TunnelRequestTaskManagerDelete(ctx *edge.RequestCtx, replicaSet uint64, req *DeleteRequest, res *Bool, kvs ...*rony.KeyValue) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_TaskManagerDelete, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_Bool:
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
func TunnelRequestTaskManagerList(ctx *edge.RequestCtx, replicaSet uint64, req *ListRequest, res *TaskViewMany, kvs ...*rony.KeyValue) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_TaskManagerList, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_TaskViewMany:
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

type TaskManagerClient struct {
	c edgec.Client
}

func NewTaskManagerClient(ec edgec.Client) *TaskManagerClient {
	return &TaskManagerClient{
		c: ec,
	}
}
func (c *TaskManagerClient) Create(req *CreateRequest, kvs ...*rony.KeyValue) (*TaskView, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_TaskManagerCreate, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_TaskView:
		x := &TaskView{}
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
func (c *TaskManagerClient) Get(req *GetRequest, kvs ...*rony.KeyValue) (*TaskView, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_TaskManagerGet, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_TaskView:
		x := &TaskView{}
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
func (c *TaskManagerClient) Delete(req *DeleteRequest, kvs ...*rony.KeyValue) (*Bool, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_TaskManagerDelete, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_Bool:
		x := &Bool{}
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
func (c *TaskManagerClient) List(req *ListRequest, kvs ...*rony.KeyValue) (*TaskViewMany, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_TaskManagerList, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_TaskViewMany:
		x := &TaskViewMany{}
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

func prepareTaskManagerCommand(cmd *cobra.Command, c edgec.Client) (*TaskManagerClient, error) {
	// Bind current flags to the registered flags in config package
	err := config.BindCmdFlags(cmd)
	if err != nil {
		return nil, err
	}

	return NewTaskManagerClient(c), nil
}

var genTaskManagerCreateCmd = func(h ITaskManagerCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareTaskManagerCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Create(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.StringFlag("title", "", ""),
		config.StringFlag("tODOs", "", ""),
		config.Int64Flag("dueDate", 0, ""),
	)
	return cmd
}
var genTaskManagerGetCmd = func(h ITaskManagerCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareTaskManagerCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Get(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.Int64Flag("taskID", 0, ""),
	)
	return cmd
}
var genTaskManagerDeleteCmd = func(h ITaskManagerCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareTaskManagerCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Delete(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.Int64Flag("taskID", 0, ""),
	)
	return cmd
}
var genTaskManagerListCmd = func(h ITaskManagerCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareTaskManagerCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.List(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.Int32Flag("offset", 0, ""),
		config.Int32Flag("limit", 0, ""),
	)
	return cmd
}

type ITaskManagerCli interface {
	Create(cli *TaskManagerClient, cmd *cobra.Command, args []string) error
	Get(cli *TaskManagerClient, cmd *cobra.Command, args []string) error
	Delete(cli *TaskManagerClient, cmd *cobra.Command, args []string) error
	List(cli *TaskManagerClient, cmd *cobra.Command, args []string) error
}

func RegisterTaskManagerCli(h ITaskManagerCli, c edgec.Client, rootCmd *cobra.Command) {
	subCommand := &cobra.Command{
		Use: "TaskManager",
	}
	subCommand.AddCommand(
		genTaskManagerCreateCmd(h, c),
		genTaskManagerGetCmd(h, c),
		genTaskManagerDeleteCmd(h, c),
		genTaskManagerListCmd(h, c),
	)

	rootCmd.AddCommand(subCommand)
}

var _ = bytes.MinRead

func init() {

}
