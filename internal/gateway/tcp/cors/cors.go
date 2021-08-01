package cors

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"strings"
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
		c.origins = strings.Join(config.AllowedOrigins, ",")
	}
	if len(config.AllowedHeaders) == 0 {
		c.origins = "*"
	} else {
		c.origins = strings.Join(config.AllowedHeaders, ",")
	}
	if len(config.AllowedMethods) == 0 {
		c.origins = strings.Join([]string{
			fasthttp.MethodGet, fasthttp.MethodHead, fasthttp.MethodPost,
			fasthttp.MethodPatch, fasthttp.MethodConnect, fasthttp.MethodDelete,
			fasthttp.MethodTrace, fasthttp.MethodOptions,
		}, ",")
	} else {
		c.origins = strings.Join(config.AllowedMethods, ",")
	}
	return c
}

func (c *CORS) Handle(reqCtx *fasthttp.RequestCtx) bool {
	// ByPass CORS (Cross Origin Resource Sharing) check
	reqCtx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, c.origins)
	reqCtx.Response.Header.Set(fasthttp.HeaderAccessControlRequestMethod, c.methods)
	reqCtx.Response.Header.Set(fasthttp.HeaderAccessControlAllowHeaders, c.headers)

	if reqCtx.Request.Header.IsOptions() {
		reqCtx.SetStatusCode(http.StatusOK)
		reqCtx.SetConnectionClose()
		return true
	}
	return false
}
