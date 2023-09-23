package directed

import (
	"reflect"
	"testing"
)

func TestGraph(t *testing.T) {
	graph := NewDirectedGraph()

	// Create test nodes
	nodeA := &Node{Value: "A"}
	nodeB := &Node{Value: "B"}
	nodeC := &Node{Value: "C"}

	// Add nodes to directed
	_ = graph.AddNode("A", nodeA)
	_ = graph.AddNode("B", nodeB)
	_ = graph.AddNode("C", nodeC)

	// Create test edges
	edgeAB := &Edge{From: nodeA, To: nodeB, Value: "Edge A-B"}
	edgeBC := &Edge{From: nodeB, To: nodeC, Value: "Edge B-C"}

	// Add edges to directed
	_ = graph.AddEdge("A", "B", "T1", edgeAB.Value)
	_ = graph.AddEdge("B", "C", "T2", edgeBC.Value)

	// Test node count
	expectedNodeCount := 3
	actualNodeCount := len(graph.Nodes)
	if actualNodeCount != expectedNodeCount {
		t.Errorf("Unexpected node count. Expected: %d, Actual: %d", expectedNodeCount, actualNodeCount)
	}

	// Test edge count for each node
	tests := []struct {
		node       *Node
		expected   int
		edgeValues []interface{}
	}{
		{graph.Nodes[0], 1, []interface{}{"Edge A-B"}},
		{graph.Nodes[1], 1, []interface{}{"Edge B-C"}},
		{graph.Nodes[2], 0, []interface{}{}},
	}

	for _, test := range tests {
		actualEdgeCount := len(test.node.Edges)
		if actualEdgeCount != test.expected {
			t.Errorf("Unexpected edge count for node %v. Expected: %d, Actual: %d", test.node.Value, test.expected, actualEdgeCount)
		}

		for i, edge := range test.node.Edges {
			expectedValue := test.edgeValues[i]
			actualValue := edge.Value
			if !reflect.DeepEqual(actualValue, expectedValue) {
				t.Errorf("Unexpected edge value for node %v. Expected: %v, Actual: %v", test.node.Value, expectedValue, actualValue)
			}
		}
	}
}
