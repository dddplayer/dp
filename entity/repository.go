package entity

import (
	"github.com/dddplayer/core/valueobject"
)

type Repository interface {
	Find(id valueobject.Identifier) DomainObject
	Insert(obj DomainObject) error
	Walk(cb func(obj DomainObject) error)
}
