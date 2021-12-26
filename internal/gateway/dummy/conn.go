package dummyGateway

import (
	"mime/multipart"
	"sync"

	"github.com/ronaksoft/rony/errors"
)

/*
   Creation Time: 2020 - Oct - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Conn struct {
	id         uint64
	clientIP   string
	persistent bool
	mtx        sync.Mutex
	kv         map[string]interface{}
	onMessage  func(connID uint64, streamID int64, data []byte, hdr map[string]string)

	// extra data for REST
	status  int
	httpHdr map[string]string
	method  string
	path    string
	body    []byte
}

func NewConn(id uint64) *Conn {
	return &Conn{
		id:         id,
		kv:         map[string]interface{}{},
		httpHdr:    map[string]string{},
		persistent: false,
	}
}

func (c *Conn) WithHandler(onMessage func(connID uint64, streamID int64, data []byte, hdr map[string]string)) *Conn {
	c.onMessage = onMessage

	return c
}

func (c *Conn) ReadHeader(key string) string {
	return c.httpHdr[key]
}

func (c *Conn) Redirect(statusCode int, newHostPort string) {
	// TODO:: implement it
}

func (c *Conn) RedirectURL(statusCode int, url string) {
	// TODO:: implement it
}

func (c *Conn) WriteStatus(status int) {
	c.status = status
}

func (c *Conn) WriteHeader(key, value string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.httpHdr[key] = value
}

func (c *Conn) MultiPart() (*multipart.Form, error) {
	return nil, errors.New(errors.NotImplemented, "Multipart")
}

func (c *Conn) Method() string {
	return c.method
}

func (c *Conn) Path() string {
	return c.path
}

func (c *Conn) Body() []byte {
	return c.body
}

func (c *Conn) Get(key string) interface{} {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.kv[key]
}

func (c *Conn) Set(key string, val interface{}) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.kv[key] = val
}

func (c *Conn) Walk(f func(k string, v interface{}) bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for k, v := range c.kv {
		if !f(k, v) {
			return
		}
	}
}

func (c *Conn) ConnID() uint64 {
	return c.id
}

func (c *Conn) ClientIP() string {
	return c.clientIP
}

func (c *Conn) WriteBinary(streamID int64, data []byte) error {
	c.onMessage(c.id, streamID, data, c.httpHdr)

	return nil
}

func (c *Conn) Persistent() bool {
	return c.persistent
}

func (c *Conn) SetPersistent(b bool) *Conn {
	c.persistent = b

	return c
}
