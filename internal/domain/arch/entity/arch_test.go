package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"github.com/dddplayer/dp/internal/domain/code"
	"github.com/dddplayer/dp/pkg/datastructure/directed"
	"testing"
)

func TestArch_ObjectHandler(t *testing.T) {
	// Create an instance of Arch
	arc := &Arch{
		CodeHandler:     &valueobject.CodeHandler{},
		relationDigraph: &RelationDigraph{},
		directory:       &Directory{},
	}

	// Call the ObjectHandler method
	objectHandler := arc.ObjectHandler()

	// Check if the returned object is of type code.Handler
	_, ok := objectHandler.(code.Handler)
	if !ok {
		t.Errorf("Expected ObjectHandler to return an object of type code.Handler")
	}
}

func TestArch_BuildDirectory(t *testing.T) {
	claObj1 := newMockObjectWithId("test/cmd", "cla1", 1)
	claObj2 := newMockObjectWithId("test/internal/domain/testdomain", "cla2", 1)
	claObj3 := newMockObjectWithId("test/pkg", "cla3", 1)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(claObj1)
	_ = mockRepo.Insert(claObj2)
	_ = mockRepo.Insert(claObj3)

	arc := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			ObjRepo: mockRepo,
		},
		relationDigraph: &RelationDigraph{},
		directory:       nil, // Initialize directory as nil for testing purposes
	}

	// Call the buildDirectory method
	err := arc.buildDirectory()

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the directory was built correctly
	if arc.directory == nil {
		t.Error("Expected directory to be non-nil after calling buildDirectory")
	}
}

func TestBuildOriginGraph(t *testing.T) {
	claObj1 := newMockObject(1)
	claObj2 := newMockObject(2)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(claObj1)
	_ = mockRepo.Insert(claObj2)

	mockRelRepo := &MockRelationRepository{
		relations: make([]arch.Relation, 0),
	}
	// Create mock Relation objects for testing
	depRel := &MockDependenceRelation{
		from:      claObj1,
		dependsOn: claObj2,
	}
	compRel := &MockCompositionRelation{
		from:  claObj1,
		child: claObj2,
	}
	_ = mockRelRepo.Insert(depRel)
	_ = mockRelRepo.Insert(compRel)

	// Create an instance of Arch
	arc := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			ObjRepo: mockRepo,
			RelRepo: mockRelRepo,
		},
		relationDigraph: nil,
		directory:       nil,
	}

	err := arc.buildOriginGraph()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(arc.relationDigraph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, but got %d", len(arc.relationDigraph.Nodes))
	}
}

func TestBuildPlain(t *testing.T) {
	claObj1 := newMockObjectWithId("test/cmd", "cla1", 1)
	claObj2 := newMockObjectWithId("test/internal/domain/testdomain", "cla2", 1)
	claObj3 := newMockObjectWithId("test/pkg", "cla3", 1)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(claObj1)
	_ = mockRepo.Insert(claObj2)
	_ = mockRepo.Insert(claObj3)

	mockRelRepo := &MockRelationRepository{
		relations: make([]arch.Relation, 0),
	}
	// Create mock Relation objects for testing
	depRel := &MockDependenceRelation{
		from:      claObj1,
		dependsOn: claObj2,
	}
	compRel := &MockCompositionRelation{
		from:  claObj1,
		child: claObj2,
	}
	_ = mockRelRepo.Insert(depRel)
	_ = mockRelRepo.Insert(compRel)

	arc := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			ObjRepo: mockRepo,
			RelRepo: mockRelRepo,
		},
		relationDigraph: nil,
		directory:       nil,
	}

	err := arc.BuildPlain()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestBuildHexagon(t *testing.T) {
	claObj1 := newMockObjectWithId("test/cmd", "cla1", 1)
	claObj2 := newMockObjectWithId("test/internal/domain/testdomain", "cla2", 1)
	claObj3 := newMockObjectWithId("test/pkg", "cla3", 1)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(claObj1)
	_ = mockRepo.Insert(claObj2)
	_ = mockRepo.Insert(claObj3)

	mockRelRepo := &MockRelationRepository{
		relations: make([]arch.Relation, 0),
	}
	// Create mock Relation objects for testing
	depRel := &MockDependenceRelation{
		from:      claObj1,
		dependsOn: claObj2,
	}
	compRel := &MockCompositionRelation{
		from:  claObj1,
		child: claObj2,
	}
	_ = mockRelRepo.Insert(depRel)
	_ = mockRelRepo.Insert(compRel)

	arc := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			ObjRepo: mockRepo,
			RelRepo: mockRelRepo,
		},
		relationDigraph: nil,
		directory:       nil,
	}

	err := arc.BuildHexagon()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestGenerateCombinations(t *testing.T) {
	// Create an array of mock objects for testing
	mockObjects := []arch.Object{
		newMockObject(1),
		newMockObject(2),
		newMockObject(3),
	}

	// Generate combinations
	combinations := generateCombinations(mockObjects)

	// Verify the number of combinations
	expectedCombinations := len(mockObjects) * (len(mockObjects) - 1)
	if len(combinations) != expectedCombinations {
		t.Errorf("Expected %d combinations, but got %d", expectedCombinations, len(combinations))
	}

	// Verify that each combination contains different objects
	for i, combo1 := range combinations {
		for j, combo2 := range combinations {
			if i != j {
				if combo1.First == combo2.First && combo1.Second == combo2.Second {
					t.Errorf("Duplicate combination found: %v", combo1)
				}
			}
		}
	}
}

func TestSummaryDomainComponentRelations(t *testing.T) {
	name := "TestDiagram"
	mockDiagram, err := NewDiagram(name, arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	// Create some mock objects
	domain := "TestDomain"
	funcObj0 := newMockDomainFunction(domain, newMockObjectFunction(0))
	funcObj1 := newMockDomainFunction(domain, newMockObjectFunction(1))
	funcObj2 := newMockDomainFunction(domain, newMockObjectFunction(2))
	funcObj3 := newMockDomainFunction(domain, newMockObjectFunction(3))

	// Add the mock objects to the Diagram
	_ = mockDiagram.AddObj(funcObj0)
	_ = mockDiagram.AddObj(funcObj1)
	_ = mockDiagram.AddObj(funcObj2)
	_ = mockDiagram.AddObj(funcObj3)
	mockDiagram.objs = append(mockDiagram.objs, funcObj0, funcObj1, funcObj2, funcObj3)

	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}
	// Add the mock ObjIdentifier objects to the graph
	err = g.AddObj(funcObj0.OriginIdentifier())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(funcObj1.OriginIdentifier())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(funcObj2.OriginIdentifier())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(funcObj3.OriginIdentifier())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create edges for testing
	_ = g.AddEdge(funcObj0.OriginIdentifier().ID(), funcObj1.OriginIdentifier().ID(), arch.RelationTypeComposition, valueobject.NewEmptyRelationPos())
	_ = g.AddEdge(funcObj1.OriginIdentifier().ID(), funcObj2.OriginIdentifier().ID(), arch.RelationTypeEmbedding, valueobject.NewEmptyRelationPos())
	_ = g.AddEdge(funcObj2.OriginIdentifier().ID(), funcObj3.OriginIdentifier().ID(), arch.RelationTypeAssociation, valueobject.NewEmptyRelationPos())

	mockArch := &Arch{
		CodeHandler:     &valueobject.CodeHandler{},
		relationDigraph: g,
		directory:       nil,
	}

	err = mockArch.summaryDomainComponentRelations(mockDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	var totalEdges int
	for _, n := range mockDiagram.Nodes {
		totalEdges += len(n.Edges)
	}
	if totalEdges != (4 + 3 + 1 + 0) {
		t.Errorf("Expected 8 edges, but got %d", totalEdges)
	}
}

func TestBuildStrategicArchGraph(t *testing.T) {
	mockDirectory, objs := newMockDirectoryWithDomainObjs()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
	}

	mockArch := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			Scope:   "TestDomain",
			ObjRepo: mockRepo,
		},
		relationDigraph: &RelationDigraph{Graph: directed.NewDirectedGraph()},
		directory:       mockDirectory,
	}

	diagram, err := mockArch.buildStrategicArchGraph()

	if err == nil {
		t.Errorf("Expected error occurs, but got nil")
	} else if err.Error() != "aggregate test has no entity" {
		t.Errorf("Expected error message: aggregate test has no entity, but got: %v", err.Error())
	}

	mockDirectory, objs = newMockDirectoryWithAggregate()
	mockRepo = &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
		_ = mockArch.relationDigraph.AddObj(mockObj.Identifier())
	}
	mockArch.directory = mockDirectory
	mockArch.ObjRepo = mockRepo

	diagram, err = mockArch.buildStrategicArchGraph()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	expectedNodeCount := 3
	if len(diagram.Nodes) != expectedNodeCount {
		t.Errorf("Expected %d nodes, but got %d", expectedNodeCount, len(diagram.Nodes))
	}

	expectedEdgeCount := 2
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}
}

func TestStrategicGraph(t *testing.T) {
	claObj1 := newMockObjectWithId("test/cmd", "cla1", 1)
	claObj2 := newMockObjectWithId("test/internal/domain/testdomain", "cla2", 1)
	claObj20 := newMockClassWithName("test/internal/domain/testdomain/entity", "testdomain")
	claObj3 := newMockObjectWithId("test/pkg", "cla3", 1)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(claObj1)
	_ = mockRepo.Insert(claObj2)
	_ = mockRepo.Insert(claObj20)
	_ = mockRepo.Insert(claObj3)

	mockRelRepo := &MockRelationRepository{
		relations: make([]arch.Relation, 0),
	}
	// Create mock Relation objects for testing
	depRel := &MockDependenceRelation{
		from:      claObj1,
		dependsOn: claObj2,
	}
	compRel := &MockCompositionRelation{
		from:  claObj1,
		child: claObj2,
	}
	_ = mockRelRepo.Insert(depRel)
	_ = mockRelRepo.Insert(compRel)

	mockArch := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			ObjRepo: mockRepo,
			RelRepo: mockRelRepo,
		},
		relationDigraph: nil,
		directory:       nil,
	}

	// Call the StrategicGraph function
	diagram, err := mockArch.StrategicGraph()

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	expectedEdgeCount := 1
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}
}

func TestDomainComponentRelations(t *testing.T) {
	name := "TestDiagram"
	mockDiagram, err := NewDiagram(name, arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	// Create some mock objects
	domain := "TestDomain"
	funcObj0 := newMockDomainFunction(domain, newMockObjectFunction(0))
	funcObj1 := newMockDomainFunction(domain, newMockObjectFunction(1))
	funcObj2 := newMockDomainFunction(domain, newMockObjectFunction(2))
	funcObj3 := newMockDomainFunction(domain, newMockObjectFunction(3))

	// Add the mock objects to the Diagram
	_ = mockDiagram.AddObj(funcObj0)
	_ = mockDiagram.AddObj(funcObj1)
	_ = mockDiagram.AddObj(funcObj2)
	_ = mockDiagram.AddObj(funcObj3)
	mockDiagram.objs = append(mockDiagram.objs, funcObj0, funcObj1, funcObj2, funcObj3)

	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}
	// Add the mock ObjIdentifier objects to the graph
	err = g.AddObj(funcObj0.OriginIdentifier())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(funcObj1.OriginIdentifier())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(funcObj2.OriginIdentifier())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = g.AddObj(funcObj3.OriginIdentifier())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create edges for testing
	_ = g.AddEdge(funcObj0.OriginIdentifier().ID(), funcObj1.OriginIdentifier().ID(), arch.RelationTypeComposition, valueobject.NewEmptyRelationPos())
	_ = g.AddEdge(funcObj1.OriginIdentifier().ID(), funcObj2.OriginIdentifier().ID(), arch.RelationTypeEmbedding, valueobject.NewEmptyRelationPos())
	_ = g.AddEdge(funcObj2.OriginIdentifier().ID(), funcObj3.OriginIdentifier().ID(), arch.RelationTypeAssociation, valueobject.NewEmptyRelationPos())

	mockArch := &Arch{
		CodeHandler:     &valueobject.CodeHandler{},
		relationDigraph: g,
		directory:       nil,
	}

	err = mockArch.domainComponentRelations(mockDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	var totalEdges int
	for _, n := range mockDiagram.Nodes {
		totalEdges += len(n.Edges)
	}
	if totalEdges != 3 {
		t.Errorf("Expected 8 edges, but got %d", totalEdges)
	}
}

func TestBuildTacticArchGraph(t *testing.T) {
	mockDirectory, objs := newMockDirectoryWithDomainObjs()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
	}

	mockArch := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			Scope:   "TestDomain",
			ObjRepo: mockRepo,
		},
		relationDigraph: &RelationDigraph{Graph: directed.NewDirectedGraph()},
		directory:       mockDirectory,
	}

	diagram, err := mockArch.buildTacticArchGraph()

	if err == nil {
		t.Errorf("Expected error occurs, but got nil")
	} else if err.Error() != "aggregate test has no entity" {
		t.Errorf("Expected error message: aggregate test has no entity, but got: %v", err.Error())
	}

	mockDirectory, objs = newMockDirectoryWithAggregate()
	mockRepo = &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
		_ = mockArch.relationDigraph.AddObj(mockObj.Identifier())
	}
	mockArch.directory = mockDirectory
	mockArch.ObjRepo = mockRepo

	diagram, err = mockArch.buildTacticArchGraph()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	expectedNodeCount := 5
	if len(diagram.Nodes) != expectedNodeCount {
		t.Errorf("Expected %d nodes, but got %d", expectedNodeCount, len(diagram.Nodes))
	}

	expectedEdgeCount := 3
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}
}

func TestTacticGraph(t *testing.T) {
	claObj1 := newMockObjectWithId("test/cmd", "cla1", 1)
	claObj2 := newMockObjectWithId("test/internal/domain/testdomain", "cla2", 1)
	claObj20 := newMockClassWithName("test/internal/domain/testdomain/entity", "testdomain")
	claObj3 := newMockObjectWithId("test/pkg", "cla3", 1)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(claObj1)
	_ = mockRepo.Insert(claObj2)
	_ = mockRepo.Insert(claObj20)
	_ = mockRepo.Insert(claObj3)

	mockRelRepo := &MockRelationRepository{
		relations: make([]arch.Relation, 0),
	}
	// Create mock Relation objects for testing
	depRel := &MockDependenceRelation{
		from:      claObj1,
		dependsOn: claObj2,
	}
	compRel := &MockCompositionRelation{
		from:  claObj1,
		child: claObj2,
	}
	_ = mockRelRepo.Insert(depRel)
	_ = mockRelRepo.Insert(compRel)

	mockArch := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			ObjRepo: mockRepo,
			RelRepo: mockRelRepo,
		},
		relationDigraph: nil,
		directory:       nil,
	}

	// Call the StrategicGraph function
	diagram, err := mockArch.TacticGraph(&MockOptions{
		ShowAllRel: true,
	})

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	expectedEdgeCount := 2
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}
}
