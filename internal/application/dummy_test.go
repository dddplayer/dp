package application

import (
	"errors"
	"github.com/dddplayer/dp/internal/domain/arch"
)

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

var main = `package main

import (
	"fmt"
	"module/cmd"
	"module/internal/domain/test/entity"
	"module/internal/domain/test/valueobject"
	"module/pkg"
)

func main() {
	cmd.Func1()
	pkg.Func1()
	t := &entity.Test{}
	vo := &valueobject.VO{}
	fmt.Println(t, vo)
}`
