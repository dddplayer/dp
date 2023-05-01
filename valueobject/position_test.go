package valueobject

import (
	"testing"
)

func TestNewPosition(t *testing.T) {
	filename := "test.go"
	offset := 100
	line := 5
	column := 10

	vp := &PositionMock{
		filename: filename,
		offset:   offset,
		line:     line,
		column:   column,
	}

	p := NewPosition(vp)

	if p.Filename != filename {
		t.Errorf("Expected filename '%s', but got '%s'", filename, p.Filename)
	}

	if p.Offset != offset {
		t.Errorf("Expected offset '%d', but got '%d'", offset, p.Offset)
	}

	if p.Line != line {
		t.Errorf("Expected line '%d', but got '%d'", line, p.Line)
	}

	if p.Column != column {
		t.Errorf("Expected column '%d', but got '%d'", column, p.Column)
	}
}

type PositionMock struct {
	filename string
	offset   int
	line     int
	column   int
}

func (p *PositionMock) Filename() string { return p.filename }
func (p *PositionMock) Offset() int      { return p.offset }
func (p *PositionMock) Line() int        { return p.line }
func (p *PositionMock) Column() int      { return p.column }
