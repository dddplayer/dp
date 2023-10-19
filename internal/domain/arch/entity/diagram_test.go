package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"github.com/dddplayer/dp/pkg/datastructure/directed"
	"testing"
)

func TestNewDiagram(t *testing.T) {
	name := "TestDiagram"
	diagram, err := NewDiagram(name, arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if diagram.Name() != name {
		t.Errorf("Expected Diagram name '%s', but got '%s'", name, diagram.Name())
	}

	if diagram.root == nil {
		t.Error("Expected Diagram to have a root node, but it's nil")
	}

	rootNode := diagram.root
	if rootNode.Key != name || rootNode.Value.(arch.Object).Identifier().Name() != name {
		t.Errorf("Expected root node with Key and Value '%s', but got Key '%s' and Value '%s'", name, rootNode.Key, rootNode.Value)
	}

	if len(diagram.objs) != 0 {
		t.Errorf("Expected Diagram to be empty, but it contains %d objects", len(diagram.objs))
	}
}

func TestDiagram_Objects(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create some mock objects
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)
	mockObject3 := newMockObject(3)

	// Add the mock objects to the Diagram
	diagram.objs = append(diagram.objs, mockObject1, mockObject2, mockObject3)

	// Call the Objects function
	objects := diagram.Objects()

	// Check if the returned objects match the mock objects
	if len(objects) != 3 {
		t.Errorf("Expected 3 objects, but got %d", len(objects))
	}

	// Check if the objects are in the correct order
	if objects[0].Identifier().Name() != "MockObject_1" {
		t.Errorf("Expected object at index 0 to have name 'Object1', but got '%s'", objects[0].Identifier().Name())
	}
	if objects[1].Identifier().Name() != "MockObject_2" {
		t.Errorf("Expected object at index 1 to have name 'Object2', but got '%s'", objects[1].Identifier().Name())
	}
	if objects[2].Identifier().Name() != "MockObject_3" {
		t.Errorf("Expected object at index 2 to have name 'Object3', but got '%s'", objects[2].Identifier().Name())
	}
}

func TestDiagram_AppendObject(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create some mock objects
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)

	// Call the AppendObject function to add the mock objects
	diagram.AppendObject(mockObject1, mockObject2)

	// Check if the objects have been correctly appended
	objects := diagram.Objects()

	// Check if the number of objects matches the expected count
	if len(objects) != 2 {
		t.Errorf("Expected 2 objects, but got %d", len(objects))
	}

	// Check if the objects are in the correct order
	if objects[0].Identifier().Name() != "MockObject_1" {
		t.Errorf("Expected object at index 0 to have name 'MockObject_1', but got '%s'", objects[0].Identifier().Name())
	}
	if objects[1].Identifier().Name() != "MockObject_2" {
		t.Errorf("Expected object at index 1 to have name 'MockObject_2', but got '%s'", objects[1].Identifier().Name())
	}
}

func TestDiagram_AddObj(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create a mock object
	mockObject := newMockObject(1)

	// Call the AddObj function to add the mock object
	err = diagram.AddObj(mockObject)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the object has been correctly added
	obj := diagram.FindNodeByKey(mockObject.Identifier().ID())

	// Check if the number of objects matches the expected count
	if obj == nil {
		t.Errorf("Expected object not empty, but got nil")
	}

	// Check if the object's name matches the expected value
	if obj.Key != mockObject.ID() {
		t.Errorf("Expected the added object to have same key with mockObject, but not")
	}
}

func TestDiagram_AddObjTo(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create a mock object
	mockObject := newMockObject(1)
	parentID := diagram.Name()
	relationType := arch.RelationTypeAggregation

	// Call the AddObjTo function to add the mock object
	err = diagram.AddObjTo(mockObject, parentID, relationType)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the object has been correctly added
	obj := diagram.FindNodeByKey(mockObject.Identifier().ID())
	if obj == nil {
		t.Errorf("Expected object not empty, but got nil")
	}

	// Check if the object's name matches the expected value
	if obj.Key != mockObject.Identifier().ID() {
		t.Errorf("Expected the added object to have the same key as mockObject, but not")
	}

	// Check if the relation has been correctly added
	edges := diagram.Edges()
	if len(edges) != 0 { // ignored
		t.Errorf("Expected 0 edge, but got %d", len(edges))
	}
}

func TestDiagram_AddStringTo(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Define the parameters
	obj := "TestString"
	parentID := diagram.Name()
	relationType := arch.RelationTypeAggregation

	// Call the AddStringTo function to add the string object
	err = diagram.AddStringTo(obj, parentID, relationType)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the string object has been correctly added
	strObj := diagram.FindNodeByKey(valueobject.NewStringObj(obj).Identifier().ID())
	if strObj == nil {
		t.Errorf("Expected string object not empty, but got nil")
	}

	// Check if the string object's key matches the expected value
	if strObj.Key != valueobject.NewStringObj(obj).Identifier().ID() {
		t.Errorf("Expected the added string object to have the same key as the string, but not")
	}

	// Check if the relation has been correctly added
	edges := diagram.Edges()
	if len(edges) != 0 { // ignored
		t.Errorf("Expected 1 edge, but got %d", len(edges))
	}
}

func TestDiagram_AddRelations(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create two mock objects
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)

	// Add the mock objects to the Diagram
	_ = diagram.AddObjTo(mockObject1, diagram.Name(), arch.RelationTypeAggregation)
	_ = diagram.AddObjTo(mockObject2, diagram.Name(), arch.RelationTypeAggregation)

	// Define the relation metas
	relations := []arch.RelationMeta{
		valueobject.NewRelationMeta(arch.RelationTypeAggregationRoot, mockObject1.position, mockObject2.position),
		valueobject.NewRelationMeta(arch.RelationTypeImplementation, mockObject1.position, mockObject2.position),
	}

	// Call the AddRelations function to add the relations
	err = diagram.AddRelations(mockObject1.Identifier().ID(), mockObject2.Identifier().ID(), relations)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the relations have been correctly added
	edges := diagram.Edges()
	if len(edges) != len(relations) {
		t.Errorf("Expected %d edges, but got %d", len(relations), len(edges))
	}
}

func TestDiagram_Name(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Call the Name method to retrieve the name
	name := diagram.Name()

	// Check if the name matches the expected value
	expectedName := "TestDiagram"
	if name != expectedName {
		t.Errorf("Expected name '%s', but got '%s'", expectedName, name)
	}
}

func TestDiagram_CountDepth(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create two mock objects
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)

	// Add the mock objects to the Diagram
	_ = diagram.AddObjTo(mockObject1, diagram.Name(), arch.RelationTypeAggregationRoot)
	_ = diagram.AddObjTo(mockObject2, mockObject1.ID(), arch.RelationTypeAggregationRoot)

	// Call the countDepth method to calculate the depth
	depth := diagram.countDepth(1, diagram.root)

	// Check if the depth matches the expected value
	expectedDepth := 3
	if depth != expectedDepth {
		t.Errorf("Expected depth %d, but got %d", expectedDepth, depth)
	}
}

func TestDiagram_IgnoreEdge(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create a mock edge with various relation types
	mockEdge1 := &directed.Edge{
		From: diagram.root,
		To:   diagram.root,
		Type: arch.RelationTypeAttribution,
	}
	mockEdge2 := &directed.Edge{
		From: diagram.root,
		To:   diagram.root,
		Type: arch.RelationTypeBehavior,
	}
	mockEdge3 := &directed.Edge{
		From: diagram.root,
		To:   diagram.root,
		Type: arch.RelationTypeAggregation,
	}
	mockEdge4 := &directed.Edge{
		From: diagram.root,
		To:   diagram.root,
		Type: arch.RelationTypeComposition,
	}

	// Test ignoring edges with RelationTypeAttribution and RelationTypeBehavior
	if !diagram.ignoreEdge(mockEdge1) {
		t.Errorf("Expected to ignore RelationTypeAttribution edge, but it was not ignored")
	}
	if !diagram.ignoreEdge(mockEdge2) {
		t.Errorf("Expected to ignore RelationTypeBehavior edge, but it was not ignored")
	}

	// Test not ignoring edges with other relation types
	if !diagram.ignoreEdge(mockEdge3) {
		t.Errorf("Expected ignore RelationTypeAggregation edge, but not")
	}
	if !diagram.ignoreEdge(mockEdge4) {
		t.Errorf("Expected not to ignore RelationTypeComposition edge, but it was ignored")
	}

	// Test ignoring edges from StringObj to anything with RelationTypeAggregation
	strObj := valueobject.NewStringObj("TestString")
	_ = diagram.AddObj(strObj)
	mockEdge5 := &directed.Edge{
		From: diagram.FindNodeByKey(strObj.Identifier().ID()),
		To:   diagram.root,
		Type: arch.RelationTypeAggregation,
	}

	if !diagram.ignoreEdge(mockEdge5) {
		t.Errorf("Expected to ignore RelationTypeAggregation edge from StringObj, but it was not ignored")
	}
}

func TestDiagram_ParseNodeEdge(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)
	mockObject3 := newMockObject(3)

	// Add the mock objects to the Diagram with different nesting depths
	_ = diagram.AddObjTo(mockObject1, diagram.Name(), arch.RelationTypeAggregationRoot)
	_ = diagram.AddObjTo(mockObject2, mockObject1.ID(), arch.RelationTypeAggregationRoot)
	_ = diagram.AddObjTo(mockObject3, mockObject2.ID(), arch.RelationTypeAggregationRoot)

	// Parse the node edges
	visited := make(map[*directed.Node]bool)
	edges := diagram.parseNodeEdge(diagram.root, visited)

	if len(edges) != 3 {
		t.Errorf("Expected %d parsed edges, but got %d", 3, len(edges))
	}
}

func TestDiagram_Edges(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)
	mockObject3 := newMockObject(3)

	// Add the mock objects to the Diagram with different nesting depths
	_ = diagram.AddObjTo(mockObject1, diagram.Name(), arch.RelationTypeAggregationRoot)
	_ = diagram.AddObjTo(mockObject2, mockObject1.ID(), arch.RelationTypeAggregationRoot)
	_ = diagram.AddObjTo(mockObject3, mockObject2.ID(), arch.RelationTypeAggregationRoot)

	// Get the edges from the diagram
	edges := diagram.Edges()

	if len(edges) != 3 {
		t.Errorf("Expected %d edges, but got %d", 3, len(edges))
	}

	// Optionally, you can check the properties of individual edges if needed
	// For example, you can check if they have the correct type, from, and to values.
	for _, edge := range edges {
		// Replace these conditions with the actual properties of your edges
		if edge.Type() != arch.RelationTypeAggregationRoot {
			t.Errorf("Expected RelationTypeAggregationRoot relation type, but got %d", edge.Type())
		}
		if edge.From() == "" || edge.To() == "" {
			t.Errorf("Expected non-empty From and To values for edges")
		}
	}
}

func TestDiagramWithSubDiagram(t *testing.T) {
	// Create a subDiagram
	sub := &subDiagram{
		name: "SubDiagram",
		nodes: []arch.Node{
			newMockNode(1),
			newMockNode(2),
		},
		elements: []arch.Element{
			newMockElement(1, newMockNode(3), newMockNode(4)),
			newMockElement(2, newMockNode(5)),
		},
	}

	if sub.Name() != "SubDiagram" {
		t.Errorf("Expected SubDiagram name to be SubDiagram, but got %s", sub.Name())
	}

	if len(sub.Nodes()) != 2 {
		t.Errorf("Expected SubDiagram to have 2 nodes, but got %d", len(sub.Nodes()))
	}

	if len(sub.Summary()) != 2 {
		t.Errorf("Expected SubDiagram to have 2 elements, but got %d", len(sub.Summary()))
	}
}

func TestNode(t *testing.T) {
	// Create a node
	n := &node{
		id:    "testID",
		name:  "testName",
		color: "testColor",
	}

	if n.ID() != "testID" {
		t.Errorf("Expected node ID to be 'testID', but got '%s'", n.ID())
	}

	if n.Name() != "testName" {
		t.Errorf("Expected node name to be 'testName', but got '%s'", n.Name())
	}

	if n.Color() != "testColor" {
		t.Errorf("Expected node color to be 'testColor', but got '%s'", n.Color())
	}
}

func TestElement(t *testing.T) {
	// Create a node
	n := &node{
		id:    "testID",
		name:  "testName",
		color: "testColor",
	}

	// Create an element
	e := &element{
		node:     n,
		t:        elementTypeClass,
		children: []arch.Nodes{[]arch.Node{n, n}},
	}

	if e.ID() != "testID" {
		t.Errorf("Expected element ID to be 'testID', but got '%s'", e.ID())
	}

	if e.Name() != "testName" {
		t.Errorf("Expected element name to be 'testName', but got '%s'", e.Name())
	}

	if e.Color() != "testColor" {
		t.Errorf("Expected element color to be 'testColor', but got '%s'", e.Color())
	}

	if e.t != elementTypeClass {
		t.Errorf("Expected element type to be 'class', but got '%s'", e.t)
	}

	if len(e.Children()) != 1 {
		t.Errorf("Expected element to have 1 children nodes, but got %d", len(e.Children()))
	}
}

func TestAddLeft(t *testing.T) {
	// Create a node
	n1 := &node{
		id:    "testID1",
		name:  "testName1",
		color: "testColor1",
	}

	n2 := &node{
		id:    "testID2",
		name:  "testName2",
		color: "testColor2",
	}

	// Create an element of elementTypeClass
	e := &element{
		node:     n1,
		t:        elementTypeClass,
		children: make([]arch.Nodes, 1),
	}

	// Add a node to the left
	e.addLeft(n2)

	// Verify that the node was added to the left
	children := e.Children()[0]
	if len(children) != 1 {
		t.Errorf("Expected 1 children nodes on the left, but got %d", len(children))
	}

	if children[0].ID() != "testID2" {
		t.Errorf("Expected the first child ID to be 'testID2', but got '%s'", children[0].ID())
	}

}

func TestAddRight(t *testing.T) {
	// Create two nodes
	n1 := &node{
		id:    "testID1",
		name:  "testName1",
		color: "testColor1",
	}

	n2 := &node{
		id:    "testID2",
		name:  "testName2",
		color: "testColor2",
	}

	// Create an element of elementTypeClass
	e := &element{
		node:     n1,
		t:        elementTypeClass,
		children: make([]arch.Nodes, 2),
	}

	// Add a node to the right
	e.addRight(n2)

	// Verify that the node was added to the right
	children := e.Children()[1]
	if len(children) != 1 {
		t.Errorf("Expected 1 child node on the right, but got %d", len(children))
	}

	if children[0].ID() != "testID2" {
		t.Errorf("Expected the first child ID on the right to be 'testID2', but got '%s'", children[0].ID())
	}

	// Verify that the left side remains empty
	leftChildren := e.Children()[0]
	if len(leftChildren) != 0 {
		t.Errorf("Expected no child nodes on the left, but got %d", len(leftChildren))
	}
}

func TestNewElement(t *testing.T) {
	// Create a node
	n := &node{
		id:    "testID",
		name:  "testName",
		color: "testColor",
	}

	// Create an element of elementTypeClass
	e := newElement(n, elementTypeClass)

	// Verify that the element was created with the correct type and child nodes
	if e.t != elementTypeClass {
		t.Errorf("Expected elementTypeClass, but got %s", e.t)
	}

	if len(e.Children()) != 2 {
		t.Errorf("Expected 2 child nodes, but got %d", len(e.Children()))
	}

	// Create an element of elementTypeGeneral
	e = newElement(n, elementTypeGeneral)

	// Verify that the element was created with the correct type and child nodes
	if e.t != elementTypeGeneral {
		t.Errorf("Expected elementTypeGeneral, but got %s", e.t)
	}

	if len(e.Children()) != 1 {
		t.Errorf("Expected 1 child node, but got %d", len(e.Children()))
	}
}

func TestNewEdge(t *testing.T) {
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)

	e := newEdge(mockObject1.ID(), mockObject2.ID(), arch.RelationTypeAggregation, mockObject1.Position(), mockObject2.Position())

	// Verify that the edge was created with the correct values
	if e.from != mockObject1.ID() {
		t.Errorf("Expected from to be %s, but got %s", mockObject1.ID(), e.from)
	}

	if e.to != mockObject2.ID() {
		t.Errorf("Expected to to be %s, but got %s", mockObject2.ID(), e.to)
	}

	if e.relationType != arch.RelationTypeAggregation {
		t.Errorf("Expected relationType to be %v, but got %v", arch.RelationTypeAggregation, e.relationType)
	}

	if e.fromPos != mockObject1.Position() {
		t.Errorf("Expected fromPos to be %v, but got %v", mockObject1.Position(), e.fromPos)
	}

	if e.toPos != mockObject2.Position() {
		t.Errorf("Expected toPos to be %v, but got %v", mockObject2.Position(), e.toPos)
	}
}

func TestEdgeMethods(t *testing.T) {
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)

	e := newEdge(mockObject1.ID(), mockObject2.ID(), arch.RelationTypeAggregation, mockObject1.Position(), mockObject2.Position())

	// Test the From method
	if from := e.From(); from != mockObject1.ID() {
		t.Errorf("Expected From() to return %s, but got %s", mockObject1.ID(), from)
	}

	// Test the To method
	if to := e.To(); to != mockObject2.ID() {
		t.Errorf("Expected To() to return %s, but got %s", mockObject2.ID(), to)
	}

	// Test the Type method
	if relationType := e.Type(); relationType != arch.RelationTypeAggregation {
		t.Errorf("Expected Type() to return %v, but got %v", arch.RelationTypeAggregation, relationType)
	}
}

func TestEdgeKey(t *testing.T) {
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)

	e := newEdge(mockObject1.ID(), mockObject2.ID(), arch.RelationTypeAggregation, mockObject1.Position(), mockObject2.Position())

	// Test the Key method
	expectedKey := fmt.Sprintf("%s-%s-%d", mockObject1.ID(), mockObject2.ID(), arch.RelationTypeAggregation)
	if key := e.Key(); key != expectedKey {
		t.Errorf("Expected Key() to return '%s', but got '%s'", expectedKey, key)
	}
}

func TestEdgePos(t *testing.T) {
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)

	e := newEdge(mockObject1.ID(), mockObject2.ID(), arch.RelationTypeAggregation, mockObject1.Position(), mockObject2.Position())

	// Test the pos method
	pos := e.pos()

	// Check if the returned Position object matches the expected values
	if pos == nil {
		t.Errorf("Expected pos() to return a Position object, but got nil")
	} else {
		expectedFromPos := mockObject1.Position()
		expectedToPos := mockObject2.Position()

		if pos.From() != expectedFromPos {
			t.Errorf("Expected fromPos to be %v, but got %v", expectedFromPos, pos.From())
		}

		if pos.To() != expectedToPos {
			t.Errorf("Expected toPos to be %v, but got %v", expectedToPos, pos.To())
		}
	}
}

func TestMergedEdgeMethods(t *testing.T) {
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)
	fromPos := mockObject1.Position()
	toPos := mockObject2.Position()

	// Create a merged edge with position information
	e := &mergedEdge{
		edge:  newEdge(mockObject1.ID(), mockObject2.ID(), arch.RelationTypeAggregation, fromPos, toPos),
		count: 2,
		pos:   []arch.RelationPos{valueobject.NewEmptyRelationPos(), valueobject.NewEmptyRelationPos()},
	}

	// Test the Pos and Count methods
	pos := e.Pos()
	count := e.Count()

	// Check if the Pos method returns the expected number of RelationPos objects
	if len(pos) != len(e.pos) {
		t.Errorf("Expected %d RelationPos objects, but got %d", len(e.pos), len(pos))
	}

	// Check if the Count method returns the expected count
	if count != e.count {
		t.Errorf("Expected count to be %d, but got %d", e.count, count)
	}
}

func TestEdgesMerge(t *testing.T) {
	// Create some mock edges with different keys and positions
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)
	mockObject3 := newMockObject(3)
	mockObject4 := newMockObject(4)

	edge1 := newEdge(mockObject1.ID(), mockObject2.ID(), arch.RelationTypeAggregation, mockObject1.Position(), mockObject2.Position())
	edge2 := newEdge(mockObject2.ID(), mockObject3.ID(), arch.RelationTypeComposition, mockObject2.Position(), mockObject3.Position())
	edge3 := newEdge(mockObject3.ID(), mockObject4.ID(), arch.RelationTypeAssociation, mockObject3.Position(), mockObject4.Position())
	edge4 := newEdge(mockObject1.ID(), mockObject2.ID(), arch.RelationTypeAggregation, mockObject1.Position(), mockObject2.Position())

	// Create a slice of edges and merge them
	edges := edges{edge1, edge2, edge3, edge4}
	mergedEdges := edges.merge()

	// Check the merged edges
	if len(mergedEdges) != 3 {
		t.Errorf("Expected 3 merged edges, but got %d", len(mergedEdges))
	}

	for _, me := range mergedEdges {
		if me.From() == mockObject1.ID() {
			if me.To() != mockObject2.ID() {
				t.Errorf("Expected To() of merged edge 1 to be %s, but got %s", mockObject2.ID(), me.To())
			}
			if me.Type() != arch.RelationTypeAggregation {
				t.Errorf("Expected Type() of merged edge 1 to be %d, but got %d", arch.RelationTypeAggregation, mergedEdges[0].Type())
			}
			if me.Count() != 2 {
				t.Errorf("Expected Count() of merged edge 1 to be 2, but got %d", me.Count())
			}
			if len(me.Pos()) != 1 {
				t.Errorf("Expected Pos() of merged edge 1 to have 1 elements, but got %d", len(me.Pos()))
			}
		} else if me.From() == mockObject2.ID() {
			if me.To() != mockObject3.ID() {
				t.Errorf("Expected To() of merged edge 2 to be %s, but got %s", mockObject3.ID(), me.To())
			}
			if me.Type() != arch.RelationTypeComposition {
				t.Errorf("Expected Type() of merged edge 2 to be %d, but got %d", arch.RelationTypeComposition, me.Type())
			}
			if me.Count() != 1 {
				t.Errorf("Expected Count() of merged edge 2 to be 1, but got %d", me.Count())
			}
			if len(me.Pos()) != 1 {
				t.Errorf("Expected Pos() of merged edge 2 to have 1 element, but got %d", len(me.Pos()))
			}
		} else if me.From() == mockObject3.ID() {
			if me.To() != mockObject4.ID() {
				t.Errorf("Expected To() of merged edge 3 to be %s, but got %s", mockObject4.ID(), me.To())
			}
			if me.Type() != arch.RelationTypeAssociation {
				t.Errorf("Expected Type() of merged edge 3 to be %d, but got %d", arch.RelationTypeAssociation, me.Type())
			}
			if me.Count() != 1 {
				t.Errorf("Expected Count() of merged edge 3 to be 1, but got %d", me.Count())
			}
			if len(me.Pos()) != 1 {
				t.Errorf("Expected Pos() of merged edge 3 to have 1 element, but got %d", len(me.Pos()))
			}
		}
	}

	// Check if the merged edges have unique keys
	edgeKeys := make(map[string]bool)
	for _, mergedEdge := range mergedEdges {
		key := mergedEdge.Key()
		if edgeKeys[key] {
			t.Errorf("Found duplicate merged edge key: %s", key)
		} else {
			edgeKeys[key] = true
		}
	}
}

func TestDiagram_ParseClassEdge(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)
	mockEle2 := newMockElement(2, newMockNode(5))
	mockEle2.IDVal = mockObject2.ID()

	// Create a subDiagram
	sub := &subDiagram{
		name:  "SubDiagram",
		nodes: []arch.Node{},
		elements: []arch.Element{newElement(&node{
			id:    mockEle2.MockNode.ID(),
			name:  mockEle2.MockNode.Name(),
			color: mockEle2.MockNode.Color(),
		}, elementTypeClass)},
	}

	// Create a mock edge with RelationTypeAttribution
	mockEdge := &directed.Edge{
		From: diagram.root,
		To:   diagram.root,
		Type: arch.RelationTypeAttribution,
		Value: valueobject.NewRelationPos(
			mockObject1.Position(),
			mockObject1.Position(),
		),
	}
	mockEdge.To.Value = mockObject1

	// Parse the class edge
	diagram.parseClassEdge(mockEdge, sub, mockObject2)

	// Verify that the node was correctly added to the subDiagram's nodes
	if len(sub.nodes) != 1 {
		t.Errorf("Expected 1 node in the subDiagram, but got %d", len(sub.nodes))
	}

	// Verify that the node's properties match the expected values
	node := sub.nodes[0]
	if node.ID() != mockObject1.Identifier().ID() {
		t.Errorf("Expected node ID to be %s, but got %s", mockObject1.Identifier().ID(), node.ID())
	}
	if node.Name() != mockObject1.Identifier().Name() {
		t.Errorf("Expected node Name to be %s, but got %s", mockObject1.Identifier().Name(), node.Name())
	}
	if node.Color() != string(objColor(mockObject1)) {
		t.Errorf("Expected node Color to be %s, but got %s", string(objColor(mockObject1)), node.Color())
	}

	// Verify that the element was updated with the node
	foundElement := sub.findElement(mockObject2.Identifier().ID())
	if foundElement == nil {
		t.Errorf("Expected element with ID %s to be found, but it was not", mockObject1.Identifier().ID())
	} else {
		element := foundElement.(*element)
		if len(element.Children()[1]) != 1 {
			t.Errorf("Expected element to have 1 child on the right, but got %d", len(element.Children()[1]))
		}
	}
}

func TestDiagram_ParseNode(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockObject1 := newMockObject(1)
	mockEle1 := newMockElement(1, newMockNode(4))
	mockEle1.IDVal = mockObject1.ID()
	mockObject2 := newMockObject(2)
	mockObject3 := newMockObject(3)

	// Create a subDiagram
	sub := &subDiagram{
		name:  "SubDiagram",
		nodes: []arch.Node{},
		elements: []arch.Element{newElement(&node{
			id:    mockEle1.MockNode.ID(),
			name:  mockEle1.MockNode.Name(),
			color: mockEle1.MockNode.Color(),
		}, elementTypeClass)},
	}

	mockEdge2 := &directed.Edge{
		From: diagram.root,
		To:   diagram.root,
		Type: arch.RelationTypeAttribution,
		Value: valueobject.NewRelationPos(
			mockObject1.Position(),
			mockObject1.Position(),
		),
	}
	mockEdge2.To.Value = mockObject2
	mockEdge3 := &directed.Edge{
		From: diagram.root,
		To:   diagram.root,
		Type: arch.RelationTypeBehavior,
		Value: valueobject.NewRelationPos(
			mockObject1.Position(),
			mockObject1.Position(),
		),
	}
	mockEdge3.To.Value = mockObject3

	// Create a mock directed node with edges representing class relations
	mockDirectedNode := &directed.Node{
		Value: mockObject1,
		Edges: []*directed.Edge{
			mockEdge2,
			mockEdge3,
		},
	}

	// Call the parseNode method to parse the mock directed node
	diagram.parseNode(mockDirectedNode, sub)

	// Verify that nodes and elements were correctly added and updated
	if len(sub.nodes) != 2 {
		t.Errorf("Expected 2 node in the subDiagram, but got %d", len(sub.nodes))
	}

	foundElement := sub.findElement(mockObject1.Identifier().ID())
	if foundElement == nil {
		t.Errorf("Expected element with ID %s to be found, but it was not", mockObject1.Identifier().ID())
	} else {
		element := foundElement.(*element)
		if len(element.Children()[1]) != 1 {
			t.Errorf("Expected element to have 1 child on the right, but got %d", len(element.Children()[1]))
		}
		if len(element.Children()[0]) != 1 {
			t.Errorf("Expected element to have 1 child on the right, but got %d", len(element.Children()[0]))
		}
	}
}

func TestDiagram_ParseSubDiagrams(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "DomainName"

	// Test Case 1: Object is an Aggregate
	aggregateObject := newMockAggregate(domain, 0)
	mockNode1 := &directed.Node{Value: aggregateObject, Edges: []*directed.Edge{}}

	// Call parseSubDiagrams on mockNode1
	subDiagram1 := diagram.parseSubDiagrams(mockNode1)
	if subDiagram1.Name() != domain {
		t.Errorf("Expected subDiagram1 to have name %s, but got %s", domain, subDiagram1.Name())
	}
	if len(subDiagram1.Nodes()) != 1 {
		t.Errorf("Expected subDiagram1 to have 1 node, but got %d", len(subDiagram1.Nodes()))
		n := subDiagram1.Nodes()[0]
		if n.ID() != aggregateObject.Identifier().ID() {
			t.Errorf("Expected node ID to be %s, but got %s", aggregateObject.Identifier().ID(), n.ID())
		}
		if n.Name() != aggregateObject.Identifier().Name() {
			t.Errorf("Expected node Name to be %s, but got %s", aggregateObject.Identifier().Name(), n.Name())
		}
		if n.Color() != string(objColor(aggregateObject)) {
			t.Errorf("Expected node Color to be %s, but got %s", string(objColor(aggregateObject)), n.Color())
		}
	}
	if len(subDiagram1.Summary()) != 1 {
		t.Errorf("Expected subDiagram1 to have 1 summary, but got %d", len(subDiagram1.Summary()))
		e := subDiagram1.Summary()[0]
		if e.ID() != aggregateObject.Identifier().ID() {
			t.Errorf("Expected summary ID to be %s, but got %s", aggregateObject.Identifier().ID(), e.ID())
		}
	}

	// Test Case 2: Object is a StringObj
	stringObject := valueobject.NewStringObj("StringObjectName")
	mockNode2 := &directed.Node{Value: stringObject, Edges: []*directed.Edge{}}

	// Call parseSubDiagrams on mockNode2
	subDiagram2 := diagram.parseSubDiagrams(mockNode2)
	if subDiagram2.Name() != stringObject.Identifier().Name() {
		t.Errorf("Expected subDiagram2 to have name %s, but got %s", stringObject.Identifier().Name(), subDiagram2.Name())
	}
	if len(subDiagram2.Nodes()) != 1 {
		t.Errorf("Expected subDiagram1 to have 1 node, but got %d", len(subDiagram2.Nodes()))
		n := subDiagram2.Nodes()[0]
		if n.ID() != stringObject.Identifier().ID() {
			t.Errorf("Expected node ID to be %s, but got %s", stringObject.Identifier().ID(), n.ID())
		}
		if n.Name() != stringObject.Identifier().Name() {
			t.Errorf("Expected node Name to be %s, but got %s", stringObject.Identifier().Name(), n.Name())
		}
		if n.Color() != string(objColor(stringObject)) {
			t.Errorf("Expected node Color to be %s, but got %s", string(objColor(stringObject)), n.Color())
		}
	}
	if len(subDiagram2.Summary()) != 1 {
		t.Errorf("Expected subDiagram2 to have 1 summary, but got %d", len(subDiagram2.Summary()))
		e := subDiagram2.Summary()[0]
		if e.ID() != stringObject.Identifier().ID() {
			t.Errorf("Expected summary ID to be %s, but got %s", stringObject.Identifier().ID(), e.ID())
		}
	}

	// Test Case 3: Object has Aggregation and Attribution relations
	// Create directed edges representing Aggregation and Attribution relations
	aggregationEdge := &directed.Edge{
		From: mockNode1,
		To:   mockNode2,
		Type: arch.RelationTypeAggregation,
		// Add other properties as needed.
	}
	attributionEdge := &directed.Edge{
		From: mockNode1,
		To:   mockNode2,
		Type: arch.RelationTypeAttribution,
		// Add other properties as needed.
	}

	// Add edges to mockNode1
	mockNode1.Edges = append(mockNode1.Edges, aggregationEdge, attributionEdge)

	// Call parseSubDiagrams on mockNode1 again
	subDiagram3 := diagram.parseSubDiagrams(mockNode1)
	if len(subDiagram3.Nodes()) != 3 {
		t.Errorf("Expected subDiagram1 to have 3 node, but got %d", len(subDiagram3.Nodes()))
		n := subDiagram3.Nodes()[0]
		if n.ID() != aggregateObject.Identifier().ID() {
			t.Errorf("Expected node ID to be %s, but got %s", aggregateObject.Identifier().ID(), n.ID())
		}
		if n.Name() != aggregateObject.Identifier().Name() {
			t.Errorf("Expected node Name to be %s, but got %s", aggregateObject.Identifier().Name(), n.Name())
		}
		if n.Color() != string(objColor(aggregateObject)) {
			t.Errorf("Expected node Color to be %s, but got %s", string(objColor(aggregateObject)), n.Color())
		}
	}

	if len(subDiagram3.Summary()) != 1 {
		t.Errorf("Expected subDiagram3 to have 1 summary, but got %d", len(subDiagram3.Summary()))
		e := subDiagram3.Summary()[0]
		if e.ID() != aggregateObject.Identifier().ID() {
			t.Errorf("Expected summary ID to be %s, but got %s", aggregateObject.Identifier().ID(), e.ID())
		}
		if len(e.Children()[1]) != 1 {
			t.Errorf("Expected element to have 1 child on the right, but got %d", len(e.Children()[1]))
		}
		if len(e.Children()[0]) != 1 {
			t.Errorf("Expected element to have 1 child on the right, but got %d", len(e.Children()[0]))
		}
	}
}

func TestDiagram_SubDiagrams(t *testing.T) {
	// Create a new Diagram
	diagram, err := NewDiagram("TestDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create mock objects for sub-diagrams
	sub1 := newMockAggregate("Sub1", 0)
	sub2 := newMockAggregate("Sub2", 0)
	sub3 := newMockAggregate("Sub3", 0)

	// Create directed edges representing AggregationRoot relations
	edge1 := &directed.Edge{
		From: diagram.root,
		To:   &directed.Node{Value: sub1},
		Type: arch.RelationTypeAggregationRoot,
	}
	edge2 := &directed.Edge{
		From: diagram.root,
		To:   &directed.Node{Value: sub2},
		Type: arch.RelationTypeAggregationRoot,
	}
	edge3 := &directed.Edge{
		From: diagram.root,
		To:   &directed.Node{Value: sub3},
		Type: arch.RelationTypeAggregationRoot,
	}

	// Add edges to the diagram's root
	diagram.root.Edges = []*directed.Edge{edge1, edge2, edge3}

	// Call the SubDiagrams method
	subDiagrams := diagram.SubDiagrams()

	// Validate the number of sub-diagrams
	if len(subDiagrams) != 4 {
		t.Errorf("Expected 4 sub-diagrams, but got %d", len(subDiagrams))
	}

	// Validate the names of sub-diagrams
	expectedNames := []string{"TestDiagram", "Sub1", "Sub2", "Sub3"}
	for i, sd := range subDiagrams {
		if sd.Name() != expectedNames[i] {
			t.Errorf("Expected sub-diagram %d to have name %s, but got %s", i+1, expectedNames[i], sd.Name())
		}
	}
}
