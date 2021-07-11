package service

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/tools"
	"time"
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

// Sample implements ISample interface.
type Sample struct {
	ServerID string
}

func (h *Sample) Set(ctx *edge.RequestCtx, req *SetRequest, res *SetResponse) {

}

func (h *Sample) Get(ctx *edge.RequestCtx, req *GetRequest, res *GetResponse) {

}

func (h *Sample) EchoInternal(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) {
	panic("implement me")
}

func (h *Sample) Echo(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID
}

func (h *Sample) EchoDelay(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID
	time.Sleep(time.Second * 1)
}

func (h *Sample) EchoTunnel(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID

}

var EchoRest = edge.NewRestProxy(
	func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
		req := &EchoRequest{}
		err := req.UnmarshalJSON(conn.Body())
		if err != nil {
			return err
		}
		ctx.FillEnvelope(conn.ConnID(), C_SampleEcho, req)
		return nil
	},
	func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
		envelope := ctx.BufferPop()
		if envelope == nil {
			return errors.ErrInternalServer
		}
		switch envelope.Constructor {
		case C_EchoResponse:
			x := &EchoResponse{}
			_ = x.Unmarshal(envelope.Message)
			b, err := x.MarshalJSON()
			if err != nil {
				return err
			}
			return conn.WriteBinary(ctx.StreamID(), b)
		case rony.C_Error:
			x := &rony.Error{}
			_ = x.Unmarshal(envelope.Message)

		default:
			return errors.ErrUnexpectedResponse

		}
		return errors.ErrInternalServer
	},
)

var EchoRestBinding = edge.NewRestProxy(
	func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
		req := &EchoRequest{}
		req.Int = tools.StrToInt64(tools.GetString(conn.Get("value"), "0"))
		req.Timestamp = tools.StrToInt64(tools.GetString(conn.Get("ts"), "0"))
		ctx.FillEnvelope(conn.ConnID(), C_SampleEcho, req)
		return nil
	},
	func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
		envelope := ctx.BufferPop()
		if envelope == nil {
			return errors.ErrInternalServer
		}
		switch envelope.Constructor {
		case C_EchoResponse:
			x := &EchoResponse{}
			_ = x.Unmarshal(envelope.Message)
			b, err := x.MarshalJSON()
			if err != nil {
				return err
			}
			return conn.WriteBinary(ctx.StreamID(), b)
		case rony.C_Error:
			x := &rony.Error{}
			_ = x.Unmarshal(envelope.Message)

		default:
			return errors.ErrUnexpectedResponse

		}
		return errors.ErrInternalServer
	},
)
