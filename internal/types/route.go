package types

type RouteInfo struct {
	Method      string              `json:"method"`
	Path        string              `json:"path"`
	Handler     string              `json:"handler"`
	FilePath    string              `json:"filePath,omitempty"`
	Imports     []Import            `json:"imports,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

type Import struct {
	Name string `json:"name"`
	From string `json:"from"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required,omitempty"`
	Schema      Schema `json:"schema"`
	Description string `json:"description,omitempty"`
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
	Type       string              `json:"type,omitempty"`
	Properties map[string]Schema   `json:"properties,omitempty"`
	Items      *Schema             `json:"items,omitempty"`
	Example    interface{}         `json:"example,omitempty"`
}
