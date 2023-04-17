package entity

import "github.com/dddplayer/core/dot/entity"

type dotGraph struct {
	name  string
	nodes []entity.DotNode
}

func (dg *dotGraph) Name() string {
	return dg.name
}

func (dg *dotGraph) Nodes() []entity.DotNode {
	return dg.nodes
}

type dotNode struct {
	name string
}

func (dn *dotNode) Name() string {
	return dn.name
}
