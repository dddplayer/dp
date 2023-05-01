package entity

import (
	"github.com/dddplayer/core/valueobject"
	"strings"
)

func IsEntity(path string) bool {
	return strings.Contains(path, string(KindEntity))
}

func IsValueObject(path string) bool {
	return strings.Contains(path, string(KindValueObject))
}

func IsFactory(path string) bool {
	return strings.Contains(path, string(KindFactory))
}

func IsService(path string) bool {
	return strings.Contains(path, string(KindService))
}

func NewEntity(id valueobject.Identifier, pos valueobject.Position) *Entity {
	return &Entity{
		keyObj: &keyObj{
			Class: &Class{
				obj: &obj{
					id:  &id,
					pos: &pos,
				},
				attrs:   []*valueobject.Identifier{},
				methods: []*valueobject.Identifier{},
			},
		},
	}
}

func NewValueObject(id valueobject.Identifier, pos valueobject.Position) *ValueObject {
	return &ValueObject{
		keyObj: &keyObj{
			Class: &Class{
				obj: &obj{
					id:  &id,
					pos: &pos,
				},
				attrs:   []*valueobject.Identifier{},
				methods: []*valueobject.Identifier{},
			},
		},
	}
}

func NewService(id valueobject.Identifier, pos valueobject.Position) *Service {
	return &Service{
		Function: *NewFunction(id, pos),
	}
}

func NewFactory(id valueobject.Identifier, pos valueobject.Position) *Factory {
	return &Factory{
		Function: *NewFunction(id, pos),
	}
}

func NewInterface(id valueobject.Identifier, pos valueobject.Position) *Interface {
	return &Interface{
		obj: &obj{
			id:  &id,
			pos: &pos,
		},
		Methods: []*valueobject.Identifier{},
	}
}

func NewInterfaceMethod(id valueobject.Identifier, pos valueobject.Position) *InterfaceMethod {
	return &InterfaceMethod{
		obj: &obj{
			id:  &id,
			pos: &pos,
		},
	}
}

func NewFunction(id valueobject.Identifier, pos valueobject.Position) *Function {
	return &Function{
		obj: &obj{
			id:  &id,
			pos: &pos,
		},
	}
}

type GeneralComponent struct {
	*General
}

func (g *GeneralComponent) Kind() ComponentKind {
	return KindGeneral
}

func NewGeneralComponent(g *General) *GeneralComponent {
	return &GeneralComponent{
		General: g,
	}
}

type ClassComponent struct {
	*Class
}

func (c *ClassComponent) Kind() ComponentKind {
	return KindClass
}

func NewClassComponent(c *Class) *ClassComponent {
	return &ClassComponent{
		Class: c,
	}
}

type FunctionComponent struct {
	*Function
}

func (f *FunctionComponent) Kind() ComponentKind {
	return KindFunc
}

func NewFunctionComponent(f *Function) *FunctionComponent {
	return &FunctionComponent{
		Function: f,
	}
}
