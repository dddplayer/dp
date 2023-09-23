package valueobject

import (
	"path"
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
