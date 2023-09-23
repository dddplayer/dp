package directed

import "fmt"

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
	// 检查节点键是否冲突
	for _, node := range g.Nodes {
		if node.Key == key {
			return fmt.Errorf("节点键冲突：%s", key)
		}
	}

	// 创建新节点
	newNode := &Node{
		Key:   key,
		Value: value,
		Edges: []*Edge{},
	}

	// 将新节点添加到图中
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
