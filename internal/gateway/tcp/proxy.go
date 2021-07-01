package tcpGateway

import (
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

// HttpProxy help to provide RESTFull wrappers around RPC handlers.
type HttpProxy struct {
	handler gateway.MessageHandler
	routes  map[string]*trie.Trie
}

func (hp *HttpProxy) Set(method, path string, f gateway.ProxyFactory) {
	if hp.routes == nil {
		hp.routes = make(map[string]*trie.Trie)
	}
	if _, ok := hp.routes[method]; !ok {
		hp.routes[method] = trie.NewTrie()
	}
	hp.routes[method].Insert(path, trie.WithTag(method), trie.WithProxyFactory(f))
}

func (hp *HttpProxy) Handle(conn *httpConn, ctx *gateway.RequestCtx) {
	bw := gateway.NewBodyWriter()
	conn.proxy.OnRequest(conn, ctx, bw)
	hp.handler(conn, int64(ctx.ConnID()), *bw.Bytes())
	bw.Release()
}

func (hp *HttpProxy) Search(method, path string, conn *httpConn) gateway.ProxyFactory {
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
