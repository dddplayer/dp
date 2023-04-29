package entity

import (
	"errors"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"path"
	"strings"
)

type expression struct {
	expr  ast.Expr
	pkg   *packages.Package
	infos []*exprInfo
}

type exprInfo struct {
	sel  string
	val  string
	ship RelationShip
}

func (e *expression) visit(cb func(path, name string, ship RelationShip)) {
	infos, err := getExprsInfo(e.expr)
	if err != nil {
		panic(err)
	}
	for _, info := range infos {
		var exprPath string
		if info.sel == "" {
			if isBasicTypes(info.val) {
				continue
			}
			exprPath = e.pkg.ID
		} else {
			exprPath = getPath(findFile(e.pkg, e.expr).Imports, info.sel)
		}

		if exprPath != "" {
			cb(exprPath, info.val, info.ship)
		}
	}
}

func getExprsInfo(expr ast.Expr) (exprInfos []*exprInfo, err error) {
	switch x := expr.(type) {
	case *ast.MapType:
		for _, exp := range []ast.Expr{x.Key, x.Value} {
			info, err := getExprInfo(exp)
			if err != nil {
				return nil, err
			}
			exprInfos = append(exprInfos, info)
		}

	default:
		info, err := getExprInfo(expr)
		if err != nil {
			return nil, err
		}
		exprInfos = append(exprInfos, info)
	}

	return
}

func getExprInfo(expr ast.Expr) (*exprInfo, error) {
	sel, val, ship, err := extractExpr(expr)
	if err != nil {
		return nil, err
	}
	return &exprInfo{
		sel:  sel,
		val:  val,
		ship: ship,
	}, nil
}

func extractExpr(expr ast.Expr) (sel, val string, ship RelationShip, err error) {
	switch x := expr.(type) {
	case *ast.StarExpr:
		exp := x.X
		switch y := exp.(type) {
		case *ast.SelectorExpr:
			if id, ok := y.X.(*ast.Ident); ok {
				sel = id.Name
			}
			val = y.Sel.Name
			ship = OneOne
		case *ast.Ident:
			val = y.Name
			ship = OneOne
		}
	case *ast.SelectorExpr:
		if id, ok := x.X.(*ast.Ident); ok {
			sel = id.Name
		}
		val = x.Sel.Name
		ship = OneOne
	case *ast.Ident:
		val = x.Name
		ship = OneOne
	case *ast.ArrayType:
		exp := x.Elt
		switch y := exp.(type) {
		case *ast.Ident:
			val = y.Name
			ship = OneMany
		}
	case *ast.MapType:
		err = errors.New("map nesting currently not supported yet")
	}

	return
}

func isBasicTypes(name string) bool {
	return strings.Contains("string int any byte uint", name)
}

func getPath(imports []*ast.ImportSpec, name string) string {
	for _, spec := range imports {
		if spec.Name != nil && spec.Name.Name == name {
			return trimDoubleQuote(spec.Path.Value)
		} else if spec.Name == nil && path.Base(trimDoubleQuote(spec.Path.Value)) == name {
			return trimDoubleQuote(spec.Path.Value)
		}
	}
	return ""
}

func trimDoubleQuote(str string) string {
	reg := "\""
	return strings.TrimPrefix(strings.TrimSuffix(str, reg), reg)
}
