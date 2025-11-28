package middleware

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/yourusername/atomicdocs/internal/openapi"
	"github.com/yourusername/atomicdocs/internal/registry"
	"github.com/yourusername/atomicdocs/internal/types"
)

type Handler struct {
	registry *registry.Registry
}

type RegistrationPayload struct {
	Routes []types.RouteInfo `json:"routes"`
	Port   int               `json:"port"`
}

func NewHandler(reg *registry.Registry) *Handler {
	return &Handler{registry: reg}
}

func (h *Handler) RegisterRoutes(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")
	
	if string(ctx.Method()) == "OPTIONS" {
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}
	
	var payload RegistrationPayload
	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(`{"error":"invalid JSON"}`)
		return
	}
	
	h.registry.RegisterApp(payload.Port, payload.Routes)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(`{"status":"registered"}`)
}

func (h *Handler) GetSpec(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	
	portHeader := string(ctx.Request.Header.Peek("X-App-Port"))
	if portHeader == "" {
		portHeader = "3000"
	}
	
	routes := h.registry.GetByPort(portHeader)
	spec := openapi.Generate(routes, "http://localhost:"+portHeader)
	
	data, _ := json.Marshal(spec)
	ctx.SetBody(data)
}

func (h *Handler) ServeUI(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	ctx.SetBodyString(swaggerHTML)
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>API Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = () => {
            SwaggerUIBundle({
                url: '/docs/json',
                dom_id: '#swagger-ui',
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
                layout: "BaseLayout",
                deepLinking: true,
                tryItOutEnabled: true
            });
        };
    </script>
</body>
</html>`
