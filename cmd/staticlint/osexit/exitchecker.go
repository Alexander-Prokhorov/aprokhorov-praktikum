// Пакет osexit предназначен для проверки использования os.Exit() в функции main() пакета main.
package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	mainName = "main"
	pkg      = "os"
	cmd      = "Exit"
)

var OSExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for use os.Exit in func main()",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	isPackageMain := func(x *ast.File) bool {
		return x.Name.String() == mainName
	}

	isFuncMain := func(x *ast.FuncDecl) bool {
		return x.Name.String() == mainName
	}

	/*
		isOsExit := func(x *ast.CallExpr) bool {
			if x, ok := x.Fun.(*ast.SelectorExpr); ok {
				if y, ok := x.X.(*ast.Ident); ok {
					if y.Name == pkg && x.Sel.String() == cmd {
						return true
					}
				}
			}
			return false
		}
	*/

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if x, ok := node.(*ast.File); ok && isPackageMain(x) {
				ast.Inspect(x, func(node ast.Node) bool {
					if x, ok := node.(*ast.FuncDecl); ok && isFuncMain(x) {
						ast.Inspect(x, func(node ast.Node) bool {
							callExpr, ok := node.(*ast.CallExpr)
							if !ok {
								return true
							}

							selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
							if !ok {
								return true
							}

							ident, ok := selectorExpr.X.(*ast.Ident)
							if !ok {
								return true
							}

							if ident.Name == pkg && selectorExpr.Sel.String() == cmd {
								pass.Reportf(node.Pos(), "os.Exit used in func main of package main")
							}

							return true
						})
					}
					return true
				})
			}
			return true
		})
	}

	return nil, nil
}
