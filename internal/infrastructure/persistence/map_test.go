package persistence

import (
	"bytes"
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"io"
	"os"
	"strings"
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

func TestRelations_Walk_ErrorHandling(t *testing.T) {
	t.Run("Handle Error from Walker Function", func(t *testing.T) {
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

		// 保存当前 os.Stdout，以便后面恢复
		originalStdout := os.Stdout

		// 创建一个新的 *os.File 来捕获 stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// 在测试结束时恢复原始的 os.Stdout
		defer func() {
			os.Stdout = originalStdout
		}()

		// 使用通道读取捕获的 stdout 内容
		capturedOutput := make(chan string)
		go func() {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			capturedOutput <- buf.String()
		}()

		// 创建一个函数来模拟 walker 函数返回错误
		errorOccurred := false
		kv.Walk(func(rel arch.Relation) error {
			if !errorOccurred {
				errorOccurred = true
				return fmt.Errorf("an error occurred")
			}
			return nil
		})

		// 关闭管道，确保捕获完成
		w.Close()

		// 获取捕获的 stdout 内容
		captured := <-capturedOutput

		// 验证是否打印了错误信息并继续遍历
		if !errorOccurred {
			t.Errorf("Expected the walker function to return an error and continue, but it didn't")
		}
		if !strings.HasPrefix(captured, "relations Walk error:  an error occurred") {
			t.Errorf("Expected 'relations Walk error:  an error occurred', but got '%s'", captured)
		}
	})
}
