package parser

import (
	"regexp"
	"strings"
	"github.com/yourusername/atomicdocs/internal/types"
)

func ParseSchemaFromFile(schemaName string, fileContent string) map[string]types.Schema {
	// Try Zod first
	if schema := parseZodSchema(schemaName, fileContent); schema != nil {
		return schema
	}
	
	// Try Yup
	if schema := parseYupSchema(schemaName, fileContent); schema != nil {
		return schema
	}
	
	// Try Joi
	if schema := parseJoiSchema(schemaName, fileContent); schema != nil {
		return schema
	}
	
	return nil
}

func parseZodSchema(schemaName string, content string) map[string]types.Schema {
	// Pattern: export const SignupSchema = z.object({ ... })
	// Use (?s) for multiline matching
	pattern := regexp.MustCompile(`(?s)(?:export\s+)?(?:const|let|var)\s+` + schemaName + `\s*=\s*z\.object\(\{(.+?)\}\)`)
	matches := pattern.FindStringSubmatch(content)
	
	if len(matches) < 2 {
		return nil
	}
	
	schemaBody := matches[1]
	properties := make(map[string]types.Schema)
	
	// Extract fields: email: z.string(), name: z.number()
	fieldPattern := regexp.MustCompile(`(\w+)\s*:\s*z\.(\w+)\(\)`)
	fieldMatches := fieldPattern.FindAllStringSubmatch(schemaBody, -1)
	
	for _, match := range fieldMatches {
		fieldName := match[1]
		zodType := match[2]
		
		properties[fieldName] = types.Schema{
			Type: zodTypeToOpenAPI(zodType),
		}
	}
	
	return properties
}

func parseYupSchema(schemaName string, content string) map[string]types.Schema {
	// Pattern: export const SignupSchema = yup.object({ ... })
	pattern := regexp.MustCompile(`(?:export\s+)?(?:const|let|var)\s+` + schemaName + `\s*=\s*yup\.object\(\{([^}]+)\}\)`)
	matches := pattern.FindStringSubmatch(content)
	
	if len(matches) < 2 {
		return nil
	}
	
	schemaBody := matches[1]
	properties := make(map[string]types.Schema)
	
	// Extract fields: email: yup.string()
	fieldPattern := regexp.MustCompile(`(\w+)\s*:\s*yup\.(\w+)\(\)`)
	fieldMatches := fieldPattern.FindAllStringSubmatch(schemaBody, -1)
	
	for _, match := range fieldMatches {
		fieldName := match[1]
		yupType := match[2]
		
		properties[fieldName] = types.Schema{
			Type: yupTypeToOpenAPI(yupType),
		}
	}
	
	return properties
}

func parseJoiSchema(schemaName string, content string) map[string]types.Schema {
	// Pattern: export const SignupSchema = Joi.object({ ... })
	pattern := regexp.MustCompile(`(?:export\s+)?(?:const|let|var)\s+` + schemaName + `\s*=\s*Joi\.object\(\{([^}]+)\}\)`)
	matches := pattern.FindStringSubmatch(content)
	
	if len(matches) < 2 {
		return nil
	}
	
	schemaBody := matches[1]
	properties := make(map[string]types.Schema)
	
	// Extract fields: email: Joi.string()
	fieldPattern := regexp.MustCompile(`(\w+)\s*:\s*Joi\.(\w+)\(\)`)
	fieldMatches := fieldPattern.FindAllStringSubmatch(schemaBody, -1)
	
	for _, match := range fieldMatches {
		fieldName := match[1]
		joiType := match[2]
		
		properties[fieldName] = types.Schema{
			Type: joiTypeToOpenAPI(joiType),
		}
	}
	
	return properties
}

func zodTypeToOpenAPI(zodType string) string {
	switch strings.ToLower(zodType) {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		return "array"
	case "object":
		return "object"
	default:
		return "string"
	}
}

func yupTypeToOpenAPI(yupType string) string {
	switch strings.ToLower(yupType) {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		return "array"
	case "object":
		return "object"
	default:
		return "string"
	}
}

func joiTypeToOpenAPI(joiType string) string {
	switch strings.ToLower(joiType) {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		return "array"
	case "object":
		return "object"
	default:
		return "string"
	}
}
