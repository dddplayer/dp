package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"path"
	"strings"
)

type domainObj struct {
	*obj
	domain string
}

func (do *domainObj) Identifier() arch.ObjIdentifier {
	return &domainIdent{
		ident:  do.obj.id,
		Domain: do.domain,
	}
}

func (do *domainObj) OriginIdentifier() arch.ObjIdentifier {
	return do.obj.Identifier()
}

func (do *domainObj) Domain() string {
	return do.domain
}

type domainIdent struct {
	*ident
	Domain string
}

func (di *domainIdent) ID() string {
	return path.Join(path.Base(di.Domain), strings.TrimPrefix(di.ident.ID(), di.Domain))
}

type DomainGeneral struct {
	*domainObj
}

type DomainFunction struct {
	*domainObj
}

type DomainInterface struct {
	*domainObj
	Methods []*DomainFunction
}

type DomainAttr struct {
	*domainObj
}

type DomainClass struct {
	*domainObj
	Attributes []*DomainAttr
	Methods    []*DomainFunction
}

type Entity struct {
	*DomainClass
}

type ValueObject struct {
	*DomainClass
}

type Aggregate struct {
	*Entity
	Name string
}

func NewDomainClass(cla *Class, d string, attrs []*DomainAttr, methods []*DomainFunction) *DomainClass {
	return &DomainClass{
		domainObj: &domainObj{
			obj:    cla.obj,
			domain: d,
		},
		Attributes: attrs,
		Methods:    methods,
	}
}

func NewDomainAttr(attr *Attr, d string) *DomainAttr {
	return &DomainAttr{
		domainObj: &domainObj{
			obj:    attr.obj,
			domain: d,
		},
	}
}

func NewDomainFunction(f *Function, d string) *DomainFunction {
	return &DomainFunction{
		domainObj: &domainObj{
			obj:    f.obj,
			domain: d,
		},
	}
}

func NewDomainInterface(i *Interface, d string, methods []*DomainFunction) *DomainInterface {
	return &DomainInterface{
		domainObj: &domainObj{
			obj:    i.obj,
			domain: d,
		},
		Methods: methods,
	}
}

func NewDomainGeneral(g *General, d string) *DomainGeneral {
	return &DomainGeneral{
		domainObj: &domainObj{
			obj:    g.obj,
			domain: d,
		},
	}
}

func NewEntity(cla *DomainClass) *Entity {
	return &Entity{
		DomainClass: cla,
	}
}

func NewAggregate(e *Entity, name string) *Aggregate {
	return &Aggregate{
		Entity: e,
		Name:   name,
	}
}

func NewValueObject(cla *DomainClass) *ValueObject {
	return &ValueObject{
		DomainClass: cla,
	}
}
