package intimacy

import (
	"reflect"
	"sort"
	"testing"
)

func TestPair(t *testing.T) {
	first := &node{name: "first"}
	second := &node{name: "second"}
	p := &pair{first, second}

	// Test first and second fields of the pair
	if p.first != first {
		t.Errorf("Expected first field to be %v, but got %v", first, p.first)
	}
	if p.second != second {
		t.Errorf("Expected second field to be %v, but got %v", second, p.second)
	}
}

func TestEdge(t *testing.T) {
	// Create nodes
	node1 := &node{name: "node1"}
	node2 := &node{name: "node2"}

	// Create edge and pair
	pair := &pair{first: node1, second: node2}
	edge := &edge{value: 1, pair: pair}

	// Test value and pair fields of the edge
	if edge.value != 1 {
		t.Errorf("Expected value field to be 1, but got %v", edge.value)
	}
	if edge.pair != pair {
		t.Errorf("Expected pair field to be %v, but got %v", pair, edge.pair)
	}
}

func TestNode(t *testing.T) {
	// Create node
	node := &node{name: "node"}

	// Test name and edges fields of the node
	if node.name != "node" {
		t.Errorf("Expected name field to be 'node', but got '%v'", node.name)
	}
	if node.edges != nil {
		t.Errorf("Expected edges field to not be nil, but it was nil")
	}
}

func TestNewGraph(t *testing.T) {
	g := NewGraph()

	// Check if the directed has an initial node
	if g.init == nil {
		t.Error("NewGraph() failed: directed init node is nil")
	}
	// Check if the initial node has the correct name
	if g.init.name != "intimacy_tree_root_node" {
		t.Errorf("NewGraph() failed: expected node name 'intimacy_tree_root_node', got '%s'", g.init.name)
	}
	// Check if the initial node has an empty edges slice
	if len(g.init.edges) != 0 {
		t.Errorf("NewGraph() failed: expected initial node edges to be empty, got %d edges", len(g.init.edges))
	}
}

func TestIntersect(t *testing.T) {
	// 创建测试用的节点
	n1 := &node{name: "node1", edges: nil}
	n2 := &node{name: "node2", edges: nil}
	n3 := &node{name: "node3", edges: nil}
	n4 := &node{name: "node4", edges: nil}
	n5 := &node{name: "node5", edges: nil}
	n6 := &node{name: "node6", edges: nil}

	// 构造测试用的数组
	arr1 := []*node{n1, n2, n3}
	arr2 := []*node{n3, n4, n5}
	arr3 := []*node{n1, n3, n5, n6}

	// 执行测试
	result := intersect(arr1, arr2)

	// 验证结果
	if len(result) != 1 || result[0] != n3 {
		t.Errorf("Expected intersection of arr1 and arr2 to be [%v], but got %v", n3, result)
	}

	result = intersect(arr2, arr3)
	if len(result) != 2 || (result[0] != n3 && result[1] != n5) {
		t.Errorf("Expected intersection of arr2 and arr3 to be [%v %v], but got %v", n3, n5, result)
	}

	result = intersect(arr1, arr3)
	if len(result) != 2 || (result[0] != n1 && result[1] != n3) {
		t.Errorf("Expected intersection of arr1 and arr3 to be [%v %v], but got %v", n1, n3, result)
	}

	// 测试空数组
	var arr4 []*node
	var arr5 []*node
	result = intersect(arr4, arr5)
	if len(result) != 0 {
		t.Errorf("Expected intersection of empty arrays to be [], but got %v", result)
	}
}

func TestUnion(t *testing.T) {
	n1 := &node{name: "node1", edges: nil}
	n2 := &node{name: "node2", edges: nil}
	n3 := &node{name: "node3", edges: nil}
	n4 := &node{name: "node4", edges: nil}

	// Test empty array
	var a1 []*node
	var a2 []*node
	var want nodes
	got := union(a1, a2)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("union(%v, %v) = %v; want %v", a1, a2, got, want)
	}

	// Test arrays with one common element
	a1 = []*node{n1, n2}
	a2 = []*node{n2, n3}
	want = []*node{n1, n2, n3}
	got = union(a1, a2)
	sort.Sort(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("union(%v, %v) = %v; want %v", a1, a2, got, want)
	}

	// Test arrays with all elements in common
	a1 = []*node{n1, n2, n3}
	a2 = []*node{n1, n2, n3}
	want = []*node{n1, n2, n3}
	got = union(a1, a2)
	sort.Sort(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("union(%v, %v) = %v; want %v", a1, a2, got, want)
	}

	// Test arrays with no common element
	a1 = []*node{n1, n2}
	a2 = []*node{n3, n4}
	want = []*node{n1, n2, n3, n4}
	got = union(a1, a2)
	sort.Sort(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("union(%v, %v) = %v; want %v", a1, a2, got, want)
	}
}

func TestAdd(t *testing.T) {
	g := NewGraph()

	n1 := g.add("node1")
	if len(g.init.edges) != 1 {
		t.Errorf("Expected g.init.edges length to be 1 but got %d", len(g.init.edges))
	}
	if g.init.edges[0].pair.first != g.init {
		t.Errorf("Expected g.init.edges[0].pair.first to be g.init but got %v", g.init.edges[0].pair.first)
	}
	if g.init.edges[0].pair.second != n1 {
		t.Errorf("Expected g.init.edges[0].pair.second to be %v but got %v", n1, g.init.edges[0].pair.second)
	}
	if n1.edges != nil {
		t.Errorf("Expected n1.edges to be nil but got %v", n1.edges)
	}

	n2 := g.add("node2")
	if len(g.init.edges) != 2 {
		t.Errorf("Expected g.init.edges length to be 2 but got %d", len(g.init.edges))
	}
	if g.init.edges[1].pair.first != g.init {
		t.Errorf("Expected g.init.edges[1].pair.first to be g.init but got %v", g.init.edges[1].pair.first)
	}
	if g.init.edges[1].pair.second != n2 {
		t.Errorf("Expected g.init.edges[1].pair.second to be %v but got %v", n2, g.init.edges[1].pair.second)
	}
	if n2.edges != nil {
		t.Errorf("Expected n2.edges to be nil but got %v", n2.edges)
	}
}

func TestGraphFind(t *testing.T) {
	g := NewGraph()
	a := g.add("a")
	b := g.add("b")
	c := g.add("c")

	foundA := g.find("a")
	if foundA != a {
		t.Errorf("Expected to find node 'a', but got %+v", foundA)
	}

	foundB := g.find("b")
	if foundB != b {
		t.Errorf("Expected to find node 'b', but got %+v", foundB)
	}

	foundC := g.find("c")
	if foundC != c {
		t.Errorf("Expected to find node 'c', but got %+v", foundC)
	}

	foundD := g.find("d")
	if foundD != nil {
		t.Errorf("Expected to not find node 'd', but got %+v", foundD)
	}
}

func TestGraph_IntimacyPlusOne(t *testing.T) {
	g := NewGraph()
	err := g.IntimacyPlusOne("A", "B")
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	// check if the edge was added correctly
	fn := g.find("A")
	if len(fn.edges) != 1 {
		t.Errorf("expected 1 edge for node A but got %d", len(fn.edges))
	}
	sn := g.find("B")
	if len(sn.edges) != 1 {
		t.Errorf("expected 1 edge for node B but got %d", len(sn.edges))
	}
	eg := fn.edges[0]
	if eg.pair.first != fn || eg.pair.second != sn || eg.value != 1 {
		t.Errorf("unexpected edge value, expected {%s, %s}:1 but got {%s, %s}:%d", fn.name, sn.name, eg.pair.first.name, eg.pair.second.name, eg.value)
	}

	// add another edge
	err = g.IntimacyPlusOne("B", "C")
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	// check if the second edge was added correctly
	fn = g.find("B")
	if len(fn.edges) != 2 {
		t.Errorf("expected 2 edges for node B but got %d", len(fn.edges))
	}
	sn = g.find("C")
	if len(sn.edges) != 1 {
		t.Errorf("expected 1 edge for node C but got %d", len(sn.edges))
	}
	for _, eg := range fn.edges {
		if eg.pair.another(fn).name == sn.name {
			if eg.value != 1 {
				t.Errorf("unexpected edge value, expected {%s, %s}:1 but got {%s, %s}:%d", fn.name, sn.name, eg.pair.first.name, eg.pair.second.name, eg.value)
			}
			return
		}
	}
	t.Errorf("expected an edge between nodes B and C but found none")
}

func TestIntimacy(t *testing.T) {
	g := NewGraph()

	_ = g.add("Alice")
	_ = g.add("Bob")
	_ = g.add("Charlie")
	_ = g.add("David")
	_ = g.add("Susan")

	_ = g.IntimacyPlusOne("Alice", "Bob")
	_ = g.IntimacyPlusOne("Alice", "Charlie")
	_ = g.IntimacyPlusOne("Bob", "David")
	_ = g.IntimacyPlusOne("Charlie", "David")
	_ = g.IntimacyPlusOne("Charlie", "Susan")

	tests := []struct {
		name1    string
		name2    string
		expected float64
	}{
		{"Alice", "Bob", 1},
		{"Alice", "Charlie", 1},
		{"Alice", "David", 1},
		{"Alice", "Susan", 0.5},
		{"Bob", "Charlie", 2.0 / 3},
		{"Charlie", "David", 1},
	}

	for _, test := range tests {
		got := g.Intimacy(test.name1, test.name2)
		if got != test.expected {
			t.Errorf("Intimacy(%q, %q) = %v, want %v", test.name1, test.name2, got, test.expected)
		}
	}
}
