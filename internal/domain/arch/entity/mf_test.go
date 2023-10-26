package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"testing"
)

func TestDirFilter_IsValid(t *testing.T) {
	sf := &DirFilter{
		pkgSet: []string{
			"example.com/pkg1",
			"example.com/pkg2",
			"example.com/pkg3",
		},
		// 其他字段的值可以根据需要进行设置
	}

	// 编写测试用例，包括输入和预期输出
	testCases := []struct {
		input    string
		expected bool
	}{
		{"example.com/pkg1", true},  // 预期返回 true
		{"example.com", true},       // 预期返回 true
		{"example.com/pkg4", false}, // 预期返回 false
		{"another.com/pkg1", false}, // 预期返回 false
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := sf.IsValid(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %v for input %s, but got %v", tc.expected, tc.input, actual)
			}
		})
	}
}

func TestDirFilter_IsExist(t *testing.T) {
	obj1 := newMockObject(1)
	obj2 := newMockObject(2)
	obj3 := newMockObject(3)
	obj4 := newMockObject(4)
	obj5 := newMockObject(5)
	sf := &DirFilter{
		objs: []arch.ObjIdentifier{
			obj1.Identifier(),
			obj2.Identifier(),
			obj3.Identifier(),
		},
	}

	testCases := []struct {
		input    arch.ObjIdentifier
		expected bool
	}{
		{obj1.Identifier(), true},
		{obj2.Identifier(), true},
		{obj4.Identifier(), false},
		{obj5.Identifier(), false},
	}

	for _, tc := range testCases {
		t.Run(tc.input.ID(), func(t *testing.T) {
			actual := sf.isExist(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %v for input %s, but got %v", tc.expected, tc.input.ID(), actual)
			}
		})
	}
}

func TestGetReceiverFromSourceData(t *testing.T) {
	receiver2 := newMockObject(2)
	receiver4 := newMockObject(4)
	sourceData := []arch.Object{
		newMockObject(1),
		newMockObject(2),
		newMockObject(3),
	}

	testCase1 := struct {
		receiver   arch.ObjIdentifier
		sourceData []arch.Object
		expected   arch.Object
	}{
		receiver:   receiver2.Identifier(),
		sourceData: sourceData,
		expected:   sourceData[1],
	}

	testCase2 := struct {
		receiver   arch.ObjIdentifier
		sourceData []arch.Object
		expected   arch.Object
	}{
		receiver:   receiver4.Identifier(),
		sourceData: sourceData,
		expected:   nil,
	}

	t.Run("Test Case 1", func(t *testing.T) {
		actual := getReceiverFromSourceData(testCase1.receiver, testCase1.sourceData)
		if actual != testCase1.expected {
			t.Errorf("Expected %v, but got %v", testCase1.expected, actual)
		}
	})

	t.Run("Test Case 2", func(t *testing.T) {
		actual := getReceiverFromSourceData(testCase2.receiver, testCase2.sourceData)
		if actual != testCase2.expected {
			t.Errorf("Expected %v, but got %v", testCase2.expected, actual)
		}
	})
}

func TestDirFilter_FilterObjs(t *testing.T) {
	obj1 := newMockObject(1)
	obj2 := newMockObject(2)
	obj3 := newMockObject(3)

	sf := &DirFilter{
		objs: []arch.ObjIdentifier{
			obj1.Identifier(),
			obj2.Identifier(),
			obj3.Identifier(),
		},
	}

	func1 := valueobject.NewFunction(obj1, obj2)
	func2 := valueobject.NewFunction(obj2, nil)

	// 创建一些模拟的 sourceData
	sourceData := []arch.Object{
		func1, func2,
		valueobject.NewGeneral(obj3),
	}

	// 调用被测试的函数
	filteredObjs := sf.FilterObjs(sourceData)

	// 验证筛选结果是否符合预期
	expectedResult := []arch.Object{
		func1, func2,
	}
	if !equalSlices(filteredObjs, expectedResult) {
		t.Errorf("Filtered objects do not match the expected result.")
	}
}

// 用于比较两切片是否相等的辅助函数
func equalSlices(slice1, slice2 []arch.Object) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, obj := range slice1 {
		if obj != slice2[i] {
			return false
		}
	}
	return true
}

func TestMessageFlow_NewDirFilter(t *testing.T) {
	mf := &MessageFlow{
		relationDigraph: NewRelationDigraph(),
		mainPkgPath:     "/path/to",
		endPkgPath:      "/path/to/sub",
	}

	// 调用被测试的函数
	dirFilter, err := mf.newDirFilter()

	// 验证返回的结果是否符合预期
	if err != nil {
		t.Errorf("Expected no error, but got error: %v", err)
	}

	if dirFilter == nil {
		t.Error("Expected a DirFilter, but got nil")
	}

	// 验证生成的 DirFilter 是否符合预期
	expectedPkgSet := []string{"/path/to", "/path/to/sub"}
	if !equalStringSlices(dirFilter.pkgSet, expectedPkgSet) {
		t.Errorf("Expected pkgSet %v, but got %v", expectedPkgSet, dirFilter.pkgSet)
	}

	// 测试无法找到主函数的情况
	mf.mainPkgPath = "example.com/mainpkg"
	dirFilter, err = mf.newDirFilter()

	if err == nil || err.Error() != "main func not found" {
		t.Errorf("Expected an error with 'main func not found' message, but got error: %v", err)
	}

	if dirFilter != nil {
		t.Error("Expected nil DirFilter, but got a non-nil DirFilter")
	}
}

// 用于比较两字符串切片是否相等的辅助函数
func equalStringSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, str := range slice1 {
		if str != slice2[i] {
			return false
		}
	}
	return true
}

func TestBuildDiagram(t *testing.T) {
	mockDirectory, objs := newMockMfDirectoryWithObjs()
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}
	for _, mockObj := range objs {
		_ = mockRepo.Insert(mockObj)
	}

	mf := &MessageFlow{
		directory:       mockDirectory,
		objRepo:         mockRepo,
		relationDigraph: NewRelationDigraph(),
		mainPkgPath:     "/path/to",
		endPkgPath:      "/path/to/sub",
		modulePath:      "/path/to",
	}

	diagram, err := mf.buildDiagram()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	expectedNodeCount := 7
	if len(diagram.Nodes) != expectedNodeCount {
		t.Errorf("Expected %d nodes, but got %d", expectedNodeCount, len(diagram.Nodes))
	}

	expectedEdgeCount := 5
	if len(diagram.Edges()) != expectedEdgeCount {
		t.Errorf("Expected edges to be %d, but got %d", expectedEdgeCount, len(diagram.Edges()))
	}
}
