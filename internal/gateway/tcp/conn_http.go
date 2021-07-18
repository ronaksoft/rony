package tcpGateway

import (
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/tools"
	"github.com/valyala/fasthttp"
	"mime/multipart"
)

/*
   Creation Time: 2019 - Nov - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	strLocation = []byte(fasthttp.HeaderLocation)
)

// httpConn
type httpConn struct {
	gateway    *Gateway
	ctx        *fasthttp.RequestCtx
	clientIP   []byte
	clientType []byte
	mtx        tools.SpinLock
	kv         map[string]interface{}
}

func (c *httpConn) Get(key string) interface{} {
	c.mtx.Lock()
	v := c.kv[key]
	c.mtx.Unlock()

	return v
}

func (c *httpConn) Set(key string, val interface{}) {
	c.mtx.Lock()
	c.kv[key] = val
	c.mtx.Unlock()
}

func (c *httpConn) ConnID() uint64 {
	return c.ctx.ConnID()
}

func (c *httpConn) ClientIP() string {
	return string(c.clientIP)
}

func (c *httpConn) SetClientIP(ip []byte) {
	c.clientIP = append(c.clientIP[:0], ip...)
}

func (c *httpConn) SetClientType(ct []byte) {
	c.clientType = append(c.clientType[:0], ct...)
}

func (c *httpConn) WriteStatus(status int) {
	c.ctx.SetStatusCode(status)
}

func (c *httpConn) WriteBinary(streamID int64, data []byte) (err error) {
	_, err = c.ctx.Write(data)
	metrics.IncCounter(metrics.CntGatewayOutgoingHttpMessage)

	return err
}

func (c *httpConn) WriteHeader(key, value string) {
	c.ctx.Response.Header.Add(key, value)
}

func (c *httpConn) MultiPart() (*multipart.Form, error) {
	return c.ctx.MultipartForm()
}

func (c *httpConn) Method() string {
	return tools.B2S(c.ctx.Method())
}

func (c *httpConn) Path() string {
	return tools.B2S(c.ctx.Request.URI().PathOriginal())
}

func (c *httpConn) Body() []byte {
	return c.ctx.PostBody()
}

func (c *httpConn) Persistent() bool {
	return false
}

func (c *httpConn) Redirect(statusCode int, newHostPort string) {
	u := fasthttp.AcquireURI()
	c.ctx.URI().CopyTo(u)
	u.SetHost(newHostPort)
	c.ctx.Response.Header.SetCanonical(strLocation, u.FullURI())
	c.ctx.Response.SetStatusCode(statusCode)
	fasthttp.ReleaseURI(u)
}
