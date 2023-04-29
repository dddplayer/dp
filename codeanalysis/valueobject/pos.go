package valueobject

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

type Position interface {
	Filename() string
	Offset() int
	Line() int
	Column() int
}

func NewPosition(f string, o, l, c int) Position {
	return &pos{
		filename: f,
		offset:   o,
		line:     l,
		column:   c,
	}
}

type pos struct {
	filename string
	offset   int
	line     int
	column   int
}

func (p *pos) Filename() string { return p.filename }
func (p *pos) Offset() int      { return p.offset }
func (p *pos) Line() int        { return p.line }
func (p *pos) Column() int      { return p.column }

func AstPosition(pkg *packages.Package, node ast.Node) Position {
	return getPosition(pkg.Fset, node.Pos())
}

func SsaPosition(pkg *ssa.Package, t *types.TypeName) Position {
	return getPosition(pkg.Prog.Fset, t.Pos())
}

func SsaFuncPosition(pkg *ssa.Package, f *ssa.Function) Position {
	return getPosition(pkg.Prog.Fset, f.Pos())
}

func SsaInstructionPosition(pkg *ssa.Package, f ssa.Instruction) Position {
	return getPosition(pkg.Prog.Fset, f.Pos())
}

func getPosition(f *token.FileSet, pos token.Pos) Position {
	position := f.Position(pos)
	if position.IsValid() != true {
		return nil
	}
	return NewPosition(position.Filename, position.Offset, position.Line, position.Column)
}
