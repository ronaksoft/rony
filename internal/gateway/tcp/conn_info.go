package tcpGateway

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/tools"
	"sync"
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
	metaPool       sync.Pool
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
	kvs        []*rony.KeyValue
	clientIP   []byte
	clientType []byte
}

func newConnInfo() *connInfo {
	return &connInfo{
		kvs: make([]*rony.KeyValue, 0, 8),
	}
}

func (m *connInfo) AppendKV(key, value string) {
	kv := rony.PoolKeyValue.Get()
	kv.Key = key
	kv.Value = value
	m.kvs = append(m.kvs, kv)
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

func acquireConnInfo(reqCtx *gateway.RequestCtx) *connInfo {
	mt, ok := metaPool.Get().(*connInfo)
	if !ok {
		mt = newConnInfo()
	}
	reqCtx.Request.Header.VisitAll(func(key, value []byte) {
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
	})
	mt.SetClientIP(tools.S2B(reqCtx.RemoteIP().To4().String()), true)


	return mt
}

func releaseConnInfo(m *connInfo) {
	for _, kv := range m.kvs {
		rony.PoolKeyValue.Put(kv)
	}
	m.clientIP = m.clientIP[:0]
	m.clientType = m.clientType[:0]
	metaPool.Put(m)
}
