package generator

import (
	"fmt"
	"os"
	"text/template"
)

type CodeGenerator struct {
	APIDesign *APIDesigner
}

func NewCodeGenerator(apiDesign *APIDesigner) *CodeGenerator {
	return &CodeGenerator{
		APIDesign: apiDesign,
	}
}

func (cg *CodeGenerator) GenerateAPICode() error {
	// Generate main.go
	if err := cg.generateMainFile(); err != nil {
		return fmt.Errorf("error generating main.go: %v", err)
	}

	// Generate handlers.go
	if err := cg.generateHandlersFile(); err != nil {
		return fmt.Errorf("error generating handlers.go: %v", err)
	}

	return nil
}

func (cg *CodeGenerator) generateMainFile() error {
	mainTemplate := `package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	{{range .Endpoints}}
	r.{{.Method}}("{{.Path}}", {{.FunctionName}}Handler)
	{{end}}

	r.Run(":8080")
}
`

	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create("generated_main.go")
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, cg.APIDesign)
}

func (cg *CodeGenerator) generateHandlersFile() error {
	handlersTemplate := `package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

{{range .Endpoints}}
func {{.FunctionName}}Handler(c *gin.Context) {
	{{range .Parameters}}
	var {{.Name}} {{.Type}}
	{{if eq .Location "body"}}
	if err := c.ShouldBindJSON(&{{.Name}}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	{{else if eq .Location "path"}}
	{{.Name}} = c.Param("{{.Name}}")
	{{else if eq .Location "query"}}
	{{.Name}} = c.Query("{{.Name}}")
	{{end}}
	{{end}}

	// TODO: Implement {{.FunctionName}} logic here

	c.JSON(http.StatusOK, gin.H{
		"message": "{{.FunctionName}} executed successfully",
	})
}
{{end}}
`

	tmpl, err := template.New("handlers").Parse(handlersTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create("generated_handlers.go")
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, cg.APIDesign)
}

func (cg *CodeGenerator) GenerateGoModFile() error {
	goModContent := `module soft-crusher-api

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
)
`

	return os.WriteFile("go.mod", []byte(goModContent), 0644)
}
