package valueobject

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"testing"
)

type position struct {
	filename string
	offset   int
	line     int
	column   int
}

func TestGetPosition(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", -1, 100)

	// 创建一个无效的位置（位置超出文件范围）
	invalidPos := token.Pos(200)

	// 测试无效位置的情况
	result := getPosition(fset, invalidPos)
	if result != nil {
		t.Errorf("getPosition(%v, %v) = %v, expected nil", fset, invalidPos, result)
	}

	pos := file.Pos(50)
	expected := &position{
		filename: "test.go",
		offset:   50,
		line:     1,
		column:   51,
	}

	result = getPosition(fset, pos)

	if result == nil {
		t.Errorf("getPosition(%v, %v) = nil, expected %v", fset, pos, expected)
		return
	}

	if result.Filename() != expected.filename {
		t.Errorf("getPosition(%v, %v).filename = %v, expected %v", fset, pos, result.Filename(), expected.filename)
	}

	if result.Offset() != expected.offset {
		t.Errorf("getPosition(%v, %v).offset = %v, expected %v", fset, pos, result.Offset(), expected.offset)
	}

	if result.Line() != expected.line {
		t.Errorf("getPosition(%v, %v).line = %v, expected %v", fset, pos, result.Line(), expected.line)
	}

	if result.Column() != expected.column {
		t.Errorf("getPosition(%v, %v).column = %v, expected %v", fset, pos, result.Column(), expected.column)
	}
}

func TestSsaPosition(t *testing.T) {
	pkg := &ssa.Package{
		Prog: &ssa.Program{
			Fset: token.NewFileSet(),
		},
	}
	file := pkg.Prog.Fset.AddFile("test.go", -1, 100)
	pos := file.Pos(50)

	// 使用 NewTypeName 方法创建虚拟的 TypeName 对象
	typeName := types.NewTypeName(pos, nil, "TypeName", types.Typ[types.Int])

	// 调用 SsaPosition 函数
	result := SsaPosition(pkg, typeName)

	// 期望的结果
	expected := &position{
		filename: "test.go", // 假设位置为 example.go
		offset:   50,        // 假设位置偏移为 0
		line:     1,         // 假设该位置在文件的第一行
		column:   51,        // 假设该位置在文件的第一行第一列
	}

	// 检查结果是否与预期相符
	if result.Filename() != expected.filename {
		t.Errorf("Expected filename: %v, got: %v", expected.filename, result.Filename())
	}

	if result.Offset() != expected.offset {
		t.Errorf("Expected offset: %v, got: %v", expected.offset, result.Offset())
	}

	if result.Line() != expected.line {
		t.Errorf("Expected line: %v, got: %v", expected.line, result.Line())
	}

	if result.Column() != expected.column {
		t.Errorf("Expected column: %v, got: %v", expected.column, result.Column())
	}
}

func TestAstPosition(t *testing.T) {
	// 创建一个虚拟的 packages.Package
	pkg := &packages.Package{
		Fset: token.NewFileSet(),
		// 其他属性根据需要设置
	}
	file := pkg.Fset.AddFile("test.go", -1, 100)
	pos := file.Pos(50)

	starExpr := &ast.StarExpr{
		Star: pos,
		X: &ast.Ident{
			NamePos: pos,
			Name:    "x",
		},
	}

	// 调用 AstPosition 函数
	result := AstPosition(pkg, starExpr)
	if result == nil {
		t.Errorf("AstPosition(%v, %v) = nil, expected non-nil", pkg, starExpr)
	}

	// 期望的结果
	expected := &position{
		filename: "test.go", // 假设位置为 example.go
		offset:   50,        // 假设位置偏移为 0
		line:     1,         // 假设该位置在文件的第一行
		column:   51,        // 假设该位置在文件的第一行第一列
	}

	// 检查结果是否与预期相符
	if result.Filename() != expected.filename {
		t.Errorf("Expected filename: %v, got: %v", expected.filename, result.Filename())
	}

	if result.Offset() != expected.offset {
		t.Errorf("Expected offset: %v, got: %v", expected.offset, result.Offset())
	}

	if result.Line() != expected.line {
		t.Errorf("Expected line: %v, got: %v", expected.line, result.Line())
	}

	if result.Column() != expected.column {
		t.Errorf("Expected column: %v, got: %v", expected.column, result.Column())
	}
}
