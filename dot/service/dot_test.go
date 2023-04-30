package service

import (
	"bytes"
	"github.com/dddplayer/core/dot/entity"
	"strings"
	"testing"
)

func TestWriteDot(t *testing.T) {
	// Create a new instance of MyDotGraph and add some nodes and edges
	graph := NewMyDotGraph()
	graph.AddNode("node1", []entity.DotElement{
		NewMyDotElement("name1", "", ""),
		NewMyDotElement("name2", "", "")})
	graph.AddNode("node2", []entity.DotElement{
		NewMyDotElement("name3", "", "")})
	graph.AddEdge("node1", "node2")

	// Write the Dot representation of the graph to a buffer
	var buf bytes.Buffer
	err := WriteDot(graph, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the Dot representation matches the expected output
	actualOutput := strings.TrimSpace(buf.String())
	if strings.Contains(actualOutput, "node1 [label=<") == false ||
		strings.Contains(actualOutput, "node2 [label=<") == false ||
		strings.Contains(actualOutput, "node1 -> node2") == false {
		t.Errorf("unexpected result. got: %s", actualOutput)
	}
}

// NewMyDotGraph returns a new instance of MyDotGraph
func NewMyDotGraph() *MyDotGraph {
	return &MyDotGraph{
		name:  "MyGraph",
		nodes: make(map[string]entity.DotNode),
		edges: make(map[string]map[string]bool),
	}
}

// MyDotGraph is a simple implementation of the DotGraph interface
type MyDotGraph struct {
	name  string
	nodes map[string]entity.DotNode
	edges map[string]map[string]bool
}

func (g *MyDotGraph) Name() string {
	return g.name
}

func (g *MyDotGraph) Nodes() []entity.DotNode {
	nodes := make([]entity.DotNode, 0, len(g.nodes))
	for _, node := range g.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (g *MyDotGraph) Edges() []entity.DotEdge {
	edges := make([]entity.DotEdge, 0, len(g.edges))
	for from, targets := range g.edges {
		for to := range targets {
			edges = append(edges, NewMyDotEdge(from, to))
		}
	}
	return edges
}

func (g *MyDotGraph) AddNode(name string, elements []entity.DotElement) {
	g.nodes[name] = NewMyDotNode(name, elements)
}

func (g *MyDotGraph) AddEdge(from, to string) {
	if _, ok := g.edges[from]; !ok {
		g.edges[from] = make(map[string]bool)
	}
	g.edges[from][to] = true
}

// NewMyDotNode returns a new instance of MyDotNode
func NewMyDotNode(name string, elements []entity.DotElement) *MyDotNode {
	return &MyDotNode{
		name:     name,
		elements: elements,
	}
}

// MyDotNode is a simple implementation of the DotNode interface
type MyDotNode struct {
	name     string
	elements []entity.DotElement
}

func (n *MyDotNode) Name() string {
	return n.name
}

func (n *MyDotNode) Elements() []entity.DotElement {
	return n.elements
}

// NewMyDotElement returns a new instance of MyDotElement
func NewMyDotElement(name, port, color string) *MyDotElement {
	return &MyDotElement{
		name:  name,
		port:  port,
		color: color,
	}
}

// MyDotElement is a simple implementation of the DotElement interface
type MyDotElement struct {
	name       string
	port       string
	color      string
	attributes []any
}

func (e *MyDotElement) Name() string {
	return e.name
}

func (e *MyDotElement) Port() string {
	return e.port
}

func (e *MyDotElement) Color() string {
	return e.color
}

func (e *MyDotElement) Attributes() []interface{} {
	return e.attributes
}

// NewMyDotEdge returns a new instance of MyDotEdge
func NewMyDotEdge(from, to string) *MyDotEdge {
	return &MyDotEdge{
		from: from,
		to:   to,
	}
}

// MyDotEdge is a simple implementation of the DotEdge interface
type MyDotEdge struct {
	from string
	to   string
}

func (e *MyDotEdge) From() string {
	return e.from
}

func (e *MyDotEdge) To() string {
	return e.to
}
