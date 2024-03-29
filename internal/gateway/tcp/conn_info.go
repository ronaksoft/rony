package tcpGateway

import (
	"bytes"
	"sync"

	"github.com/ronaksoft/rony/tools"
	"github.com/valyala/fasthttp"
)

/*
   Creation Time: 2021 - Mar - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	ronyUpgrade    = []byte("Upgrade")
	connInfoPool   sync.Pool
	ignoredHeaders = map[string]bool{
		"Host":                     true,
		"Upgrade":                  true,
		"Connection":               true,
		"Sec-Websocket-Version":    true,
		"Sec-Websocket-Protocol":   true,
		"Sec-Websocket-Extensions": true,
		"Sec-Websocket-Key":        true,
		"Sec-Websocket-Accept":     true,
	}
)

type connInfo struct {
	kvs        map[string]string
	clientIP   []byte
	clientType []byte
	upgrade    bool
}

func newConnInfo() *connInfo {
	return &connInfo{
		kvs: make(map[string]string, 8),
	}
}

func (m *connInfo) AppendKV(key, value string) {
	m.kvs[key] = value
}

func (m *connInfo) SetClientIP(clientIP []byte, onlyNew bool) {
	if onlyNew && len(m.clientIP) > 0 {
		return
	}
	m.clientIP = append(m.clientIP[:0], clientIP...)
}

func (m *connInfo) SetClientType(clientType []byte) {
	m.clientType = append(m.clientType[:0], clientType...)
}

func (m *connInfo) Upgrade() bool {
	return m.upgrade
}

func acquireConnInfo(reqCtx *fasthttp.RequestCtx) *connInfo {
	mt, ok := connInfoPool.Get().(*connInfo)
	if !ok {
		mt = newConnInfo()
	}
	reqCtx.Request.Header.VisitAll(
		func(key, value []byte) {
			switch tools.B2S(key) {
			case "Cf-Connecting-Ip":
				mt.SetClientIP(value, false)
			case "X-Forwarded-For", "X-Real-Ip", "Forwarded", "X-Forwarded":
				mt.SetClientIP(value, true)
			case "X-Client-Type":
				mt.SetClientType(value)
			default:
				if !ignoredHeaders[tools.B2S(key)] {
					mt.AppendKV(string(key), string(value))
				}
			}
		},
	)

	reqCtx.PostArgs().VisitAll(
		func(key, value []byte) {
			mt.AppendKV(string(key), string(value))
		},
	)

	reqCtx.QueryArgs().VisitAll(
		func(key, value []byte) {
			mt.AppendKV(string(key), string(value))
		},
	)

	mt.SetClientIP(tools.S2B(reqCtx.RemoteIP().To4().String()), true)

	if reqCtx.Request.Header.ConnectionUpgrade() {
		mt.upgrade = true
	} else if bytes.Equal(reqCtx.Request.Header.Peek("X-Rony-Upgrade"), ronyUpgrade) {
		reqCtx.Request.Header.SetBytesKV(tools.S2B(fasthttp.HeaderConnection), ronyUpgrade)
		mt.upgrade = true
	}

	return mt
}

func releaseConnInfo(m *connInfo) {
	for k := range m.kvs {
		delete(m.kvs, k)
	}
	m.upgrade = false
	m.clientIP = m.clientIP[:0]
	m.clientType = m.clientType[:0]
	connInfoPool.Put(m)
}
