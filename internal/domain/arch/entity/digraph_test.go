package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"github.com/dddplayer/dp/pkg/datastructure/directed"
	"golang.org/x/exp/slices"
	"testing"
)

func TestSummaryNodeEdges(t *testing.T) {
	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}

	// Create nodes for testing
	nodeA := directed.NewNode("A", nil)
	nodeB := directed.NewNode("B", nil)
	nodeC := directed.NewNode("C", nil)
	nodeD := directed.NewNode("D", nil)

	// Create edges for testing
	edgeAB := &directed.Edge{From: nodeA, To: nodeB, Type: arch.RelationTypeComposition}
	edgeBC := &directed.Edge{From: nodeB, To: nodeC, Type: arch.RelationTypeEmbedding}
	edgeCD := &directed.Edge{From: nodeC, To: nodeD, Type: arch.RelationTypeAssociation}

	// Add edges to nodes
	nodeA.Edges = append(nodeA.Edges, edgeAB)
	nodeB.Edges = append(nodeB.Edges, edgeBC)
	nodeC.Edges = append(nodeC.Edges, edgeCD)

	// Add nodes to the graph
	g.Nodes = append(g.Nodes, nodeA, nodeB, nodeC, nodeD)

	// Call the method being tested
	edges, err := g.obtainEdgesFromNodeTree(nodeA)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert the expected edges
	expectedEdges := []*directed.Edge{edgeAB, edgeBC, edgeCD}
	if len(edges) != len(expectedEdges) {
		t.Errorf("Expected %d edges, but got %d", len(expectedEdges), len(edges))
	}

	for i := range edges {
		if !slices.Contains(expectedEdges, edges[i]) {
			t.Errorf("Expected edge %v not found", edges[i])
		}
	}
}

func TestAddObj(t *testing.T) {
	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}

	// Create a mock ObjIdentifier for testing
	objID := &MockObjIdentifier{
		id:   "123",
		name: "TestObject",
		dir:  "/path/to/object",
	}

	// Call the method being tested
	err := g.AddObj(objID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Find the added node in the graph
	node := g.FindNodeByKey(objID.ID())
	if node == nil {
		t.Error("Node not found in the graph")
	}

	// Assert the properties of the added node
	if node.Key != objID.ID() {
		t.Errorf("Expected node key: %s, but got: %s", objID.ID(), node.Key)
	}

	if node.Value != objID {
		t.Error("Expected node value to be the ObjIdentifier object")
	}

	// Attempt to add the same ObjIdentifier again
	err = g.AddObj(objID)
	if err == nil {
		t.Error("Expected error when adding duplicate ObjIdentifier, but got nil")
	}
}

func TestSummaryRelations(t *testing.T) {
	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}

	// Create mock ObjIdentifier objects for testing
	objA := &MockObjIdentifier{
		id:   "A",
		name: "ObjectA",
		dir:  "/path/to/objectA",
	}
	objB := &MockObjIdentifier{
		id:   "B",
		name: "ObjectB",
		dir:  "/path/to/objectB",
	}
	objC := &MockObjIdentifier{
		id:   "C",
		name: "ObjectC",
		dir:  "/path/to/objectC",
	}
	objD := &MockObjIdentifier{
		id:   "D",
		name: "ObjectD",
		dir:  "/path/to/objectD",
	}

	// Add the mock ObjIdentifier objects to the graph
	err := g.AddObj(objA)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(objB)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(objC)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(objD)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create edges for testing
	_ = g.AddEdge(objA.ID(), objB.ID(), arch.RelationTypeComposition, valueobject.NewEmptyRelationPos())
	_ = g.AddEdge(objB.ID(), objC.ID(), arch.RelationTypeEmbedding, valueobject.NewEmptyRelationPos())
	_ = g.AddEdge(objC.ID(), objD.ID(), arch.RelationTypeAssociation, valueobject.NewEmptyRelationPos())

	// Call the method being tested
	relations, err := g.SummaryRelationMetas(objA, objD)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert the expected relations
	expectedRelations := []arch.RelationType{arch.RelationTypeAssociation}
	if len(relations) != len(expectedRelations) {
		t.Errorf("Expected %d relations, but got %d", len(expectedRelations), len(relations))
	}

	for i := range relations {
		if relations[i].Type() != expectedRelations[i] {
			t.Errorf("Expected relation %v, but got %v", expectedRelations[i], relations[i])
		}
	}
}

func TestAddRelation(t *testing.T) {
	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}

	// Create mock ObjIdentifier objects for testing
	mockIdentifier1 := &MockObjIdentifier{
		id:   "A",
		name: "ObjectA",
		dir:  "/path/to/objectA",
	}
	mockPosition1 := &MockPosition{FilenameVal: "file1", OffsetVal: 10, LineVal: 5, ColumnVal: 2}
	objA := MockObject{id: mockIdentifier1, position: mockPosition1}

	mockIdentifier2 := &MockObjIdentifier{
		id:   "B",
		name: "ObjectB",
		dir:  "/path/to/objectB",
	}
	mockPosition2 := &MockPosition{FilenameVal: "file2", OffsetVal: 10, LineVal: 5, ColumnVal: 2}
	objB := MockObject{id: mockIdentifier2, position: mockPosition2}

	// Add the mock ObjIdentifier objects to the graph
	err := g.AddObj(objA)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(objB)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create mock Relation objects for testing
	depRel := &MockDependenceRelation{
		from:      objA,
		dependsOn: objB,
	}
	compRel := &MockCompositionRelation{
		from:  objA,
		child: objB,
	}
	embRel := &MockEmbeddingRelation{
		from:     objA,
		embedded: objB,
	}
	implRel := &MockImplementationRelation{
		from:       objA,
		implements: []arch.Object{objB},
	}
	assocRel := &MockAssociationRelation{
		from:            objA,
		refer:           objB,
		associationType: arch.RelationTypeAssociationOneOne,
	}

	// Call the method being tested
	err = g.AddRelation(depRel)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddRelation(compRel)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddRelation(embRel)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddRelation(implRel)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddRelation(assocRel)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Retrieve the added edges
	edgesA, _ := g.obtainEdgesFromNodeTree(g.FindNodeByKey(objA.ID()))
	edgesB, _ := g.obtainEdgesFromNodeTree(g.FindNodeByKey(objB.ID()))

	// Assert the number of edges for node A
	expectedEdgesA := 5
	if len(edgesA) != expectedEdgesA {
		t.Errorf("Expected %d edges for node A, but got %d", expectedEdgesA, len(edgesA))
	}

	// Assert the number of edges for node B
	expectedEdgesB := 0
	if len(edgesB) != expectedEdgesB {
		t.Errorf("Expected %d edges for node B, but got %d", expectedEdgesB, len(edgesB))
	}
}
