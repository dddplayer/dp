package radix

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLongestPrefix(t *testing.T) {
	testCases := []struct {
		k1       string
		k2       string
		expected int
	}{
		{"hello", "help", 3},
		{"world", "word", 3},
		{"apple", "orange", 0},
		{"", "hello", 0},
		{"hello", "", 0},
		{"", "", 0},
	}

	for _, tc := range testCases {
		if got := longestPrefix(tc.k1, tc.k2); got != tc.expected {
			t.Errorf("longestPrefix(%q, %q) = %d, expected %d", tc.k1, tc.k2, got, tc.expected)
		}
	}
}

func TestSortAscending(t *testing.T) {
	testCases := []struct {
		es       edges
		expected edges
	}{
		{
			edges{
				&edge{"apple", nil, nil},
				&edge{"cat", nil, nil},
				&edge{"banana", nil, nil},
			},
			edges{
				&edge{"apple", nil, nil},
				&edge{"banana", nil, nil},
				&edge{"cat", nil, nil},
			},
		},
		{
			edges{
				&edge{"dog", nil, nil},
				&edge{"fish", nil, nil},
				&edge{"bird", nil, nil},
			},
			edges{
				&edge{"bird", nil, nil},
				&edge{"dog", nil, nil},
				&edge{"fish", nil, nil},
			},
		},
		{
			edges{
				&edge{"hello", nil, nil},
				&edge{"world", nil, nil},
				&edge{"goodbye", nil, nil},
			},
			edges{
				&edge{"goodbye", nil, nil},
				&edge{"hello", nil, nil},
				&edge{"world", nil, nil},
			},
		},
	}

	for _, tc := range testCases {
		sortAscending(tc.es)
		if !reflect.DeepEqual(tc.es, tc.expected) {
			t.Errorf("sortAscending(%v) = %v, expected %v", tc.es, tc.es, tc.expected)
		}
	}
}

func TestAddEdge(t *testing.T) {
	n1 := &node{val: "a"}
	n2 := &node{val: "b"}

	e1 := &edge{"apple", nil, nil}
	e2 := &edge{"banana", nil, nil}
	e3 := &edge{"cat", nil, nil}

	// Add first edge to node n1
	n1.addEdge(e1)

	// Make sure the edge is properly added to node n1
	if n1.suffixes[0] != e1 {
		t.Errorf("Expected edge %v to be added to node %v, but got %v instead",
			e1, n1, n1.suffixes[0])
	}

	// Add second edge to node n1
	n1.addEdge(e2)

	// Make sure both edges are properly added to node n1
	if n1.suffixes[0] != e1 || n1.suffixes[1] != e2 {
		t.Errorf("Expected edges %v and %v to be added to node %v, but got %v instead",
			e1, e2, n1, n1.suffixes)
	}

	// Add edge to node n2
	n2.addEdge(e3)

	// Make sure the edge is properly added to node n2
	if n2.suffixes[0] != e3 {
		t.Errorf("Expected edge %v to be added to node %v, but got %v instead",
			e3, n2, n2.suffixes[0])
	}
}

func TestTree_Insert(t *testing.T) {
	tree := &Tree{root: &node{}}

	// insert a node to the root
	if ok := tree.Insert("a", 1); !ok {
		t.Errorf("failed to insert a node to the root")
	}

	// insert a child node to an existing node
	if ok := tree.Insert("ab", 2); !ok {
		t.Errorf("failed to insert a child node to an existing node")
	}

	// update an existing node
	if ok := tree.Insert("a", 3); !ok {
		t.Errorf("failed to update an existing node")
	}

	// insert a long key
	if ok := tree.Insert("hello, world!", 4); !ok {
		t.Errorf("failed to insert a long key")
	}

	// insert an empty key should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("inserting an empty key should panic")
		}
	}()
	tree.Insert("", 5)
}

func TestInsert_RightDepthSplit(t *testing.T) {
	tree := &Tree{root: &node{}}

	tree.Insert("abc", 1)
	tree.Insert("abx", 2)

	childNode := tree.root

	var commonEdge *edge
	for _, e := range childNode.suffixes {
		if e.name == "ab" {
			commonEdge = e
			break
		}
	}

	if commonEdge == nil {
		t.Errorf("Expected a common edge with name 'c' in childNode.suffixes, but not found")
	}

	if commonEdge.end.val != nil {
		t.Errorf("Expected end node value nil, but got %v", commonEdge.end.val)
	}

	if len(commonEdge.end.suffixes) != 2 {
		t.Errorf("Expected empty suffixes for end node, but got %v", len(commonEdge.end.suffixes))
	}
}

func TestInsert_RightDepthSplitCommonEdge(t *testing.T) {
	tree := &Tree{root: &node{}}

	tree.Insert("ab", 1)
	tree.Insert("abx", 2)

	childNode := tree.root

	var commonEdge *edge
	for _, e := range childNode.suffixes {
		if e.name == "ab" {
			commonEdge = e
			break
		}
	}

	if commonEdge == nil {
		t.Errorf("Expected a common edge with name 'c' in childNode.suffixes, but not found")
	}

	if commonEdge.end.val != 1 {
		t.Errorf("Expected end node value nil, but got %v", commonEdge.end.val)
	}

	if len(commonEdge.end.suffixes) != 1 {
		t.Errorf("Expected empty suffixes for end node, but got %v", len(commonEdge.end.suffixes))
	}
}

func TestGetFirstByte(t *testing.T) {
	s1 := "apple"
	s2 := "book"
	s3 := ""

	if getFirstByte(s1) != 'a' {
		t.Errorf("getFirstByte(%s) = %c; want a", s1, getFirstByte(s1))
	}

	if getFirstByte(s2) != 'b' {
		t.Errorf("getFirstByte(%s) = %c; want b", s2, getFirstByte(s2))
	}

	if getFirstByte(s3) != 0 {
		t.Errorf("getFirstByte(%s) = %c; want 0", s3, getFirstByte(s3))
	}
}

func TestFindEdgeWithSamePrefix(t *testing.T) {
	es := edges{
		&edge{name: "apple"},
		&edge{name: "banana"},
		&edge{name: "avocado"},
		&edge{name: "orange"},
	}

	// test case 0: empty edges
	e := findEdgeWithSamePrefix('a', nil)
	if e != nil {
		t.Errorf("expected no edge found, but got: %v", e)
	}

	// test case 1: firstByte exists in one of the edges
	e = findEdgeWithSamePrefix('a', es)
	if e == nil || e.name != "apple" {
		t.Errorf("expected edge with prefix 'a', but got: %v", e)
	}

	// test case 2: firstByte exists in multiple edges
	es = append(es, &edge{name: "apricot"})
	e = findEdgeWithSamePrefix('a', es)
	if e == nil || e.name != "apple" {
		t.Errorf("expected edge with prefix 'apple', but got: %v", e)
	}

	// test case 3: firstByte does not exist in any edge
	e = findEdgeWithSamePrefix('z', es)
	if e != nil {
		t.Errorf("expected no edge found, but got: %v", e)
	}
}

func TestTree_Get(t *testing.T) {
	tree := NewTree()

	// test empty tree
	v, ok := tree.Get("foo")
	if v != nil || ok {
		t.Errorf("Expected (nil, false), but got (%v, %v)", v, ok)
	}

	// insert a node and test a non-existing key
	tree.Insert("bar", "baz")
	v, ok = tree.Get("foo")
	if v != nil || ok {
		t.Errorf("Expected (nil, false), but got (%v, %v)", v, ok)
	}

	// test an existing key
	v, ok = tree.Get("bar")
	if v != "baz" || !ok {
		t.Errorf("Expected (baz, true), but got (%v, %v)", v, ok)
	}
}

func TestDelEdge(t *testing.T) {
	rootNode := &node{}
	edge1 := &edge{name: "Edge1", start: rootNode, end: nil}
	edge2 := &edge{name: "Edge2", start: nil, end: rootNode}
	rootNode.suffixes = append(rootNode.suffixes, edge1)
	rootNode.delEdge(edge1)

	if len(rootNode.suffixes) != 0 {
		t.Errorf("Expected no edges in suffixes, but got %d", len(rootNode.suffixes))
	}

	if rootNode.delEdge(edge2); len(rootNode.suffixes) != 0 {
		t.Errorf("Expected no change when deleting an edge not connected to the node")
	}
}

func TestWalkNode(t *testing.T) {
	tree := &Tree{root: &node{}}

	tree.Insert("abc", 1)
	tree.Insert("abx", 2)
	tree.Insert("abcd", 3)
	tree.Insert("abxy", 4)

	var results []string

	walker := func(prefix string, v any, ws WalkState) WalkStatus {
		results = append(results, fmt.Sprintf("%s:%v:%v", prefix, v, ws))
		return WalkContinue
	}

	walkNode("", tree.root, walker)

	expectedResults := []string{
		":<nil>:1",
		"ab:<nil>:1",
		"abc:1:1",
		"abcd:3:1",
		"abcd:3:2",
		"abc:1:2",
		"abx:2:1",
		"abxy:4:1",
		"abxy:4:2",
		"abx:2:2",
		"ab:<nil>:2",
		":<nil>:2",
	}

	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, but got %d", len(expectedResults), len(results))
		return
	}

	for i, expected := range expectedResults {
		if results[i] != expected {
			t.Errorf("Result %d: Expected '%s', but got '%s'", i, expected, results[i])
		}
	}
}
