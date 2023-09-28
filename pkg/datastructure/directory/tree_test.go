package directory

import (
	"golang.org/x/exp/slices"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDirectoryTree(t *testing.T) {
	filePaths := []string{
		"/path/to/file1.txt",
		"/path/to/file2.txt",
		"/path/to/nested/file3.txt",
		"/path/to/nested/file4.txt",
	}

	rootNode := BuildDirectoryTree(filePaths)

	// Validate root node name
	if rootNode.Name != "/path/to" {
		t.Errorf("Root node name doesn't match. Expected: /path/to, Actual: %s", rootNode.Name)
	}

	// Validate the number of children nodes
	if len(rootNode.Children) != 1 {
		t.Errorf("Number of children nodes is incorrect. Expected: 1, Actual: %d", len(rootNode.Children))
	}

	// Validate the name of the child node
	childNode := rootNode.Children["nested"]
	if childNode == nil {
		t.Errorf("Child node 'nested' not found")
	}
}

func TestFindCommonRootDirectory(t *testing.T) {
	filePaths := []string{
		"/path/to/file1.txt",
		"/path/to/file2.txt",
		"/path/to/nested/file3.txt",
		"/path/to/nested/file4.txt",
		"/path/file5.txt",
	}

	rootDirectory := FindCommonRootDirectory(filePaths)

	// Validate the root directory
	expectedRoot := "/path"
	if rootDirectory != expectedRoot {
		t.Errorf("Common root directory doesn't match. Expected: %s, Actual: %s", expectedRoot, rootDirectory)
	}
}

func TestFindCommonRootDirectory_EmptyPaths(t *testing.T) {
	var filePaths []string
	rootDirectory := FindCommonRootDirectory(filePaths)

	// Validate the root directory
	expectedRoot := ""
	if rootDirectory != expectedRoot {
		t.Errorf("Common root directory doesn't match. Expected: %s, Actual: %s", expectedRoot, rootDirectory)
	}
}

func TestAddPath(t *testing.T) {
	node := &TreeNode{
		Name:     "/path/to",
		Children: make(map[string]*TreeNode),
	}

	// Add path "/path/to/nested"
	node.AddPath("nested")

	// Validate the number of children nodes
	if len(node.Children) != 1 {
		t.Errorf("Number of children nodes is incorrect. Expected: 1, Actual: %d", len(node.Children))
	}

	// Validate the name of the child node
	childNode := node.Children["nested"]
	if childNode == nil {
		t.Errorf("Child node 'nested' not found")
	}
}

func TestTrimFilePath(t *testing.T) {
	filePath := "/path/to/nested/file.txt"

	rootDirectory := "/path/to"
	trimmedPath := strings.TrimPrefix(filepath.Dir(filePath), rootDirectory)
	trimmedPath = strings.TrimPrefix(trimmedPath, string(filepath.Separator))
	trimmedPath = strings.TrimSuffix(trimmedPath, string(filepath.Separator))

	// Validate the trimmed path
	expectedPath := "nested"
	if trimmedPath != expectedPath {
		t.Errorf("Trimmed path doesn't match. Expected: %s, Actual: %s", expectedPath, trimmedPath)
	}
}

func TestTreeNode_AddValue(t *testing.T) {
	// Create a root TreeNode
	root := &TreeNode{
		Name:     "github.com/dddplayer/markdown",
		Children: make(map[string]*TreeNode),
		Value:    nil,
	}

	// Add a child TreeNode
	child := &TreeNode{
		Name:     "child",
		Children: make(map[string]*TreeNode),
		Value:    nil,
	}
	root.Children["child"] = child

	err := root.AddValue("nonexistent", "value")
	if err == nil {
		t.Error("Expected an error for adding a value to a non-existent node, but got no error")
	} else if err.Error() != "node not found: nonexistent" {
		t.Errorf("Expected error message 'node not found: nonexistent', but got: %v", err)
	}

	err = root.AddValue("child", "child_value")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if child.Value != "child_value" {
		t.Errorf("Expected child node value to be 'child_value', but got: %v", child.Value)
	}
}

func TestTreeNode_AddValue_Recursive(t *testing.T) {
	// 创建一个根 TreeNode
	root := &TreeNode{
		Name:     "root",
		Children: make(map[string]*TreeNode),
		Value:    nil,
	}

	// 添加一些子节点
	child1 := &TreeNode{
		Name:     "child1",
		Children: make(map[string]*TreeNode),
		Value:    nil,
	}
	root.Children["child1"] = child1

	child2 := &TreeNode{
		Name:     "child2",
		Children: make(map[string]*TreeNode),
		Value:    nil,
	}
	child1.Children["child2"] = child2

	child3 := &TreeNode{
		Name:     "child3",
		Children: make(map[string]*TreeNode),
		Value:    nil,
	}
	child2.Children["child3"] = child3

	// 在路径中包含多个路径段，进行递归调用
	err := root.AddValue("child1/child2/child3", "recursive_value")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 验证递归路径下的节点值是否正确设置
	if child1.Children["child2"].Children["child3"].Value != "recursive_value" {
		t.Errorf("Expected 'recursive_value', but got: %v", child1.Children["child2"].Children["child3"].Value)
	}
}

func TestGetValue(t *testing.T) {
	// 创建一个根 TreeNode
	root := &TreeNode{Name: "root", Children: make(map[string]*TreeNode), Value: nil}

	// 添加子节点
	root.Children["folder1"] = &TreeNode{Name: "folder1", Children: make(map[string]*TreeNode), Value: nil}
	root.Children["folder2"] = &TreeNode{Name: "folder2", Children: make(map[string]*TreeNode), Value: "folder2_value"}

	// 正确的路径，成功获取值
	value, err := root.GetValue("folder2")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if value != "folder2_value" {
		t.Errorf("Expected 'folder2_value', but got %v", value)
	}

	// 路径中包含不存在的节点，返回错误
	_, err = root.GetValue("folder3/file2.txt")
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	} else if err.Error() != "node not found: folder3/file2.txt" {
		t.Errorf("Expected error message 'node not found: folder3/file2.txt', but got: %v", err)
	}

	// 路径为空，返回错误
	_, err = root.GetValue("")
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	} else if err.Error() != "node not found: " {
		t.Errorf("Expected error message 'node not found: ', but got: %v", err)
	}

	// 路径为单个节点名称，与根节点的子节点匹配，成功获取值
	value, err = root.GetValue("folder1")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if value != nil {
		t.Errorf("Expected nil value, but got %v", value)
	}
}

func TestGetNode(t *testing.T) {
	// 创建一个根 TreeNode
	root := &TreeNode{Name: "root", Children: make(map[string]*TreeNode), Value: nil}

	// 添加子节点
	root.Children["folder1"] = &TreeNode{Name: "folder1", Children: make(map[string]*TreeNode), Value: nil}
	root.Children["folder2"] = &TreeNode{Name: "folder2", Children: make(map[string]*TreeNode), Value: "folder2_value"}

	// 正确的路径，成功获取节点
	node := root.GetNode("folder2")
	if node == nil {
		t.Errorf("Expected a non-nil node, but got nil")
	}
	if node.Value != "folder2_value" {
		t.Errorf("Expected 'folder2_value', but got %v", node.Value)
	}

	// 路径中包含不存在的节点，返回 nil
	node = root.GetNode("folder3/file2.txt")
	if node != nil {
		t.Errorf("Expected nil node, but got a non-nil node")
	}

	// 路径为空，返回 nil
	node = root.GetNode("")
	if node != nil {
		t.Errorf("Expected nil node, but got a non-nil node")
	}

	// 路径为单个节点名称，与根节点的子节点匹配，成功获取节点
	node = root.GetNode("folder1")
	if node == nil {
		t.Errorf("Expected a non-nil node, but got nil")
	}
	if node.Value != nil {
		t.Errorf("Expected nil value, but got %v", node.Value)
	}
}

func TestWalk(t *testing.T) {
	// 创建一个简单的目录树
	root := &TreeNode{Name: "root", Children: make(map[string]*TreeNode), Value: nil}
	root.Children["folder1"] = &TreeNode{Name: "folder1", Children: make(map[string]*TreeNode), Value: nil}
	root.Children["folder2"] = &TreeNode{Name: "folder2", Children: make(map[string]*TreeNode), Value: "folder2_value"}
	root.Children["folder3"] = &TreeNode{Name: "folder3", Children: make(map[string]*TreeNode), Value: nil}
	root.Children["folder1"].Children["subfolder"] = &TreeNode{Name: "subfolder", Children: make(map[string]*TreeNode), Value: "subfolder_value"}

	// 创建一个用于记录遍历结果的切片
	var traversalResult []string

	// 定义回调函数，记录遍历结果
	walkFunc := func(dir string, value any) error {
		traversalResult = append(traversalResult, dir)
		return nil
	}

	// 调用 Walk 函数进行遍历
	Walk(root, walkFunc)

	// 验证遍历结果是否符合预期
	expectedResults := []string{
		"root",
		"root/folder1",
		"root/folder1/subfolder",
		"root/folder2",
		"root/folder3",
	}

	if len(traversalResult) != len(expectedResults) {
		t.Errorf("Expected %d traversal results, but got %d", len(expectedResults), len(traversalResult))
		return
	}

	for i, expected := range expectedResults {
		if !slices.Contains(traversalResult, expected) {
			t.Errorf("Result %d: Expected '%s', not in '%s'", i, expected, traversalResult)
		}
	}
}

func TestWalkRecursive(t *testing.T) {
	// 创建一个简单的目录树
	root := &TreeNode{Name: "root", Children: make(map[string]*TreeNode), Value: nil}
	root.Children["folder1"] = &TreeNode{Name: "folder1", Children: make(map[string]*TreeNode), Value: nil}
	root.Children["folder2"] = &TreeNode{Name: "folder2", Children: make(map[string]*TreeNode), Value: "folder2_value"}
	root.Children["folder3"] = &TreeNode{Name: "folder3", Children: make(map[string]*TreeNode), Value: nil}
	root.Children["folder1"].Children["subfolder"] = &TreeNode{Name: "subfolder", Children: make(map[string]*TreeNode), Value: "subfolder_value"}

	// 创建一个用于记录遍历结果的切片
	var traversalResult []string

	// 定义回调函数，记录遍历结果
	walkFunc := func(dir string, value any) error {
		traversalResult = append(traversalResult, dir)
		return nil
	}

	// 调用 walkRecursive 进行遍历
	if err := root.walkRecursive("", walkFunc); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 验证遍历结果是否符合预期
	expectedResults := []string{
		"root",
		"root/folder1",
		"root/folder1/subfolder",
		"root/folder2",
		"root/folder3",
	}

	if len(traversalResult) != len(expectedResults) {
		t.Errorf("Expected %d traversal results, but got %d", len(expectedResults), len(traversalResult))
		return
	}

	for i, expected := range expectedResults {
		if !slices.Contains(traversalResult, expected) {
			t.Errorf("Result %d: Expected '%s', not in '%s'", i, expected, traversalResult)
		}
	}
}
