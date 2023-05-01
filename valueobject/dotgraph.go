package valueobject

import "github.com/dddplayer/core/dot/entity"

func NewDotGraph(name string) *DotGraph {
	return &DotGraph{
		name:  name,
		nodes: []entity.DotNode{},
		edges: []entity.DotEdge{},
	}
}

func NewDotEdge(from, to string) *DotEdge {
	return &DotEdge{
		from: from,
		to:   to,
	}
}

func NewNode(name string) *DotNode {
	return &DotNode{
		name:     name,
		elements: []entity.DotElement{},
	}
}

func NewElement(id, name string) *DotElement {
	return &DotElement{
		id:   id,
		name: name,
	}
}

func NewAttribute(id, name, color string) *DotAttribute {
	return &DotAttribute{
		id:    id,
		name:  name,
		color: color,
	}
}

type (
	DotGraph struct {
		name  string
		nodes []entity.DotNode
		edges []entity.DotEdge
	}

	DotNode struct {
		name     string
		elements []entity.DotElement
	}

	DotElement struct {
		id    string
		name  string
		color string
		attrs []any
	}

	DotAttribute struct {
		id    string
		name  string
		color string
	}

	DotEdge struct {
		from string
		to   string
	}
)

func (dg *DotGraph) Name() string                { return dg.name }
func (dg *DotGraph) Nodes() []entity.DotNode     { return dg.nodes }
func (dg *DotGraph) Edges() []entity.DotEdge     { return dg.edges }
func (dg *DotGraph) AppendNode(n entity.DotNode) { dg.nodes = append(dg.nodes, n) }
func (dg *DotGraph) AppendEdge(e entity.DotEdge) { dg.edges = append(dg.edges, e) }

func (dn *DotNode) Name() string                  { return dn.name }
func (dn *DotNode) Elements() []entity.DotElement { return dn.elements }
func (dn *DotNode) Append(es ...entity.DotElement) {
	for _, e := range es {
		dn.elements = append(dn.elements, e)
	}
}

func (de *DotElement) Attributes() []any { return de.attrs }
func (de *DotElement) Name() string      { return de.name }
func (de *DotElement) Port() string      { return de.id }
func (de *DotElement) Color() string     { return de.color }
func (de *DotElement) SetColor(c string) { de.color = c }
func (de *DotElement) Append(attr any)   { de.attrs = append(de.attrs, attr) }

func (da *DotAttribute) Name() string  { return da.name }
func (da *DotAttribute) Color() string { return da.color }
func (da *DotAttribute) Port() string  { return da.id }

func (de *DotEdge) From() string { return de.from }
func (de *DotEdge) To() string   { return de.to }
