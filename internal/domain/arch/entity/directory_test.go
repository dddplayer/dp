package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/pkg/datastructure/directory"
	"reflect"
	"testing"
)

func TestNewDirectory(t *testing.T) {
	// Define your mock directory tree structure
	mockPaths := []string{
		"/path/to/file1.txt",
		"/path/to/file2.txt",
		"/path/to/nested/file3.txt",
		"/path/to/nested/file4.txt",
	}

	// Call the function to be tested
	dir := NewDirectory(mockPaths)

	// Verify that the root node is empty
	if dir.root == nil {
		t.Errorf("Expected root node to be empty, but it's not")
	}
}

// MockTreeNode 模拟的 TreeNode 结构体
type MockTreeNode struct {
	NameVal     string
	ChildrenVal map[string]*MockTreeNode
	ValueVal    interface{}
}

func TestIsHexagon(t *testing.T) {
	// Create a mock instance of Directory
	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "root",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): &directory.TreeNode{},
				string(arch.HexagonDirectoryPkg): &directory.TreeNode{},
				string(arch.HexagonDirectoryInternal): &directory.TreeNode{
					Children: map[string]*directory.TreeNode{
						string(arch.HexagonDirectoryDomain): &directory.TreeNode{
							Children: map[string]*directory.TreeNode{
								string(arch.HexagonDirectoryEntity):      nil, // ValueObject missing
								string(arch.HexagonDirectoryValueObject): nil, // Entity missing
							},
						},
					},
				},
			},
		},
	}

	// Call the function to be tested
	isHexagon := mockDirectory.isHexagon()

	// Verify that isHexagon is true
	if !isHexagon {
		t.Errorf("Expected isHexagon to be true, but it's not")
	}
}

func TestIsHexagonWithInValidStructure(t *testing.T) {
	// Create a mock instance of Directory
	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "root",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): &directory.TreeNode{},
				string(arch.HexagonDirectoryPkg): &directory.TreeNode{},
				string(arch.HexagonDirectoryInternal): &directory.TreeNode{
					Children: map[string]*directory.TreeNode{
						string(arch.HexagonDirectoryDomain): &directory.TreeNode{
							Children: map[string]*directory.TreeNode{
								string(arch.HexagonDirectoryEntity):      &directory.TreeNode{}, // Entity present
								string(arch.HexagonDirectoryValueObject): &directory.TreeNode{}, // ValueObject present
							},
						},
					},
				},
			},
		},
	}

	// Call the function to be tested
	isHexagon := mockDirectory.isHexagon()

	// Verify that isHexagon is false
	if isHexagon {
		t.Errorf("Expected isHexagon to be false, but it's false")
	}
}

func TestArchDesignPatternWithHexagon(t *testing.T) {
	// Create a mock instance of Directory with isHexagon returning true
	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "root",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): &directory.TreeNode{},
				string(arch.HexagonDirectoryPkg): &directory.TreeNode{},
				string(arch.HexagonDirectoryInternal): &directory.TreeNode{
					Children: map[string]*directory.TreeNode{
						string(arch.HexagonDirectoryDomain): &directory.TreeNode{
							Children: map[string]*directory.TreeNode{
								string(arch.HexagonDirectoryEntity):      nil, // ValueObject missing
								string(arch.HexagonDirectoryValueObject): nil, // Entity missing
							},
						},
					},
				},
			},
		},
	}

	// Call the function to be tested
	designPattern := mockDirectory.ArchDesignPattern()

	// Verify that designPattern is DesignPatternHexagon
	if designPattern != arch.DesignPatternHexagon {
		t.Errorf("Expected designPattern to be DesignPatternHexagon, but it's not")
	}
}

func TestArchDesignPatternWithPlain(t *testing.T) {
	// Create a mock instance of Directory with isHexagon returning false
	mockDirectory := &Directory{
		root: &directory.TreeNode{
			Name: "root",
			Children: map[string]*directory.TreeNode{
				string(arch.HexagonDirectoryCmd): &directory.TreeNode{},
				string(arch.HexagonDirectoryPkg): &directory.TreeNode{},
				string(arch.HexagonDirectoryInternal): &directory.TreeNode{
					Children: map[string]*directory.TreeNode{
						string(arch.HexagonDirectoryDomain): &directory.TreeNode{
							Children: map[string]*directory.TreeNode{
								string(arch.HexagonDirectoryEntity):      &directory.TreeNode{}, // Entity present
								string(arch.HexagonDirectoryValueObject): &directory.TreeNode{}, // ValueObject present
							},
						},
					},
				},
			},
		},
	}

	// Call the function to be tested
	designPattern := mockDirectory.ArchDesignPattern()

	// Verify that designPattern is DesignPatternPlain
	if designPattern != arch.DesignPatternPlain {
		t.Errorf("Expected designPattern to be DesignPatternPlain, but it's not")
	}
}

func TestHexagonDirectory(t *testing.T) {
	// 创建一个模拟的 Directory 实例
	dir := &Directory{}

	// 测试各种情况
	testCases := []struct {
		inputDir        string
		expectedHexagon arch.HexagonDirectory
	}{
		{"domain", arch.HexagonDirectoryDomain},
		{"entity", arch.HexagonDirectoryEntity},
		{"valueobject", arch.HexagonDirectoryValueObject},
		{"repository", arch.HexagonDirectoryRepository},
		{"factory", arch.HexagonDirectoryFactory},
		{"domain/aggregate", arch.HexagonDirectoryAggregate},
		{"invalid", arch.HexagonDirectoryInvalid},
	}

	for _, testCase := range testCases {
		result := dir.HexagonDirectory(testCase.inputDir)

		if result != testCase.expectedHexagon {
			t.Errorf("Expected %v for input directory %s, got %v", testCase.expectedHexagon, testCase.inputDir, result)
		}
	}
}

func TestGetTargetDir(t *testing.T) {
	// 创建一个模拟的 Directory 实例
	dir := &Directory{
		root: &directory.TreeNode{
			Name: "/root",
		},
	}

	// 测试有效目录
	validDir := "/root/target"
	expectedTargetDir := "target"
	result, err := dir.getTargetDir(validDir)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result != expectedTargetDir {
		t.Errorf("Expected target directory %s, got %s", expectedTargetDir, result)
	}

	// 测试无效目录
	invalidDir := "/invalid/target"
	_, err = dir.getTargetDir(invalidDir)

	if err == nil {
		t.Errorf("Expected an error for invalid directory, got no error")
	}
}

func TestAddObjs(t *testing.T) {
	// 创建一个模拟的 Directory 实例
	dir := &Directory{
		root: &directory.TreeNode{
			Name: "/root",
			Children: map[string]*directory.TreeNode{
				"target": &directory.TreeNode{
					Name: "target",
				},
			},
		},
	}

	// 测试根目录
	rootDir := "/root"
	rootObjs := []arch.ObjIdentifier{newMockObject(1)}
	err := dir.AddObjs(rootDir, rootObjs)

	if err != nil {
		t.Errorf("Expected no error for root directory, got: %v", err)
	}

	// 测试有效目录
	validDir := "/root/target"
	validObjs := []arch.ObjIdentifier{newMockObject(2)}
	err = dir.AddObjs(validDir, validObjs)

	if err != nil {
		t.Errorf("Expected no error for valid directory, got: %v", err)
	}

	// 测试无效目录
	invalidDir := "/invalid/target"
	invalidObjs := []arch.ObjIdentifier{newMockObject(3)}
	err = dir.AddObjs(invalidDir, invalidObjs)

	if err == nil {
		t.Errorf("Expected an error for invalid directory, got no error")
	}
}

func TestGetObjs(t *testing.T) {
	// 创建一个模拟的 Directory 实例
	dir := &Directory{
		root: &directory.TreeNode{
			Name: "/root",
			Children: map[string]*directory.TreeNode{
				"target": &directory.TreeNode{
					Name: "target",
				},
			},
		},
	}

	// 准备模拟的对象标识符切片
	objs := []arch.ObjIdentifier{newMockObject(1)}

	// 添加对象到目录
	targetDir := "/root/target"
	err := dir.AddObjs(targetDir, objs)
	if err != nil {
		t.Fatalf("Failed to add objects to directory: %v", err)
	}

	// 测试获取对象
	resultObjs, err := dir.GetObjs("target")

	if err != nil {
		t.Errorf("Expected no error for getting objects, got: %v", err)
	}

	// 检查获取的对象是否与预期相符
	if !reflect.DeepEqual(resultObjs, objs) {
		t.Errorf("Expected objects %v, got %v", objs, resultObjs)
	}

	// 测试不存在的目录
	invalidDir := "/invalid/target"
	_, err = dir.GetObjs(invalidDir)

	if err == nil {
		t.Errorf("Expected an error for invalid directory, got no error")
	}
}

func TestParentDir(t *testing.T) {
	// 创建一个模拟的 Directory 实例
	dir := &Directory{}

	// 测试目录
	testCases := []struct {
		inputDir          string
		expectedParentDir string
	}{
		{"/root/target", "root"},
		{"/root", "/"},
		{"/", "/"},
		{"/root/dir/file.txt", "dir"},
	}

	for _, testCase := range testCases {
		result := dir.ParentDir(testCase.inputDir)

		if result != testCase.expectedParentDir {
			t.Errorf("Expected parent directory %s for input directory %s, got %s", testCase.expectedParentDir, testCase.inputDir, result)
		}
	}
}
