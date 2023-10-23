package directed

import (
	"fmt"
	"strings"
)

type Edge struct {
	From  *Node
	To    *Node
	Type  any
	Value any
}

type Node struct {
	Key   string
	Value any
	Edges []*Edge
}

type Graph struct {
	Nodes []*Node
}

func NewDirectedGraph() *Graph {
	return &Graph{
		Nodes: []*Node{},
	}
}

func NewNode(k string, v any) *Node {
	return &Node{
		Key:   k,
		Value: v,
		Edges: []*Edge{},
	}
}

func (g *Graph) AddNode(key string, value any) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	for _, node := range g.Nodes {
		if node.Key == key {
			return fmt.Errorf("key conflict: %s", key)
		}
	}

	newNode := &Node{
		Key:   key,
		Value: value,
		Edges: []*Edge{},
	}

	g.Nodes = append(g.Nodes, newNode)
	return nil
}

func (g *Graph) FindNodeByKey(key string) *Node {
	for _, node := range g.Nodes {
		if node.Key == key {
			return node
		}
	}
	return nil
}

func (g *Graph) AddEdge(from string, to string, t any, value any) error {
	fromNode := g.FindNodeByKey(from)
	toNode := g.FindNodeByKey(to)
	if fromNode == nil {
		return fmt.Errorf("from node: %s not found in Digraph", from)
	}
	if toNode == nil {
		return fmt.Errorf("to node: %s not found in Digraph", to)
	}

	edge := &Edge{
		From:  fromNode,
		To:    toNode,
		Type:  t,
		Value: value,
	}
	fromNode.Edges = append(fromNode.Edges, edge)

	return nil
}

func (g *Graph) FindPathsToPrefix(startKey, endKeyPrefix string) [][]*Node {
	startNode := g.FindNodeByKey(startKey)

	if startNode == nil {
		return nil // 无效的起点
	}

	var paths [][]*Node
	currentPath := []*Node{startNode}

	visited := make(map[*Node]bool)
	g.findAllPathsToPrefix(startNode, endKeyPrefix, &paths, currentPath, visited)

	return paths
}

func (g *Graph) findAllPathsToPrefix(node *Node, endKeyPrefix string, paths *[][]*Node, currentPath []*Node, visited map[*Node]bool) {
	if strings.HasPrefix(node.Key, endKeyPrefix) {
		// 找到一条路径，将其添加到结果中
		*paths = append(*paths, append([]*Node(nil), currentPath...))
		return
	}

	visited[node] = true
	for _, edge := range node.Edges {
		nextNode := edge.To
		if !visited[nextNode] {
			// 递归探索下一个节点
			g.findAllPathsToPrefix(nextNode, endKeyPrefix, paths, append(currentPath, nextNode), visited)
		}
	}
	visited[node] = false
}
