package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FunctionInfo struct {
	Name       string
	Receiver   string
	Parameters []ParameterInfo
	Results    []ParameterInfo
	IsMethod   bool
	IsGeneric  bool
}

type ParameterInfo struct {
	Name string
	Type string
}

type FunctionAnalyzer struct {
	Functions []FunctionInfo
}

func NewFunctionAnalyzer() *FunctionAnalyzer {
	return &FunctionAnalyzer{
		Functions: make([]FunctionInfo, 0),
	}
}

func (fa *FunctionAnalyzer) AnalyzeFile(filePath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			funcInfo := fa.analyzeFuncDecl(x)
			fa.Functions = append(fa.Functions, funcInfo)
		}
		return true
	})

	return nil
}

func (fa *FunctionAnalyzer) analyzeFuncDecl(funcDecl *ast.FuncDecl) FunctionInfo {
	funcInfo := FunctionInfo{
		Name:      funcDecl.Name.Name,
		IsMethod:  funcDecl.Recv != nil,
		IsGeneric: funcDecl.Type.TypeParams != nil,
	}

	if funcInfo.IsMethod {
		funcInfo.Receiver = fa.extractReceiver(funcDecl.Recv)
	}

	funcInfo.Parameters = fa.extractFieldList(funcDecl.Type.Params)
	funcInfo.Results = fa.extractFieldList(funcDecl.Type.Results)

	return funcInfo
}

func (fa *FunctionAnalyzer) extractReceiver(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}
	return fa.typeToString(recv.List[0].Type)
}

func (fa *FunctionAnalyzer) extractFieldList(fieldList *ast.FieldList) []ParameterInfo {
	var params []ParameterInfo
	if fieldList == nil {
		return params
	}
	for _, field := range fieldList.List {
		typeStr := fa.typeToString(field.Type)
		if len(field.Names) == 0 {
			params = append(params, ParameterInfo{"", typeStr})
		} else {
			for _, name := range field.Names {
				params = append(params, ParameterInfo{name.Name, typeStr})
			}
		}
	}
	return params
}

func (fa *FunctionAnalyzer) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + fa.typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + fa.typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + fa.typeToString(t.Key) + "]" + fa.typeToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func()"
	case *ast.SelectorExpr:
		return fa.typeToString(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

func (fa *FunctionAnalyzer) AnalyzeDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			return fa.AnalyzeFile(path)
		}
		return nil
	})
}
