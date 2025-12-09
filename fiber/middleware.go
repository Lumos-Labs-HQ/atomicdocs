package atomicdocs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type RouteInfo struct {
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Handler     string                 `json:"handler"`
	FilePath    string                 `json:"filePath,omitempty"`
	Imports     []Import               `json:"imports,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Parameters  []Parameter            `json:"parameters,omitempty"`
	RequestBody *RequestBody           `json:"requestBody,omitempty"`
	Responses   map[string]Response    `json:"responses,omitempty"`
}

type Import struct {
	Name string `json:"name"`
	From string `json:"from"`
}

type Parameter struct {
	Name     string      `json:"name"`
	In       string      `json:"in"`
	Required bool        `json:"required,omitempty"`
	Schema   Schema      `json:"schema"`
}

type RequestBody struct {
	Required bool                       `json:"required,omitempty"`
	Content  map[string]MediaTypeObject `json:"content"`
}

type MediaTypeObject struct {
	Schema Schema `json:"schema"`
}

type Response struct {
	Description string                     `json:"description"`
	Content     map[string]MediaTypeObject `json:"content,omitempty"`
}

type Schema struct {
	Type       string            `json:"type,omitempty"`
	Format     string            `json:"format,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
}

type RegistrationPayload struct {
	Routes      []RouteInfo       `json:"routes"`
	Port        int               `json:"port"`
	SchemaFiles map[string]string `json:"schemaFiles"`
}

func New(port int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		
		if path == "/docs" || path == "/docs/json" {
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:6174%s", path), nil)
			if err != nil {
				return c.Status(503).SendString("Request creation failed")
			}
			
			req.Header.Set("X-App-Port", fmt.Sprintf("%d", port))
			
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return c.Status(503).SendString(fmt.Sprintf("AtomicDocs server not reachable: %v", err))
			}
			defer resp.Body.Close()
			
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return c.Status(503).SendString("Failed to read response")
			}
			
			c.Set("Content-Type", resp.Header.Get("Content-Type"))
			return c.Status(resp.StatusCode).Send(body)
		}
		
		return c.Next()
	}
}

func Register(app *fiber.App, port int) {
	routes := extractRoutes(app)
	
	payload := RegistrationPayload{
		Routes:      routes,
		Port:        port,
		SchemaFiles: make(map[string]string), // Empty for now
	}
	
	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:6174/api/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("✗ AtomicDocs: Registration failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("✓ AtomicDocs: Registered %d routes\n", len(routes))
	fmt.Printf("📚 Docs: http://localhost:%d/docs\n", port)
}

func extractRoutes(app *fiber.App) []RouteInfo {
	routes := []RouteInfo{}
	allRoutes := app.GetRoutes()
	
	fmt.Printf("Debug: Found %d total routes\n", len(allRoutes))
	
	for _, route := range allRoutes {
		fmt.Printf("Debug: Processing route %s %s\n", route.Method, route.Path)
		
		if strings.HasPrefix(route.Path, "/docs") {
			fmt.Printf("Debug: Skipping docs route\n")
			continue
		}
		
		if len(route.Handlers) == 0 {
			fmt.Printf("Debug: No handlers for route\n")
			continue
		}
		
		handlerName := "unknown"
		if len(route.Handlers) > 0 {
			handlerName = runtime.FuncForPC(reflect.ValueOf(route.Handlers[len(route.Handlers)-1]).Pointer()).Name()
		}
		
		info := RouteInfo{
			Method:      route.Method,
			Path:        route.Path,
			Handler:     handlerName,
			Summary:     route.Method + " " + route.Path,
			Description: fmt.Sprintf("Handler: %s", handlerName),
			Tags:        extractTags(route.Path),
			Parameters:  extractParameters(route.Path),
			Responses:   generateResponses(route.Method, route.Path),
		}
		
		if route.Method == "POST" || route.Method == "PUT" || route.Method == "PATCH" {
			info.RequestBody = generateRequestBody(route.Path)
		}
		
		routes = append(routes, info)
		fmt.Printf("Debug: Added route %s %s\n", info.Method, info.Path)
	}
	
	fmt.Printf("Debug: Returning %d routes\n", len(routes))
	return routes
}

func extractTags(path string) []string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && parts[0] != "" {
		tag := parts[0]
		// Capitalize first letter
		if len(tag) > 0 {
			tag = strings.ToUpper(tag[:1]) + tag[1:]
		}
		return []string{tag}
	}
	return []string{"API"}
}

func generateResponses(method, path string) map[string]Response {
	responses := map[string]Response{
		"200": {
			Description: "Successful response",
			Content: map[string]MediaTypeObject{
				"application/json": {
					Schema: Schema{Type: "object"},
				},
			},
		},
	}
	
	if method == "POST" {
		responses["201"] = Response{
			Description: "Created successfully",
			Content: map[string]MediaTypeObject{
				"application/json": {Schema: Schema{Type: "object"}},
			},
		}
	}
	
	if strings.Contains(path, ":") {
		responses["404"] = Response{Description: "Resource not found"}
	}
	
	responses["400"] = Response{Description: "Bad request"}
	responses["500"] = Response{Description: "Internal server error"}
	
	return responses
}

func extractParameters(path string) []Parameter {
	params := []Parameter{}
	re := regexp.MustCompile(`:(\w+)`)
	for _, match := range re.FindAllStringSubmatch(path, -1) {
		params = append(params, Parameter{
			Name:     match[1],
			In:       "path",
			Required: true,
			Schema:   Schema{Type: "string"},
		})
	}
	return params
}

func generateRequestBody(path string) *RequestBody {
	schema := Schema{
		Type: "object",
		Properties: map[string]Schema{
			"name":  {Type: "string"},
			"email": {Type: "string"},
		},
	}
	
	if strings.Contains(path, "product") {
		schema.Properties = map[string]Schema{
			"name":  {Type: "string"},
			"price": {Type: "number"},
			"stock": {Type: "integer"},
		}
	}
	
	return &RequestBody{
		Required: true,
		Content: map[string]MediaTypeObject{
			"application/json": {Schema: schema},
		},
	}
}

func generateSummary(method, path string) string {
	action := map[string]string{
		"GET": "Get", "POST": "Create", "PUT": "Update", "DELETE": "Delete", "PATCH": "Update",
	}[method]
	
	resource := "resource"
	if strings.Contains(path, "user") { resource = "user" }
	if strings.Contains(path, "product") { resource = "product" }
	if strings.Contains(path, "auth") { resource = "authentication" }
	
	return fmt.Sprintf("%s %s", action, resource)
}

func generateDescription(path string) string {
	if strings.Contains(path, "login") { return "Authenticate user and return JWT token" }
	if strings.Contains(path, "register") { return "Register new user account" }
	if strings.Contains(path, ":id") { return "Operation on specific resource by ID" }
	return "API endpoint operation"
}

func generateDetailedResponses(method, path string) map[string]Response {
	responses := map[string]Response{
		"200": {
			Description: "Success",
			Content: map[string]MediaTypeObject{
				"application/json": {Schema: getResponseSchema(path)},
			},
		},
		"400": {Description: "Bad Request - Invalid input"},
		"500": {Description: "Internal Server Error"},
	}
	
	if method == "POST" {
		responses["201"] = Response{
			Description: "Created successfully",
			Content: map[string]MediaTypeObject{
				"application/json": {Schema: getResponseSchema(path)},
			},
		}
	}
	
	if strings.Contains(path, ":") {
		responses["404"] = Response{Description: "Resource not found"}
	}
	
	return responses
}

func generateDetailedRequestBody(path string) *RequestBody {
	return &RequestBody{
		Required: true,
		Content: map[string]MediaTypeObject{
			"application/json": {Schema: getRequestSchema(path)},
		},
	}
}

func getRequestSchema(path string) Schema {
	if strings.Contains(path, "user") {
		return Schema{
			Type: "object",
			Properties: map[string]Schema{
				"name":  {Type: "string"},
				"email": {Type: "string", Format: "email"},
				"age":   {Type: "integer"},
			},
		}
	}
	if strings.Contains(path, "product") {
		return Schema{
			Type: "object", 
			Properties: map[string]Schema{
				"name":     {Type: "string"},
				"price":    {Type: "number"},
				"stock":    {Type: "integer"},
				"category": {Type: "string"},
			},
		}
	}
	return Schema{Type: "object"}
}

func getResponseSchema(path string) Schema {
	if strings.Contains(path, "user") {
		return Schema{
			Type: "object",
			Properties: map[string]Schema{
				"id":    {Type: "integer"},
				"name":  {Type: "string"},
				"email": {Type: "string"},
				"age":   {Type: "integer"},
			},
		}
	}
	if strings.Contains(path, "product") {
		return Schema{
			Type: "object",
			Properties: map[string]Schema{
				"id":       {Type: "integer"},
				"name":     {Type: "string"},
				"price":    {Type: "number"},
				"stock":    {Type: "integer"},
				"category": {Type: "string"},
			},
		}
	}
	return Schema{Type: "object"}
}
