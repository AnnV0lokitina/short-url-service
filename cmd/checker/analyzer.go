package main

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var NoExitAnalyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "check if os.Exit don't use",
	Run:  run,
}

func parseMainFunction(pass *analysis.Pass, fun *ast.FuncDecl) {
	ast.Inspect(fun, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			if s, ok := c.Fun.(*ast.SelectorExpr); ok {
				if fmt.Sprint(s.X) == "os" && s.Sel.Name == "Exit" {
					pass.Reportf(s.X.Pos(), "os.Exit in main()")
				}
			}
		}
		return true
	})
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.Contains(filename, "test") || !strings.HasSuffix(filename, ".go") {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if f, ok := node.(*ast.FuncDecl); ok {
				if f.Name.Name == "main" {
					parseMainFunction(pass, f)
				}
			}
			return true
		})
	}
	return nil, nil
}
