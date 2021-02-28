package tcpGateway

import (
	"fmt"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/internal/gateway/tcp/radix"
	"github.com/ronaksoft/rony/tools"
	"github.com/savsgio/gotils/bytes"
	"github.com/savsgio/gotils/strconv"
	gstrings "github.com/savsgio/gotils/strings"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"strings"
)

/*
   Creation Time: 2021 - Feb - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// MethodWild wild HTTP method
const MethodWild = "*"

var (
	questionMark = byte('?')

	// MatchedRoutePathParam is the param name under which the path of the matched
	// route is stored, if Router.SaveMatchedRoutePath is set.
	MatchedRoutePathParam = fmt.Sprintf("__matchedRoutePath::%s__", bytes.Rand(make([]byte, 15)))
)

type Mux struct {
	trees              []*radix.Tree
	customMethodsIndex map[string]int
	registeredPaths    map[string][]string

	// If enabled, adds the matched route path onto the ctx.UserValue context
	// before invoking the handler.
	// The matched route path is only added to handlers of routes that were
	// registered when this option was enabled.
	SaveMatchedRoutePath bool

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 308 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 308 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// An optional gateway.ProxyHandler that is called on automatic OPTIONS requests.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS gateway.ProxyHandler

	// Configurable gateway.ProxyHandler which is called when no matching route is
	// found. If it is not set, default NotFound is used.
	NotFound gateway.ProxyHandler

	// Configurable ProxyHandler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, ctx.Error with fasthttp.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed gateway.ProxyHandler

	// Cached value of global (*) allowed methods
	globalAllowed string
}

func NewMux(addrs string) *Mux {
	p := &Mux{}
	return p
}

func (r *Mux) saveMatchedRoutePath(path string, handler gateway.ProxyHandler) gateway.ProxyHandler {
	return func(ctx *gateway.RequestCtx, f gateway.ProxyFunc) {
		ctx.SetUserValue(MatchedRoutePathParam, path)
		handler(ctx, f)
	}
}

func (r *Mux) methodIndexOf(method string) int {
	switch method {
	case fasthttp.MethodGet:
		return 0
	case fasthttp.MethodHead:
		return 1
	case fasthttp.MethodPost:
		return 2
	case fasthttp.MethodPut:
		return 3
	case fasthttp.MethodPatch:
		return 4
	case fasthttp.MethodDelete:
		return 5
	case fasthttp.MethodConnect:
		return 6
	case fasthttp.MethodOptions:
		return 7
	case fasthttp.MethodTrace:
		return 8
	}

	if i, ok := r.customMethodsIndex[method]; ok {
		return i
	}

	return -1
}

// List returns all registered routes grouped by method
func (r *Mux) List() map[string][]string {
	return r.registeredPaths
}

// GET is a shortcut for router.Handle(fasthttp.MethodGet, path, handler)
func (r *Mux) GET(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodGet, path, handler)
}

// HEAD is a shortcut for router.Handle(fasthttp.MethodHead, path, handler)
func (r *Mux) HEAD(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodHead, path, handler)
}

// POST is a shortcut for router.Handle(fasthttp.MethodPost, path, handler)
func (r *Mux) POST(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodPost, path, handler)
}

// PUT is a shortcut for router.Handle(fasthttp.MethodPut, path, handler)
func (r *Mux) PUT(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodPut, path, handler)
}

// PATCH is a shortcut for router.Handle(fasthttp.MethodPatch, path, handler)
func (r *Mux) PATCH(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodPatch, path, handler)
}

// DELETE is a shortcut for router.Handle(fasthttp.MethodDelete, path, handler)
func (r *Mux) DELETE(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodDelete, path, handler)
}

// CONNECT is a shortcut for router.Handle(fasthttp.MethodConnect, path, handler)
func (r *Mux) CONNECT(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodConnect, path, handler)
}

// OPTIONS is a shortcut for router.Handle(fasthttp.MethodOptions, path, handler)
func (r *Mux) OPTIONS(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodOptions, path, handler)
}

// TRACE is a shortcut for router.Handle(fasthttp.MethodTrace, path, handler)
func (r *Mux) TRACE(path string, handler gateway.ProxyHandler) {
	r.Handle(fasthttp.MethodTrace, path, handler)
}

// ANY is a shortcut for router.Handle(router.MethodWild, path, handler)
//
// WARNING: Use only for routes where the request method is not important
func (r *Mux) ANY(path string, handler gateway.ProxyHandler) {
	r.Handle(MethodWild, path, handler)
}

// Handle registers a new request handler with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Mux) Handle(method, path string, handler gateway.ProxyHandler) {
	switch {
	case len(method) == 0:
		panic("method must not be empty")
	case len(path) < 1 || path[0] != '/':
		panic("path must begin with '/' in path '" + path + "'")
	case path == "/":
		panic("you must not alter the root handler")
	case handler == nil:
		panic("handler must not be nil")
	}

	r.registeredPaths[method] = append(r.registeredPaths[method], path)

	methodIndex := r.methodIndexOf(method)
	if methodIndex == -1 {
		tree := radix.New()
		tree.Mutable = false

		r.trees = append(r.trees, tree)
		methodIndex = len(r.trees) - 1
		r.customMethodsIndex[method] = methodIndex
	}

	tree := r.trees[methodIndex]
	if tree == nil {
		tree = radix.New()
		tree.Mutable = false

		r.trees[methodIndex] = tree
		r.globalAllowed = r.allowed("*", "")
	}

	if r.SaveMatchedRoutePath {
		handler = r.saveMatchedRoutePath(path, handler)
	}

	optionalPaths := getOptionalPaths(path)

	// if not has optional paths, adds the original
	if len(optionalPaths) == 0 {
		tree.Add(path, handler)
	} else {
		for _, p := range optionalPaths {
			tree.Add(p, handler)
		}
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handler function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Mux) Lookup(ctx *fasthttp.RequestCtx) (gateway.ProxyHandler, bool) {
	path := strconv.B2S(ctx.Request.URI().PathOriginal())
	methodIndex := r.methodIndexOf(tools.ByteToStr(ctx.Request.Header.Method()))
	if methodIndex == -1 {
		return nil, false
	}

	if tree := r.trees[methodIndex]; tree != nil {
		handler, tsr := tree.Get(path, ctx)
		if handler != nil || tsr {
			return handler, tsr
		}
	}

	if tree := r.trees[r.methodIndexOf(MethodWild)]; tree != nil {
		return tree.Get(path, ctx)
	}

	return nil, false
}

func (r *Mux) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" || path == "/*" { // server-wide{ // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.registeredPaths {
				if method == fasthttp.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else { // specific path
		for method := range r.registeredPaths {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == fasthttp.MethodOptions {
				continue
			}

			handle, _ := r.trees[r.methodIndexOf(method)].Get(path, nil)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, fasthttp.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}
	return
}

func (r *Mux) tryRedirect(ctx *fasthttp.RequestCtx, tree *radix.Tree, tsr bool, method, path string) bool {
	// Moved Permanently, request with GET method
	code := fasthttp.StatusMovedPermanently
	if method != fasthttp.MethodGet {
		// Permanent Redirect, request with same method
		code = fasthttp.StatusPermanentRedirect
	}

	if tsr && r.RedirectTrailingSlash {
		uri := bytebufferpool.Get()

		if len(path) > 1 && path[len(path)-1] == '/' {
			uri.SetString(path[:len(path)-1])
		} else {
			uri.SetString(path)
			uri.WriteString("/")
		}

		queryBuf := ctx.URI().QueryString()
		if len(queryBuf) > 0 {
			uri.WriteByte(questionMark)
			uri.Write(queryBuf)
		}

		ctx.Redirect(uri.String(), code)

		bytebufferpool.Put(uri)

		return true
	}

	// Try to fix the request path
	if r.RedirectFixedPath {
		path := strconv.B2S(ctx.Request.URI().Path())

		uri := bytebufferpool.Get()
		found := tree.FindCaseInsensitivePath(
			cleanPath(path),
			r.RedirectTrailingSlash,
			uri,
		)

		if found {
			queryBuf := ctx.URI().QueryString()
			if len(queryBuf) > 0 {
				uri.WriteByte(questionMark)
				uri.Write(queryBuf)
			}

			ctx.RedirectBytes(uri.Bytes(), code)

			bytebufferpool.Put(uri)

			return true
		}
	}

	return false
}

// Handler makes the router implement the http.Handler interface.
func (r *Mux) Handler(ctx *gateway.RequestCtx, f gateway.ProxyFunc) {
	path := strconv.B2S(ctx.Request.URI().PathOriginal())
	method := strconv.B2S(ctx.Request.Header.Method())
	methodIndex := r.methodIndexOf(method)

	if methodIndex > -1 {
		if tree := r.trees[methodIndex]; tree != nil {
			if handler, tsr := tree.Get(path, ctx); handler != nil {
				handler(ctx, f)
				return
			} else if method != fasthttp.MethodConnect && path != "/" {
				if ok := r.tryRedirect(ctx, tree, tsr, method, path); ok {
					return
				}
			}
		}
	}

	// Try to search in the wild method tree
	if tree := r.trees[r.methodIndexOf(MethodWild)]; tree != nil {
		if handler, tsr := tree.Get(path, ctx); handler != nil {
			handler(ctx, f)
			return
		} else if method != fasthttp.MethodConnect && path != "/" {
			if ok := r.tryRedirect(ctx, tree, tsr, method, path); ok {
				return
			}
		}
	}

	if r.HandleOPTIONS && method == fasthttp.MethodOptions {
		// Handle OPTIONS requests

		if allow := r.allowed(path, fasthttp.MethodOptions); allow != "" {
			ctx.Response.Header.Set("Allow", allow)
			if r.GlobalOPTIONS != nil {
				r.GlobalOPTIONS(ctx, f)
			}
			return
		}
	} else if r.HandleMethodNotAllowed {
		// Handle 405

		if allow := r.allowed(path, method); allow != "" {
			ctx.Response.Header.Set("Allow", allow)
			if r.MethodNotAllowed != nil {
				r.MethodNotAllowed(ctx, f)
			} else {
				ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
				ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
			}
			return
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound(ctx, f)
	} else {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	}
}

// cleanPath removes the '.' if it is the last character of the route
func cleanPath(path string) string {
	lenPath := len(path)

	if path[lenPath-1] == '.' {
		path = path[:lenPath-1]
	}

	return path
}

// getOptionalPaths returns all possible paths when the original path
// has optional arguments
func getOptionalPaths(path string) []string {
	paths := make([]string, 0)

	start := 0
walk:
	for {
		if start >= len(path) {
			return paths
		}

		c := path[start]
		start++

		if c != '{' {
			continue
		}

		newPath := ""
		hasRegex := false
		questionMarkIndex := -1

		brackets := 0

		for end, c := range []byte(path[start:]) {
			switch c {
			case '{':
				brackets++

			case '}':
				if brackets > 0 {
					brackets--
					continue
				} else if questionMarkIndex == -1 {
					continue walk
				}

				end++
				newPath += path[questionMarkIndex+1 : start+end]

				path = path[:questionMarkIndex] + path[questionMarkIndex+1:] // remove '?'
				paths = append(paths, newPath)
				start += end - 1

				continue walk

			case ':':
				hasRegex = true

			case '?':
				if hasRegex {
					continue
				}

				questionMarkIndex = start + end
				newPath += path[:questionMarkIndex]

				if len(path[:start-2]) == 0 {
					// include the root slash because the param is in the first segment
					paths = append(paths, "/")

				} else if !gstrings.Include(paths, path[:start-2]) {
					// include the path without the wildcard
					// -2 due to remove the '/' and '{'
					paths = append(paths, path[:start-2])
				}
			}
		}
	}
}
