package persistence

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

func TestRelations(t *testing.T) {
	t.Run("Insert and Walk", func(t *testing.T) {
		kv := &Relations{}

		mockObject1 := &MockObject{id: &MockIdentifier{IDVal: "id1"}}
		mockObject2 := &MockObject{id: &MockIdentifier{IDVal: "id2"}}

		rel1 := &MockRelation{
			relationType: arch.RelationTypeAssociation,
			fromObject:   mockObject1,
		}
		rel2 := &MockRelation{
			relationType: arch.RelationTypeComposition,
			fromObject:   mockObject2,
		}

		if err := kv.Insert(rel1); err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if err := kv.Insert(rel2); err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		if len(kv.relations) != 2 {
			t.Errorf("Expected 2 relations, but got %d", len(kv.relations))
		}

		var visitedRelations []arch.Relation
		kv.Walk(func(rel arch.Relation) error {
			visitedRelations = append(visitedRelations, rel)
			return nil
		})

		if len(visitedRelations) != 2 {
			t.Errorf("Expected to visit 2 relations, but visited %d", len(visitedRelations))
		}
	})
}

func TestRadixTree_All(t *testing.T) {
	t.Run("Return All Identifiers", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := &RadixTree{}

		// 创建一些 MockIdentifier 并插入到 RadixTree 中
		mockID1 := &MockIdentifier{IDVal: "id1"}
		mockID2 := &MockIdentifier{IDVal: "id2"}
		mockID3 := &MockIdentifier{IDVal: "id3"}

		r.objIds = []arch.ObjIdentifier{mockID1, mockID2, mockID3}

		// 调用 All 函数
		identifiers := r.All()

		// 验证返回的切片是否包含了所有插入的标识符
		if len(identifiers) != 3 {
			t.Errorf("Expected 3 identifiers, but got %d", len(identifiers))
		}

		// 验证返回的标识符切片中是否包含了所有插入的标识符
		for _, id := range r.objIds {
			found := false
			for _, returnedID := range identifiers {
				if id.ID() == returnedID.ID() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find identifier with ID '%s', but not found", id.ID())
			}
		}
	})
}

func TestRelations_Walk(t *testing.T) {
	t.Run("Walk Function Called for Each Relation", func(t *testing.T) {
		// 创建 Relations 实例
		kv := &Relations{}

		// 创建一些 MockRelation
		mockRel1 := &MockRelation{
			relationType: arch.RelationTypeAssociation,
			fromObject:   &MockObject{},
		}
		mockRel2 := &MockRelation{
			relationType: arch.RelationTypeComposition,
			fromObject:   &MockObject{},
		}
		mockRel3 := &MockRelation{
			relationType: arch.RelationTypeNone,
			fromObject:   &MockObject{},
		}

		// 将关系添加到 Relations 实例
		kv.relations = []arch.Relation{mockRel1, mockRel2, mockRel3}

		// 创建一个函数来测试 walker
		var visitedRelations []arch.Relation
		kv.Walk(func(rel arch.Relation) error {
			visitedRelations = append(visitedRelations, rel)
			return nil
		})

		// 验证 walker 是否按顺序调用了所有关系
		if len(visitedRelations) != 3 {
			t.Errorf("Expected to visit 3 relations, but visited %d", len(visitedRelations))
		}
	})
}
