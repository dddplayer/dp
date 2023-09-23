package directory

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// TreeNode 表示目录树的节点
type TreeNode struct {
	Name     string               // 节点名称
	Children map[string]*TreeNode // 子节点
	Value    interface{}          // 存储任意类型的值
}

// AddPath 将目录路径添加到目录树中
func (node *TreeNode) AddPath(path string) {
	components := strings.Split(path, string(filepath.Separator))

	currNode := node
	for _, comp := range components {
		if _, exists := currNode.Children[comp]; !exists {
			currNode.Children[comp] = &TreeNode{
				Name:     comp,
				Children: make(map[string]*TreeNode),
			}
		}
		currNode = currNode.Children[comp]
	}
}

// AddValue 添加值到目标节点
func (node *TreeNode) AddValue(childPath string, value interface{}) error {
	pathSegments := strings.Split(childPath, string(filepath.Separator))
	childName := pathSegments[0]

	child, ok := node.Children[childName]
	if !ok {
		return fmt.Errorf("node not found: %s", childPath)
	}

	if len(pathSegments) > 1 {
		return child.AddValue(strings.Join(pathSegments[1:], "/"), value)
	}

	child.Value = value
	return nil
}

// GetValue 通过路径获取值
func (node *TreeNode) GetValue(path string) (interface{}, error) {
	pathSegments := strings.Split(path, string(filepath.Separator))
	childName := pathSegments[0]

	child, ok := node.Children[childName]
	if !ok {
		// 如果子节点不存在，则返回错误信息
		return nil, fmt.Errorf("node not found: %s", path)
	}

	if len(pathSegments) > 1 {
		// 递归调用GetValue方法，将剩余路径传递给子节点
		return child.GetValue(strings.Join(pathSegments[1:], "/"))
	}

	// 当路径已经遍历完时，返回当前节点的值
	return child.Value, nil
}

func (node *TreeNode) GetNode(path string) *TreeNode {
	pathSegments := strings.Split(path, string(filepath.Separator))
	childName := pathSegments[0]

	child, ok := node.Children[childName]
	if !ok {
		return nil
	}

	if len(pathSegments) > 1 {
		return child.GetNode(strings.Join(pathSegments[1:], "/"))
	}

	return child
}

// FindCommonRootDirectory 查找所有路径的共同根目录
func FindCommonRootDirectory(filePaths []string) string {
	if len(filePaths) == 0 {
		return ""
	}

	// 使用第一个文件路径作为起始值
	rootPath := filepath.Dir(filePaths[0])

	// 对比所有路径，找出共同的根目录
	for _, filePath := range filePaths {
		for !strings.HasPrefix(filePath, rootPath) {
			rootPath = filepath.Dir(rootPath)
		}
	}

	return filepath.Clean(rootPath)
}

// BuildDirectoryTree 构建目录树
func BuildDirectoryTree(filePaths []string) *TreeNode {
	rootDirectory := FindCommonRootDirectory(filePaths)

	rootNode := &TreeNode{
		Name:     rootDirectory,
		Children: make(map[string]*TreeNode),
	}

	for _, filePath := range filePaths {
		dir := strings.TrimPrefix(filepath.Dir(filePath), rootDirectory)
		dir = strings.TrimPrefix(dir, string(filepath.Separator))
		dir = strings.TrimSuffix(dir, string(filepath.Separator))
		if dir != "" {
			rootNode.AddPath(dir)
		}
	}

	return rootNode
}

// WalkFunc 是遍历回调函数的类型
type WalkFunc func(dir string, value any) error

// WalkFrom root
// WalFrom Node

// Walk 递归遍历目录树
func Walk(node *TreeNode, cb WalkFunc) {
	if err := node.walkRecursive("", cb); err != nil {
		fmt.Println(err.Error())
	}
}

// 辅助函数，实现递归遍历
func (node *TreeNode) walkRecursive(currentDir string, cb WalkFunc) error {
	dir := path.Join(currentDir, node.Name)
	if err := cb(dir, node.Value); err != nil {
		return err
	}
	for _, child := range node.Children {
		if err := child.walkRecursive(dir, cb); err != nil {
			return err
		}
	}
	return nil
}
