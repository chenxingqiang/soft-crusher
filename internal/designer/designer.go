package designer

import (
	"fmt"
	"strings"
	"unicode"
)

type APIEndpoint struct {
	Method       string
	Path         string
	FunctionName string
	Parameters   []Parameter
	Responses    []Response
}

type Parameter struct {
	Name     string
	Type     string
	Location string // "path", "query", or "body"
}

type Response struct {
	StatusCode int
	Type       string
}

type APIDesigner struct {
	Endpoints []APIEndpoint
}

func NewAPIDesigner() *APIDesigner {
	return &APIDesigner{
		Endpoints: make([]APIEndpoint, 0),
	}
}

func (ad *APIDesigner) DesignAPI(functions []FunctionInfo) {
	for _, fn := range functions {
		endpoint := APIEndpoint{
			Method:       ad.inferHTTPMethod(fn.Name),
			Path:         ad.generatePath(fn.Name),
			FunctionName: fn.Name,
			Parameters:   ad.generateParameters(fn.Args),
			Responses:    ad.generateResponses(fn.Returns),
		}
		ad.Endpoints = append(ad.Endpoints, endpoint)
	}
}

func (ad *APIDesigner) inferHTTPMethod(funcName string) string {
	lowercaseName := strings.ToLower(funcName)
	switch {
	case strings.HasPrefix(lowercaseName, "get"):
		return "GET"
	case strings.HasPrefix(lowercaseName, "create") || strings.HasPrefix(lowercaseName, "add"):
		return "POST"
	case strings.HasPrefix(lowercaseName, "update"):
		return "PUT"
	case strings.HasPrefix(lowercaseName, "delete"):
		return "DELETE"
	default:
		return "POST"
	}
}

func (ad *APIDesigner) generatePath(funcName string) string {
	// Convert camelCase to kebab-case
	var result strings.Builder
	for i, r := range funcName {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return "/" + result.String()
}

func (ad *APIDesigner) generateParameters(args []string) []Parameter {
	parameters := make([]Parameter, 0)
	for _, arg := range args {
		parts := strings.Split(arg, " ")
		if len(parts) == 2 {
			param := Parameter{
				Name:     parts[0],
				Type:     parts[1],
				Location: "body", // Default to body, can be refined later
			}
			parameters = append(parameters, param)
		}
	}
	return parameters
}

func (ad *APIDesigner) generateResponses(returns []string) []Response {
	responses := []Response{{StatusCode: 200, Type: "OK"}}
	if len(returns) > 0 {
		responses[0].Type = strings.Join(returns, ", ")
	}
	return responses
}

func (ad *APIDesigner) PrintAPIDesign() {
	for _, endpoint := range ad.Endpoints {
		fmt.Printf("Endpoint: %s %s\n", endpoint.Method, endpoint.Path)
		fmt.Printf("  Function: %s\n", endpoint.FunctionName)
		fmt.Println("  Parameters:")
		for _, param := range endpoint.Parameters {
			fmt.Printf("    - %s (%s): %s\n", param.Name, param.Type, param.Location)
		}
		fmt.Println("  Responses:")
		for _, resp := range endpoint.Responses {
			fmt.Printf("    - %d: %s\n", resp.StatusCode, resp.Type)
		}
		fmt.Println()
	}
}
