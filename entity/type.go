package entity

import "github.com/dddplayer/core/valueobject"

type ComponentKind string

const (
	KindEntity      ComponentKind = "entity"
	KindValueObject ComponentKind = "valueobject"
	KindFactory     ComponentKind = "factory"
	KindService     ComponentKind = "service"
	KindRepository  ComponentKind = "repository"
	KindGeneral     ComponentKind = "general"
	KindClass       ComponentKind = "class"
	KindFunc        ComponentKind = "function"
)

type ObjColor string

const (
	ColorEntity      ObjColor = "#ffe599ff"
	ColorValueObject ObjColor = "#a2c4c9ff"
	ColorFactory     ObjColor = "#cfe2f3ff"
	ColorService     ObjColor = "#e69138ff"
	ColorRepository  ObjColor = "#ffffffff"
	ColorMethod      ObjColor = "#a4c2f4ff"
	ColorAttribute   ObjColor = "#f3f3f3ff"
	ColorInterface   ObjColor = "#9fc5e8ff"
	ColorClass       ObjColor = "#b4a7d6ff"
	ColorGeneral     ObjColor = "#f4ccccff"
	ColorFunc        ObjColor = "#ead1dcff"
)

type DomainObject interface {
	Identifier() *valueobject.Identifier
	Position() *valueobject.Position
}

type DomainComponent interface {
	DomainObject
	Kind() ComponentKind
}

type DomainClass interface {
	DomainObject
	Commands() []*valueobject.Identifier
	AppendCommand(id *valueobject.Identifier)
	Attributes() []*valueobject.Identifier
	AppendAttribute(id *valueobject.Identifier)
}

type DomainKeyObj interface {
	DomainComponent
	Commands() []*valueobject.Identifier
	AppendCommand(id *valueobject.Identifier)
	Attributes() []*valueobject.Identifier
	AppendAttribute(id *valueobject.Identifier)
}

type obj struct {
	id  *valueobject.Identifier
	pos *valueobject.Position
}

func (o *obj) Identifier() *valueobject.Identifier { return o.id }
func (o *obj) Position() *valueobject.Position     { return o.pos }

type (
	General struct {
		*obj
	}

	Class struct {
		*obj
		attrs   []*valueobject.Identifier
		methods []*valueobject.Identifier
	}

	Function struct {
		*obj
	}

	Attr struct {
		*obj
	}

	keyObj struct {
		*Class
	}

	Entity struct {
		*keyObj
	}

	ValueObject struct {
		*keyObj
	}

	Service struct {
		Function
	}

	Factory struct {
		Function
	}

	Interface struct {
		*obj
		Methods []*valueobject.Identifier
	}

	InterfaceMethod struct {
		*obj
	}
)

func (k *Class) Attributes() []*valueobject.Identifier      { return k.attrs }
func (k *Class) AppendAttribute(id *valueobject.Identifier) { k.attrs = append(k.attrs, id) }
func (k *Class) Commands() []*valueobject.Identifier        { return k.methods }
func (k *Class) AppendCommand(id *valueobject.Identifier)   { k.methods = append(k.methods, id) }

func (i *Interface) Append(m *InterfaceMethod) { i.Methods = append(i.Methods, m.id) }
func (e *Entity) Kind() ComponentKind          { return KindEntity }
func (v *ValueObject) Kind() ComponentKind     { return KindValueObject }
func (s *Service) Kind() ComponentKind         { return KindService }
func (s *Factory) Kind() ComponentKind         { return KindFactory }
