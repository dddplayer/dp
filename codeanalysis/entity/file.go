package entity

import (
	"go/ast"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
)

func findFile(pkg *packages.Package, expr ast.Expr) *ast.File {
	if expr == nil {
		return nil
	}

	filename := pkg.Fset.Position(expr.Pos()).Filename
	if filename != "" {
		idx := slices.Index(pkg.CompiledGoFiles, filename)
		if idx != -1 {
			return pkg.Syntax[idx]
		}
	}
	return nil
}
