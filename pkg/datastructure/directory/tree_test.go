package directory

import (
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
	}

	rootDirectory := FindCommonRootDirectory(filePaths)

	// Validate the root directory
	expectedRoot := "/path/to"
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
