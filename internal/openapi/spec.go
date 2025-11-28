package openapi

import (
	"github.com/yourusername/atomicdocs/internal/types"
)

type Spec struct {
	OpenAPI string                `json:"openapi"`
	Info    Info                  `json:"info"`
	Servers []Server              `json:"servers"`
	Paths   map[string]PathItem   `json:"paths"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
}

type Operation struct {
	Summary     string                     `json:"summary,omitempty"`
	Description string                     `json:"description,omitempty"`
	Tags        []string                   `json:"tags,omitempty"`
	Parameters  []types.Parameter          `json:"parameters,omitempty"`
	RequestBody *types.RequestBody         `json:"requestBody,omitempty"`
	Responses   map[string]types.Response  `json:"responses"`
}

func Generate(routes []types.RouteInfo, baseURL string) *Spec {
	paths := make(map[string]PathItem)
	
	for _, route := range routes {
		if paths[route.Path].Get == nil && paths[route.Path].Post == nil &&
			paths[route.Path].Put == nil && paths[route.Path].Delete == nil &&
			paths[route.Path].Patch == nil {
			paths[route.Path] = PathItem{}
		}
		
		op := &Operation{
			Summary:     route.Summary,
			Description: route.Description,
			Tags:        route.Tags,
			Parameters:  route.Parameters,
			RequestBody: route.RequestBody,
			Responses:   route.Responses,
		}
		
		if op.Responses == nil {
			op.Responses = map[string]types.Response{
				"200": {Description: "Successful response"},
			}
		}
		
		item := paths[route.Path]
		switch route.Method {
		case "GET":
			item.Get = op
		case "POST":
			item.Post = op
		case "PUT":
			item.Put = op
		case "DELETE":
			item.Delete = op
		case "PATCH":
			item.Patch = op
		}
		paths[route.Path] = item
	}
	
	return &Spec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       "API Documentation",
			Description: "Auto-generated API documentation",
			Version:     "1.0.0",
		},
		Servers: []Server{{URL: baseURL}},
		Paths:   paths,
	}
}
