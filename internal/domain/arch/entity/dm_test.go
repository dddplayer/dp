package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"path"
	"testing"
)

func TestNewDomainModel(t *testing.T) {
	mockRepo := &MockObjectRepository{}
	mockDirectory := newMockEmptyDirectory()

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if model == nil {
		t.Errorf("Expected non-nil model, but got nil")
	}
	if model.repo != mockRepo {
		t.Errorf("Expected repo to be mockRepo, but got %v", model.repo)
	}
	if len(model.aggregates) != 0 {
		t.Errorf("Expected aggregate group to be 0, but got %v", len(model.aggregates))
	}
}

func TestDomainModel_DomainName(t *testing.T) {
	mockRepo := &MockObjectRepository{}
	mockDirectory := newMockEmptyDirectory()

	dm, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domainName, err := dm.DomainName()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	expectedDomainName := "root/internal/domain"
	if domainName != expectedDomainName {
		t.Errorf("Expected domain name %s, but got: %s", expectedDomainName, domainName)
	}
}

func TestDomainModel_DomainName_Error(t *testing.T) {
	mockRepo := &MockObjectRepository{}
	mockDirectory := newMockInvalidEmptyDirectory()

	dm, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	_, err = dm.DomainName()
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildAbstractDomainComponent(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test case: Unsupported component type
	err = model.buildAbstractComponent(mockDiagram, mockGroup, mockGroup.Name(), "InvalidComponentType")
	if err == nil {
		t.Errorf("Expected an error for unsupported component type, but got nil")
	}
	expectedErrMsg := "unsupported objs type: InvalidComponentType"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%v'", expectedErrMsg, err)
	}

	err = model.buildAbstractComponent(mockDiagram, mockGroup, mockGroup.Name(), valueobject.GeneralComponent)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.buildAbstractComponent(mockDiagram, mockGroup, mockGroup.Name(), valueobject.FunctionComponent)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the relationships in the diagram
	// Get the list of objects from the diagram
	objects := mockDiagram.Objects()

	// Verify the number of objects in the diagram
	expectedNodeCount := 2
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	// Get the list of edges from the diagram
	edges := mockDiagram.Edges()

	// Verify the number of edges in the diagram
	expectedEdgeCount := 3 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestBuildAbstractDomainComponent_Error(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	componentKey := path.Join(mockGroup.Name(), string(valueobject.GeneralComponent))
	_ = mockDiagram.AddStringTo(componentKey, mockGroup.Name(), arch.RelationTypeAbstraction)
	err = model.buildAbstractComponent(mockDiagram, mockGroup, mockGroup.Name(), valueobject.GeneralComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	funcObj := newMockObjectFunction(0)
	df := newMockDomainFunction("testdomain", funcObj)
	_ = mockDiagram.AddObjTo(df, componentKey, arch.RelationTypeAggregation)
	err = model.buildAbstractComponent(mockDiagram, mockGroup, mockGroup.Name(), valueobject.FunctionComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildDMComponents(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.buildComponents(mockDiagram, mockGroup, mockGroup.Name(), valueobject.ClassComponent)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.buildComponents(mockDiagram, mockGroup, mockGroup.Name(), valueobject.InterfaceComponent)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test case: Unsupported component type
	err = model.buildComponents(mockDiagram, mockGroup, mockGroup.Name(), "InvalidComponentType")
	if err == nil {
		t.Errorf("Expected an error for unsupported component type, but got nil")
	}
	expectedErrMsg := "unsupported objs type: InvalidComponentType"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%v'", expectedErrMsg, err)
	}

	// Verify the relationships in the diagram
	// Get the list of objects from the diagram
	objects := mockDiagram.Objects()

	// Verify the number of objects in the diagram
	expectedNodeCount := 5
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	// Get the list of edges from the diagram
	edges := mockDiagram.Edges()

	// Verify the number of edges in the diagram
	expectedEdgeCount := 1 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestBuildDMComponents_Error(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"
	claObj := newMockObject(0)
	claAttrObj := newMockObjectAttribute(0)
	claMethodObj := newMockObjectMethod(0)
	df := valueobject.NewDomainFunction(valueobject.NewFunction(claMethodObj, claObj.Identifier()), domain)
	mockClass := newMockDomainClass(domain, claObj, claAttrObj, claMethodObj)

	_ = mockDiagram.AddObjTo(df, mockClass.Identifier().ID(), arch.RelationTypeBehavior)
	err = model.buildComponents(mockDiagram, mockGroup, mockGroup.Name(), valueobject.ClassComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
	err = model.buildComponents(mockDiagram, mockGroup, mockGroup.Name(), valueobject.ClassComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	intfObj := newMockObjectInterface(0)
	intfMethodObj := newMockObjectInterfaceMethod(0)
	dif := valueobject.NewDomainFunction(valueobject.NewFunction(intfMethodObj, intfObj.Identifier()), domain)
	mockItf := newMockDomainInterface(domain, intfObj, []*MockObject{intfMethodObj})

	_ = mockDiagram.AddObjTo(dif, mockItf.Identifier().ID(), arch.RelationTypeBehavior)
	err = model.buildComponents(mockDiagram, mockGroup, mockGroup.Name(), valueobject.InterfaceComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
	err = model.buildComponents(mockDiagram, mockGroup, mockGroup.Name(), valueobject.InterfaceComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildDomainComponents(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.buildDomainComponents(mockDiagram, mockGroup, mockGroup.Name())
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the relationships in the diagram
	// Get the list of objects from the diagram
	objects := mockDiagram.Objects()

	// Verify the number of objects in the diagram
	expectedNodeCount := 7
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	// Get the list of edges from the diagram
	edges := mockDiagram.Edges()

	// Verify the number of edges in the diagram
	expectedEdgeCount := 3 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestBuildDomainComponents_General_Error(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	componentKey := path.Join(mockGroup.Name(), string(valueobject.GeneralComponent))
	_ = mockDiagram.AddStringTo(componentKey, mockGroup.Name(), arch.RelationTypeAbstraction)
	err = model.buildDomainComponents(mockDiagram, mockGroup, mockGroup.Name())
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildDomainComponents_Function_Error(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	componentKey := path.Join(mockGroup.Name(), string(valueobject.GeneralComponent))
	funcObj := newMockObjectFunction(0)
	mockFunc := newMockFunction(funcObj)
	dif := valueobject.NewDomainFunction(valueobject.NewFunction(funcObj, mockFunc.Identifier()), "testdomain")
	_ = mockDiagram.AddObjTo(dif, componentKey, arch.RelationTypeAggregation)

	err = model.buildDomainComponents(mockDiagram, mockGroup, mockGroup.Name())
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildDomainComponents_Class_Error(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"
	claObj := newMockObject(0)
	claAttrObj := newMockObjectAttribute(0)
	claMethodObj := newMockObjectMethod(0)
	df := valueobject.NewDomainFunction(valueobject.NewFunction(claMethodObj, claObj.Identifier()), domain)
	mockClass := newMockDomainClass(domain, claObj, claAttrObj, claMethodObj)

	_ = mockDiagram.AddObjTo(df, mockClass.Identifier().ID(), arch.RelationTypeAggregation)

	err = model.buildDomainComponents(mockDiagram, mockGroup, mockGroup.Name())
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildDomainComponents_Interface_Error(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"
	intfObj := newMockObjectInterface(0)
	intfMethodObj := newMockObjectInterfaceMethod(0)
	dif := valueobject.NewDomainFunction(valueobject.NewFunction(intfMethodObj, intfObj.Identifier()), domain)
	mockItf := newMockDomainInterface(domain, intfObj, []*MockObject{intfMethodObj})

	_ = mockDiagram.AddObjTo(dif, mockItf.Identifier().ID(), arch.RelationTypeBehavior)

	err = model.buildDomainComponents(mockDiagram, mockGroup, mockGroup.Name())
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestAddDomainClass(t *testing.T) {
	mockObject := newMockObject(0)
	mockAttr := newMockObjectAttribute(0)
	mockMethod := newMockObjectMethod(0)

	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(mockObject)
	_ = mockRepo.Insert(mockAttr)
	_ = mockRepo.Insert(mockMethod)

	domain := "testdomain"
	mockClass := newMockDomainClass(domain, mockObject, mockAttr, mockMethod)

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if err := g.AddObjTo(mockClass, g.Name(), arch.RelationTypeAggregation); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.addDomainClass(g, mockClass)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the relationships in the diagram
	// Get the list of objects from the diagram
	objects := g.Objects()

	// Verify the number of objects in the diagram
	expectedNodeCount := 3 // Class node + attribute nodes + method nodes
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	// Get the list of edges from the diagram
	edges := g.Edges()

	// Verify the number of edges in the diagram
	expectedEdgeCount := 0 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestAddDomainClass_Error(t *testing.T) {
	mockObject := newMockObject(0)
	mockAttr := newMockObjectAttribute(0)
	mockMethod := newMockObjectMethod(0)

	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(mockObject)
	_ = mockRepo.Insert(mockAttr)
	_ = mockRepo.Insert(mockMethod)

	domain := "testdomain"
	mockClass := newMockDomainClass(domain, mockObject, mockAttr, mockMethod)

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if err := g.AddObjTo(mockClass, g.Name(), arch.RelationTypeAggregation); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	df := valueobject.NewDomainFunction(valueobject.NewFunction(mockMethod, mockObject.Identifier()), domain)
	_ = g.AddObjTo(df, mockClass.Identifier().ID(), arch.RelationTypeAggregation)

	err = model.addDomainClass(g, mockClass)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestDomainModel_AddAggregateToDiagram(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"

	mockAggregate := newMockAggregate(domain, 0)
	mockAggregateGroup := newMockAggregateGroup(domain)

	err = model.addAggregateToDiagram(g, mockAggregateGroup, mockAggregate)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the relationships in the diagram
	// Get the list of objects from the diagram
	objects := g.Objects()

	// Verify the number of objects in the diagram
	expectedNodeCount := 3 // Class node + attribute nodes + method nodes
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	// Get the list of edges from the diagram
	edges := g.Edges()

	// Verify the number of edges in the diagram
	expectedEdgeCount := 1 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestDomainModel_AddAggregateToDiagram_Error(t *testing.T) {
	mockGroup := newMockDomainGroup("testGroup", 0)
	mockDirectory := newMockEmptyDirectory()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"

	mockAggregate := newMockAggregate(domain, 0)
	mockAggregateGroup := newMockAggregateGroup(domain)

	claObj := newMockObject(0)
	claAttrObj := newMockObjectAttribute(0)
	claMethodObj := newMockObjectMethod(0)
	df := valueobject.NewDomainFunction(valueobject.NewFunction(claMethodObj, claObj.Identifier()), domain)
	mockClass := newMockDomainClass(domain, claObj, claAttrObj, claMethodObj)

	_ = g.AddObjTo(df, mockClass.Identifier().ID(), arch.RelationTypeAggregation)

	err = model.addAggregateToDiagram(g, mockAggregateGroup, mockAggregate)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	_ = g.AddObjTo(mockAggregate, g.Name(), arch.RelationTypeAggregationRoot)
	err = model.addAggregateToDiagram(g, mockAggregateGroup, mockAggregate)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestDomainModel_AddNodeToAggregate(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"
	mockAggregate := newMockAggregate(domain, 0)
	if err := g.AddObjTo(mockAggregate, g.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	key := "someKey"
	err = model.addNodeToAggregate(g, key, mockAggregate)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	objects := g.Objects()
	expectedNodeCount := 1
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	edges := g.Edges()
	expectedEdgeCount := 2
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestDomainModel_AddVOToNode(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"
	mockVO := newMockValueObject(domain, 0)

	err = model.addVOToNode(g, g.Name(), mockVO)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	objects := g.Objects()
	expectedNodeCount := 3 // Class node + attribute nodes + method nodes
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	edges := g.Edges()
	expectedEdgeCount := 0 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestDomainModel_AddVOToNode_Error(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"
	mockVO := newMockValueObject(domain, 0)

	_ = g.AddObjTo(mockVO, g.Name(), arch.RelationTypeAggregation)

	err = model.addVOToNode(g, g.Name(), mockVO)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestDomainModel_AddEntityToNode(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"
	mockEntity := newMockEntity(domain, 0)

	err = model.addEntityToNode(g, g.Name(), mockEntity)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	objects := g.Objects()
	expectedNodeCount := 3 // Class node + attribute nodes + method nodes
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	edges := g.Edges()
	expectedEdgeCount := 0 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestDomainModel_AddEntityToNode_Error(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "testdomain"
	mockEntity := newMockEntity(domain, 0)

	_ = g.AddObjTo(mockEntity, g.Name(), arch.RelationTypeAggregation)

	err = model.addEntityToNode(g, g.Name(), mockEntity)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestDomainModel_FindAggregateGroup(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	aggregateGroup1 := newMockAggregateGroupWithName("domain", "Group1")
	aggregateGroup2 := newMockAggregateGroupWithName("domain", "Group2")
	aggregateGroup3 := newMockAggregateGroupWithName("domain", "Group3")

	model.aggregates = append(model.aggregates, aggregateGroup1, aggregateGroup2, aggregateGroup3)

	foundGroup := model.FindAggregateGroup("Group2")
	if foundGroup == nil {
		t.Errorf("Expected to find AggregateGroup 'Group2', but got nil")
	} else if foundGroup.Name() != "Group2" {
		t.Errorf("Expected AggregateGroup with name 'Group2', but got '%s'", foundGroup.Name())
	}

	notFoundGroup := model.FindAggregateGroup("NonExistentGroup")
	if notFoundGroup != nil {
		t.Errorf("Expected not to find AggregateGroup 'NonExistentGroup', but got a match")
	}
}

func TestDomainModel_GetClass(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	claObj0 := newMockObject(0)
	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	cla0 := newMockClass(claObj0, claAttrObj0, claMethodObj0)

	claObj1 := newMockObject(1)
	claAttrObj1 := newMockObjectAttribute(1)
	claMethodObj1 := newMockObjectMethod(1)
	cla1 := newMockClass(claObj1, claAttrObj1, claMethodObj1)

	claObj2 := newMockObject(2)
	claAttrObj2 := newMockObjectAttribute(2)
	claMethodObj2 := newMockObjectMethod(2)
	cla2 := newMockClass(claObj2, claAttrObj2, claMethodObj2)

	_ = mockRepo.Insert(cla0)
	_ = mockRepo.Insert(cla1)
	_ = mockRepo.Insert(cla2)

	objIds := []arch.ObjIdentifier{claObj0.Identifier(), claObj1.Identifier()}
	classes, err := model.getClass(objIds)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(classes) != 2 {
		t.Errorf("Expected 2 classes in the result, but got %d", len(classes))
	}

	classNames := make(map[string]bool)
	for _, class := range classes {
		if cla, ok := class.(*valueobject.Class); ok {
			classNames[cla.Identifier().ID()] = true
		}
	}
	if !classNames[claObj0.Identifier().ID()] {
		t.Errorf("Expected class '%s' in the result, but it's missing", claObj0.Identifier().ID())
	}
	if !classNames[claObj1.Identifier().ID()] {
		t.Errorf("Expected class '%s' in the result, but it's missing", claObj1.Identifier().ID())
	}
}

func TestDomainModel_GetClass_Error(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	claObj0 := newMockObject(0)

	objIds := []arch.ObjIdentifier{claObj0.Identifier()}
	_, err = model.getClass(objIds)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestDomainModel_ProcessComponent(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "TestDomain"
	aggregateGroup := newMockAggregateGroup(domain)
	mockEntity1 := newMockEntity(domain, 1)
	mockEntity2 := newMockEntity(domain, 2)
	mockVO1 := newMockValueObject(domain, 1)
	mockVO2 := newMockValueObject(domain, 2)

	// 调用 processComponent 函数，将一组对象传入并指定组件类型为 EntityComponent
	err = model.processComponent(aggregateGroup, valueobject.EntityComponent, []arch.Object{mockEntity1, mockEntity2})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	var entityGroups []valueobject.Group
	for _, sg := range aggregateGroup.SubGroups() {
		switch sg.(type) {
		case *valueobject.EntityGroup:
			entityGroups = append(entityGroups, sg)
		}
	}

	if len(entityGroups) != 1 {
		t.Errorf("Expected 1 EntityGroup in the AggregateGroup, but got %d", len(entityGroups))
	}
	entityGroup := entityGroups[0]
	if entityGroup.(valueobject.DomainGroup).Domain() != domain {
		t.Errorf("Expected EntityGroup with name 'TestDomain', but got '%s'", entityGroup.Name())
	}

	err = model.processComponent(aggregateGroup, valueobject.VOComponent, []arch.Object{mockVO1, mockVO2})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	var voGroups []valueobject.Group

	for _, sg := range aggregateGroup.SubGroups() {
		switch sg.(type) {
		case *valueobject.VOGroup:
			voGroups = append(voGroups, sg)
		}
	}

	if len(voGroups) != 1 {
		t.Errorf("Expected 1 VOGroup in the AggregateGroup, but got %d", len(voGroups))
	}
	voGroup := voGroups[0]
	if voGroup.(valueobject.DomainGroup).Domain() != domain {
		t.Errorf("Expected VOGroup with name 'TestDomain', but got '%s'", voGroup.Name())
	}
}

func TestDomainModel_ProcessClasses(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "TestDomain"
	aggregateGroup := newMockAggregateGroup(domain)
	model.aggregates = append(model.aggregates, aggregateGroup)

	claObj0 := newMockObject(0)
	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	cla0 := newMockClass(claObj0, claAttrObj0, claMethodObj0)

	claObj1 := newMockObject(1)
	claAttrObj1 := newMockObjectAttribute(1)
	claMethodObj1 := newMockObjectMethod(1)
	cla1 := newMockClass(claObj1, claAttrObj1, claMethodObj1)

	_ = mockRepo.Insert(cla0)
	_ = mockRepo.Insert(cla1)

	dir := "root/aggregate0/testdir"
	err = model.processClasses([]arch.ObjIdentifier{cla0.Identifier(), cla1.Identifier()}, valueobject.EntityComponent, dir)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	var entityGroups []valueobject.Group
	for _, sg := range aggregateGroup.SubGroups() {
		switch sg.(type) {
		case *valueobject.EntityGroup:
			entityGroups = append(entityGroups, sg)
		}
	}

	if len(entityGroups) != 1 {
		t.Errorf("Expected 1 EntityGroup in the AggregateGroup, but got %d", len(entityGroups))
	}
	entityGroup := entityGroups[0]
	if entityGroup.(valueobject.DomainGroup).Domain() != domain {
		t.Errorf("Expected EntityGroup with name 'TestDomain', but got '%s'", entityGroup.Name())
	}
}

func TestDomainModel_ProcessClasses_Error(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "TestDomain"
	aggregateGroup := newMockAggregateGroup(domain)
	model.aggregates = append(model.aggregates, aggregateGroup)

	claObj0 := newMockObject(0)
	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	cla0 := newMockClass(claObj0, claAttrObj0, claMethodObj0)

	dir := "root/aggregate0/testdir"
	err = model.processClasses([]arch.ObjIdentifier{cla0.Identifier()}, valueobject.EntityComponent, dir)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestDomainModel_ProcessObjects(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "TestDomain"
	aggregateGroup := newMockAggregateGroup(domain)
	model.aggregates = append(model.aggregates, aggregateGroup)

	claObj0 := newMockObject(0)
	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	cla0 := newMockClass(claObj0, claAttrObj0, claMethodObj0)

	claObj1 := newMockObject(1)
	claAttrObj1 := newMockObjectAttribute(1)
	claMethodObj1 := newMockObjectMethod(1)
	cla1 := newMockClass(claObj1, claAttrObj1, claMethodObj1)

	_ = mockRepo.Insert(cla0)
	_ = mockRepo.Insert(cla1)

	dir := "root/aggregate0/testdir"
	err = model.processObjects([]arch.ObjIdentifier{cla0.Identifier(), cla1.Identifier()}, valueobject.EntityComponent, dir)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	var entityGroups []valueobject.Group
	for _, sg := range aggregateGroup.SubGroups() {
		switch sg.(type) {
		case *valueobject.EntityGroup:
			entityGroups = append(entityGroups, sg)
		}
	}

	if len(entityGroups) != 1 {
		t.Errorf("Expected 1 EntityGroup in the AggregateGroup, but got %d", len(entityGroups))
	}
	entityGroup := entityGroups[0]
	if entityGroup.(valueobject.DomainGroup).Domain() != domain {
		t.Errorf("Expected EntityGroup with name 'TestDomain', but got '%s'", entityGroup.Name())
	}
}

func TestDomainModel_ProcessObjects_Error(t *testing.T) {
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	mockDirectory := newMockEmptyDirectory()
	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	domain := "TestDomain"
	aggregateGroup := newMockAggregateGroup(domain)
	model.aggregates = append(model.aggregates, aggregateGroup)

	claObj0 := newMockObject(0)
	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	cla0 := newMockClass(claObj0, claAttrObj0, claMethodObj0)

	dir := "root/aggregate0/testdir"
	err = model.processObjects([]arch.ObjIdentifier{cla0.Identifier()}, valueobject.EntityComponent, dir)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestGeneralModel_StrategicGrouping(t *testing.T) {
	mockDirectory, objs := newMockDirectoryWithDomainObjs()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.StrategicGrouping()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the number of objects in the diagram
	expectedSubGroupCount := 1
	if len(model.aggregates) != expectedSubGroupCount {
		t.Errorf("Expected %d subgroups in the group, but got %d", expectedSubGroupCount, len(model.aggregates))
	}

	if len(model.aggregates) > 0 {
		agg := model.aggregates[0]
		for _, sg := range agg.SubGroups() {
			switch sg.(type) {
			case valueobject.DomainGroup:
				domainGroup := sg.(valueobject.DomainGroup)
				if len(sg.Objects()) != 0 {
					t.Errorf("Expected 1 object in the DomainGroup, but got %d", len(sg.Objects()))
				}
				if domainGroup.Domain() != "root/internal/domain/test" {
					t.Errorf("Expected DomainGroup with name 'TestDomain', but got '%s'", domainGroup.Domain())
				}
			default:
				t.Errorf("Expected DomainGroup, but got %T", sg)
			}
		}
	}
}

func TestGeneralModel_TacticGrouping(t *testing.T) {
	mockDirectory, objs := newMockDirectoryWithDomainObjs()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
	}

	model, err := NewDomainModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.TacticGrouping()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the number of objects in the diagram
	expectedSubGroupCount := 1
	if len(model.aggregates) != expectedSubGroupCount {
		t.Errorf("Expected %d subgroups in the group, but got %d", expectedSubGroupCount, len(model.aggregates))
	}

	if len(model.aggregates) > 0 {
		agg := model.aggregates[0]
		for _, sg := range agg.SubGroups() {
			switch sg.(type) {
			case valueobject.DomainGroup:
				domainGroup := sg.(valueobject.DomainGroup)
				if len(sg.Objects()) != 4 {
					t.Errorf("Expected 1 object in the DomainGroup, but got %d", len(sg.Objects()))
				}
				if domainGroup.Domain() != "root/internal/domain/test" {
					t.Errorf("Expected DomainGroup with name 'TestDomain', but got '%s'", domainGroup.Domain())
				}
			default:
				t.Errorf("Expected DomainGroup, but got %T", sg)
			}
		}
	}
}
