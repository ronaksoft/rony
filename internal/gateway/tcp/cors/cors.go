package cors

import (
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

/*
   Creation Time: 2021 - Aug - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Config struct {
	AllowedHeaders []string
	AllowedMethods []string
	AllowedOrigins []string
}

type CORS struct {
	headers string
	methods string
	origins string
}

func New(config Config) *CORS {
	c := &CORS{}
	if len(config.AllowedOrigins) == 0 {
		c.origins = "*"
	} else {
		c.origins = strings.Join(config.AllowedOrigins, ", ")
	}
	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = []string{
			"Origin", "Accept", "Content-Type",
			"X-Requested-With", "X-Auth-Tokens", "Authorization",
		}
	}
	c.headers = strings.Join(config.AllowedHeaders, ",")
	if len(config.AllowedMethods) == 0 {
		c.methods = strings.Join([]string{
			fasthttp.MethodGet, fasthttp.MethodHead, fasthttp.MethodPost,
			fasthttp.MethodPatch, fasthttp.MethodConnect, fasthttp.MethodDelete,
			fasthttp.MethodTrace, fasthttp.MethodOptions,
		}, ", ")
	} else {
		c.methods = strings.Join(config.AllowedMethods, ", ")
	}

	return c
}

func (c *CORS) Handle(reqCtx *fasthttp.RequestCtx) bool {
	// ByPass CORS (Cross Origin Resource Sharing) check
	if c.origins == "*" {
		reqCtx.Response.Header.SetBytesV(
			fasthttp.HeaderAccessControlAllowOrigin,
			reqCtx.Request.Header.Peek(fasthttp.HeaderOrigin),
		)
	} else {
		reqCtx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, c.origins)
	}

	reqCtx.Response.Header.Set(fasthttp.HeaderAccessControlRequestMethod, c.methods)
	reqCtx.Response.Header.Add("Vary", fasthttp.HeaderOrigin)

	if reqCtx.Request.Header.IsOptions() {
		reqCtx.Response.Header.Add("Vary", fasthttp.HeaderAccessControlRequestMethod)
		reqCtx.Response.Header.Add("Vary", fasthttp.HeaderAccessControlRequestHeaders)
		reqHeaders := reqCtx.Request.Header.Peek(fasthttp.HeaderAccessControlRequestHeaders)
		if len(reqHeaders) > 0 {
			reqCtx.Response.Header.SetBytesV(fasthttp.HeaderAccessControlAllowHeaders, reqHeaders)
		} else {
			reqCtx.Response.Header.Set(fasthttp.HeaderAccessControlAllowHeaders, c.headers)
		}

		reqCtx.SetStatusCode(http.StatusNoContent)
		reqCtx.SetConnectionClose()

		return true
	}

	return false
}
