package persistence

import (
	"github.com/dddplayer/dp/pkg/datastructure/radix"
	"testing"
)

func TestRadixTree_Find(t *testing.T) {
	t.Run("Object Found", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := &RadixTree{Tree: radix.NewTree()}

		// 创建一个 MockObject 和 MockIdentifier 并插入到 RadixTree 中
		mockID := &MockIdentifier{IDVal: "id1"}
		mockObject := &MockObject{id: mockID}
		r.Tree.Insert(mockID.ID(), mockObject)

		// 查找对象
		foundObject := r.Find(mockID)

		// 验证是否找到对象
		if foundObject == nil {
			t.Errorf("Expected to find the object, but got nil")
		}

		// 验证找到的对象是否正确
		if foundObject != mockObject {
			t.Errorf("Expected the found object to be the same as the inserted object")
		}
	})

	t.Run("Object Not Found", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := &RadixTree{Tree: radix.NewTree()}

		// 创建一个 MockIdentifier，但不插入到 RadixTree 中
		mockID := &MockIdentifier{IDVal: "id2"}

		// 查找对象
		foundObject := r.Find(mockID)

		// 验证是否未找到对象
		if foundObject != nil {
			t.Errorf("Expected not to find the object, but found one")
		}
	})
}

func TestRadixTree_Insert(t *testing.T) {
	t.Run("Insert Successful", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := &RadixTree{Tree: radix.NewTree()}

		// 创建一个 MockObject 和 MockIdentifier
		mockID := &MockIdentifier{IDVal: "id1"}
		mockObject := &MockObject{id: mockID}

		// 插入对象
		err := r.Insert(mockObject)

		// 验证插入是否成功
		if err != nil {
			t.Errorf("Expected insert to be successful, but got error: %v", err)
		}

		// 验证对象是否插入到树中
		foundObject, ok := r.Tree.Get(mockID.ID())
		if !ok {
			t.Errorf("Expected to find the inserted object, but not found")
		}
		if foundObject != mockObject {
			t.Errorf("Expected the found object to be the same as the inserted object")
		}
	})

	t.Run("Insert Duplicate ID", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := &RadixTree{Tree: radix.NewTree()}

		// 创建两个 MockObject，但它们具有相同的 ID
		mockID := &MockIdentifier{IDVal: "id2"}
		mockObject1 := &MockObject{id: mockID}
		mockObject2 := &MockObject{id: mockID}

		// 插入第一个对象
		err1 := r.Insert(mockObject1)

		// 验证第一个插入是否成功
		if err1 != nil {
			t.Errorf("Expected first insert to be successful, but got error: %v", err1)
		}

		// 尝试插入第二个对象，应该失败
		err2 := r.Insert(mockObject2)

		if err2 != nil {
			t.Errorf("Expected second insert to be successful, but got error: %v", err2)
		}

		foundObject, ok := r.Tree.Get(mockID.ID())
		if !ok {
			t.Errorf("Expected to find the inserted object, but not found")
		}
		if foundObject != mockObject2 {
			t.Errorf("Expected the found object to be the same as the inserted object")
		}
	})
}
