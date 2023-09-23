package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
)

func objColor(object arch.Object) arch.ObjColor {
	switch object.(type) {
	case *valueobject.Aggregate:
		return arch.ColorAggregate
	case *valueobject.Entity:
		return arch.ColorEntity
	case *valueobject.ValueObject:
		return arch.ColorValueObject
	case *valueobject.Class:
		return arch.ColorClass
	case *valueobject.DomainInterface, *valueobject.Interface:
		return arch.ColorInterface
	case *valueobject.DomainAttr, *valueobject.Attr:
		return arch.ColorAttribute
	case *valueobject.Function, *valueobject.DomainFunction:
		return arch.ColorFunc
	case *valueobject.StringObj:
		return arch.ColorWhite
	default:
		return arch.ColorGeneral
	}
}

func objColorWithParent(object, parent arch.Object) arch.ObjColor {
	switch parent.(type) {
	case *valueobject.Aggregate, *valueobject.Entity, *valueobject.ValueObject, *valueobject.Class:
		switch object.(type) {
		case *valueobject.DomainFunction, *valueobject.Function:
			return arch.ColorMethod
		}
	}
	return arch.ColorFunc
}
