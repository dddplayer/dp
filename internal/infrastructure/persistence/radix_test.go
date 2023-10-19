package persistence

import (
	"errors"
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

func TestRadixTree_Find(t *testing.T) {
	t.Run("Object Found", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := NewRadixTree()

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
		r := NewRadixTree()

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
		r := NewRadixTree()

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
		r := NewRadixTree()

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

func TestRadixTree_All(t *testing.T) {
	t.Run("Return All Identifiers", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := &RadixTree{}

		// 创建一些 MockIdentifier 并插入到 RadixTree 中
		mockID1 := &MockIdentifier{IDVal: "id1"}
		mockID2 := &MockIdentifier{IDVal: "id2"}
		mockID3 := &MockIdentifier{IDVal: "id3"}

		r.objIds = map[string]arch.ObjIdentifier{
			mockID1.ID(): mockID1,
			mockID2.ID(): mockID2,
			mockID3.ID(): mockID3,
		}

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

func TestRadixTree_Walk(t *testing.T) {
	t.Run("Walk with Objects", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := NewRadixTree()

		// 创建一些 MockObject 和 MockIdentifier 并插入到 RadixTree 中
		mockID1 := &MockIdentifier{IDVal: "id1"}
		mockObject1 := &MockObject{id: mockID1}
		r.Tree.Insert(mockID1.ID(), mockObject1)

		mockID2 := &MockIdentifier{IDVal: "id2"}
		mockObject2 := &MockObject{id: mockID2}
		r.Tree.Insert(mockID2.ID(), mockObject2)

		// 创建回调函数，用于记录回调的次数
		var callbackCount int

		// 调用 Walk 函数
		r.Walk(func(obj arch.Object) error {
			callbackCount++
			return nil
		})

		// 验证回调函数是否按预期执行了两次（每个对象一次）
		if callbackCount != 2 {
			t.Errorf("Expected callback to be called 2 times, but got %d times", callbackCount)
		}
	})

	t.Run("Walk with Error", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := NewRadixTree()

		// 创建一个 MockObject 并插入到 RadixTree 中
		mockID := &MockIdentifier{IDVal: "id1"}
		mockObject := &MockObject{id: mockID}
		r.Tree.Insert(mockID.ID(), mockObject)

		// 创建回调函数，用于返回错误
		callbackError := errors.New("callback error")

		// 调用 Walk 函数
		r.Walk(func(obj arch.Object) error {
			return callbackError
		})

	})

	t.Run("Walk with No Objects", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := NewRadixTree()

		// 创建回调函数，用于记录回调的次数
		var callbackCount int

		// 调用 Walk 函数
		r.Walk(func(obj arch.Object) error {
			callbackCount++
			return nil
		})

		// 验证回调函数是否没有被执行（因为没有对象）
		if callbackCount != 0 {
			t.Errorf("Expected callback to be called 0 times, but got %d times", callbackCount)
		}
	})
}

func TestRadixTree_GetObjects(t *testing.T) {
	t.Run("Get Objects with Valid IDs", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := NewRadixTree()

		// 创建一些 MockObject 和 MockIdentifier 并插入到 RadixTree 中
		mockID1 := &MockIdentifier{IDVal: "id1"}
		mockObject1 := &MockObject{id: mockID1}
		r.Tree.Insert(mockID1.ID(), mockObject1)

		mockID2 := &MockIdentifier{IDVal: "id2"}
		mockObject2 := &MockObject{id: mockID2}
		r.Tree.Insert(mockID2.ID(), mockObject2)

		// 获取对象的标识符列表
		ids := []arch.ObjIdentifier{mockID1, mockID2}

		// 调用 GetObjects 函数
		objs, err := r.GetObjects(ids)

		// 验证是否没有返回错误
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		// 验证返回的对象切片是否包含了两个对象
		if len(objs) != 2 {
			t.Errorf("Expected 2 objects, but got %d", len(objs))
		}

		// 验证返回的对象切片中是否包含了预期的对象
		if objs[0] != mockObject1 {
			t.Errorf("Expected the first object to be the same as mockObject1")
		}
		if objs[1] != mockObject2 {
			t.Errorf("Expected the second object to be the same as mockObject2")
		}
	})

	t.Run("Get Objects with Invalid IDs", func(t *testing.T) {
		// 创建 RadixTree 实例
		r := NewRadixTree()

		// 创建一些 MockObject 和 MockIdentifier 并插入到 RadixTree 中
		mockID1 := &MockIdentifier{IDVal: "id1"}
		mockObject1 := &MockObject{id: mockID1}
		r.Tree.Insert(mockID1.ID(), mockObject1)

		mockID2 := &MockIdentifier{IDVal: "id2"}

		// 获取对象的标识符列表，其中一个标识符无效
		ids := []arch.ObjIdentifier{mockID1, mockID2}

		// 调用 GetObjects 函数
		objs, err := r.GetObjects(ids)

		// 验证是否返回了错误
		if err == nil {
			t.Error("Expected an error, but got no error")
		} else {
			// 验证错误消息是否包含了无效的标识符
			expectedErrorMsg := "object id2 not found"
			if err.Error() != expectedErrorMsg {
				t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
			}

			// 验证返回的对象切片是否为空
			if objs != nil {
				t.Errorf("Expected nil object slice, but got %v", objs)
			}
		}
	})
}
