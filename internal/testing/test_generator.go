package testing

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type TestingSuiteGenerator struct {
	APIDesign *APIDesigner
}

func NewTestingSuiteGenerator(apiDesign *APIDesigner) *TestingSuiteGenerator {
	return &TestingSuiteGenerator{
		APIDesign: apiDesign,
	}
}

func (tsg *TestingSuiteGenerator) GenerateTests() error {
	testTemplate := `package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	{{range .Endpoints}}
	r.{{.Method}}("{{.Path}}", {{.FunctionName}}Handler)
	{{end}}
	return r
}

{{range .Endpoints}}
func Test{{.FunctionName}}(t *testing.T) {
	router := SetupRouter()

	{{if eq .Method "GET"}}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("{{.Method}}", "{{.Path}}", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "{{.FunctionName}} executed successfully")
	{{else}}
	// Sample request body
	body := map[string]interface{}{
		{{range .Parameters}}
		"{{.Name}}": "sample_{{.Name}}",
		{{end}}
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("{{.Method}}", "{{.Path}}", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "{{.FunctionName}} executed successfully")
	{{end}}
}
{{end}}
`

	tmpl, err := template.New("tests").Parse(testTemplate)
	if err != nil {
		return fmt.Errorf("error parsing test template: %v", err)
	}

	f, err := os.Create("generated_handlers_test.go")
	if err != nil {
		return fmt.Errorf("error creating test file: %v", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, tsg.APIDesign)
	if err != nil {
		return fmt.Errorf("error executing test template: %v", err)
	}

	fmt.Println("Test suite generated successfully: generated_handlers_test.go")
	return nil
}

func (tsg *TestingSuiteGenerator) UpdateGoModFile() error {
	goModContent := `module soft-crusher-api

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/stretchr/testify v1.7.0
)
`

	return os.WriteFile("go.mod", []byte(goModContent), 0644)
}
