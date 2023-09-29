package valueobject

import (
	"path"
	"reflect"
	"strings"
	"testing"
)

func TestDomainObjMethods(t *testing.T) {
	domain := "example.com"
	obj := &obj{
		id:  &ident{name: "test", pkg: "testpkg"},
		pos: &pos{},
	}
	domainObj := &domainObj{
		obj:    obj,
		domain: domain,
	}

	identifier := domainObj.Identifier()
	if domainIdentifier, ok := identifier.(*domainIdent); !ok {
		t.Errorf("Expected Identifier to return a domainIdent, but got %T", identifier)
	} else {
		if domainIdentifier.Domain != domain {
			t.Errorf("Expected domainIdent's Domain to be %s, but got %s", domain, domainIdentifier.Domain)
		}
	}

	originIdentifier := domainObj.OriginIdentifier()
	if originIdentifier != obj.Identifier() {
		t.Errorf("Expected OriginIdentifier to return the same identifier, but got different ones")
	}

	if domainObj.Domain() != domain {
		t.Errorf("Expected Domain to return %s, but got %s", domain, domainObj.Domain())
	}
}

func TestDomainIdentIDMethod(t *testing.T) {
	domain := "example.com"
	baseName := "test"
	pkg := "testpkg"
	ident := &ident{
		name: baseName,
		pkg:  pkg,
	}
	domainIdent := &domainIdent{
		ident:  ident,
		Domain: domain,
	}

	expectedID := path.Join(path.Base(domain), strings.TrimPrefix(ident.ID(), domain))
	if domainIdent.ID() != expectedID {
		t.Errorf("Expected ID to be %s, but got %s", expectedID, domainIdent.ID())
	}
}

func TestDomainTypes(t *testing.T) {
	domain := "example.com"
	obj := &obj{
		id:  &ident{name: "test", pkg: "testpkg"},
		pos: &pos{},
	}
	domainObj := &domainObj{
		obj:    obj,
		domain: domain,
	}

	domainGeneral := &DomainGeneral{
		domainObj: domainObj,
	}

	domainFunction := &DomainFunction{
		domainObj: domainObj,
	}

	domainInterface := &DomainInterface{
		domainObj: domainObj,
		Methods: []*DomainFunction{
			domainFunction,
		},
	}

	domainAttr := &DomainAttr{
		domainObj: domainObj,
	}

	if domainGeneral.Identifier().ID() != domainObj.Identifier().ID() {
		t.Errorf("DomainGeneral's Identifier method is not working as expected")
	}

	if domainFunction.Identifier().ID() != domainObj.Identifier().ID() {
		t.Errorf("DomainFunction's Identifier method is not working as expected")
	}

	if domainInterface.Identifier().ID() != domainObj.Identifier().ID() {
		t.Errorf("DomainInterface's Identifier method is not working as expected")
	}

	if domainAttr.Identifier().ID() != domainObj.Identifier().ID() {
		t.Errorf("DomainAttr's Identifier method is not working as expected")
	}
}

func TestNewDomainClass(t *testing.T) {
	classObj := &obj{
		id:  &ident{name: "TestClass", pkg: "testpkg"},
		pos: &pos{},
	}
	class := &Class{
		obj: classObj,
	}

	var attrs []*DomainAttr
	var methods []*DomainFunction

	domain := "example.com"
	domainClass := NewDomainClass(class, domain, attrs, methods)

	if domainClass.domainObj.obj != classObj {
		t.Errorf("Expected domainObj.obj to be the same as classObj, but got different ones")
	}

	if domainClass.domainObj.domain != domain {
		t.Errorf("Expected domainObj.domain to be %s, but got %s", domain, domainClass.domainObj.domain)
	}
}

func TestDomainClassTypes(t *testing.T) {
	classObj := &obj{
		id:  &ident{name: "TestClass", pkg: "testpkg"},
		pos: &pos{},
	}

	var attrs []*DomainAttr
	var methods []*DomainFunction

	domain := "example.com"
	domainClass := &DomainClass{
		domainObj: &domainObj{
			obj:    classObj,
			domain: domain,
		},
		Attributes: attrs,
		Methods:    methods,
	}

	entity := &Entity{
		DomainClass: domainClass,
	}

	valueObject := &ValueObject{
		DomainClass: domainClass,
	}

	aggregateName := "TestAggregate"
	aggregate := &Aggregate{
		Entity: entity,
		Name:   aggregateName,
	}

	if entity.DomainClass != domainClass {
		t.Errorf("Entity's DomainClass is not set correctly")
	}

	if valueObject.DomainClass != domainClass {
		t.Errorf("ValueObject's DomainClass is not set correctly")
	}

	if aggregate.Entity != entity {
		t.Errorf("Aggregate's Entity is not set correctly")
	}
	if aggregate.Name != aggregateName {
		t.Errorf("Aggregate's Name is not set correctly")
	}
}

func TestNewDomainAttr(t *testing.T) {
	// 创建一个模拟的属性对象和领域标识符
	mockAttr := &Attr{
		obj: NewObj(&obj{
			id:  &ident{name: "TestAttr", pkg: "testpkg"},
			pos: &pos{},
		}),
	}

	domain := "example.com"

	// 调用 NewDomainAttr 函数创建 DomainAttr
	domainAttr := NewDomainAttr(mockAttr, domain)

	// 验证 DomainAttr 对象的属性是否正确设置
	if domainAttr.domainObj.obj != mockAttr.obj {
		t.Errorf("Expected domainObj.obj to be the same as mockAttr.obj, but got different ones")
	}

	if domainAttr.domainObj.domain != domain {
		t.Errorf("Expected domainObj.domain to be %s, but got %s", domain, domainAttr.domainObj.domain)
	}
}

func TestNewDomainFunction(t *testing.T) {
	// 创建一个模拟的函数对象和领域字符串
	mockFunction := &Function{
		obj: NewObj(&obj{
			id:  &ident{name: "TestFunction", pkg: "testpkg"},
			pos: &pos{},
		}),
	}

	domain := "example.com"

	// 调用 NewDomainFunction 函数创建 DomainFunction
	domainFunction := NewDomainFunction(mockFunction, domain)

	// 验证 DomainFunction 对象的属性是否正确设置
	if domainFunction.domainObj.obj != mockFunction.obj {
		t.Errorf("Expected domainObj.obj to be the same as mockFunction.obj, but got different ones")
	}

	if domainFunction.domainObj.domain != domain {
		t.Errorf("Expected domainObj.domain to be %s, but got %s", domain, domainFunction.domainObj.domain)
	}
}

func TestNewDomainInterface(t *testing.T) {
	// 创建一个模拟的接口对象、领域字符串和方法列表
	mockInterface := &Interface{
		obj: NewObj(&obj{
			id:  &ident{name: "TestInterface", pkg: "testpkg"},
			pos: &pos{},
		}),
	}
	domain := "example.com"
	mockMethods := []*DomainFunction{
		NewDomainFunction(&Function{
			obj: NewObj(&obj{
				id:  &ident{name: "Method1", pkg: "testpkg"},
				pos: &pos{},
			}),
		}, domain),
		NewDomainFunction(&Function{
			obj: NewObj(&obj{
				id:  &ident{name: "Method2", pkg: "testpkg"},
				pos: &pos{},
			}),
		}, domain),
	}

	// 调用 NewDomainInterface 函数创建 DomainInterface
	domainInterface := NewDomainInterface(mockInterface, domain, mockMethods)

	// 验证 DomainInterface 对象的属性是否正确设置
	if domainInterface.domainObj.obj != mockInterface.obj {
		t.Errorf("Expected domainObj.obj to be the same as mockInterface.obj, but got different ones")
	}

	if domainInterface.domainObj.domain != domain {
		t.Errorf("Expected domainObj.domain to be %s, but got %s", domain, domainInterface.domainObj.domain)
	}

	if !reflect.DeepEqual(domainInterface.Methods, mockMethods) {
		t.Error("Expected Methods to be the same as mockMethods, but they are different")
	}
}

func TestNewDomainGeneral(t *testing.T) {
	// 创建一个模拟的通用对象和领域字符串
	mockGeneral := &General{
		obj: NewObj(&obj{
			id:  &ident{name: "TestGeneral", pkg: "testpkg"},
			pos: &pos{},
		}),
	}
	domain := "example.com"

	// 调用 NewDomainGeneral 函数创建 DomainGeneral
	domainGeneral := NewDomainGeneral(mockGeneral, domain)

	// 验证 DomainGeneral 对象的属性是否正确设置
	if domainGeneral.domainObj.obj != mockGeneral.obj {
		t.Errorf("Expected domainObj.obj to be the same as mockGeneral.obj, but got different ones")
	}

	if domainGeneral.domainObj.domain != domain {
		t.Errorf("Expected domainObj.domain to be %s, but got %s", domain, domainGeneral.domainObj.domain)
	}
}

func TestNewEntity(t *testing.T) {
	// 创建一个模拟的领域类对象
	mockClass := &DomainClass{
		domainObj: &domainObj{
			obj: &obj{
				id:  &ident{name: "TestClass", pkg: "testpkg"},
				pos: &pos{},
			},
			domain: "example.com",
		},
		Attributes: []*DomainAttr{},
		Methods:    []*DomainFunction{},
	}

	// 调用 NewEntity 函数创建 Entity 对象
	entity := NewEntity(mockClass)

	// 验证 Entity 对象的属性是否正确设置
	if entity.DomainClass != mockClass {
		t.Errorf("Expected Entity's DomainClass to be the same as mockClass, but got different ones")
	}
}

func TestNewAggregate(t *testing.T) {
	// 创建一个模拟的 Entity 对象和名称
	mockEntity := &Entity{
		DomainClass: &DomainClass{
			domainObj: &domainObj{
				obj: &obj{
					id:  &ident{name: "TestClass", pkg: "testpkg"},
					pos: &pos{},
				},
				domain: "example.com",
			},
			Attributes: []*DomainAttr{},
			Methods:    []*DomainFunction{},
		},
	}

	aggregateName := "TestAggregate"

	// 调用 NewAggregate 函数创建 Aggregate 对象
	aggregate := NewAggregate(mockEntity, aggregateName)

	// 验证 Aggregate 对象的属性是否正确设置
	if aggregate.Entity != mockEntity {
		t.Errorf("Expected Aggregate's Entity to be the same as mockEntity, but got different ones")
	}

	if aggregate.Name != aggregateName {
		t.Errorf("Expected Aggregate's Name to be %s, but got %s", aggregateName, aggregate.Name)
	}
}

func TestNewValueObject(t *testing.T) {
	// 创建一个模拟的 DomainClass 对象
	mockDomainClass := &DomainClass{
		domainObj: &domainObj{
			obj: &obj{
				id:  &ident{name: "TestClass", pkg: "testpkg"},
				pos: &pos{},
			},
			domain: "example.com",
		},
		Attributes: []*DomainAttr{},
		Methods:    []*DomainFunction{},
	}

	// 调用 NewValueObject 函数创建 ValueObject 对象
	valueObject := NewValueObject(mockDomainClass)

	// 验证 ValueObject 对象的属性是否正确设置
	if valueObject.DomainClass != mockDomainClass {
		t.Errorf("Expected ValueObject's DomainClass to be the same as mockDomainClass, but got different ones")
	}
}
