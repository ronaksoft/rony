package service

import (
	"time"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2020 - Mar - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. -I=../../../.. --go_out=. service.proto
//go:generate protoc -I=. -I=../../../.. --gorony_out=. service.proto
//go:generate protoc -I=. -I=../../../.. --gorony_out=rony_opt=open_api:. service.proto

// Sample implements ISample interface.
type Sample struct {
	ServerID string
}

func (h *Sample) Set(ctx *edge.RequestCtx, req *SetRequest, res *SetResponse) *rony.Error {
	return nil
}

func (h *Sample) Get(ctx *edge.RequestCtx, req *GetRequest, res *GetResponse) *rony.Error {
	return nil
}

func (h *Sample) EchoInternal(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) *rony.Error {
	panic("implement me")
}

func (h *Sample) Echo(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) *rony.Error {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID

	return nil
}

func (h *Sample) EchoDelay(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) *rony.Error {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID
	time.Sleep(time.Second * 1)

	return nil
}

func (h *Sample) EchoTunnel(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) *rony.Error {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID

	return nil
}

var EchoRest = edge.NewRestProxy(
	func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
		req := &EchoRequest{}
		err := req.UnmarshalJSON(conn.Body())
		if err != nil {
			return err
		}
		ctx.Fill(conn.ConnID(), C_SampleEcho, req)

		return nil
	},
	func(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
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
			case rony.C_Error:
				x := &rony.Error{}
				_ = x.Unmarshal(envelope.Message)
				err = x
			default:
				err = errors.ErrUnexpectedResponse
			}
		}) {
			err = errors.ErrInternalServer
		}

		return
	},
)

var EchoRestBinding = edge.NewRestProxy(
	func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
		req := &EchoRequest{}
		req.Int = tools.StrToInt64(tools.GetString(conn.Get("value"), "0"))
		req.Timestamp = tools.StrToInt64(tools.GetString(conn.Get("ts"), "0"))
		ctx.Fill(conn.ConnID(), C_SampleEcho, req)

		return nil
	},
	func(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
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
			}
			err = errors.ErrUnexpectedResponse
		}) {
			err = errors.ErrInternalServer
		}

		return
	},
)
