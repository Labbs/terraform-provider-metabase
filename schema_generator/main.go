package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/labbs/terraform-provider-metabase/schema_generator/normalizer"
	"github.com/labbs/terraform-provider-metabase/schema_generator/oapi_codegen"
)

// Structs for OpenAPI JSON
type OpenAPI struct {
	Paths      map[string]map[string]Operation `json:"paths"`
	Components Components                      `json:"components"`
	Info       Info                            `json:"info"`
	OpenAPI    string                          `json:"openapi"`
}

type Operation struct {
	Summary     string              `json:"summary,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
}

type Parameter struct {
	Name     string                 `json:"name"`
	In       string                 `json:"in"`
	Required bool                   `json:"required,omitempty"`
	Schema   map[string]interface{} `json:"schema"`
}

type Response struct {
	Description string             `json:"description"`
	Content     map[string]Content `json:"content"`
}

type Content struct {
	Schema map[string]interface{} `json:"schema"`
}

type Components struct {
	Schemas map[string]interface{} `json:"schemas,omitempty"`
}

type RequestBody struct {
	Description   *string                `json:"description,omitempty"`
	Content       map[string]MediaType   `json:"content"` // Required.
	Required      *bool                  `json:"required,omitempty"`
	MapOfAnything map[string]interface{} `json:"-"` // Key must match pattern: `^x-`.
}

type MediaType struct {
	Schema Schema `json:"schema,omitempty"`
	// MapOfAnything map[string]interface{} `json:"-"` // Key must match pattern: `^x-`.
}

type Schema struct {
	Properties map[string]Propertie `json:"properties"`
	Required   []string             `json:"required"`
	Type       string               `json:"type"`
}

type Propertie struct {
	AllOf       []Propertie `json:"allOf,omitempty"`
	Description *string     `json:"description,omitempty"`
	Ref         string      `json:"$ref,omitempty"`
	Type        string      `json:"type,omitempty"`
	Items       *Items      `json:"items,omitempty"`
}

type Items struct {
	Type       string               `json:"type,omitempty"`
	Properties map[string]Propertie `json:"properties"`
	Required   []string             `json:"required,omitempty"`
}

type Reference struct {
	Ref         string  `json:"$ref"`
	Summary     *string `json:"summary,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

func sanitizeJSON(data []byte) []byte {
	stringData := string(data)

	// Replace "type": "null" with "nullable": true
	stringData = strings.ReplaceAll(stringData, `"type": "null"`, `"nullable": true`)
	stringData = strings.ReplaceAll(stringData, `"type":"null"`, `"nullable": true`)

	stringData = strings.ReplaceAll(stringData, "null%", "null")
	stringData = strings.ReplaceAll(stringData, "nil%", "null")

	return []byte(stringData)
}

func main() {
	resp, err := http.Get("http://127.0.0.1:3002/api/docs/openapi.json")
	if err != nil {
		log.Fatalf("Error retrieving OpenAPI JSON from URL : %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body : %v", err)
	}

	body = sanitizeJSON(body)
	if err != nil {
		log.Fatalf("Error sanitizing JSON: %v", err)
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		log.Fatalf("Error initially decoding JSON : %v", err)
	}

	var mapper *normalizer.SchemaNameMapper
	if components, ok := rawData["components"].(map[string]interface{}); ok {
		if schemas, ok := components["schemas"].(map[string]interface{}); ok {
			var normalizedSchemas map[string]interface{}
			normalizedSchemas, mapper = normalizer.NormalizeSchemas(schemas)
			components["schemas"] = normalizedSchemas
		}
	}

	// Normalize refs
	normalizedData := normalizer.NormalizeRefs(rawData, mapper)

	// Convert normalized data back to JSON
	intermediateJSON, err := json.Marshal(normalizedData)
	if err != nil {
		log.Fatalf("Error converting normalized data : %v", err)
	}

	var openapi OpenAPI
	if err := json.Unmarshal(intermediateJSON, &openapi); err != nil {
		log.Fatalf("Error while decoding final JSON : %v", err)
	}

	if openapi.Components.Schemas == nil {
		openapi.Components.Schemas = make(map[string]interface{})
	}

	// Add missing definitions if needed
	openapi.Components.Schemas["metabase_util_cron_ScheduleMap"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"schedule_day": map[string]interface{}{
				"type": "string",
			},
			"schedule_frame": map[string]interface{}{
				"type": "string",
			},
			"schedule_hour": map[string]interface{}{
				"type": "integer",
			},
			"schedule_minute": map[string]interface{}{
				"type": "integer",
			},
			"schedule_type": map[string]interface{}{
				"type": "string",
			},
		},
		"required": []string{"schedule_type"},
	}

	paramPattern := regexp.MustCompile(`\{([^}]+)\}`)

	for path, operations := range openapi.Paths {
		expectedParams := paramPattern.FindAllStringSubmatch(path, -1)

		for method, operation := range operations {
			if operation.Responses == nil {
				operation.Responses = map[string]Response{
					"200": {
						Description: "Default response",
						Content: map[string]Content{
							"application/json": {
								Schema: map[string]interface{}{
									"type":                 "object",
									"additionalProperties": true,
								},
							},
						},
					},
				}
			}

			existingParams := map[string]bool{}
			for _, param := range operation.Parameters {
				existingParams[param.Name] = true
			}

			for _, match := range expectedParams {
				paramName := match[1]
				if !existingParams[paramName] {
					operation.Parameters = append(operation.Parameters, Parameter{
						Name:     paramName,
						In:       "path",
						Required: true,
						Schema: map[string]interface{}{
							"type": "string",
						},
					})
				}
			}
			openapi.Paths[path][method] = operation
		}
	}

	// Create version directory
	version := strings.Join(strings.Split(openapi.Info.Version, ".")[:2], ".")
	var path string = "../metabase/" + strings.ReplaceAll(version, ".", "_")

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Failed to create version directory: %v", err)
	}

	updatedJSON, err := json.MarshalIndent(openapi, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode json : %v", err)
	}

	if err := os.WriteFile(path+"/updated-openapi.json", updatedJSON, 0644); err != nil {
		log.Fatalf("Failed to save the updated OpenAPI json file: %v", err)
	}

	oapi_codegen.NewOapiCodegenConfig(path)

	oapi_codegen.ExecuteCodegen(path)
}
