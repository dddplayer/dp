package valueobject

import "github.com/dddplayer/dp/internal/domain/arch"

type relationPos struct {
	from arch.Position
	to   arch.Position
}

func (rp *relationPos) From() arch.Position { return rp.from }
func (rp *relationPos) To() arch.Position   { return rp.to }

func NewRelationPos(f, t arch.Position) arch.RelationPos {
	return &relationPos{
		from: f,
		to:   t,
	}
}

func NewEmptyRelationPos() arch.RelationPos {
	return &relationPos{}
}
