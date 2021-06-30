package gateway

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/pools"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2019 - Aug - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Gateway defines the gateway interface where clients could connect
// and communicate with the edge server
type Gateway interface {
	Start()
	Run()
	Shutdown()
	GetConn(connID uint64) rony.Conn
	Addr() []string
	Protocol() Protocol
}

type (
	RequestCtx     = fasthttp.RequestCtx
	ConnectHandler = func(c rony.Conn, kvs ...*rony.KeyValue)
	MessageHandler = func(c rony.Conn, streamID int64, data []byte, bypass bool)
	CloseHandler   = func(c rony.Conn)
)

// ProxyHandle defines the interface for proxy handlers. This is used in 'rest' package to support
// writing restfull wrappers for the RPC handlers.
type ProxyHandle interface {
	OnRequest(conn rony.Conn, ctx *RequestCtx) []byte
	OnResponse(data []byte, bodyWriter BodyWriter, hdrWriter *HeaderWriter)
	Release()
}

// ProxyFactory is the factory which generates ProxyHandle.
type ProxyFactory interface {
	Get() ProxyHandle
}

// BodyWriter is implemented by bodyWriter internally. This interface is used by the developer to write desired
// output to the REST http response. WriteProto is a helper function for marshalling and then writing the proto message
// to the body of http response.
type BodyWriter interface {
	Write(data []byte)
	WriteProto(p proto.Message) error
}

type bodyWriter struct {
	buf *pools.ByteBuffer
}

func NewBodyWriter() *bodyWriter {
	bw := &bodyWriter{}
	return bw
}

func (bw *bodyWriter) WriteProto(m proto.Message) error {
	if bw.buf == nil {
		bw.buf = pools.Buffer.FromProto(m)
		return nil
	}
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	bw.buf.AppendFrom(b)
	return nil
}

func (bw *bodyWriter) Write(data []byte) {
	if bw.buf == nil {
		bw.buf = pools.Buffer.FromBytes(data)
		return
	}
	bw.buf.AppendFrom(data)
}

func (bw *bodyWriter) Bytes() *[]byte {
	if bw.buf == nil {
		return nil
	}
	return bw.buf.Bytes()
}

func (bw *bodyWriter) Release() {
	if bw.buf == nil {
		return
	}
	pools.Buffer.Put(bw.buf)
}

// HeaderWriter is the object responsible for writing response header.
type HeaderWriter map[string]string

func NewHeaderWriter() *HeaderWriter {
	return &HeaderWriter{}
}

func (hw HeaderWriter) Set(key, value string) {
	hw[key] = value
}

func (hw HeaderWriter) Release() {

}
