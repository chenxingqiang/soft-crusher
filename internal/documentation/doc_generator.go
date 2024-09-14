package documentation

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type DocumentationGenerator struct {
	APIDesign *APIDesigner
}

func NewDocumentationGenerator(apiDesign *APIDesigner) *DocumentationGenerator {
	return &DocumentationGenerator{
		APIDesign: apiDesign,
	}
}

func (dg *DocumentationGenerator) GenerateSwaggerDoc() error {
	swagger := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Soft-Crusher Generated API",
			Version: "1.0.0",
		},
		Paths: openapi3.Paths{},
	}

	for _, endpoint := range dg.APIDesign.Endpoints {
		path := endpoint.Path
		method := strings.ToLower(endpoint.Method)

		operation := &openapi3.Operation{
			Summary:     fmt.Sprintf("%s operation", endpoint.FunctionName),
			Description: fmt.Sprintf("Endpoint for %s", endpoint.FunctionName),
			Parameters:  []*openapi3.ParameterRef{},
			Responses:   openapi3.Responses{},
		}

		for _, param := range endpoint.Parameters {
			paramLocation := openapi3.ParameterInQuery
			if param.Location == "path" {
				paramLocation = openapi3.ParameterInPath
			}

			parameter := &openapi3.Parameter{
				Name:        param.Name,
				In:          paramLocation,
				Description: fmt.Sprintf("Parameter %s", param.Name),
				Required:    paramLocation == openapi3.ParameterInPath,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: dg.convertGoTypeToSwaggerType(param.Type),
					},
				},
			}
			operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{Value: parameter})
		}

		for _, response := range endpoint.Responses {
			operation.Responses[fmt.Sprintf("%d", response.StatusCode)] = &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: &response.Type,
				},
			}
		}

		if swagger.Paths[path] == nil {
			swagger.Paths[path] = &openapi3.PathItem{}
		}
		switch method {
		case "get":
			swagger.Paths[path].Get = operation
		case "post":
			swagger.Paths[path].Post = operation
		case "put":
			swagger.Paths[path].Put = operation
		case "delete":
			swagger.Paths[path].Delete = operation
		}
	}

	return dg.writeSwaggerJSON(swagger)
}

func (dg *DocumentationGenerator) convertGoTypeToSwaggerType(goType string) string {
	switch goType {
	case "int", "int32", "int64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		return "string"
	}
}

func (dg *DocumentationGenerator) writeSwaggerJSON(swagger *openapi3.T) error {
	data, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling Swagger JSON: %v", err)
	}

	err = os.WriteFile("swagger.json", data, 0644)
	if err != nil {
		return fmt.Errorf("error writing Swagger JSON file: %v", err)
	}

	fmt.Println("Swagger documentation generated successfully: swagger.json")
	return nil
}
