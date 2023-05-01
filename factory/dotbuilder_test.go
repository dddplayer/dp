package factory

import (
	"fmt"
	dot "github.com/dddplayer/core/dot/entity"
	"github.com/dddplayer/core/entity"
	"github.com/dddplayer/core/valueobject"
	"reflect"
	"testing"
)

func newEntity() *entity.Entity {
	pos := &valueobject.Position{Filename: "entity1.go", Offset: 10, Line: 5, Column: 15}
	e1 := &valueobject.Identifier{Name: "1", Path: "domain/test_domain/entity"}
	f1 := &valueobject.Identifier{Name: "func1", Path: "domain/test_domain/entity/1"}
	a1 := &valueobject.Identifier{Name: "attr1", Path: "domain/test_domain/entity/1"}
	et1 := entity.NewEntity(*e1, *pos)
	et1.AppendCommand(f1)
	et1.AppendAttribute(a1)
	return et1
}

func newValueObject() *entity.ValueObject {
	pos := &valueobject.Position{Filename: "vo1.go", Offset: 10, Line: 5, Column: 15}
	v1 := &valueobject.Identifier{Name: "1", Path: "domain/test_domain/valueobject"}
	f1 := &valueobject.Identifier{Name: "func1", Path: "domain/test_domain/valueobject/1"}
	a1 := &valueobject.Identifier{Name: "attr1", Path: "domain/test_domain/valueobject/1"}
	vo1 := entity.NewValueObject(*v1, *pos)
	vo1.AppendCommand(f1)
	vo1.AppendAttribute(a1)
	return vo1
}

func newService() *entity.Service {
	pos := &valueobject.Position{Filename: "service.go", Offset: 10, Line: 5, Column: 15}
	s1 := &valueobject.Identifier{Name: "1", Path: "domain/test_domain/service"}
	vo1 := entity.NewService(*s1, *pos)
	return vo1
}

func newFactory() *entity.Factory {
	pos := &valueobject.Position{Filename: "factory.go", Offset: 10, Line: 5, Column: 15}
	id := &valueobject.Identifier{Name: "1", Path: "domain/test_domain/factory"}
	f1 := entity.NewFactory(*id, *pos)
	return f1
}

func newInterfaceMethod() *entity.InterfaceMethod {
	pos := &valueobject.Position{Filename: "interface.go", Offset: 10, Line: 5, Column: 15}
	v1 := &valueobject.Identifier{Name: "method", Path: "domain/test_domain/entities/interface"}
	vo1 := entity.NewInterfaceMethod(*v1, *pos)
	return vo1
}

func newInterface() *entity.Interface {
	pos := &valueobject.Position{Filename: "interface.go", Offset: 10, Line: 5, Column: 15}
	v1 := &valueobject.Identifier{Name: "interface", Path: "domain/test_domain/entities"}
	vo1 := entity.NewInterface(*v1, *pos)
	vo1.Append(newInterfaceMethod())
	return vo1
}

func newTestDomain() *entity.Domain {
	d := &entity.Domain{
		Name:       "domain/test_domain",
		Components: []entity.DomainComponent{newEntity(), newValueObject(), newService(), newFactory()},
		Interfaces: []*entity.Interface{newInterface()},
	}

	return d
}

func TestDotBuilder_buildNode(t *testing.T) {
	db := &DotBuilder{}
	d := newTestDomain()
	node := db.buildNode(d)

	expectedNodeName := "test_domain"
	if node.Name() != expectedNodeName {
		t.Errorf("Unexpected node name: got %s, expected %s", node.Name(), expectedNodeName)
	}

	expectedNodeElementsLength := 5
	if len(node.Elements()) != expectedNodeElementsLength {
		t.Errorf("Unexpected node elements length: got %d, expected %d", len(node.Elements()), expectedNodeElementsLength)
	}

	// Check if all components are present in the node's elements
	service1 := newService()
	factory1 := newFactory()
	entity1 := newEntity()
	valueObject1 := newValueObject()
	interface1 := newInterface()

	for i, ele := range node.Elements() {
		switch i {
		case 0: // Service node
			if ele.Name() != "service" {
				t.Errorf("Expected service with name %s in node elements, got: %s", "service", ele.Name())
			}
			if len(ele.Attributes()) != 1 {
				t.Errorf("Expected 1 attribute in service node elements, got: %d", len(ele.Attributes()))
			}
			if a, ok := ele.Attributes()[0].(dot.DotAttribute); ok {
				if a.Name() != service1.Identifier().Name {
					t.Errorf("Expected service attribute with name %s, got: %s", service1.Identifier().Name, a.Name())
				}
			} else {
				t.Errorf("Expected service attribute with type DotAttribute, got: %v", reflect.TypeOf(ele.Attributes()[0]))
			}
		case 1: // Factory node
			if ele.Name() != "factory" {
				t.Errorf("Expected factory with name %s in node elements, got: %s", "factory", ele.Name())
			}
			if len(ele.Attributes()) != 1 {
				t.Errorf("Expected 1 attribute in factory node elements, got: %d", len(ele.Attributes()))
			}
			if a, ok := ele.Attributes()[0].(dot.DotAttribute); ok {
				if a.Name() != service1.Identifier().Name {
					t.Errorf("Expected factory attribute with name %s, got: %s", factory1.Identifier().Name, a.Name())
				}
			} else {
				t.Errorf("Expected factory attribute with type DotAttribute, got: %v", reflect.TypeOf(ele.Attributes()[0]))
			}
		case 2: // Entity node
			if ele.Name() != entity1.Identifier().Name {
				t.Errorf("Expected entity with name %s in node elements, got: %s", entity1.Identifier().Name, ele.Name())
			}
			if len(ele.Attributes()) != 2 {
				t.Errorf("Expected 2 attribute in entity node elements, got: %d", len(ele.Attributes()))
			}
			if a, ok := ele.Attributes()[0].([]dot.DotAttribute); ok {
				if a[0].Name() != entity1.Commands()[0].Name {
					t.Errorf("Expected entity attribute with name %s, got: %s", entity1.Commands()[0].Name, a[0].Name())
				}
			} else {
				t.Errorf("Expected entity attribute with type []DotAttribute, got: %v", reflect.TypeOf(ele.Attributes()[0]))
			}
			if a, ok := ele.Attributes()[1].([]dot.DotAttribute); ok {
				if a[0].Name() != entity1.Attributes()[0].Name {
					t.Errorf("Expected entity attribute with name %s, got: %s", entity1.Attributes()[0].Name, a[0].Name())
				}
			} else {
				t.Errorf("Expected entity attribute with type []DotAttribute, got: %v", reflect.TypeOf(ele.Attributes()[1]))
			}

		case 3: // ValueObject node
			if ele.Name() != valueObject1.Identifier().Name {
				t.Errorf("Expected valueObject with name %s in node elements, got: %s", valueObject1.Identifier().Name, ele.Name())
			}
			if len(ele.Attributes()) != 2 {
				t.Errorf("Expected 2 attribute in valueObject node elements, got: %d", len(ele.Attributes()))
			}
			if a, ok := ele.Attributes()[0].([]dot.DotAttribute); ok {
				if a[0].Name() != valueObject1.Commands()[0].Name {
					t.Errorf("Expected valueObject attribute with name %s, got: %s", valueObject1.Commands()[0].Name, a[0].Name())
				}
			} else {
				t.Errorf("Expected valueObject attribute with type []DotAttribute, got: %v", reflect.TypeOf(ele.Attributes()[0]))
			}
			if a, ok := ele.Attributes()[1].([]dot.DotAttribute); ok {
				if a[0].Name() != valueObject1.Attributes()[0].Name {
					t.Errorf("Expected valueObject attribute with name %s, got: %s", valueObject1.Attributes()[0].Name, a[0].Name())
				}
			} else {
				t.Errorf("Expected valueObject attribute with type []DotAttribute, got: %v", reflect.TypeOf(ele.Attributes()[1]))
			}

		case 4: // Interface node
			if ele.Name() != interface1.Identifier().Name {
				t.Errorf("Expected Interface with name %s in node elements, got: %s", interface1.Identifier().Name, ele.Name())
			}
			if len(ele.Attributes()) != 1 {
				t.Errorf("Expected 1 attribute in valueObject node elements, got: %d", len(ele.Attributes()))
			}
			if a, ok := ele.Attributes()[0].([]dot.DotAttribute); ok {
				if a[0].Name() != interface1.Methods[0].Name {
					t.Errorf("Expected Interface attribute with name %s, got: %s", interface1.Methods[0].Name, a[0].Name())
				}
			} else {
				t.Errorf("Expected valueObject attribute with type []DotAttribute, got: %v", reflect.TypeOf(ele.Attributes()[0]))
			}
		}
	}
}

func TestDotBuilder_buildNodes(t *testing.T) {
	db := &DotBuilder{}
	d := &entity.Domain{
		Name: "domain/test_domain",
		SubDomains: []*entity.Domain{
			{Name: "domain/test_domain/1"},
			{
				Name: "domain/test_domain/2",
				SubDomains: []*entity.Domain{
					{
						Name:       "domain/test_domain/2/1",
						SubDomains: []*entity.Domain{},
					},
					{
						Name:       "domain/test_domain/2/2",
						SubDomains: []*entity.Domain{},
					},
				},
			},
		},
	}

	g := &valueobject.DotGraph{}
	if err := db.buildNodes(d, g); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedNodes := 5
	actualNodes := len(g.Nodes())
	if actualNodes != expectedNodes {
		t.Errorf("unexpected number of nodes, expected %d but got %d", expectedNodes, actualNodes)
	}
}

func TestDotBuilder_buildEdges(t *testing.T) {
	db := &DotBuilder{}
	from := valueobject.Identifier{Name: "from", Path: "test/entity/fromPath"}
	to := valueobject.Identifier{Name: "to", Path: "test/valueobject/toPath"}
	db.relations = []*valueobject.Relation{
		{From: &from, To: &to},
	}
	g := &valueobject.DotGraph{}
	err := db.buildEdges(g)
	if err != nil {
		t.Errorf("buildEdges returned unexpected error: %v", err)
	}
	expectedEdge := valueobject.NewDotEdge("fromPath:fromPath.from", "toPath:toPath.to")
	if len(g.Edges()) != 1 || g.Edges()[0].From() != expectedEdge.From() || g.Edges()[0].To() != expectedEdge.To() {
		t.Errorf("buildEdges got from: %s, to: %s, Expected: from: %s, to: %s",
			g.Edges()[0].From(), g.Edges()[0].To(), expectedEdge.From(), expectedEdge.To())
	}
}

func TestDotBuilder_buildGraph(t *testing.T) {
	// setup
	db := &DotBuilder{
		Domain: &entity.Domain{
			Name: "test",
			SubDomains: []*entity.Domain{
				{Name: "sub1"},
				{Name: "sub2"},
			},
		},
	}
	from := valueobject.Identifier{Name: "from", Path: "test/entity/fromPath"}
	to := valueobject.Identifier{Name: "to", Path: "test/valueobject/toPath"}
	db.relations = []*valueobject.Relation{
		{From: &from, To: &to},
	}

	// execute
	got, err := db.buildGraph()

	// verify
	if err != nil {
		t.Errorf("buildGraph returned unexpected error: %v", err)
	}
	if len(got.Nodes()) != 3 {
		t.Errorf("buildGraph got %d nodes, expected %d", len(got.Nodes()), 3)
	}
	if len(got.Edges()) != 1 {
		t.Errorf("buildGraph got %d edges, expected %d", len(got.Edges()), 1)
	}
}

func TestDotBuilder_buildDomain(t *testing.T) {
	db := &DotBuilder{
		Domain: entity.NewDomain("domain"),
	}

	pos := &valueobject.Position{Filename: "entity1.go", Offset: 10, Line: 5, Column: 15}
	id := &valueobject.Identifier{Name: "entity1", Path: "domain/test_domain"}
	comp := entity.NewEntity(*id, *pos)

	db.buildDomain(comp)

	// Ensure that the component was added to the domain's components slice
	if len(db.Domain.SubDomains) != 1 {
		t.Fatalf("expected 1 subdomain, got %d", len(db.Domain.SubDomains))
	}
	subdomain := db.Domain.SubDomains[0]
	if subdomain.Name != "domain/test_domain" {
		t.Errorf("expected subdomain name to be 'domain/test_domain', got '%s'", subdomain.Name)
	}
	if len(subdomain.Components) != 1 {
		t.Fatalf("expected subdomain to have 1 component, got %d", len(subdomain.Components))
	}
	if subdomain.Components[0].Identifier().Name != "entity1" {
		t.Errorf("expected component name to be 'entity1', got '%s'", subdomain.Components[0].Identifier().Name)
	}

	// Test adding an interface
	ipos := &valueobject.Position{Filename: "interface.go", Offset: 10, Line: 5, Column: 15}
	iid := &valueobject.Identifier{Name: "interface1", Path: "domain/test_domain"}
	iface := entity.NewInterface(*iid, *ipos)
	db.buildDomain(iface)
	if len(subdomain.Interfaces) != 1 {
		t.Fatalf("expected subdomain to have 1 interface, got %d", len(subdomain.Interfaces))
	}
	if subdomain.Interfaces[0].Identifier().Name != "interface1" {
		t.Errorf("expected interface name to be 'interface1', got '%s'", subdomain.Interfaces[0].Identifier().Name)
	}
}

func TestDotBuilder_Build(t *testing.T) {
	// Define mock repository with test data
	mockData := map[valueobject.Identifier]entity.DomainObject{
		{Path: "domain/component1"}: entity.NewEntity(
			valueobject.Identifier{Name: "entity1", Path: "domain/component1"},
			valueobject.Position{Filename: "entity1.go", Offset: 10, Line: 5, Column: 15}),
		{Path: "domain/component2"}: entity.NewEntity(
			valueobject.Identifier{Name: "entity2", Path: "domain/component2"},
			valueobject.Position{Filename: "entity2.go", Offset: 10, Line: 5, Column: 15}),
		{Path: "domain/component3"}: entity.NewEntity(
			valueobject.Identifier{Name: "entity3", Path: "domain/component3"},
			valueobject.Position{Filename: "entity3.go", Offset: 10, Line: 5, Column: 15}),
		{Path: "domain/relation1"}: &valueobject.Relation{
			From: &valueobject.Identifier{Name: "entity1", Path: "domain/component1"},
			To:   &valueobject.Identifier{Name: "entity2", Path: "domain/component2"},
		},
		{Path: "domain/relation2"}: &valueobject.Relation{
			From: &valueobject.Identifier{Name: "entity1", Path: "domain/component1"},
			To:   &valueobject.Identifier{Name: "entity3", Path: "domain/component3"},
		},
	}
	mockRepo := &mockRepository{data: mockData}

	// Create DotBuilder and build graph
	db := &DotBuilder{Repo: mockRepo, Domain: entity.NewDomain("domain")}
	g, err := db.Build()
	if err != nil {
		t.Fatalf("Build returned unexpected error: %v", err)
	}

	// Ensure graph has the expected number of nodes and edges
	expectedNodes := 4 // domain, component1, component2, component3
	expectedEdges := 2 // relation1, relation2
	if len(g.Nodes()) != expectedNodes {
		t.Errorf("expected %d nodes, got %d", expectedNodes, len(g.Nodes()))
	}
	if len(g.Edges()) != expectedEdges {
		t.Errorf("expected %d edges, got %d", expectedEdges, len(g.Edges()))
	}

	// Ensure nodes are correct
	expectedNodeNames := []string{"component1", "component2", "component3"}
	for _, nodeName := range expectedNodeNames {
		if !containsNodeName(g.Nodes(), nodeName) {
			t.Errorf("expected node '%s' to be in graph", nodeName)
		}
	}

	// Ensure edges are correct
	expectedEdgePaths := []string{"component1:component1.entity1 - component2:component2.entity2", "component1:component1.entity1 - component3:component3.entity3"}
	for _, edgePath := range expectedEdgePaths {
		if !containsEdgePath(g.Edges(), edgePath) {
			t.Errorf("expected edge '%s' to be in graph", edgePath)
		}
	}
}

// Utility function to check if a node with the given name exists in the slice of nodes
func containsNodeName(nodes []dot.DotNode, name string) bool {
	for _, node := range nodes {
		if node.Name() == name {
			return true
		}
	}
	return false
}

// Utility function to check if an edge with the given path exists in the slice of edges
func containsEdgePath(edges []dot.DotEdge, path string) bool {
	for _, edge := range edges {
		if fmt.Sprintf("%s - %s", edge.From(), edge.To()) == path {
			return true
		}
	}
	return false
}
