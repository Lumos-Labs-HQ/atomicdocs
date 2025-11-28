package main

import (
	"github.com/valyala/fasthttp"
	"github.com/yourusername/atomicdocs/internal/middleware"
	"github.com/yourusername/atomicdocs/internal/registry"
)

func main() {
	reg := registry.New()
	handler := middleware.NewHandler(reg)
	
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		
		switch path {
		case "/api/register":
			handler.RegisterRoutes(ctx)
		case "/docs":
			handler.ServeUI(ctx)
		case "/docs/json":
			handler.GetSpec(ctx)
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	}
	
	fasthttp.ListenAndServe(":6174", requestHandler)
}
