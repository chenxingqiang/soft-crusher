package generator

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/chenxingqiang/soft-crusher/internal/analyzer"
)

type APIGenerator struct {
	Functions []analyzer.FunctionInfo
}

func NewAPIGenerator(functions []analyzer.FunctionInfo) *APIGenerator {
	return &APIGenerator{
		Functions: functions,
	}
}

func (ag *APIGenerator) GenerateAPI() (string, error) {
	apiTemplate := `
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

{{range .Functions}}
type {{.Name}}Request struct {
	{{range .Parameters}}
	{{.Name}} {{.Type}} ` + "`" + `json:"{{.Name | toLower}}" validate:"required"` + "`" + `
	{{end}}
}

func {{.Name}}Handler(c *gin.Context) {
	var req {{.Name}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement {{.Name}} logic here

	c.JSON(http.StatusOK, gin.H{"message": "{{.Name}} executed successfully"})
}
{{end}}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	{{range .Functions}}
	r.POST("/{{.Name | toLower}}", {{.Name}}Handler)
	{{end}}

	return r
}

func main() {
	r := SetupRouter()
	r.Run(":8080")
}
`

	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
	}

	tmpl, err := template.New("api").Funcs(funcMap).Parse(apiTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing API template: %v", err)
	}

	var result strings.Builder
	err = tmpl.Execute(&result, ag)
	if err != nil {
		return "", fmt.Errorf("error executing API template: %v", err)
	}

	return result.String(), nil
}
