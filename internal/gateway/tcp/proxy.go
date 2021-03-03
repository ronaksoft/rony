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

type HttpProxy struct {
	handler gateway.MessageHandler
	routes  map[string]*trie.Trie
}

func (hp *HttpProxy) CreateHandle(
	onRequest func(conn rony.Conn, data []byte) []byte,
	onResponse func(data []byte) ([]byte, map[string]string),
) gateway.ProxyHandle {
	return &simpleProxy{
		onRequestFunc:  onRequest,
		onResponseFunc: onResponse,
	}
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
	hp.handler(conn, int64(ctx.ConnID()), conn.proxy.OnRequest(conn, ctx.PostBody()))
}

func (hp *HttpProxy) search(method, path string, conn *httpConn) gateway.ProxyHandle {
	n := hp.routes[method].Search(path, conn)
	if n == nil {
		return nil
	}

	return n.Proxy
}

type simpleProxy struct {
	onRequestFunc  func(conn rony.Conn, date []byte) []byte
	onResponseFunc func(data []byte) ([]byte, map[string]string)
}

func (s *simpleProxy) OnRequest(conn rony.Conn, data []byte) []byte {
	return s.onRequestFunc(conn, data)
}

func (s *simpleProxy) OnResponse(data []byte) ([]byte, map[string]string) {
	return s.onResponseFunc(data)
}
