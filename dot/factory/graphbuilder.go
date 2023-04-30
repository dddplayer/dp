package factory

import (
	"fmt"
	"github.com/dddplayer/core/datastructure/intimacy"
	"github.com/dddplayer/core/dot/entity"
	"github.com/dddplayer/core/dot/valueobject"
	"golang.org/x/exp/slices"
	"sort"
	"strings"
)

func NewGraphBuilder(dg entity.DotGraph) *GraphBuilder {
	return &GraphBuilder{
		dotGraph: dg,
	}
}

type GraphBuilder struct {
	dotGraph entity.DotGraph
}

func (gb *GraphBuilder) Build() (*entity.Digraph, error) {
	digraph := &entity.Digraph{
		Name:  gb.dotGraph.Name(),
		Label: fmt.Sprintf("\n\n%s\nDomain Model\n\nPowered by DDD Player", gb.dotGraph.Name()),
		Nodes: []*entity.Node{},
		Edges: []*entity.Edge{},
	}
	for _, n := range gb.dotGraph.Nodes() {
		if n := gb.buildNode(n); n != nil {
			digraph.Nodes = append(digraph.Nodes, n)
		}
	}
	for _, e := range gb.dotGraph.Edges() {
		if e := gb.buildEdge(e); e != nil {
			digraph.Edges = append(digraph.Edges, e)
		}
	}

	return digraph, nil
}

func portStr(name string) string {
	ps := name
	for _, old := range []string{".", "-"} {
		ps = strings.ReplaceAll(ps, old, entity.DotJoiner)
	}
	return ps
}

func (gb *GraphBuilder) buildEdge(e entity.DotEdge) *entity.Edge {
	return &entity.Edge{
		From: portStr(e.From()),
		To:   portStr(e.To()),
	}
}

func (gb *GraphBuilder) buildNode(n entity.DotNode) *entity.Node {
	if sortedEls := gb.getNodeSortedElements(n); sortedEls != nil {

		n := &entity.Node{
			Name: portStr(n.Name()),
			Rows: []*entity.Row{},
		}
		if err := n.Build(sortedEls); err != nil {
			fmt.Println(err)
			return nil
		}

		return n
	}
	return nil
}

func (gb *GraphBuilder) getNodeSortedElements(n entity.DotNode) entity.DotElements {
	els := entity.DotElements(n.Elements())
	return gb.getNodeSortedElementsRecursive(n, els)
}

func (gb *GraphBuilder) getNodeSortedElementsRecursive(n entity.DotNode, els entity.DotElements) entity.DotElements {
	if len(els) == 0 {
		return nil
	}

	if len(els) == 1 {
		return els
	}

	edges := gb.getNodeInternalEdges(n)
	ig := gb.buildElementsIntimacy(n.Name(), edges, els)
	sortedEls := gb.sortElementsWithIntimacy(els, ig)
	sorted := entity.DotElements{}
	if len(sortedEls) > 1 {
		sorted = append(sorted, sortedEls.First())

		rest := sortedEls[1:]
		sortedEls = gb.getNodeSortedElementsRecursive(n, rest)
		for _, se := range sortedEls {
			sorted = append(sorted, se)
		}
		return sorted
	}

	return nil
}

func (gb *GraphBuilder) sortElementsWithIntimacy(els entity.DotElements, ig *intimacy.Graph) entity.DotElements {
	sorted := entity.DotElements{}

	intimacyEls := valueobject.IntimacyElements{}
	firstEle := els.First()
	if firstEle != nil {
		sorted = append(sorted, firstEle)

		for i := 1; i < len(els); i++ {
			intimacyEls = append(intimacyEls, &valueobject.IntimacyElement{
				Key: i,
				Val: ig.Intimacy(els[i].Name(), firstEle.Name()),
			})
		}
		sort.Sort(intimacyEls)
		for _, ie := range intimacyEls {
			sorted = append(sorted, els[ie.Key])
		}
	}

	return sorted
}

func (gb *GraphBuilder) getNodeInternalEdges(n entity.DotNode) []entity.DotEdge {
	var edges []entity.DotEdge
	for _, e := range gb.dotGraph.Edges() {
		if strings.HasPrefix(e.From(), n.Name()) && strings.HasPrefix(e.To(), n.Name()) {
			edges = append(edges, e)
		}
	}
	return edges
}

func (gb *GraphBuilder) buildElementsIntimacy(nodeName string, edges []entity.DotEdge, els entity.DotElements) *intimacy.Graph {
	firstEle := els.First()
	if firstEle != nil {
		firstElePorts := gb.getElementPorts(nodeName, firstEle)

		iG := intimacy.NewGraph()
		for i := 1; i < len(els); i++ {
			secondEle := els[i]
			secondElePorts := gb.getElementPorts(nodeName, secondEle)

			for _, e := range edges {
				if e.From() == e.To() {
					continue
				}

				if (slices.Contains(firstElePorts, e.From()) && slices.Contains(secondElePorts, e.To())) ||
					(slices.Contains(firstElePorts, e.To()) && slices.Contains(secondElePorts, e.From())) {
					if err := iG.IntimacyPlusOne(firstEle.Name(), secondEle.Name()); err != nil {
						return nil
					}
				}
			}
		}
		return iG
	}

	return nil
}

func (gb *GraphBuilder) buildNodeElementsIntimacy(n entity.DotNode) *intimacy.Graph {
	edges := gb.getNodeInternalEdges(n)
	els := entity.DotElements(n.Elements())
	return gb.buildElementsIntimacy(n.Name(), edges, els)
}

func (gb *GraphBuilder) getElementPorts(nodeName string, e entity.DotElement) []string {
	var ports []string
	ports = append(ports, entity.DotAttrPort(nodeName, e.Port()))

	for _, a := range e.Attributes() {
		switch as := a.(type) {
		case []entity.DotAttribute:
			for _, a := range as {
				ports = append(ports, entity.DotAttrPort(nodeName, a.Port()))
			}
		case entity.DotAttribute:
			ports = append(ports, entity.DotAttrPort(nodeName, as.Port()))
		}
	}

	return ports
}
