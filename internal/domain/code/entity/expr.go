package entity

import (
	"errors"
	"github.com/dddplayer/dp/internal/domain/code"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"path"
	"strings"
)

type expression struct {
	expr  ast.Expr
	pkg   *packages.Package
	infos []*exprInfo
	errs  []error
}

type exprInfo struct {
	sel  string
	val  string
	ship code.RelationShip
}

func (e *expression) visit(cb func(path, name string, ship code.RelationShip)) {
	infos, err := getExprsInfo(e.expr)
	if err != nil {
		e.errs = append(e.errs, err)
		return
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

func extractExpr(expr ast.Expr) (sel, val string, ship code.RelationShip, err error) {
	switch x := expr.(type) {
	case *ast.StarExpr:
		exp := x.X
		switch y := exp.(type) {
		case *ast.SelectorExpr:
			if id, ok := y.X.(*ast.Ident); ok {
				sel = id.Name
			}
			val = y.Sel.Name
			ship = code.OneOne
		case *ast.Ident:
			val = y.Name
			ship = code.OneOne
		}
	case *ast.SelectorExpr:
		if id, ok := x.X.(*ast.Ident); ok {
			sel = id.Name
		}
		val = x.Sel.Name
		ship = code.OneOne
	case *ast.Ident:
		val = x.Name
		ship = code.OneOne
	case *ast.ArrayType:
		exp := x.Elt
		switch y := exp.(type) {
		case *ast.Ident:
			val = y.Name
			ship = code.OneMany
		}
	case *ast.MapType:
		err = errors.New("map currently not supported yet, ignore at this time")
	}

	return
}

func isBasicTypes(name string) bool {
	return strings.Contains("bool string int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64 uintptr "+
		"byte rune float32 float64 complex64 complex128 any interface{} error", name)
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
