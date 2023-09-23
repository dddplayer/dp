package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/code"
)

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
func (p pos) Valid() bool       { return p.line != -1 }

func newPosition(id code.Position) *pos {
	return &pos{
		filename: id.Filename(),
		offset:   id.Offset(),
		line:     id.Line(),
		column:   id.Column(),
	}
}

func emptyPosition() *pos {
	return &pos{
		filename: "",
		offset:   0,
		line:     -1,
		column:   0,
	}
}
