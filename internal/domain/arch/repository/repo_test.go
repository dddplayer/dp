package repository

import (
	"errors"
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

type MockObject struct {
	id       arch.ObjIdentifier
	position arch.Position
}

func (mo MockObject) Identifier() arch.ObjIdentifier {
	return mo.id
}

func (mo MockObject) Position() arch.Position {
	return mo.position
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
	return nil, errors.New("some objects not found in mock repository")
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

type MockIdentifier struct {
	IDVal                  string
	NameVal                string
	NameSeparatorLengthVal int
	DirVal                 string
}

func (mi *MockIdentifier) ID() string {
	return mi.IDVal
}

func (mi *MockIdentifier) Name() string {
	return mi.NameVal
}

func (mi *MockIdentifier) NameSeparatorLength() int {
	return mi.NameSeparatorLengthVal
}

func (mi *MockIdentifier) Dir() string {
	return mi.DirVal
}

type MockPosition struct {
	FilenameVal string
	OffsetVal   int
	LineVal     int
	ColumnVal   int
}

func (mp *MockPosition) Filename() string {
	return mp.FilenameVal
}

func (mp *MockPosition) Offset() int {
	return mp.OffsetVal
}

func (mp *MockPosition) Line() int {
	return mp.LineVal
}

func (mp *MockPosition) Column() int {
	return mp.ColumnVal
}

func TestObjectRepository(t *testing.T) {
	mockIdentifier1 := &MockIdentifier{IDVal: "id1"}
	mockPosition1 := &MockPosition{FilenameVal: "file1", OffsetVal: 10, LineVal: 5, ColumnVal: 2}
	mockObject1 := MockObject{id: mockIdentifier1, position: mockPosition1}

	mockIdentifier2 := &MockIdentifier{IDVal: "id2"}
	mockPosition2 := &MockPosition{FilenameVal: "file2", OffsetVal: 20, LineVal: 8, ColumnVal: 4}
	mockObject2 := MockObject{id: mockIdentifier2, position: mockPosition2}

	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{mockIdentifier1, mockIdentifier2},
	}
	_ = mockRepo.Insert(mockObject1)
	_ = mockRepo.Insert(mockObject2)

	t.Run("Test Find", func(t *testing.T) {
		foundObj := mockRepo.Find(mockIdentifier1)
		if foundObj == nil {
			t.Errorf("Expected to find object, but got nil")
		}
	})

	t.Run("Test GetObjects", func(t *testing.T) {
		ids := []arch.ObjIdentifier{mockIdentifier1, mockIdentifier2}
		objects, err := mockRepo.GetObjects(ids)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if len(objects) != len(ids) {
			t.Errorf("Expected %d objects, but got %d", len(ids), len(objects))
		}
	})

	t.Run("Test All", func(t *testing.T) {
		allIdents := mockRepo.All()
		if len(allIdents) != len(mockRepo.idents) {
			t.Errorf("Expected %d identifiers, but got %d", len(mockRepo.idents), len(allIdents))
		}
	})
}

// 模拟一个用于测试的 Relation
type MockRelation struct {
	relationType arch.RelationType
	fromObject   arch.Object
}

func (mr *MockRelation) Type() arch.RelationType {
	return mr.relationType
}

func (mr *MockRelation) From() arch.Object {
	return mr.fromObject
}

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

func TestRelationRepository(t *testing.T) {
	mockRepo := &MockRelationRepository{
		relations: make([]arch.Relation, 0),
	}

	mockObject1 := &MockObject{id: &MockIdentifier{IDVal: "id1"}}
	mockObject2 := &MockObject{id: &MockIdentifier{IDVal: "id2"}}

	mockRelation1 := &MockRelation{
		relationType: arch.RelationTypeAssociation,
		fromObject:   mockObject1,
	}
	mockRelation2 := &MockRelation{
		relationType: arch.RelationTypeComposition,
		fromObject:   mockObject2,
	}

	t.Run("Test Insert", func(t *testing.T) {
		err := mockRepo.Insert(mockRelation1)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if len(mockRepo.relations) != 1 {
			t.Errorf("Expected 1 relation, but got %d", len(mockRepo.relations))
		}
	})

	t.Run("Test Walk", func(t *testing.T) {
		mockRepo.relations = append(mockRepo.relations, mockRelation1, mockRelation2)

		walkedRelations := 0
		mockRepo.Walk(func(rel arch.Relation) error {
			walkedRelations++
			return nil
		})

		if walkedRelations != len(mockRepo.relations) {
			t.Errorf("Expected %d walked relations, but got %d", len(mockRepo.relations), walkedRelations)
		}
	})
}
