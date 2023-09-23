package factory

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/dot"
)

type DummyDotGraph struct {
	NameVal  string
	NodesVal []arch.Node
	EdgesVal []arch.Edge
}

func (g *DummyDotGraph) Name() string                 { return g.NameVal }
func (g *DummyDotGraph) Nodes() []arch.Node           { return g.NodesVal }
func (g *DummyDotGraph) Edges() []arch.Edge           { return g.EdgesVal }
func (g *DummyDotGraph) Summary() []arch.Element      { return nil }
func (g *DummyDotGraph) SubGraphs() []arch.SubDiagram { return nil }
func (g *DummyDotGraph) Print()                       {}

type DummyNode struct {
	name     string
	elements []arch.Element
}

func (n DummyNode) Name() string             { return n.name }
func (n DummyNode) Elements() []arch.Element { return n.elements }
func (n DummyNode) ID() string               { return n.name }
func (n DummyNode) Color() string            { return "" }

type DummyDotEdge struct {
	FromVal string
	ToVal   string
	L       string
	TT      string
}

func (e *DummyDotEdge) From() string                 { return e.FromVal }
func (e *DummyDotEdge) To() string                   { return e.ToVal }
func (e *DummyDotEdge) Label() string                { return e.L }
func (e *DummyDotEdge) Tooltip() string              { return e.TT }
func (e *DummyDotEdge) ArrowHead() dot.EdgeArrowHead { return dot.EdgeArrowHeadONormal }
func (e *DummyDotEdge) Count() int                   { return 1 }
func (e *DummyDotEdge) Pos() []arch.RelationPos      { return []arch.RelationPos{} }
func (e *DummyDotEdge) Type() arch.RelationType      { return arch.RelationTypeNone }

type DummyDotElement struct {
	NameVal       string
	PortVal       string
	ColorVal      string
	AttributesVal []arch.Nodes
}

func (e *DummyDotElement) Name() string           { return e.NameVal }
func (e *DummyDotElement) ID() string             { return e.PortVal }
func (e *DummyDotElement) Color() string          { return e.ColorVal }
func (e *DummyDotElement) Children() []arch.Nodes { return e.AttributesVal }

type DummyDotAttribute struct {
	name  string
	port  string
	color string
}

func (a DummyDotAttribute) Name() string  { return a.name }
func (a DummyDotAttribute) ID() string    { return a.port }
func (a DummyDotAttribute) Color() string { return a.color }

func generateDummyDotGraph() *DummyDotGraph {
	node1 := &DummyNode{
		name: "Node1",
		elements: []arch.Element{
			&DummyDotElement{
				NameVal:  "Elem1",
				PortVal:  "ElePort1",
				ColorVal: "red",
				AttributesVal: []arch.Nodes{
					[]arch.Node{
						&DummyDotAttribute{name: "attr1", port: "port1", color: "green"},
						&DummyDotAttribute{name: "attr2", port: "port2", color: "blue"},
					},
					[]arch.Node{
						&DummyDotAttribute{name: "attr3", port: "port3", color: "green"},
						&DummyDotAttribute{name: "attr4", port: "port4", color: "blue"},
					},
				},
			},
			&DummyDotElement{
				NameVal:       "Elem10",
				PortVal:       "port1",
				ColorVal:      "blue",
				AttributesVal: []arch.Nodes{},
			},
			&DummyDotElement{
				NameVal:  "Elem2",
				PortVal:  "port1",
				ColorVal: "blue",
				AttributesVal: []arch.Nodes{
					[]arch.Node{
						&DummyDotAttribute{name: "attr5", port: "port5", color: "green"},
						&DummyDotAttribute{name: "attr6", port: "port6", color: "blue"}},
				},
			},
		},
	}
	node2 := &DummyNode{
		name: "Node2",
		elements: []arch.Element{
			&DummyDotElement{
				NameVal:  "Elem3",
				PortVal:  "port1",
				ColorVal: "blue",
				AttributesVal: []arch.Nodes{
					[]arch.Node{
						&DummyDotAttribute{name: "attr5", port: "port5", color: "green"},
						&DummyDotAttribute{name: "attr6", port: "port6", color: "blue"},
					},
				},
			},
			&DummyDotElement{
				NameVal:  "Elem4",
				PortVal:  "port2",
				ColorVal: "red",
				AttributesVal: []arch.Nodes{
					[]arch.Node{
						&DummyDotAttribute{name: "attr7", port: "port7", color: "green"},
					},
				},
			},
		},
	}
	graph := &DummyDotGraph{
		NameVal: "TestGraph",
		NodesVal: []arch.Node{
			node1,
			node2,
		},
		EdgesVal: []arch.Edge{
			&DummyDotEdge{FromVal: "Node1:port1", ToVal: "Node1:port5"},
			&DummyDotEdge{FromVal: "Node1:port2", ToVal: "Node1:port6"},
			&DummyDotEdge{FromVal: "Node2:port5", ToVal: "Node2:port7"},
		},
	}
	return graph
}
