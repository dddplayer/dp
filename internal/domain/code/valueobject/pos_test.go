package valueobject

import (
	"go/token"
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
	pos := file.Pos(50)
	expected := &position{
		filename: "test.go",
		offset:   50,
		line:     1,
		column:   51,
	}

	result := getPosition(fset, pos)

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
