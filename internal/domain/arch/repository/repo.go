package repository

import "github.com/dddplayer/dp/internal/domain/arch"

type ObjectRepository interface {
	Find(id arch.ObjIdentifier) arch.Object
	GetObjects(ids []arch.ObjIdentifier) ([]arch.Object, error)
	All() []arch.ObjIdentifier
	Insert(obj arch.Object) error
	Walk(walker func(obj arch.Object) error)
}

type RelationRepository interface {
	Insert(rel arch.Relation) error
	Walk(walker func(rel arch.Relation) error)
}
