// Package noexit анализатор запрещает прямой вызов os.Exit в main.
package noexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer запрещает использование os.Exit в main.
var Analyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "запрещает прямой вызов os.Exit в main",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	// для пакета main
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// ищем main()
			if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "main" {

				// ищем os.Exit()
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					if call, ok := n.(*ast.CallExpr); ok {
						if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
							if pkgIdent, okident := sel.X.(*ast.Ident); okident && pkgIdent.Name == "os" && sel.Sel.Name == "Exit" {
								pass.Reportf(call.Pos(), "нельзя использовать os.Exit в main")
							}
						}
					}
					return true
				})
			}
			return true
		})

	}
	return nil, nil
}
