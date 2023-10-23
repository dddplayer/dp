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

func TestArch_BuildDirectory_Error(t *testing.T) {
	mockRepo := newMockObjectRepoWithInvalidDir()

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
	if err == nil {
		t.Error("Expected error, but got nil")
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
			Scope:   "test",
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

func TestDomainComponentRelations_DomainObjectError(t *testing.T) {
	name := "TestDiagram"

	mockDiagram, err := NewDiagram(name, arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	// Create some mock objects
	domain := "TestDomain"
	funcObj0 := newMockObjectFunction(0)
	funcObj1 := newMockDomainFunction(domain, newMockObjectFunction(1))

	// Add the mock objects to the Diagram
	_ = mockDiagram.AddObj(funcObj0)
	_ = mockDiagram.AddObj(funcObj1)
	mockDiagram.objs = append(mockDiagram.objs, funcObj0, funcObj1)

	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}

	mockArch := &Arch{
		CodeHandler:     &valueobject.CodeHandler{},
		relationDigraph: g,
		directory:       nil,
	}

	err = mockArch.domainComponentRelations(mockDiagram)
	if err == nil {
		t.Errorf("Expected error occurs, but got nil")
	}

	mockDiagram, err = NewDiagram(name, arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	// Create some mock objects
	funcObj2 := newMockDomainFunction(domain, newMockObjectFunction(2))
	funcObj3 := newMockObjectFunction(3)

	// Add the mock objects to the Diagram
	_ = mockDiagram.AddObj(funcObj2)
	_ = mockDiagram.AddObj(funcObj3)
	mockDiagram.objs = append(mockDiagram.objs, funcObj2, funcObj3)

	// Create a new RelationDigraph instance
	g = &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}

	mockArch = &Arch{
		CodeHandler:     &valueobject.CodeHandler{},
		relationDigraph: g,
		directory:       nil,
	}

	err = mockArch.domainComponentRelations(mockDiagram)
	if err == nil {
		t.Errorf("Expected error occurs, but got nil")
	}
}

func TestDomainComponentRelations_RelationMetaError(t *testing.T) {
	name := "TestDiagram"
	mockDiagram, err := NewDiagram(name, arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	// Create some mock objects
	domain := "TestDomain"
	funcObj0 := newMockDomainFunction(domain, newMockObjectFunction(0))
	funcObj1 := newMockDomainFunction(domain, newMockObjectFunction(1))

	// Add the mock objects to the Diagram
	_ = mockDiagram.AddObj(funcObj0)
	_ = mockDiagram.AddObj(funcObj1)
	mockDiagram.objs = append(mockDiagram.objs, funcObj0, funcObj1)

	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}

	mockArch := &Arch{
		CodeHandler:     &valueobject.CodeHandler{},
		relationDigraph: g,
		directory:       nil,
	}

	err = mockArch.domainComponentRelations(mockDiagram)
	if err == nil {
		t.Errorf("Expected error occurs, but got nil")
	}

}

func TestDomainComponentRelations_RelationError(t *testing.T) {
	name := "TestDiagram"
	mockDiagram, err := NewDiagram(name, arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	// Create some mock objects
	domain := "TestDomain"
	funcObj0 := newMockDomainFunction(domain, newMockObjectFunction(0))
	funcObj1 := newMockDomainFunction(domain, newMockObjectFunction(1))

	// Add the mock objects to the Diagram
	mockDiagram.objs = append(mockDiagram.objs, funcObj0, funcObj1)

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

	// Create edges for testing
	_ = g.AddEdge(funcObj0.OriginIdentifier().ID(), funcObj1.OriginIdentifier().ID(), arch.RelationTypeComposition, valueobject.NewEmptyRelationPos())

	mockArch := &Arch{
		CodeHandler:     &valueobject.CodeHandler{},
		relationDigraph: g,
		directory:       nil,
	}

	err = mockArch.domainComponentRelations(mockDiagram)
	if err == nil {
		t.Errorf("Expected error occurs, but got nil")
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
			Scope:   "test",
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

func TestGeneralGraph(t *testing.T) {
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
			Scope:   "test",
			ObjRepo: mockRepo,
			RelRepo: mockRelRepo,
		},
		relationDigraph: nil,
		directory:       nil,
	}

	// Call the StrategicGraph function
	diagram, err := mockArch.GeneralGraph(&MockOptions{
		ShowAllRel:            false,
		ShowStructEmbeddedRel: false,
	})

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	expectedEdgeCount := 6
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}

	// Call the StrategicGraph function
	diagram, err = mockArch.GeneralGraph(&MockOptions{
		ShowAllRel:            true,
		ShowStructEmbeddedRel: false,
	})

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	expectedEdgeCount = 6
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}

	// Call the StrategicGraph function
	diagram, err = mockArch.GeneralGraph(&MockOptions{
		ShowAllRel:            false,
		ShowStructEmbeddedRel: true,
	})

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	expectedEdgeCount = 6
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}
}

func TestBuildGeneralArchGraph(t *testing.T) {
	mockDirectory, objs := newMockDirectoryWithObjs()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
	}

	mockArch := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			Scope:   "testpackage",
			ObjRepo: mockRepo,
		},
		relationDigraph: &RelationDigraph{Graph: directed.NewDirectedGraph()},
		directory:       mockDirectory,
	}

	diagram, err := mockArch.buildGeneralArchGraph()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	expectedNodeCount := 7
	if len(diagram.Nodes) != expectedNodeCount {
		t.Errorf("Expected %d nodes, but got %d", expectedNodeCount, len(diagram.Nodes))
	}

	expectedEdgeCount := 6
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}
}

func TestComponentRelations(t *testing.T) {
	domain := "testpackage"
	mockDiagram, err := NewDiagram(domain, arch.PlainDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)
	mockObject3 := newMockObject(3)

	// Add the mock objects to the Diagram with different nesting depths
	_ = mockDiagram.AddObjTo(mockObject1, mockDiagram.Name(), arch.RelationTypeAggregationRoot)
	_ = mockDiagram.AddObjTo(mockObject2, mockObject1.ID(), arch.RelationTypeAggregationRoot)
	_ = mockDiagram.AddObjTo(mockObject3, mockObject2.ID(), arch.RelationTypeAggregationRoot)

	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}
	// Add the mock ObjIdentifier objects to the graph
	_ = g.AddObj(mockObject1.Identifier())
	_ = g.AddObj(mockObject2.Identifier())
	_ = g.AddObj(mockObject3.Identifier())

	// Create edges for testing
	_ = g.AddEdge(mockObject1.Identifier().ID(), mockObject2.Identifier().ID(), arch.RelationTypeComposition, valueobject.NewEmptyRelationPos())
	_ = g.AddEdge(mockObject2.Identifier().ID(), mockObject3.Identifier().ID(), arch.RelationTypeAggregation, valueobject.NewEmptyRelationPos())

	mockArch := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			Scope: domain,
		},
		relationDigraph: g,
		directory:       nil,
	}

	// 调用被测试的函数
	err = mockArch.componentRelations(mockDiagram)

	// 检查是否返回了错误
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	expectedNodeCount := 4
	if len(mockDiagram.Nodes) != expectedNodeCount {
		t.Errorf("Expected %d nodes, but got %d", expectedNodeCount, len(mockDiagram.Nodes))
	}

	expectedEdgeCount := 4
	if len(mockDiagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(mockDiagram.Edges()))
	}
}

func TestComponentAssociationRelations(t *testing.T) {
	domain := "testpackage"
	mockDiagram, err := NewDiagram(domain, arch.PlainDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	mockObject1 := newMockObject(1)
	mockObject2 := newMockObject(2)
	mockObject3 := newMockObject(3)

	// Add the mock objects to the Diagram with different nesting depths
	_ = mockDiagram.AddObjTo(mockObject1, mockDiagram.Name(), arch.RelationTypeAggregationRoot)
	_ = mockDiagram.AddObjTo(mockObject2, mockObject1.ID(), arch.RelationTypeAggregationRoot)
	_ = mockDiagram.AddObjTo(mockObject3, mockObject2.ID(), arch.RelationTypeAggregationRoot)

	// Create a new RelationDigraph instance
	g := &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}
	// Add the mock ObjIdentifier objects to the graph
	_ = g.AddObj(mockObject1.Identifier())
	_ = g.AddObj(mockObject2.Identifier())
	_ = g.AddObj(mockObject3.Identifier())

	// Create edges for testing
	_ = g.AddEdge(mockObject1.Identifier().ID(), mockObject2.Identifier().ID(), arch.RelationTypeAssociationOneOne, valueobject.NewEmptyRelationPos())
	_ = g.AddEdge(mockObject2.Identifier().ID(), mockObject3.Identifier().ID(), arch.RelationTypeAssociationOneMany, valueobject.NewEmptyRelationPos())

	mockArch := &Arch{
		CodeHandler: &valueobject.CodeHandler{
			Scope: domain,
		},
		relationDigraph: g,
		directory:       nil,
	}

	// 调用被测试的函数
	err = mockArch.componentAssociationRelations(mockDiagram)

	// 检查是否返回了错误
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	expectedNodeCount := 4
	if len(mockDiagram.Nodes) != expectedNodeCount {
		t.Errorf("Expected %d nodes, but got %d", expectedNodeCount, len(mockDiagram.Nodes))
	}

	expectedEdgeCount := 5
	if len(mockDiagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(mockDiagram.Edges()))
	}
}

func TestFilterAssociationMetas(t *testing.T) {
	metas := []arch.RelationMeta{
		NewMockRelationMeta(arch.RelationTypeAssociationOneOne, &MockRelationPos{
			fromPos: &MockPosition{FilenameVal: "/path/to/file1.txt", LineVal: 10, ColumnVal: 5},
			toPos:   &MockPosition{FilenameVal: "/path/to/file2.txt", LineVal: 20, ColumnVal: 15}}),
		NewMockRelationMeta(arch.RelationTypeAssociationOneMany, &MockRelationPos{
			fromPos: &MockPosition{FilenameVal: "/path/to/file1.txt", LineVal: 10, ColumnVal: 5},
			toPos:   &MockPosition{FilenameVal: "/path/to/file2.txt", LineVal: 20, ColumnVal: 15}}),
		NewMockRelationMeta(arch.RelationTypeBehavior, &MockRelationPos{
			fromPos: &MockPosition{FilenameVal: "/path/to/file1.txt", LineVal: 10, ColumnVal: 5},
			toPos:   &MockPosition{FilenameVal: "/path/to/file2.txt", LineVal: 20, ColumnVal: 15}}),
	}

	// 创建一个模拟的 Arch 实例
	arc := &Arch{}

	// 调用被测试的函数
	filteredMetas := arc.filterAssociationMetas(metas)

	// 检查过滤后的切片是否包含了关联关系类型的 RelationMeta 对象
	for _, meta := range filteredMetas {
		switch meta.Type() {
		case arch.RelationTypeAssociationOneOne, arch.RelationTypeAssociationOneMany, arch.RelationTypeAssociation:
			// 正确类型的 RelationMeta 对象
		default:
			t.Errorf("Unexpected RelationMeta type: %v", meta.Type())
		}
	}

	// 检查过滤后的切片长度是否与期望相符
	expectedLength := 2 // 两个关联关系类型的 RelationMeta 对象
	if len(filteredMetas) != expectedLength {
		t.Errorf("Expected %d RelationMeta objects, got %d", expectedLength, len(filteredMetas))
	}
}
