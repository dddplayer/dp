package valueobject

import "fmt"

type Relation struct {
	From *Identifier
	To   *Identifier
	Pos  *Position
	Ship RelationShip
	Type RelationType
}

func (r *Relation) Identifier() *Identifier {
	return &Identifier{
		Name: fmt.Sprintf("%s->%s", r.From.Name, r.To.Name),
		Path: fmt.Sprintf("%s->%s", r.From.Path, r.To.Path),
	}
}

func (r *Relation) Position() *Position {
	return r.Pos
}

type RelationShip int

const (
	OneOne RelationShip = 1 << iota
	OneMany
)

type RelationType int

const (
	TypeRefer RelationType = 1 << iota
	TypeImplements
	TypeCall
)
