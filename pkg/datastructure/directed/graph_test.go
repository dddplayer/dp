package directed

import (
	"fmt"
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

func TestNewNode(t *testing.T) {
	// 创建一个新节点
	key := "testKey"
	value := "testValue"
	node := NewNode(key, value)

	// 验证节点的属性是否正确设置
	if node.Key != key {
		t.Errorf("Expected node.Key to be %s, but got %s", key, node.Key)
	}

	if node.Value != value {
		t.Errorf("Expected node.Value to be %v, but got %v", value, node.Value)
	}

	// 验证节点的边列表是否为空
	if len(node.Edges) != 0 {
		t.Errorf("Expected an empty edge list, but got %d edges", len(node.Edges))
	}
}

func TestGraph_AddNode(t *testing.T) {
	// 创建一个新的图
	graph := &Graph{Nodes: []*Node{}}

	// 添加一个新节点
	key1 := "node1"
	value1 := "value1"
	err := graph.AddNode(key1, value1)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 验证节点是否被正确添加
	if len(graph.Nodes) != 1 {
		t.Errorf("Expected 1 node in the graph, but got %d", len(graph.Nodes))
	}
	if graph.Nodes[0].Key != key1 {
		t.Errorf("Expected node key to be %s, but got %s", key1, graph.Nodes[0].Key)
	}
	if graph.Nodes[0].Value != value1 {
		t.Errorf("Expected node value to be %v, but got %v", value1, graph.Nodes[0].Value)
	}

	// 尝试添加一个具有冲突键的节点
	key2 := "node1" // 与上面的节点键冲突
	value2 := "value2"
	err = graph.AddNode(key2, value2)
	if err == nil {
		t.Error("Expected an error for adding a node with a conflicting key, but got no error")
	} else if err.Error() != fmt.Sprintf("key conflict: %s", key2) {
		t.Errorf("Expected error message 'key conflict: %s', but got: %v", key2, err)
	}

	// 验证图中的节点数量没有增加
	if len(graph.Nodes) != 1 {
		t.Errorf("Expected 1 node in the graph, but got %d", len(graph.Nodes))
	}
}

func TestGraph_FindNodeByKey(t *testing.T) {
	// 创建一个新的图
	graph := &Graph{Nodes: []*Node{}}

	// 添加一些节点到图中
	node1 := &Node{Key: "node1", Value: "value1", Edges: []*Edge{}}
	node2 := &Node{Key: "node2", Value: "value2", Edges: []*Edge{}}
	node3 := &Node{Key: "node3", Value: "value3", Edges: []*Edge{}}
	graph.Nodes = append(graph.Nodes, node1, node2, node3)

	// 查找一个存在的节点键
	keyToFind := "node2"
	foundNode := graph.FindNodeByKey(keyToFind)
	if foundNode == nil {
		t.Errorf("Expected to find a node with key %s, but got nil", keyToFind)
	}
	if foundNode.Key != keyToFind {
		t.Errorf("Expected node key to be %s, but got %s", keyToFind, foundNode.Key)
	}

	// 查找一个不存在的节点键
	keyToFind = "nonexistent"
	foundNode = graph.FindNodeByKey(keyToFind)
	if foundNode != nil {
		t.Errorf("Expected to find nil for nonexistent key, but got a node with key %s", keyToFind)
	}
}

func TestGraph_AddEdge(t *testing.T) {
	// 创建一个新的图
	graph := &Graph{Nodes: []*Node{}}

	// 添加一些节点到图中
	node1 := &Node{Key: "node1", Value: "value1", Edges: []*Edge{}}
	node2 := &Node{Key: "node2", Value: "value2", Edges: []*Edge{}}
	graph.Nodes = append(graph.Nodes, node1, node2)

	// 添加一条边，确保它被正确添加到图中
	from := "node1"
	to := "node2"
	edgeType := "edgeType"
	edgeValue := "edgeValue"
	err := graph.AddEdge(from, to, edgeType, edgeValue)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 验证边是否被正确添加到起始节点
	if len(node1.Edges) != 1 {
		t.Errorf("Expected 1 edge for node1, but got %d", len(node1.Edges))
	}
	addedEdge := node1.Edges[0]
	if addedEdge.From != node1 || addedEdge.To != node2 || addedEdge.Type != edgeType || addedEdge.Value != edgeValue {
		t.Errorf("Added edge does not match the expected values")
	}

	// 尝试添加一条边，其中起始节点不存在
	from = "nonexistent"
	to = "node2"
	err = graph.AddEdge(from, to, edgeType, edgeValue)
	if err == nil {
		t.Error("Expected an error for adding an edge with a nonexistent 'from' node, but got no error")
	} else if err.Error() != fmt.Sprintf("from node: %s not found in Digraph", from) {
		t.Errorf("Expected error message 'from node: %s not found in Digraph', but got: %v", from, err)
	}

	// 尝试添加一条边，其中目标节点不存在
	from = "node1"
	to = "nonexistent"
	err = graph.AddEdge(from, to, edgeType, edgeValue)
	if err == nil {
		t.Error("Expected an error for adding an edge with a nonexistent 'to' node, but got no error")
	} else if err.Error() != fmt.Sprintf("to node: %s not found in Digraph", to) {
		t.Errorf("Expected error message 'to node: %s not found in Digraph', but got: %v", to, err)
	}
}

func TestGraph_FindPathsToPrefix(t *testing.T) {
	graph := NewDirectedGraph()

	_ = graph.AddNode("A", nil)
	_ = graph.AddNode("B", nil)
	_ = graph.AddNode("C", nil)
	_ = graph.AddNode("D1", nil)
	_ = graph.AddNode("D2", nil)
	_ = graph.AddNode("E", nil)

	_ = graph.AddEdge("A", "B", nil, nil)
	_ = graph.AddEdge("A", "C", nil, nil)
	_ = graph.AddEdge("B", "D1", nil, nil)
	_ = graph.AddEdge("C", "D2", nil, nil)
	_ = graph.AddEdge("D1", "E", nil, nil)
	_ = graph.AddEdge("D2", "E", nil, nil)

	startKey := "A"
	endKeyPrefix := "D"

	paths := graph.FindPathsToPrefix(startKey, endKeyPrefix)
	if paths != nil {
		t.Logf("Paths from %s to keys with prefix %s:", startKey, endKeyPrefix)
		for _, path := range paths {
			t.Logf("Path: ")
			for _, node := range path {
				t.Logf("%s -> ", node.Key)
			}
			t.Logf("\n")
		}
	} else {
		t.Errorf("No paths found from %s to keys with prefix %s.", startKey, endKeyPrefix)
	}
}

func TestFindAllPathsToPrefix(t *testing.T) {
	graph := NewDirectedGraph()

	_ = graph.AddNode("A", nil)
	_ = graph.AddNode("B", nil)
	_ = graph.AddNode("C", nil)
	_ = graph.AddNode("D1", nil)
	_ = graph.AddNode("D2", nil)
	_ = graph.AddNode("E", nil)

	_ = graph.AddEdge("A", "B", nil, nil)
	_ = graph.AddEdge("A", "C", nil, nil)
	_ = graph.AddEdge("B", "D1", nil, nil)
	_ = graph.AddEdge("C", "D2", nil, nil)
	_ = graph.AddEdge("D1", "E", nil, nil)
	_ = graph.AddEdge("D2", "E", nil, nil)

	startNode := graph.FindNodeByKey("A")
	endKeyPrefix := "D"

	var paths [][]*Node
	currentPath := []*Node{startNode}
	visited := make(map[*Node]bool)

	graph.findAllPathsToPrefix(startNode, endKeyPrefix, &paths, currentPath, visited)

	if paths != nil {
		t.Logf("Paths to keys with prefix %s:", endKeyPrefix)
		for _, path := range paths {
			t.Logf("Path: ")
			for _, node := range path {
				t.Logf("%s -> ", node.Key)
			}
			t.Logf("\n")
		}
	} else {
		t.Errorf("No paths found to keys with prefix %s.", endKeyPrefix)
	}
}
