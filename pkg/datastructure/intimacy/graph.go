package intimacy

import "strings"

type pair struct {
	first  *node
	second *node
}

func (p *pair) another(n *node) *node {
	if p.first == n {
		return p.second
	}
	return p.first
}

type edge struct {
	value uint
	pair  *pair
}

type node struct {
	name  string
	edges []*edge
}

type Graph struct {
	init *node
}

func NewGraph() *Graph {
	return &Graph{init: &node{
		name:  "intimacy_tree_root_node",
		edges: nil,
	}}
}

func (g *Graph) Intimacy(first, second string) float64 {
	fn := g.find(first)
	if fn == nil {
		return 0
	}
	sn := g.find(second)
	if sn == nil {
		return 0
	}

	dirRelated := 0.0
	for _, e := range fn.edges {
		if e.pair.another(fn).name == sn.name {
			dirRelated = float64(e.value)
		}
	}

	indirRelated := 0.0
	var fnNodes, snNodes []*node
	for _, e := range fn.edges {
		fnNodes = append(fnNodes, e.pair.another(fn))
	}
	for _, e := range sn.edges {
		snNodes = append(snNodes, e.pair.another(sn))
	}
	iNodes := intersect(fnNodes, snNodes)
	uNodes := union(fnNodes, snNodes)

	indirRelated = float64(len(iNodes)) / float64(len(uNodes))

	return dirRelated + indirRelated
}

// intersect 获取两个数组的合集
func intersect(arr1, arr2 []*node) []*node {
	// 定义一个空的切片，用于存储合集
	var result []*node

	// 遍历第一个数组，如果元素在第二个数组中出现，则加入合集
	for _, v := range arr1 {
		for _, w := range arr2 {
			if v == w {
				result = append(result, v)
				break
			}
		}
	}

	return result
}

type nodes []*node

func (ele nodes) Len() int           { return len(ele) }
func (ele nodes) Swap(i, j int)      { ele[i], ele[j] = ele[j], ele[i] }
func (ele nodes) Less(i, j int) bool { return strings.Compare(ele[i].name, ele[j].name) < 0 }

// union 获取两个数组的并集
func union(arr1, arr2 []*node) nodes {
	// 定义一个 map，用于存储数组中的元素
	m := make(map[*node]bool)

	// 遍历第一个数组，将元素加入 map 中
	for _, v := range arr1 {
		m[v] = true
	}

	// 遍历第二个数组，将元素加入 map 中
	for _, v := range arr2 {
		m[v] = true
	}

	// 将 map 中的元素转换为切片，即为并集
	var result []*node
	for k := range m {
		result = append(result, k)
	}

	return result
}

func (g *Graph) IntimacyPlusOne(first, second string) error {
	fn := g.find(first)
	if fn == nil {
		fn = g.add(first)
	}
	sn := g.find(second)
	if sn == nil {
		sn = g.add(second)
	}

	for _, eg := range fn.edges {
		if eg.pair.another(fn).name == sn.name {
			eg.value++
			return nil
		}
	}

	e := &edge{
		value: 1,
		pair: &pair{
			first:  fn,
			second: sn,
		},
	}
	fn.edges = append(fn.edges, e)
	sn.edges = append(sn.edges, e)

	return nil
}

func (g *Graph) add(name string) *node {
	n := &node{
		name:  name,
		edges: nil,
	}
	e := &edge{
		value: 0,
		pair: &pair{
			first:  g.init,
			second: n,
		},
	}
	g.init.edges = append(g.init.edges, e)
	return n
}

func (g *Graph) find(name string) *node {
	for _, e := range g.init.edges {
		if e.pair.first.name == name {
			return e.pair.first
		} else if e.pair.second.name == name {
			return e.pair.second
		}
	}
	return nil
}
