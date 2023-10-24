package entity

import (
	"errors"
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"github.com/dddplayer/dp/pkg/datastructure/directory"
	"path"
)

// MockObjIdentifier is a mock implementation of ObjIdentifier for testing
type MockObjIdentifier struct {
	id                     string
	name                   string
	dir                    string
	NameSeparatorLengthVal int
}

// ID returns the identifier of the object
func (m *MockObjIdentifier) ID() string {
	return m.id
}

// Name returns the name of the object
func (m *MockObjIdentifier) Name() string {
	return m.name
}

// Dir returns the directory of the object
func (m *MockObjIdentifier) Dir() string {
	return m.dir
}

func (m *MockObjIdentifier) NameSeparatorLength() int {
	return m.NameSeparatorLengthVal
}

type MockDependenceRelation struct {
	from      arch.Object
	dependsOn arch.Object
}

func (rel *MockDependenceRelation) Type() arch.RelationType {
	return arch.RelationTypeDependency
}

func (rel *MockDependenceRelation) From() arch.Object {
	return rel.from
}

func (rel *MockDependenceRelation) DependsOn() arch.Object {
	return rel.dependsOn
}

type MockCompositionRelation struct {
	from  arch.Object
	child arch.Object
}

func (rel *MockCompositionRelation) Type() arch.RelationType {
	return arch.RelationTypeComposition
}

func (rel *MockCompositionRelation) From() arch.Object {
	return rel.from
}

func (rel *MockCompositionRelation) Child() arch.Object {
	return rel.child
}

type MockEmbeddingRelation struct {
	from     arch.Object
	embedded arch.Object
}

func (rel *MockEmbeddingRelation) Type() arch.RelationType {
	return arch.RelationTypeEmbedding
}

func (rel *MockEmbeddingRelation) From() arch.Object {
	return rel.from
}

func (rel *MockEmbeddingRelation) Embedded() arch.Object {
	return rel.embedded
}

type MockImplementationRelation struct {
	from       arch.Object
	implements []arch.Object
}

func (rel *MockImplementationRelation) Type() arch.RelationType {
	return arch.RelationTypeImplementation
}

func (rel *MockImplementationRelation) From() arch.Object {
	return rel.from
}

func (rel *MockImplementationRelation) Implements() []arch.Object {
	return rel.implements
}

func (rel *MockImplementationRelation) Implemented(obj arch.Object) {
	rel.implements = append(rel.implements, obj)
}

type MockAssociationRelation struct {
	from            arch.Object
	refer           arch.Object
	associationType arch.RelationType
}

func (rel *MockAssociationRelation) Type() arch.RelationType {
	return arch.RelationTypeAssociation
}

func (rel *MockAssociationRelation) From() arch.Object {
	return rel.from
}

func (rel *MockAssociationRelation) Refer() arch.Object {
	return rel.refer
}

func (rel *MockAssociationRelation) AssociationType() arch.RelationType {
	return rel.associationType
}

type MockObject struct {
	id                     arch.ObjIdentifier
	position               arch.Position
	name                   string
	dir                    string
	NameSeparatorLengthVal int
}

// ID returns the identifier of the object
func (mo MockObject) ID() string {
	return mo.id.ID()
}

func (mo MockObject) Identifier() arch.ObjIdentifier {
	return mo.id
}

func (mo MockObject) Position() arch.Position {
	return mo.position
}

// Name returns the name of the object
func (mo MockObject) Name() string {
	return mo.name
}

// Dir returns the directory of the object
func (mo MockObject) Dir() string {
	return mo.dir
}

func (mo MockObject) NameSeparatorLength() int {
	return mo.NameSeparatorLengthVal
}

type MockPosition struct {
	FilenameVal string
	OffsetVal   int
	LineVal     int
	ColumnVal   int
}

func (mp *MockPosition) Filename() string {
	return mp.FilenameVal
}

func (mp *MockPosition) Offset() int {
	return mp.OffsetVal
}

func (mp *MockPosition) Line() int {
	return mp.LineVal
}

func (mp *MockPosition) Column() int {
	return mp.ColumnVal
}

func (mp *MockPosition) IsEqual(pos arch.Position) bool {
	return mp.Filename() == pos.Filename() &&
		mp.Offset() == pos.Offset() &&
		mp.Line() == pos.Line() &&
		mp.Column() == pos.Column()
}

type MockRelationRepository struct {
	relations []arch.Relation
}

func (mrr *MockRelationRepository) Insert(rel arch.Relation) error {
	mrr.relations = append(mrr.relations, rel)
	return nil
}

func (mrr *MockRelationRepository) Walk(walker func(rel arch.Relation) error) {
	for _, rel := range mrr.relations {
		if err := walker(rel); err != nil {
			break
		}
	}
}

func newMockObjectRepoWithInvalidDir() *MockObjectRepository {
	claObj1 := newMockObjectWithId("test", "cla1", 1)
	claObj2 := newMockObjectWithId("tex", "cla2", 1)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(claObj1)
	_ = mockRepo.Insert(claObj2)
	return mockRepo
}

func newMockObjRepoWithDuplicatedIdent() *MockObjectRepository {
	claObj1 := newMockObject(1)
	claObj2 := newMockObject(2)
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	_ = mockRepo.Insert(claObj1)
	_ = mockRepo.Insert(claObj2)
	mockRepo.idents = append(mockRepo.idents, claObj1.Identifier())
	return mockRepo
}

type MockObjectRepository struct {
	objects map[string]arch.Object
	idents  []arch.ObjIdentifier
}

func (mor *MockObjectRepository) Find(id arch.ObjIdentifier) arch.Object {
	return mor.objects[id.ID()]
}

func (mor *MockObjectRepository) GetObjects(ids []arch.ObjIdentifier) ([]arch.Object, error) {
	var result []arch.Object
	for _, id := range ids {
		obj := mor.objects[id.ID()]
		if obj != nil {
			result = append(result, obj)
		}
	}
	if len(result) == len(ids) {
		return result, nil
	}
	return nil, errors.New("some objects not found")
}

func (mor *MockObjectRepository) All() []arch.ObjIdentifier {
	return mor.idents
}

func (mor *MockObjectRepository) Insert(obj arch.Object) error {
	mor.objects[obj.Identifier().ID()] = obj
	mor.idents = append(mor.idents, obj.Identifier())
	return nil
}

func (mor *MockObjectRepository) Walk(walker func(obj arch.Object) error) {
	for _, obj := range mor.objects {
		if err := walker(obj); err != nil {
			break
		}
	}
}

type MockGroup struct {
	NameFunc          func() string
	SubGroupsFunc     func() []valueobject.Group
	AppendGroupsFunc  func(...valueobject.Group)
	ObjectsFunc       func() []arch.Object
	AppendObjectsFunc func(...arch.Object)
	ClassesFunc       func() []*valueobject.Class
	GeneralsFunc      func() []*valueobject.General
	FunctionsFunc     func() []*valueobject.Function
	InterfacesFunc    func() []*valueobject.Interface
	MockObjects       []*MockObject
}

func (m *MockGroup) AllMockObjects() []*MockObject {
	var all []*MockObject
	for _, mo := range m.MockObjects {
		all = append(all, mo)
	}
	for _, sg := range m.SubGroups() {
		if g, ok := sg.(*MockGroup); ok {
			for _, mo := range g.MockObjects {
				all = append(all, mo)
			}
		}
	}
	return all
}

func (m *MockGroup) Name() string {
	return m.NameFunc()
}

func (m *MockGroup) SubGroups() []valueobject.Group {
	return m.SubGroupsFunc()
}

func (m *MockGroup) AppendGroups(groups ...valueobject.Group) {
	m.AppendGroupsFunc(groups...)
}

func (m *MockGroup) Objects() []arch.Object {
	return m.ObjectsFunc()
}

func (m *MockGroup) AppendObjects(objects ...arch.Object) {
	m.AppendObjectsFunc(objects...)
}

func (m *MockGroup) Classes() []*valueobject.Class {
	return m.ClassesFunc()
}

func (m *MockGroup) Generals() []*valueobject.General {
	return m.GeneralsFunc()
}

func (m *MockGroup) Functions() []*valueobject.Function {
	return m.FunctionsFunc()
}

func (m *MockGroup) Interfaces() []*valueobject.Interface {
	return m.InterfacesFunc()
}

func newMockInvalidEmptyDirectory() *Directory {
	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "root",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): {},
				string(arch.HexagonDirectoryPkg): {},
			},
		},
	}
	return mockDirectory
}

func newMockEmptyDirectory() *Directory {
	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "root",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): {},
				string(arch.HexagonDirectoryPkg): {},
				string(arch.HexagonDirectoryInternal): {
					Children: map[string]*directory.TreeNode{
						string(arch.HexagonDirectoryDomain): {
							Children: map[string]*directory.TreeNode{
								string(arch.HexagonDirectoryEntity):      nil,
								string(arch.HexagonDirectoryValueObject): nil,
							},
						},
					},
				},
			},
		},
	}
	return mockDirectory
}

func newMockDirectoryWithObjs() (*Directory, []arch.Object) {
	claObj0 := newMockObject(0)
	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	cla0 := newMockClass(claObj0, claAttrObj0, claMethodObj0)

	claObj1 := newMockObject(1)
	claAttrObj1 := newMockObjectAttribute(1)
	claMethodObj1 := newMockObjectMethod(1)
	cla1 := newMockClass(claObj1, claAttrObj1, claMethodObj1)

	funcObj0 := newMockObjectFunction(0)
	funcObj1 := newMockObjectFunction(1)

	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "testpackage",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): {
					Name:  string(arch.HexagonDirectoryCmd),
					Value: []arch.ObjIdentifier{funcObj0.Identifier()},
				},
				string(arch.HexagonDirectoryPkg): {
					Name:  string(arch.HexagonDirectoryPkg),
					Value: []arch.ObjIdentifier{funcObj1.Identifier()},
				},
				string(arch.HexagonDirectoryInternal): {
					Name: string(arch.HexagonDirectoryInternal),
					Children: map[string]*directory.TreeNode{
						string(arch.HexagonDirectoryDomain): {
							Name: string(arch.HexagonDirectoryDomain),
							Children: map[string]*directory.TreeNode{
								string(arch.HexagonDirectoryEntity): {
									Name:  string(arch.HexagonDirectoryEntity),
									Value: []arch.ObjIdentifier{cla0.Identifier(), claAttrObj0, claMethodObj0},
								},
								string(arch.HexagonDirectoryValueObject): {
									Name:  string(arch.HexagonDirectoryValueObject),
									Value: []arch.ObjIdentifier{cla1.Identifier(), claAttrObj1, claMethodObj1},
								},
							},
						},
					},
				},
			},
		},
	}
	return mockDirectory, []arch.Object{
		claObj0, claAttrObj0, claMethodObj0,
		claObj1, claAttrObj1, claMethodObj1,
		funcObj0, funcObj1}
}

func newMockDirectoryWithDomainObjs() (*Directory, []arch.Object) {
	domain := "TestDomain"
	claObj0 := newMockObject(0)
	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	cla0 := newMockDomainClass(domain, claObj0, claAttrObj0, claMethodObj0)

	claObj1 := newMockObject(1)
	claAttrObj1 := newMockObjectAttribute(1)
	claMethodObj1 := newMockObjectMethod(1)
	cla1 := newMockDomainClass(domain, claObj1, claAttrObj1, claMethodObj1)

	funcObj0 := newMockDomainFunction(domain, newMockObjectFunction(0))
	funcObj1 := newMockDomainFunction(domain, newMockObjectFunction(1))
	funcObj2 := newMockDomainFunction(domain, newMockObjectFunction(2))
	funcObj3 := newMockDomainFunction(domain, newMockObjectFunction(3))

	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "root",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): {
					Name:  string(arch.HexagonDirectoryCmd),
					Value: []arch.ObjIdentifier{funcObj0.Identifier()},
				},
				string(arch.HexagonDirectoryPkg): {
					Name:  string(arch.HexagonDirectoryPkg),
					Value: []arch.ObjIdentifier{funcObj1.Identifier()},
				},
				string(arch.HexagonDirectoryInternal): {
					Name: string(arch.HexagonDirectoryInternal),
					Children: map[string]*directory.TreeNode{
						string(arch.HexagonDirectoryDomain): {
							Name: string(arch.HexagonDirectoryDomain),
							Children: map[string]*directory.TreeNode{
								"test": {
									Name: "test",
									Children: map[string]*directory.TreeNode{
										string(arch.HexagonDirectoryEntity): {
											Name: string(arch.HexagonDirectoryEntity),
											Value: []arch.ObjIdentifier{
												cla0.Identifier(), claAttrObj0.Identifier(), claMethodObj0.Identifier(),
												funcObj2.Identifier()},
										},
										string(arch.HexagonDirectoryValueObject): {
											Name: string(arch.HexagonDirectoryValueObject),
											Value: []arch.ObjIdentifier{
												cla1.Identifier(), claAttrObj1.Identifier(), claMethodObj1.Identifier(),
												funcObj3.Identifier()},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return mockDirectory, []arch.Object{
		claObj0, claAttrObj0, claMethodObj0, cla0,
		claObj1, claAttrObj1, claMethodObj1, cla1,
		funcObj0, funcObj1, funcObj2, funcObj3}
}

func newMockDirectoryWithAggregate() (*Directory, []arch.Object) {
	domain := "TestDomain"
	claObj0 := newTestMockObject()
	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	cla0 := newMockClass(claObj0, claAttrObj0, claMethodObj0)

	claObj1 := newMockObject(1)
	claAttrObj1 := newMockObjectAttribute(1)
	claMethodObj1 := newMockObjectMethod(1)
	cla1 := newMockClass(claObj1, claAttrObj1, claMethodObj1)

	funcObj0 := newMockDomainFunction(domain, newMockObjectFunction(0))
	funcObj1 := newMockDomainFunction(domain, newMockObjectFunction(1))
	funcObj2 := newMockDomainFunction(domain, newMockObjectFunction(2))
	funcObj3 := newMockDomainFunction(domain, newMockObjectFunction(3))

	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "root",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): {
					Name:  string(arch.HexagonDirectoryCmd),
					Value: []arch.ObjIdentifier{funcObj0.Identifier()},
				},
				string(arch.HexagonDirectoryPkg): {
					Name:  string(arch.HexagonDirectoryPkg),
					Value: []arch.ObjIdentifier{funcObj1.Identifier()},
				},
				string(arch.HexagonDirectoryInternal): {
					Name: string(arch.HexagonDirectoryInternal),
					Children: map[string]*directory.TreeNode{
						string(arch.HexagonDirectoryDomain): {
							Name: string(arch.HexagonDirectoryDomain),
							Children: map[string]*directory.TreeNode{
								"test": {
									Name: "test",
									Children: map[string]*directory.TreeNode{
										string(arch.HexagonDirectoryEntity): {
											Name: string(arch.HexagonDirectoryEntity),
											Value: []arch.ObjIdentifier{
												cla0.Identifier(), claAttrObj0.Identifier(), claMethodObj0.Identifier(),
												funcObj2.Identifier()},
										},
										string(arch.HexagonDirectoryValueObject): {
											Name: string(arch.HexagonDirectoryValueObject),
											Value: []arch.ObjIdentifier{
												cla1.Identifier(), claAttrObj1.Identifier(), claMethodObj1.Identifier(),
												funcObj3.Identifier()},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return mockDirectory, []arch.Object{
		claObj0, claAttrObj0, claMethodObj0, cla0,
		claObj1, claAttrObj1, claMethodObj1, cla1,
		funcObj0, funcObj1, funcObj2, funcObj3}
}

func newTestMockObject() *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   "testpackage/Test",
			name: "Test",
			dir:  "testpackage"},
		name:     "Test",
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 5, ColumnVal: 2},
	}
}

func newMockObject(id int) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:                     fmt.Sprintf("testpackage/MockObject_%d", id),
			name:                   fmt.Sprintf("MockObject_%d", id),
			dir:                    "testpackage",
			NameSeparatorLengthVal: 1,
		},
		name:     fmt.Sprintf("MockObject_%d", id),
		position: &MockPosition{FilenameVal: fmt.Sprintf("mockfile_%d", id), OffsetVal: 10, LineVal: 5, ColumnVal: 2},
	}
}

func newMockObjectWithStr(id string) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   fmt.Sprintf("testpackage/MockObject_%s", id),
			name: fmt.Sprintf("MockObject_%s", id),
			dir:  "testpackage"},
		name:     fmt.Sprintf("MockObject_%s", id),
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 5, ColumnVal: 2},
	}
}

func newMockClassWithName(dir, name string) *valueobject.Class {
	obj := &MockObject{
		id: &MockObjIdentifier{
			id:   path.Join(dir, name),
			name: name,
			dir:  dir},
		name:     name,
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 5, ColumnVal: 2},
	}

	claAttrObj0 := newMockObjectAttribute(0)
	claMethodObj0 := newMockObjectMethod(0)
	return newMockClass(obj, claAttrObj0, claMethodObj0)
}

func newMockObjectWithId(dir, name string, index int) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   path.Join(dir, name),
			name: fmt.Sprintf("%s_%d", name, index),
			dir:  dir},
		name:     fmt.Sprintf("%s_%d", name, index),
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 5, ColumnVal: 2},
	}
}

func newMockObjectGeneral(id int) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   fmt.Sprintf("testpackage/general_%d", id),
			name: fmt.Sprintf("general_%d", id),
			dir:  "testpackage"},
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 5, ColumnVal: 2},
	}
}

func newMockObjectAttribute(id int) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   fmt.Sprintf("testpackage/attribute_%d", id),
			name: fmt.Sprintf("attribute_%d", id),
			dir:  "testpackage"},
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 6, ColumnVal: 2},
	}
}

func newMockObjectMethod(id int) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   fmt.Sprintf("testpackage/method_%d", id),
			name: fmt.Sprintf("method_%d", id),
			dir:  "testpackage"},
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 7, ColumnVal: 2},
	}
}

func newMockClass(obj, attr, method *MockObject) *valueobject.Class {
	attributeIdentifiers := []arch.ObjIdentifier{attr.Identifier()}
	methodIdentifiers := []arch.ObjIdentifier{method.Identifier()}

	mockClass := valueobject.NewClass(obj, attributeIdentifiers, methodIdentifiers)

	return mockClass
}

func newMockDomainClass(d string, obj, attr, method *MockObject) *valueobject.DomainClass {
	attributeObjs := []*valueobject.DomainAttr{
		valueobject.NewDomainAttr(valueobject.NewAttr(attr), d)}
	methodObjs := []*valueobject.DomainFunction{
		valueobject.NewDomainFunction(valueobject.NewFunction(method, obj.Identifier()), d)}

	cla := newMockClass(obj, attr, method)

	mockClass := valueobject.NewDomainClass(cla, d, attributeObjs, methodObjs)

	return mockClass
}

func newMockObjectInterface(id int) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   fmt.Sprintf("testpackage/interface_%d", id),
			name: fmt.Sprintf("interface_%d", id),
			dir:  "testpackage"},
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 5, ColumnVal: 2},
	}
}

func newMockObjectInterfaceMethod(id int) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   fmt.Sprintf("testpackage/interface_method_%d", id),
			name: fmt.Sprintf("interface_method_%d", id),
			dir:  "testpackage"},
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 7, ColumnVal: 2},
	}
}

func newMockObjectFunction(id int) *MockObject {
	return &MockObject{
		id: &MockObjIdentifier{
			id:   fmt.Sprintf("testpackage/function_%d", id),
			name: fmt.Sprintf("function_%d", id),
			dir:  "testpackage"},
		position: &MockPosition{FilenameVal: "mockfile", OffsetVal: 10, LineVal: 7, ColumnVal: 2},
	}
}

func newMockInterface(obj *MockObject, methods []*MockObject) *valueobject.Interface {
	var ms []arch.Object
	for _, m := range methods {
		ms = append(ms, m)
	}
	return valueobject.NewInterface(obj, ms)
}

func newMockDomainInterface(d string, obj *MockObject, methods []*MockObject) *valueobject.DomainInterface {
	var ms []*valueobject.DomainFunction
	for _, m := range methods {
		ms = append(ms, valueobject.NewDomainFunction(
			valueobject.NewFunction(m, obj.Identifier()), d))
	}
	return valueobject.NewDomainInterface(newMockInterface(obj, methods), d, ms)
}

func newMockGeneral(obj *MockObject) *valueobject.General {
	return valueobject.NewGeneral(obj)
}

func newMockDomainGeneral(d string, obj *MockObject) *valueobject.DomainGeneral {
	return valueobject.NewDomainGeneral(newMockGeneral(obj), d)
}

func newMockFunction(obj *MockObject) *valueobject.Function {
	return valueobject.NewFunction(obj, nil)
}

func newMockDomainFunction(d string, obj *MockObject) *valueobject.DomainFunction {
	return valueobject.NewDomainFunction(newMockFunction(obj), d)
}

func newMockGroup(name string, id int) *MockGroup {
	claObj := newMockObject(id)
	claAttrObj := newMockObjectAttribute(id)
	claMethodObj := newMockObjectMethod(id)

	intfObj := newMockObjectInterface(id)
	intfMethodObj := newMockObjectInterfaceMethod(id)

	genObj := newMockObjectGeneral(id)

	funcObj := newMockObjectFunction(id)

	g := &MockGroup{
		NameFunc: func() string {
			return name
		},
		ClassesFunc: func() []*valueobject.Class {
			return []*valueobject.Class{
				newMockClass(claObj, claAttrObj, claMethodObj)}
		},
		InterfacesFunc: func() []*valueobject.Interface {
			return []*valueobject.Interface{
				newMockInterface(intfObj, []*MockObject{intfMethodObj})}
		},
		GeneralsFunc: func() []*valueobject.General {
			return []*valueobject.General{newMockGeneral(genObj)}
		},
		FunctionsFunc: func() []*valueobject.Function {
			return []*valueobject.Function{newMockFunction(funcObj)}
		},
		SubGroupsFunc: func() []valueobject.Group {
			return []valueobject.Group{}
		},
	}

	g.MockObjects = []*MockObject{
		claObj,
		claAttrObj,
		claMethodObj,
		intfObj,
		intfMethodObj,
		genObj,
		funcObj,
	}

	return g
}

func newTwoLevelMockGroup() *MockGroup {
	rootGroup := newMockGroup("root", 0)
	subGroup1 := newMockGroup("subgroup1", 1)
	subGroup2 := newMockGroup("subgroup2", 2)

	rootGroup.SubGroupsFunc = func() []valueobject.Group {
		return []valueobject.Group{subGroup1, subGroup2}
	}

	return rootGroup
}

type MockDomainGroup struct {
	*MockGroup
	domain               string
	DomainGeneralsFunc   func() []*valueobject.DomainGeneral
	DomainFunctionsFunc  func() []*valueobject.DomainFunction
	DomainClassesFunc    func() []*valueobject.DomainClass
	DomainInterfacesFunc func() []*valueobject.DomainInterface
}

func (m *MockDomainGroup) Domain() string {
	return m.domain
}

func (m *MockDomainGroup) DomainGenerals() []*valueobject.DomainGeneral {
	return m.DomainGeneralsFunc()
}

func (m *MockDomainGroup) DomainFunctions() []*valueobject.DomainFunction {
	return m.DomainFunctionsFunc()
}

func (m *MockDomainGroup) DomainClasses() []*valueobject.DomainClass {
	return m.DomainClassesFunc()
}

func (m *MockDomainGroup) DomainInterfaces() []*valueobject.DomainInterface {
	return m.DomainInterfacesFunc()
}

func newMockDomainGroup(name string, id int) *MockDomainGroup {
	domain := "testdomain"
	claObj := newMockObject(id)
	claAttrObj := newMockObjectAttribute(id)
	claMethodObj := newMockObjectMethod(id)

	intfObj := newMockObjectInterface(id)
	intfMethodObj := newMockObjectInterfaceMethod(id)

	genObj := newMockObjectGeneral(id)

	funcObj := newMockObjectFunction(id)

	g := &MockDomainGroup{
		MockGroup: &MockGroup{
			NameFunc: func() string {
				return name
			},
			SubGroupsFunc: func() []valueobject.Group {
				return []valueobject.Group{}
			},
		},
		domain: domain,
		DomainClassesFunc: func() []*valueobject.DomainClass {
			return []*valueobject.DomainClass{newMockDomainClass(domain, claObj, claAttrObj, claMethodObj)}
		},
		DomainInterfacesFunc: func() []*valueobject.DomainInterface {
			return []*valueobject.DomainInterface{newMockDomainInterface(domain, intfObj, []*MockObject{intfMethodObj})}
		},
		DomainGeneralsFunc: func() []*valueobject.DomainGeneral {
			return []*valueobject.DomainGeneral{newMockDomainGeneral(domain, genObj)}
		},
		DomainFunctionsFunc: func() []*valueobject.DomainFunction {
			return []*valueobject.DomainFunction{newMockDomainFunction(domain, funcObj)}
		},
	}

	g.MockObjects = []*MockObject{
		claObj,
		claAttrObj,
		claMethodObj,
		intfObj,
		intfMethodObj,
		genObj,
		funcObj,
	}

	return g
}

func newMockAggregate(domain string, id int) *valueobject.Aggregate {
	claObj := newMockObject(id)
	claAttrObj := newMockObjectAttribute(id)
	claMethodObj := newMockObjectMethod(id)
	cla := newMockDomainClass(domain, claObj, claAttrObj, claMethodObj)

	return valueobject.NewAggregate(valueobject.NewEntity(cla), fmt.Sprintf("aggregate%d", id))
}

func newMockAggregateGroup(domain string) *valueobject.AggregateGroup {
	return valueobject.NewAggregateGroup(newMockAggregate(domain, 0), domain)
}

func newMockAggregateGroupWithName(domain string, name string) *valueobject.AggregateGroup {
	agg := newMockAggregate(domain, 0)
	agg.Name = name
	return valueobject.NewAggregateGroup(agg, domain)
}

func newMockEntity(domain string, id int) *valueobject.Entity {
	claObj := newMockObject(id)
	claAttrObj := newMockObjectAttribute(id)
	claMethodObj := newMockObjectMethod(id)
	cla := newMockDomainClass(domain, claObj, claAttrObj, claMethodObj)

	return valueobject.NewEntity(cla)
}

func newMockValueObject(domain string, id int) *valueobject.ValueObject {
	claObj := newMockObject(id)
	claAttrObj := newMockObjectAttribute(id)
	claMethodObj := newMockObjectMethod(id)
	cla := newMockDomainClass(domain, claObj, claAttrObj, claMethodObj)

	return valueobject.NewValueObject(cla)
}

type MockNode struct {
	IDVal    string
	NameVal  string
	ColorVal string
}

func (n *MockNode) ID() string    { return n.IDVal }
func (n *MockNode) Name() string  { return n.NameVal }
func (n *MockNode) Color() string { return n.ColorVal }

func newMockNode(id int) *MockNode {
	return &MockNode{
		IDVal:    fmt.Sprintf("id%d", id),
		NameVal:  fmt.Sprintf("name%d", id),
		ColorVal: fmt.Sprintf("color%d", id),
	}
}

type MockElement struct {
	MockNode
	ChildrenVal []arch.Node
}

func (e *MockElement) Children() []arch.Nodes { return []arch.Nodes{e.ChildrenVal} }

func newMockElement(id int, children ...arch.Node) *MockElement {
	return &MockElement{
		MockNode: MockNode{
			IDVal:    fmt.Sprintf("element%d", id),
			NameVal:  fmt.Sprintf("element%d", id),
			ColorVal: fmt.Sprintf("element%d", id),
		},
		ChildrenVal: children,
	}
}

type MockOptions struct {
	ShowAllRel            bool
	ShowStructEmbeddedRel bool
}

func (o MockOptions) ShowAllRelations() bool {
	return o.ShowAllRel
}

func (o MockOptions) ShowStructEmbeddedRelations() bool {
	return o.ShowStructEmbeddedRel
}

// MockRelationMeta 是一个模拟的 RelationMeta 接口实现
type MockRelationMeta struct {
	metaType arch.RelationType
	metaPos  arch.RelationPos
}

// Type 实现 RelationMeta 接口的 Type 方法
func (m *MockRelationMeta) Type() arch.RelationType {
	return m.metaType
}

// Position 实现 RelationMeta 接口的 Position 方法
func (m *MockRelationMeta) Position() arch.RelationPos {
	return m.metaPos
}

// NewMockRelationMeta 创建一个新的 MockRelationMeta 实例
func NewMockRelationMeta(metaType arch.RelationType, metaPos arch.RelationPos) *MockRelationMeta {
	return &MockRelationMeta{
		metaType: metaType,
		metaPos:  metaPos,
	}
}

type MockRelationPos struct {
	fromPos arch.Position
	toPos   arch.Position
}

func (r *MockRelationPos) From() arch.Position { return r.fromPos }
func (r *MockRelationPos) To() arch.Position   { return r.toPos }
