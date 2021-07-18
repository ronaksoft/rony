package dummyGateway

import (
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/tools"
	"mime/multipart"
	"sync"
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

func NewConn(onMessage func(connID uint64, streamID int64, data []byte, hdr map[string]string)) *Conn {
	return &Conn{
		id:         tools.RandomUint64(0),
		kv:         make(map[string]interface{}),
		persistent: false,
		onMessage:  onMessage,
	}
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

func (c *Conn) SetPersistent(b bool) {
	c.persistent = b
}
