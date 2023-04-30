package factory

import (
	"github.com/dddplayer/core/dot/entity"
	"reflect"
	"testing"
)

type DummyDotGraph struct {
	NameVal  string
	NodesVal []entity.DotNode
	EdgesVal []entity.DotEdge
}

func (g *DummyDotGraph) Name() string {
	return g.NameVal
}

func (g *DummyDotGraph) Nodes() []entity.DotNode {
	return g.NodesVal
}

func (g *DummyDotGraph) Edges() []entity.DotEdge {
	return g.EdgesVal
}

type DummyNode struct {
	name     string
	elements []entity.DotElement
}

func (n DummyNode) Name() string {
	return n.name
}

func (n DummyNode) Elements() []entity.DotElement {
	return n.elements
}

type DummyDotEdge struct {
	FromVal string
	ToVal   string
}

func (e *DummyDotEdge) From() string {
	return e.FromVal
}

func (e *DummyDotEdge) To() string {
	return e.ToVal
}

type DummyDotElement struct {
	NameVal       string
	PortVal       string
	ColorVal      string
	AttributesVal []any
}

func (e *DummyDotElement) Name() string {
	return e.NameVal
}

func (e *DummyDotElement) Port() string {
	return e.PortVal
}

func (e *DummyDotElement) Color() string {
	return e.ColorVal
}

func (e *DummyDotElement) Attributes() []interface{} {
	return e.AttributesVal
}

type DummyDotAttribute struct {
	name  string
	port  string
	color string
}

func (a DummyDotAttribute) Name() string {
	return a.name
}

func (a DummyDotAttribute) Port() string {
	return a.port
}

func (a DummyDotAttribute) Color() string {
	return a.color
}

func TestGraphBuilder_BuildEdge(t *testing.T) {
	// 设置测试用例
	testCases := []struct {
		name         string
		dotEdge      entity.DotEdge
		expectedEdge *entity.Edge
	}{
		{
			name: "Test case 1",
			dotEdge: &DummyDotEdge{
				FromVal: "node1",
				ToVal:   "node2",
			},
			expectedEdge: &entity.Edge{
				From: "node1",
				To:   "node2",
			},
		},
		{
			name: "Test case 2",
			dotEdge: &DummyDotEdge{
				FromVal: "node3",
				ToVal:   "node4",
			},
			expectedEdge: &entity.Edge{
				From: "node3",
				To:   "node4",
			},
		},
		// 添加更多的测试用例...
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化 GraphBuilder
			gb := &GraphBuilder{
				dotGraph: &DummyDotGraph{},
			}

			// 调用 buildEdge 函数，得到实际的输出
			actualEdge := gb.buildEdge(tc.dotEdge)

			// 对比实际输出和期望输出
			if !reflect.DeepEqual(actualEdge, tc.expectedEdge) {
				t.Errorf("Unexpected result:\nexpected: %v\nactual: %v\n", tc.expectedEdge, actualEdge)
			}
		})
	}
}
func TestGraphBuilder_getElementPorts(t *testing.T) {
	graph := generateDummyDotGraph()
	expectedPorts := []string{
		"Node1:ElePort1",
		"Node1:port1",
		"Node1:port2",
		"Node1:port3",
		"Node1:port4",
	}

	// Test getElementPorts() with dummy element
	gb := GraphBuilder{}
	ports := gb.getElementPorts(graph.NodesVal[0].Name(), graph.NodesVal[0].Elements()[0])

	// Compare expected and actual result
	if !reflect.DeepEqual(ports, expectedPorts) {
		t.Errorf("Unexpected result: got %v, want %v", ports, expectedPorts)
	}
}

func TestGetNodeInternalEdges(t *testing.T) {
	graph := generateDummyDotGraph()
	builder := NewGraphBuilder(graph)

	// Test edges between Node1 elements
	node1 := graph.NodesVal[0]
	expectedEdges := []entity.DotEdge{
		&DummyDotEdge{FromVal: "Node1:port1", ToVal: "Node1:port5"},
		&DummyDotEdge{FromVal: "Node1:port2", ToVal: "Node1:port6"},
	}
	internalEdges := builder.getNodeInternalEdges(node1)
	if len(internalEdges) != len(expectedEdges) {
		t.Errorf("Expected %d internal edges, but got %d", len(expectedEdges), len(internalEdges))
	}
	for i, edge := range internalEdges {
		if edge.From() != expectedEdges[i].From() {
			t.Errorf("Expected edge %d to have From value %s, but got %s", i, expectedEdges[i].From(), edge.From())
		}
		if edge.To() != expectedEdges[i].To() {
			t.Errorf("Expected edge %d to have To value %s, but got %s", i, expectedEdges[i].To(), edge.To())
		}
	}

	// Test edges between Node2 elements
	node2 := graph.NodesVal[1]
	expectedEdges = []entity.DotEdge{
		&DummyDotEdge{FromVal: "Node2:port5", ToVal: "Node2:port7"},
	}
	internalEdges = builder.getNodeInternalEdges(node2)
	if len(internalEdges) != len(expectedEdges) {
		t.Errorf("Expected %d internal edges, but got %d", len(expectedEdges), len(internalEdges))
	}
	for i, edge := range internalEdges {
		if edge.From() != expectedEdges[i].From() {
			t.Errorf("Expected edge %d to have From value %s, but got %s", i, expectedEdges[i].From(), edge.From())
		}
		if edge.To() != expectedEdges[i].To() {
			t.Errorf("Expected edge %d to have To value %s, but got %s", i, expectedEdges[i].To(), edge.To())
		}
	}

	// Test edges for a node with no internal edges
	node3 := &DummyNode{
		name:     "Node3",
		elements: []entity.DotElement{},
	}
	internalEdges = builder.getNodeInternalEdges(node3)
	if len(internalEdges) != 0 {
		t.Errorf("Expected 0 internal edges for Node3, but got %d", len(internalEdges))
	}
}

func TestGraphBuilder_buildElementsIntimacy(t *testing.T) {
	graph := generateDummyDotGraph()
	gb := GraphBuilder{graph}

	g := gb.buildNodeElementsIntimacy(graph.NodesVal[0])

	if g == nil {
		t.Errorf("buildNodeElementsIntimacy() returned nil")
		return
	}

	if intimacy := g.Intimacy("Elem1", "Elem2"); intimacy != 2 {
		t.Errorf("buildNodeElementsIntimacy() returned a wrong intimacy: %f for Elem1 and Elem2, expected 2", intimacy)
	}
}

func TestGraphBuilder_SortElementsWithIntimacy(t *testing.T) {
	graph := generateDummyDotGraph()
	gb := GraphBuilder{graph}
	g := gb.buildNodeElementsIntimacy(graph.NodesVal[0])
	got := gb.sortElementsWithIntimacy(graph.NodesVal[0].Elements(), g)
	if len(got) != 3 {
		t.Errorf("Expected 3 elements from Node1, but got %d", len(got))
	}
	if got[0].Name() != "Elem1" {
		t.Errorf("Expected ele name is Elem1, but got %s", got[0].Name())
	}
	if got[1].Name() != "Elem2" {
		t.Errorf("Expected ele name is Elem2, but got %s", got[1].Name())
	}
	if got[2].Name() != "Elem10" {
		t.Errorf("Expected ele name is Elem10, but got %s", got[2].Name())
	}
}

func TestGraphBuilder_BuildNode(t *testing.T) {
	graph := generateDummyDotGraph()
	gb := GraphBuilder{graph}

	node1 := graph.NodesVal[0]
	got := gb.buildNode(node1)

	// Test if node is not nil
	if got == nil {
		t.Error("Expected a non-nil node, but got nil")
	}

	// Test if node name is correct
	if got.Name != node1.Name() {
		t.Errorf("Expected node name %q, but got %q", node1.Name(), got.Name)
	}
}

func generateDummyDotGraph() *DummyDotGraph {
	node1 := &DummyNode{
		name: "Node1",
		elements: []entity.DotElement{
			&DummyDotElement{
				NameVal:  "Elem1",
				PortVal:  "ElePort1",
				ColorVal: "red",
				AttributesVal: []interface{}{
					[]entity.DotAttribute{
						&DummyDotAttribute{name: "attr1", port: "port1", color: "green"},
						&DummyDotAttribute{name: "attr2", port: "port2", color: "blue"},
					},
					[]entity.DotAttribute{
						&DummyDotAttribute{name: "attr3", port: "port3", color: "green"},
						&DummyDotAttribute{name: "attr4", port: "port4", color: "blue"},
					},
				},
			},
			&DummyDotElement{
				NameVal:       "Elem10",
				PortVal:       "port1",
				ColorVal:      "blue",
				AttributesVal: []interface{}{},
			},
			&DummyDotElement{
				NameVal:  "Elem2",
				PortVal:  "port1",
				ColorVal: "blue",
				AttributesVal: []interface{}{
					&DummyDotAttribute{name: "attr5", port: "port5", color: "green"},
					&DummyDotAttribute{name: "attr6", port: "port6", color: "blue"},
				},
			},
		},
	}
	node2 := &DummyNode{
		name: "Node2",
		elements: []entity.DotElement{
			&DummyDotElement{
				NameVal:  "Elem3",
				PortVal:  "port1",
				ColorVal: "blue",
				AttributesVal: []interface{}{
					&DummyDotAttribute{name: "attr5", port: "port5", color: "green"},
					&DummyDotAttribute{name: "attr6", port: "port6", color: "blue"},
				},
			},
			&DummyDotElement{
				NameVal:  "Elem4",
				PortVal:  "port2",
				ColorVal: "red",
				AttributesVal: []interface{}{
					&DummyDotAttribute{name: "attr7", port: "port7", color: "green"},
				},
			},
		},
	}
	graph := &DummyDotGraph{
		NameVal: "TestGraph",
		NodesVal: []entity.DotNode{
			node1,
			node2,
		},
		EdgesVal: []entity.DotEdge{
			&DummyDotEdge{FromVal: "Node1:port1", ToVal: "Node1:port5"},
			&DummyDotEdge{FromVal: "Node1:port2", ToVal: "Node1:port6"},
			&DummyDotEdge{FromVal: "Node2:port5", ToVal: "Node2:port7"},
		},
	}
	return graph
}
