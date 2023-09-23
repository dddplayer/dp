package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/pkg/datastructure/directory"
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
