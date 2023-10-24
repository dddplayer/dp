package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"path"
	"testing"
)

func TestNewGeneralModel(t *testing.T) {
	mockRepo := &MockObjectRepository{}
	mockDirectory := newMockEmptyDirectory()

	model, err := NewGeneralModel(nil, mockDirectory)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
	model, err = NewGeneralModel(mockRepo, nil)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	model, err = NewGeneralModel(mockRepo, mockDirectory)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if model == nil {
		t.Errorf("Expected non-nil model, but got nil")
	}
	if model.repo != mockRepo {
		t.Errorf("Expected repo to be mockRepo, but got %v", model.repo)
	}
	if model.rootGroup != nil {
		t.Errorf("Expected rootGroup to be nil, but got %v", model.rootGroup)
	}
}

func TestFindGroup(t *testing.T) {
	mockGroup := &MockGroup{
		NameFunc: func() string {
			return "root"
		},
		SubGroupsFunc: func() []valueobject.Group {
			return []valueobject.Group{
				&MockGroup{
					NameFunc: func() string {
						return "subgroup1"
					},
					SubGroupsFunc: func() []valueobject.Group {
						return []valueobject.Group{}
					},
				},
				&MockGroup{
					NameFunc: func() string {
						return "subgroup2"
					},
					SubGroupsFunc: func() []valueobject.Group {
						return []valueobject.Group{}
					},
				},
			}
		},
	}

	model := &GeneralModel{}

	parentGroup := model.FindGroup("subgroup1", mockGroup)
	if parentGroup == nil || parentGroup.Name() != "subgroup1" {
		t.Errorf("Expected parent group to be root, but got %v", parentGroup)
	}

	parentGroup = model.FindGroup("nonexistent", mockGroup)
	if parentGroup != nil {
		t.Errorf("Expected parent group to be nil, but got %v", parentGroup)
	}
}

func TestAddClass(t *testing.T) {
	mockObject := newMockObject(0)
	mockAttr := newMockObjectAttribute(0)
	mockMethod := newMockObjectMethod(0)

	attributeIdentifiers := []arch.ObjIdentifier{mockAttr.Identifier()}
	methodIdentifiers := []arch.ObjIdentifier{mockMethod.Identifier()}

	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(mockObject)
	_ = mockRepo.Insert(mockAttr)
	_ = mockRepo.Insert(mockMethod)

	mockClass := newMockClass(mockObject, mockAttr, mockMethod)

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if err := g.AddObjTo(mockClass, g.Name(), arch.RelationTypeAggregation); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	err = model.addClass(g, mockClass)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the relationships in the diagram
	// Get the list of objects from the diagram
	objects := g.Objects()

	// Verify the number of objects in the diagram
	expectedNodeCount := 1 + len(attributeIdentifiers) + len(methodIdentifiers) // Class node + attribute nodes + method nodes
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

func TestAddClass_Attr_ErrorCase(t *testing.T) {
	// Create a mock class with attributes and methods
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

	mockClass := newMockClass(mockObject, mockAttr, mockMethod)

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if err := g.AddObjTo(mockClass, g.Name(), arch.RelationTypeAggregation); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	// Test case: Object not found in repository
	_ = g.AddObjTo(mockAttr, mockClass.Identifier().ID(), arch.RelationTypeBehavior)
	err = model.addClass(g, mockClass)
	if err == nil {
		t.Errorf("Expected an error when object not found, but got nil")
	}
	expectedErrMsg := fmt.Sprintf("key conflict: %s", mockAttr.Identifier().ID())
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%v'", expectedErrMsg, err)
	}
}

func TestAddClass_Method_ErrorCase(t *testing.T) {
	// Create a mock class with attributes and methods
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

	mockClass := newMockClass(mockObject, mockAttr, mockMethod)

	g, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if err := g.AddObjTo(mockClass, g.Name(), arch.RelationTypeAggregation); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	// Test case: Object not found in repository
	_ = g.AddObjTo(mockMethod, mockClass.Identifier().ID(), arch.RelationTypeAttribution)
	err = model.addClass(g, mockClass)
	if err == nil {
		t.Errorf("Expected an error when object not found, but got nil")
	}
	expectedErrMsg := fmt.Sprintf("key conflict: %s", mockMethod.Identifier().ID())
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%v'", expectedErrMsg, err)
	}
}

func TestBuildAttributeComponents(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.buildAttributeComponents(mockDiagram, mockGroup, valueobject.ClassComponent)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.buildAttributeComponents(mockDiagram, mockGroup, valueobject.InterfaceComponent)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test case: Unsupported component type
	err = model.buildAttributeComponents(mockDiagram, mockGroup, "InvalidComponentType")
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

func TestBuildAttributeComponents_MissingObjectClassComponent(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	// Not inserting any objects into the repository

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test case: Missing class object
	err = model.buildAttributeComponents(mockDiagram, mockGroup, valueobject.ClassComponent)
	if err == nil {
		t.Errorf("Expected an error for missing class object, but got nil")
	}
	// Ensure that err is nil because class component is ignored

	// Test case: Missing interface object
	err = model.buildAttributeComponents(mockDiagram, mockGroup, valueobject.InterfaceComponent)
	if err == nil {
		t.Errorf("Expected an error for missing interface object, but got nil")
	}
}

func TestBuildAttributeComponents_DuplicateClassInterface(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)

	claObj := newMockObject(0)
	claAttrObj := newMockObjectAttribute(0)
	claMethodObj := newMockObjectMethod(0)
	mockClass := newMockClass(claObj, claAttrObj, claMethodObj)

	intfObj := newMockObjectInterface(0)
	intfMethodObj := newMockObjectInterfaceMethod(0)
	mockItf := newMockInterface(intfObj, []*MockObject{intfMethodObj})

	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	_ = mockDiagram.AddObjTo(claMethodObj, mockClass.Identifier().ID(), arch.RelationTypeAttribution)
	err = model.buildAttributeComponents(mockDiagram, mockGroup, valueobject.ClassComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	_ = mockDiagram.AddObjTo(intfMethodObj, mockItf.Identifier().ID(), arch.RelationTypeBehavior)
	err = model.buildAttributeComponents(mockDiagram, mockGroup, valueobject.InterfaceComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildAbstractComponent(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test case: Unsupported component type
	err = model.buildAbstractComponent(mockDiagram, mockGroup, "InvalidComponentType")
	if err == nil {
		t.Errorf("Expected an error for unsupported component type, but got nil")
	}
	expectedErrMsg := "unsupported objs type: InvalidComponentType"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', but got '%v'", expectedErrMsg, err)
	}

	err = model.buildAbstractComponent(mockDiagram, mockGroup, valueobject.GeneralComponent)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.buildAbstractComponent(mockDiagram, mockGroup, valueobject.FunctionComponent)
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

func TestBuildAbstractComponent_Error(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
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
	err = model.buildAbstractComponent(mockDiagram, mockGroup, valueobject.GeneralComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	funcObj := newMockObjectFunction(0)
	mockFunc := newMockFunction(funcObj)
	_ = mockDiagram.AddObjTo(mockFunc, componentKey, arch.RelationTypeAggregation)
	err = model.buildAbstractComponent(mockDiagram, mockGroup, valueobject.FunctionComponent)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildComponents(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.buildComponents(mockDiagram, mockGroup)
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

func TestBuildComponents_InterfaceError(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	intfObj := newMockObjectInterface(0)
	intfMethodObj := newMockObjectInterfaceMethod(0)
	mockItf := newMockInterface(intfObj, []*MockObject{intfMethodObj})
	_ = mockDiagram.AddObjTo(mockItf, mockGroup.Name(), arch.RelationTypeAggregation)
	err = model.buildComponents(mockDiagram, mockGroup)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildComponents_ClassError(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	claObj := newMockObject(0)
	claAttrObj := newMockObjectAttribute(0)
	claMethodObj := newMockObjectMethod(0)
	mockClass := newMockClass(claObj, claAttrObj, claMethodObj)
	_ = mockDiagram.AddObjTo(mockClass, mockGroup.Name(), arch.RelationTypeAggregation)
	err = model.buildComponents(mockDiagram, mockGroup)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildComponents_FunctionError(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
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
	_ = mockDiagram.AddObjTo(mockFunc, componentKey, arch.RelationTypeAggregation)
	err = model.buildComponents(mockDiagram, mockGroup)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestBuildComponents_GeneralError(t *testing.T) {
	mockGroup := newMockGroup("testGroup", 0)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.MockObjects {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
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
	err = model.buildComponents(mockDiagram, mockGroup)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestAddGroupToDiagram(t *testing.T) {
	mockGroup := newTwoLevelMockGroup()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.AllMockObjects() {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.addGroupToDiagram(mockDiagram, mockGroup, mockDiagram.Name())

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the relationships in the diagram
	// Get the list of objects from the diagram
	objects := mockDiagram.Objects()

	// Verify the number of objects in the diagram
	expectedNodeCount := 21
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	// Get the list of edges from the diagram
	edges := mockDiagram.Edges()

	// Verify the number of edges in the diagram
	expectedEdgeCount := 9 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestAddGroupToDiagram_Error(t *testing.T) {
	mockGroup := newTwoLevelMockGroup()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.AllMockObjects() {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo: mockRepo,
	}

	mockDiagram, err := NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	subGroup2 := newMockGroup("subgroup2", 2)
	componentKey := path.Join(subGroup2.Name(), string(valueobject.GeneralComponent))
	_ = mockDiagram.AddStringTo(componentKey, subGroup2.Name(), arch.RelationTypeAbstraction)
	err = model.addGroupToDiagram(mockDiagram, mockGroup, mockDiagram.Name())

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	componentKey2 := path.Join(mockGroup.Name(), string(valueobject.GeneralComponent))
	_ = mockDiagram.AddStringTo(componentKey2, mockGroup.Name(), arch.RelationTypeAbstraction)
	err = model.addGroupToDiagram(mockDiagram, mockGroup, mockDiagram.Name())

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	_ = mockDiagram.AddStringTo(mockGroup.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot)
	err = model.addGroupToDiagram(mockDiagram, mockGroup, mockDiagram.Name())

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestAddRootGroupToDiagram(t *testing.T) {
	mockGroup := newTwoLevelMockGroup()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.AllMockObjects() {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo:      mockRepo,
		rootGroup: mockGroup,
	}

	mockDiagram, err := NewDiagram("root", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = model.addRootGroupToDiagram(mockDiagram)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the relationships in the diagram
	// Get the list of objects from the diagram
	objects := mockDiagram.Objects()

	// Verify the number of objects in the diagram
	expectedNodeCount := 21
	if len(objects) != expectedNodeCount {
		t.Errorf("Expected %d nodes in the diagram, but got %d", expectedNodeCount, len(objects))
	}

	// Get the list of edges from the diagram
	edges := mockDiagram.Edges()

	// Verify the number of edges in the diagram
	expectedEdgeCount := 8 // Attribute and method relationships will be ignored
	if len(edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges in the diagram, but got %d", expectedEdgeCount, len(edges))
	}
}

func TestAddRootGroupToDiagram_Error(t *testing.T) {
	mockGroup := newTwoLevelMockGroup()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range mockGroup.AllMockObjects() {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo:      mockRepo,
		rootGroup: mockGroup,
	}

	mockDiagram, err := NewDiagram("root", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	subGroup2 := newMockGroup("subgroup2", 2)
	_ = mockDiagram.AddStringTo(subGroup2.Name(), mockDiagram.Name(), arch.RelationTypeAggregationRoot)
	err = model.addRootGroupToDiagram(mockDiagram)

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	componentKey := path.Join(mockGroup.Name(), string(valueobject.GeneralComponent))
	_ = mockDiagram.AddStringTo(componentKey, mockGroup.Name(), arch.RelationTypeAbstraction)
	err = model.addRootGroupToDiagram(mockDiagram)

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

}

func TestGeneralModel_Grouping(t *testing.T) {
	mockDirectory, objs := newMockDirectoryWithObjs()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
	}

	model := &GeneralModel{
		repo:      mockRepo,
		directory: mockDirectory,
	}

	model.Grouping()

	sgs := model.rootGroup.SubGroups()

	// Verify the number of objects in the diagram
	expectedSubGroupCount := 3
	if len(sgs) != expectedSubGroupCount {
		t.Errorf("Expected %d subgroups in the group, but got %d", expectedSubGroupCount, len(sgs))
	}
}
