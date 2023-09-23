package valueobject

import "github.com/dddplayer/dp/internal/domain/arch"

type relationMeta struct {
	*relationPos
	t arch.RelationType
}

func (r *relationMeta) Type() arch.RelationType    { return r.t }
func (r *relationMeta) Position() arch.RelationPos { return r.relationPos }

func NewRelationMeta(t arch.RelationType, fp, tp arch.Position) arch.RelationMeta {
	return &relationMeta{
		relationPos: &relationPos{
			from: fp,
			to:   tp,
		},
		t: t,
	}
}
