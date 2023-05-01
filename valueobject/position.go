package valueobject

import "github.com/dddplayer/core/codeanalysis/valueobject"

type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

func NewPosition(id valueobject.Position) Position {
	return Position{
		Filename: id.Filename(),
		Offset:   id.Offset(),
		Line:     id.Line(),
		Column:   id.Column(),
	}
}
