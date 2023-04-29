package entity

type RelationShip int

const (
	OneOne RelationShip = 1 << iota
	OneMany
)

type Link struct {
	From     *Node
	To       *Node
	Relation RelationShip
}
