package router

import (
	"github.com/ronaksoft/rony/internal/gateway"
)

// Group returns a new group.
// Path auto-correction, including trailing slashes, is enabled by default.
func (g *Group) Group(path string) *Group {
	return g.router.Group(g.prefix + path)
}

// GET is a shortcut for group.Handle(fasthttp.MethodGet, path, handler)
func (g *Group) GET(path string, handler gateway.ProxyHandler) {
	g.router.GET(g.prefix+path, handler)
}

// HEAD is a shortcut for group.Handle(fasthttp.MethodHead, path, handler)
func (g *Group) HEAD(path string, handler gateway.ProxyHandler) {
	g.router.HEAD(g.prefix+path, handler)
}

// POST is a shortcut for group.Handle(fasthttp.MethodPost, path, handler)
func (g *Group) POST(path string, handler gateway.ProxyHandler) {
	g.router.POST(g.prefix+path, handler)
}

// PUT is a shortcut for group.Handle(fasthttp.MethodPut, path, handler)
func (g *Group) PUT(path string, handler gateway.ProxyHandler) {
	g.router.PUT(g.prefix+path, handler)
}

// PATCH is a shortcut for group.Handle(fasthttp.MethodPatch, path, handler)
func (g *Group) PATCH(path string, handler gateway.ProxyHandler) {
	g.router.PATCH(g.prefix+path, handler)
}

// DELETE is a shortcut for group.Handle(fasthttp.MethodDelete, path, handler)
func (g *Group) DELETE(path string, handler gateway.ProxyHandler) {
	g.router.DELETE(g.prefix+path, handler)
}

// OPTIONS is a shortcut for group.Handle(fasthttp.MethodOptions, path, handler)
func (g *Group) CONNECT(path string, handler gateway.ProxyHandler) {
	g.router.CONNECT(g.prefix+path, handler)
}

// OPTIONS is a shortcut for group.Handle(fasthttp.MethodOptions, path, handler)
func (g *Group) OPTIONS(path string, handler gateway.ProxyHandler) {
	g.router.OPTIONS(g.prefix+path, handler)
}

// OPTIONS is a shortcut for group.Handle(fasthttp.MethodOptions, path, handler)
func (g *Group) TRACE(path string, handler gateway.ProxyHandler) {
	g.router.TRACE(g.prefix+path, handler)
}

// ANY is a shortcut for group.Handle(router.MethodWild, path, handler)
//
// WARNING: Use only for routes where the request method is not important
func (g *Group) ANY(path string, handler gateway.ProxyHandler) {
	g.router.ANY(g.prefix+path, handler)
}

// Handle registers a new request handler with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (g *Group) Handle(method, path string, handler gateway.ProxyHandler) {
	g.router.Handle(method, g.prefix+path, handler)
}
