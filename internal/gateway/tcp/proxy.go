package tcpGateway

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/internal/trie"
)

/*
   Creation Time: 2021 - Mar - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// HttpProxy
type HttpProxy struct {
	handler gateway.MessageHandler
	routes  map[string]*trie.Trie
}

func (hp *HttpProxy) Set(method, path string, p gateway.ProxyHandle) {
	if hp.routes == nil {
		hp.routes = make(map[string]*trie.Trie)
	}
	if _, ok := hp.routes[method]; !ok {
		hp.routes[method] = trie.NewTrie()
	}
	hp.routes[method].Insert(path, trie.WithTag(method), trie.WithProxy(p))
}

func (hp *HttpProxy) handle(conn *httpConn, ctx *gateway.RequestCtx) {
	hp.handler(conn, int64(ctx.ConnID()), conn.proxy.OnRequest(conn, ctx))
}

func (hp *HttpProxy) search(method, path string, conn *httpConn) gateway.ProxyHandle {
	r := hp.routes[method]
	if r == nil {
		return nil
	}
	n := r.Search(path, conn)
	if n == nil {
		return nil
	}
	return n.Proxy
}

// simpleProxy
type simpleProxy struct {
	onRequestFunc  func(conn rony.Conn, reqCtx *gateway.RequestCtx) []byte
	onResponseFunc func(data []byte) ([]byte, map[string]string)
}

func (s *simpleProxy) OnRequest(conn rony.Conn, ctx *gateway.RequestCtx) []byte {
	return s.onRequestFunc(conn, ctx)
}

func (s *simpleProxy) OnResponse(data []byte) ([]byte, map[string]string) {
	return s.onResponseFunc(data)
}

func CreateHandle(
	onRequest func(conn rony.Conn, reqCtx *gateway.RequestCtx) []byte,
	onResponse func(data []byte) ([]byte, map[string]string),
) gateway.ProxyHandle {
	return &simpleProxy{
		onRequestFunc:  onRequest,
		onResponseFunc: onResponse,
	}
}
