package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

func TestGroup(t *testing.T) {
	groupName := "testGroup"
	g := &group{
		name: groupName,
	}

	if g.Name() != groupName {
		t.Errorf("Expected Name to be %s, but got %s", groupName, g.Name())
	}

	subGroupName := "subGroup"
	subGroup := &group{
		name: subGroupName,
	}
	g.AppendGroups(subGroup)

	subGroups := g.SubGroups()
	if len(subGroups) != 1 || subGroups[0] != subGroup {
		t.Errorf("Expected SubGroups to contain subGroup")
	}

	testObject := &Class{}

	g.AppendObjects(testObject)
	objects := g.Objects()
	if len(objects) != 1 || objects[0] != testObject {
		t.Errorf("Expected Objects to contain testObject")
	}
}

func TestGroupMethods(t *testing.T) {
	groupName := "testGroup"
	group := &group{
		name: groupName,
	}

	class1 := &Class{}
	class2 := &Class{}
	group.AppendObjects(class1, class2)
	classes := group.Classes()
	if len(classes) != 2 || classes[0] != class1 || classes[1] != class2 {
		t.Errorf("Expected Classes to contain class1 and class2")
	}

	general1 := &General{}
	general2 := &General{}
	group.AppendObjects(general1, general2)
	generals := group.Generals()
	if len(generals) != 2 || generals[0] != general1 || generals[1] != general2 {
		t.Errorf("Expected Generals to contain general1 and general2")
	}

	function1 := &Function{}
	function2 := &Function{}
	group.AppendObjects(function1, function2)
	functions := group.Functions()
	if len(functions) != 2 || functions[0] != function1 || functions[1] != function2 {
		t.Errorf("Expected Functions to contain function1 and function2")
	}

	interface1 := &Interface{}
	interface2 := &Interface{}
	group.AppendObjects(interface1, interface2)
	interfaces := group.Interfaces()
	if len(interfaces) != 2 || interfaces[0] != interface1 || interfaces[1] != interface2 {
		t.Errorf("Expected Interfaces to contain interface1 and interface2")
	}
}

func TestDomainGroupMethods(t *testing.T) {
	groupName := "testGroup"
	group := &group{
		name: groupName,
	}

	domain := "example.com"
	domainGroup := &domainGroup{
		group:  group,
		domain: domain,
	}

	// Add more testing code similar to the previous methods

	classObj := &obj{
		id:  &ident{name: "TestClass", pkg: "testpkg"},
		pos: &pos{},
	}
	class := &Class{
		obj: classObj,
	}
	group.AppendObjects(class)

	generalObj := &obj{
		id:  &ident{name: "TestGeneral", pkg: "testpkg"},
		pos: &pos{},
	}
	general := &General{
		obj: generalObj,
	}
	group.AppendObjects(general)

	functionObj := &obj{
		id:  &ident{name: "TestFunction", pkg: "testpkg"},
		pos: &pos{},
	}
	function := &Function{
		obj: functionObj,
	}
	group.AppendObjects(function)

	interfaceObj := &obj{
		id:  &ident{name: "TestInterface", pkg: "testpkg"},
		pos: &pos{},
	}
	iface := &Interface{
		obj: interfaceObj,
	}
	group.AppendObjects(iface)

	domainClasses := domainGroup.DomainClasses()
	if len(domainClasses) != 1 || domainClasses[0].domainObj.obj != classObj {
		t.Errorf("Expected DomainClasses to contain classObj")
	}

	domainGenerals := domainGroup.DomainGenerals()
	if len(domainGenerals) != 1 || domainGenerals[0].domainObj.obj != generalObj {
		t.Errorf("Expected DomainGenerals to contain generalObj")
	}

	domainFunctions := domainGroup.DomainFunctions()
	if len(domainFunctions) != 1 || domainFunctions[0].domainObj.obj != functionObj {
		t.Errorf("Expected DomainFunctions to contain functionObj")
	}

	domainInterfaces := domainGroup.DomainInterfaces()
	if len(domainInterfaces) != 1 || domainInterfaces[0].domainObj.obj != interfaceObj {
		t.Errorf("Expected DomainInterfaces to contain interfaceObj")
	}
}

func TestDomainGroupMethods_Getters(t *testing.T) {
	groupName := "testGroup"
	group := &group{
		name: groupName,
	}

	domain := "example.com"
	domainGroup := &domainGroup{
		group:  group,
		domain: domain,
	}

	attrObj := &obj{
		id:  &ident{name: "TestAttr", pkg: "testpkg"},
		pos: &pos{},
	}
	attr := &Attr{
		obj: attrObj,
	}
	group.AppendObjects(attr)

	functionObj := &obj{
		id:  &ident{name: "TestFunction", pkg: "testpkg"},
		pos: &pos{},
	}
	function := &Function{
		obj: functionObj,
		Receiver: &ident{
			// Add actual Receiver properties
		},
	}
	group.AppendObjects(function)

	interfaceMethodObj := &obj{
		id:  &ident{name: "TestInterfaceMethod", pkg: "testpkg"},
		pos: &pos{},
	}
	interfaceMethod := &InterfaceMethod{
		obj: interfaceMethodObj,
		// Add actual InterfaceMethod properties
	}
	group.AppendObjects(interfaceMethod)

	attrIdent := &ident{name: "TestAttr", pkg: "testpkg"}
	domainAttr := domainGroup.getDomainAttr(attrIdent)
	if domainAttr == nil || domainAttr.domainObj.obj != attrObj {
		t.Errorf("Expected getDomainAttr to return the correct DomainAttr object")
	}

	functionIdent := &ident{name: "TestFunction", pkg: "testpkg"}
	domainFunction := domainGroup.getDomainFunction(functionIdent)
	if domainFunction == nil || domainFunction.domainObj.obj != functionObj {
		t.Errorf("Expected getDomainFunction to return the correct DomainFunction object")
	}

	interfaceFunctionIdent := &ident{name: "TestInterfaceMethod", pkg: "testpkg"}
	interfaceFunction := domainGroup.getInterfaceFunction(interfaceFunctionIdent)
	if interfaceFunction == nil || interfaceFunction.domainObj.obj != interfaceMethodObj {
		t.Errorf("Expected getInterfaceFunction to return the correct DomainFunction object")
	}
}

func TestEntityGroupMethods(t *testing.T) {
	// Create a new EntityGroup
	domain := "example.com"
	entityGroup := NewEntityGroup(domain)

	// Verify that the domainGroup object has been created correctly
	if entityGroup.domainGroup == nil {
		t.Errorf("Expected domainGroup to be initialized")
	}
	if entityGroup.domainGroup.domain != domain {
		t.Errorf("Expected domainGroup's domain to be %s, but got %s", domain, entityGroup.domainGroup.domain)
	}
	if entityGroup.domainGroup.group == nil {
		t.Errorf("Expected group to be initialized in domainGroup")
	}

	// Create some mock arch.Object objects
	classObj1 := &obj{pos: &pos{}}
	classObj2 := &obj{pos: &pos{}}
	classObj3 := &obj{pos: &pos{}}
	classes := []arch.Object{classObj1, classObj2, classObj3}

	// Create a new EntityGroup with mock arch.Object objects
	entityGroupWithClasses := NewEntityGroup(domain, classes...)

	// Verify that the group's objs field contains the correct arch.Object objects
	if len(entityGroupWithClasses.objs) != len(classes) {
		t.Errorf("Expected %d arch.Object objects in entityGroupWithClasses, but got %d", len(classes), len(entityGroupWithClasses.objs))
	}
	for i, obj := range entityGroupWithClasses.objs {
		if obj != classes[i] {
			t.Errorf("Expected obj at index %d to match classes[%d]", i, i)
		}
	}
}

func TestEntitiesMethods(t *testing.T) {
	c1 := &Class{obj: &obj{id: &ident{name: "i1", pkg: "p1"}, pos: &pos{}}}
	c2 := &Class{obj: &obj{id: &ident{name: "i2", pkg: "p2"}, pos: &pos{}}}
	class1 := &DomainClass{domainObj: &domainObj{obj: c1.obj}}
	class2 := &DomainClass{domainObj: &domainObj{obj: c2.obj}}
	entity1 := &Entity{DomainClass: class1}
	entity2 := &Entity{DomainClass: class2}
	entities := Entities{entity1, entity2}

	// Create a mock EntityGroup with mock DomainClasses
	entityGroup := &EntityGroup{
		domainGroup: &domainGroup{
			group: &group{
				objs: []arch.Object{c1, c2},
			},
		},
	}

	objects := entities.Objects()
	if len(objects) != len(entities) {
		t.Errorf("Expected %d objects, but got %d", len(entities), len(objects))
	}
	for i, obj := range objects {
		if obj != entities[i] {
			t.Errorf("Expected object at index %d to match entities[%d]", i, i)
		}
	}

	// Test EntityGroup's Entities method
	entityGroupEntities := entityGroup.Entities()
	if len(entityGroupEntities) != len(entities) {
		t.Errorf("Expected %d Entity instances in entityGroupEntities, but got %d", len(entities), len(entityGroupEntities))
	}
	for i, entity := range entityGroupEntities {
		if entity.DomainClass.Identifier().ID() != entities[i].DomainClass.Identifier().ID() {
			t.Errorf("Expected DomainClass at index %d to match domainClasses[%d]", i, i)
		}
	}
}

func TestVOGroupMethods(t *testing.T) {
	// Create a new VOGroup
	domain := "example.com"
	voGroup := NewVOGroup(domain)

	// Verify that the domainGroup object has been created correctly
	if voGroup.domainGroup == nil {
		t.Errorf("Expected domainGroup to be initialized")
	}
	if voGroup.domainGroup.domain != domain {
		t.Errorf("Expected domainGroup's domain to be %s, but got %s", domain, voGroup.domainGroup.domain)
	}
	if voGroup.domainGroup.group == nil {
		t.Errorf("Expected group to be initialized in domainGroup")
	}

	// Create some mock arch.Object objects
	voObj1 := &obj{pos: &pos{}}
	voObj2 := &obj{pos: &pos{}}
	vos := []arch.Object{voObj1, voObj2}

	// Create a new VOGroup with mock arch.Object objects
	voGroupWithVOs := NewVOGroup(domain, vos...)

	// Verify that the group's objs field contains the correct arch.Object objects
	if len(voGroupWithVOs.objs) != len(vos) {
		t.Errorf("Expected %d arch.Object objects in voGroupWithVOs, but got %d", len(vos), len(voGroupWithVOs.objs))
	}
	for i, obj := range voGroupWithVOs.objs {
		if obj != vos[i] {
			t.Errorf("Expected obj at index %d to match vos[%d]", i, i)
		}
	}
}

func TestValueObjectsMethods(t *testing.T) {
	// Create some mock ValueObject instances
	vo1 := &ValueObject{}
	vo2 := &ValueObject{}
	valueObjects := ValueObjects{vo1, vo2}

	// Test ValueObjects' Objects method
	objects := valueObjects.Objects()
	if len(objects) != len(valueObjects) {
		t.Errorf("Expected %d objects, but got %d", len(valueObjects), len(objects))
	}
	for i, obj := range objects {
		if obj != valueObjects[i] {
			t.Errorf("Expected object at index %d to match valueObjects[%d]", i, i)
		}
	}

	// Add more validation...
}

func TestVOMethods(t *testing.T) {
	c1 := &Class{obj: &obj{id: &ident{name: "i1", pkg: "p1"}, pos: &pos{}}}
	c2 := &Class{obj: &obj{id: &ident{name: "i2", pkg: "p2"}, pos: &pos{}}}
	class1 := &DomainClass{domainObj: &domainObj{obj: c1.obj}}
	class2 := &DomainClass{domainObj: &domainObj{obj: c2.obj}}
	vo1 := &ValueObject{DomainClass: class1}
	vo2 := &ValueObject{DomainClass: class2}
	valueObjects := ValueObjects{vo1, vo2}

	// Create a mock VOGroup with mock DomainClasses
	voGroup := &VOGroup{
		domainGroup: &domainGroup{
			group: &group{
				objs: []arch.Object{c1, c2},
			},
		},
	}

	// Test ValueObjects' Objects method
	objects := valueObjects.Objects()
	if len(objects) != len(valueObjects) {
		t.Errorf("Expected %d objects, but got %d", len(valueObjects), len(objects))
	}
	for i, obj := range objects {
		if obj != valueObjects[i] {
			t.Errorf("Expected object at index %d to match valueObjects[%d]", i, i)
		}
	}

	// Test VOGroup's ValueObjects method
	voGroupValueObjects := voGroup.ValueObjects()
	if len(voGroupValueObjects) != len(valueObjects) {
		t.Errorf("Expected %d ValueObject instances in voGroupValueObjects, but got %d", len(valueObjects), len(voGroupValueObjects))
	}
	for i, vo := range voGroupValueObjects {
		if vo.DomainClass.Identifier().ID() != valueObjects[i].DomainClass.Identifier().ID() {
			t.Errorf("Expected DomainClass at index %d to match valueObjects[%d]", i, i)
		}
	}

	// Add more validation...
}

func TestAggregateGroupMethods(t *testing.T) {
	// Create a mock Aggregate
	aggregate := &Aggregate{Name: "TestAggregate"}
	domain := "example.com"

	// Create a new AggregateGroup
	aggregateGroup := NewAggregateGroup(aggregate, domain)

	// Verify that the domainGroup object has been created correctly
	if aggregateGroup.domainGroup == nil {
		t.Errorf("Expected domainGroup to be initialized")
	}
	if aggregateGroup.domainGroup.domain != domain {
		t.Errorf("Expected domainGroup's domain to be %s, but got %s", domain, aggregateGroup.domainGroup.domain)
	}
	if aggregateGroup.domainGroup.group == nil {
		t.Errorf("Expected group to be initialized in domainGroup")
	}
	if len(aggregateGroup.domainGroup.group.objs) != 1 {
		t.Errorf("Expected 1 object in group's objs, but got %d", len(aggregateGroup.domainGroup.group.objs))
	}
	if aggregateGroup.domainGroup.group.objs[0] != aggregate {
		t.Errorf("Expected the object in group's objs to be the Aggregate instance")
	}

	// Test the DomainName method
	domainName := aggregateGroup.DomainName()
	if domainName != domain {
		t.Errorf("Expected DomainName() to return %s, but got %s", domain, domainName)
	}

	// Add more validation...
}

func TestAggregateGroupMethodsAggregate(t *testing.T) {
	// Create a mock Aggregate
	aggregate := &Aggregate{Name: "TestAggregate"}

	// Create a mock Class
	entity := &Class{obj: &obj{id: &ident{name: "TestAggregate"}, pos: &pos{}}}

	// Create a mock EntityGroup
	entityGroup := &EntityGroup{
		domainGroup: &domainGroup{
			group: &group{
				objs: []arch.Object{entity},
			},
			domain: "example.com",
		},
	}

	// Create an AggregateGroup with the mock Aggregate and EntityGroup
	aggregateGroup := &AggregateGroup{
		domainGroup: &domainGroup{
			group: &group{
				objs:      []arch.Object{aggregate},
				subGroups: []Group{entityGroup},
			},
			domain: "example.com",
		},
	}

	// Test the Aggregate method
	agg, err := aggregateGroup.Aggregate()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if agg == nil {
		t.Errorf("Expected Aggregate instance, but got nil")
	}
	if agg != aggregate {
		t.Errorf("Expected returned Aggregate to be the mock instance")
	}
	if agg.Entity == nil {
		t.Errorf("Expected Entity to be set in the returned Aggregate")
	}
}

func TestAggregateGroupMethods_IsValid(t *testing.T) {
	// Create a mock Aggregate
	aggregate := &Aggregate{Name: "TestAggregate"}

	// Create an AggregateGroup without an associated Entity
	aggregateGroupWithoutEntity := &AggregateGroup{
		domainGroup: &domainGroup{
			group: &group{
				objs: []arch.Object{aggregate},
			},
			domain: "example.com",
		},
	}

	// Test the IsValid method for the case without an associated Entity
	isValidWithoutEntity := aggregateGroupWithoutEntity.IsValid()
	if isValidWithoutEntity {
		t.Errorf("Expected IsValid() to return false, but got true")
	}

	// Create a mock Class
	entity := &Class{obj: &obj{id: &ident{name: "TestAggregate"}, pos: &pos{}}}

	// Create a mock EntityGroup
	entityGroup := &EntityGroup{
		domainGroup: &domainGroup{
			group: &group{
				objs: []arch.Object{entity},
			},
			domain: "example.com",
		},
	}

	// Create an AggregateGroup with the mock Aggregate and EntityGroup
	aggregateGroup := &AggregateGroup{
		domainGroup: &domainGroup{
			group: &group{
				objs:      []arch.Object{aggregate},
				subGroups: []Group{entityGroup},
			},
			domain: "example.com",
		},
	}

	// Test the IsValid method
	isValid := aggregateGroup.IsValid()
	if !isValid {
		t.Errorf("Expected IsValid() to return true, but got false")
	}

	// Add more validation...
}
