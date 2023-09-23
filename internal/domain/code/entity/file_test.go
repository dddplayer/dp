package entity

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_findFile(t *testing.T) {
	testDir, err := tmpTestDir()
	defer os.RemoveAll(testDir) // 删除临时目录
	if err != nil {
		t.Fatalf("failed to create tempory test dir: %v", err)
	}

	pkg, err := packages.Load(&packages.Config{
		Mode:       packages.LoadAllSyntax,
		Tests:      false,
		Dir:        "",
		BuildFlags: build.Default.BuildTags,
	}, testDir)
	if err != nil {
		t.Fatalf("failed to load package: %v", err)
	}
	if packages.PrintErrors(pkg) > 0 {
		t.Fatalf("failed to load package: %v", err)
	}

	// Test case 1: file can be found
	var expr ast.Expr
	found := false
	packages.Visit(pkg, nil, func(pkg *packages.Package) {
		for _, f := range pkg.Syntax {
			if found {
				break
			}
			for _, decl := range f.Decls {
				if found {
					break
				}
				switch decl.(type) {
				case *ast.FuncDecl:
					f := decl.(*ast.FuncDecl)
					if f.Name.Name == "hello" {
						expr = decl.(*ast.FuncDecl).Type
						found = true
						break
					}
				}
			}
		}
	})

	file := findFile(pkg[0], expr)
	if file == nil {
		t.Fatalf("expected file to be found")
	}
	if file.Name.Name != "main" {
		t.Fatalf("unexpected file found: %s", file.Name.Name)
	}

	// Test case 2: file cannot be found
	expr, err = parser.ParseExpr(`123`)
	if err != nil {
		t.Fatalf("failed to parse expr: %v", err)
	}
	file = findFile(pkg[0], expr)
	if file != nil {
		t.Fatalf("expected file not to be found")
	}
}

func tmpTestDir() (string, error) {
	// 创建临时目录
	tmpDir, err := ioutil.TempDir(".", "test")
	if err != nil {
		fmt.Printf("Failed to create temporary directory: %v", err)
		return "", err
	}

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "main.go")
	testContent := []byte(`
package main

import "fmt"

func hello() string {
    return "hello world"
}

func main() {
    fmt.Println(hello())
}
`)
	if err := ioutil.WriteFile(testFile, testContent, 0644); err != nil {
		fmt.Printf("Failed to create test file: %v", err)
		return "", err
	}

	return tmpDir, nil
}
