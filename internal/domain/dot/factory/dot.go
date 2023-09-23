package factory

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/dot/entity"
	"github.com/dddplayer/dp/internal/domain/dot/valueobject"
	"path"
)

func NewSubGraph(sd arch.SubDiagram) *entity.SubGraph {
	g := &entity.SubGraph{
		Name:      valueobject.PortStr(sd.Name()),
		Label:     sd.Name(),
		Nodes:     []*entity.Node{},
		SubGraphs: []*entity.SubGraph{},
	}

	for _, n := range sd.Nodes() {
		node := &entity.Node{
			ID:      valueobject.PortStr(n.ID()),
			Name:    path.Base(n.Name()),
			BgColor: n.Color(),
			Table:   nil,
		}
		g.Nodes = append(g.Nodes, node)
	}

	return g
}

func NewSummarySubGraph(sd arch.SubDiagram) *entity.SubGraph {
	g := &entity.SubGraph{
		Name:      valueobject.PortStr(sd.Name()),
		Label:     sd.Name(),
		Nodes:     []*entity.Node{},
		SubGraphs: []*entity.SubGraph{},
	}

	node := &entity.Node{
		ID:    valueobject.PortStr(sd.Name()),
		Name:  path.Base(sd.Name()),
		Table: &entity.Table{Rows: []*entity.Row{}},
	}
	g.Nodes = append(g.Nodes, node)

	if err := node.Build(sd.Summary()); err != nil {
		fmt.Println(err)
		return nil
	}

	return g
}
