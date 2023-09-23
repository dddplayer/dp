package factory

import (
	"errors"
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
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
	return nil, errors.New("some objects not found in mock repository of factory")
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

func TestNewArch(t *testing.T) {
	mockScope := "testScope"
	mockObjectRepo := &MockObjectRepository{}
	mockRelationRepo := &MockRelationRepository{}

	newArch, err := NewArch(mockScope, mockObjectRepo, mockRelationRepo)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Verify that the Arch object was created with the correct values
	if newArch.Scope != mockScope {
		t.Errorf("Expected scope %s, but got %s", mockScope, newArch.Scope)
	}
	if newArch.ObjRepo != mockObjectRepo {
		t.Errorf("Expected ObjRepo to match, but it didn't")
	}
	if newArch.RelRepo != mockRelationRepo {
		t.Errorf("Expected RelRepo to match, but it didn't")
	}
}
