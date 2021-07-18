package rony

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/internal/metrics"
	"mime/multipart"
)

/*
   Creation Time: 2021 - Jan - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Conn defines the Connection interface
type Conn interface {
	ConnID() uint64
	ClientIP() string
	WriteBinary(streamID int64, data []byte) error
	// Persistent returns FALSE if this connection will be closed when edge.DispatchCtx has been done. i.e. HTTP connections
	// It returns TRUE if this connection still alive when edge.DispatchCtx has been done. i.e. WebSocket connections
	Persistent() bool
	Get(key string) interface{}
	Set(key string, val interface{})
}

// RestConn is same as Conn but it supports REST format apis.
type RestConn interface {
	Conn
	WriteStatus(status int)
	WriteHeader(key, value string)
	MultiPart() (*multipart.Form, error)
	Method() string
	Path() string
	Body() []byte
	Redirect(statusCode int, newHostPort string)
}

type LogLevel = log.Level

// SetLogLevel is used for debugging purpose
func SetLogLevel(l LogLevel) {
	log.SetLevel(l)
}

func RegisterPrometheus(registerer prometheus.Registerer) {
	metrics.Register(registerer)
}
