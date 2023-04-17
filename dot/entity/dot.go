package entity

type DotNode interface {
	Name() string
}

type DotGraph interface {
	Name() string
	Nodes() []DotNode
}
