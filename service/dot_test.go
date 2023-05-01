package service

import (
	"fmt"
	"github.com/dddplayer/core/entity"
	"github.com/dddplayer/core/valueobject"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDot(t *testing.T) {
	// create a temporary directory to hold the test package
	tempDir, err := ioutil.TempDir(".", "testpkg")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// create a test package in the temporary directory
	if err := createTestPackage(tempDir); err != nil {
		t.Fatalf("failed to create test package: %v", err)
	}

	// Create a mock repository with test data
	repo := &mockRepository{
		data: make(map[valueobject.Identifier]entity.DomainObject),
	}

	// Call the Dot function with a test domain and the mock repository
	result := Dot(tempDir, path.Join(reflect.TypeOf(mockRepository{}).PkgPath(), path.Base(tempDir)), repo)

	base := filepath.Base(tempDir)
	baseF1 := fmt.Sprintf("%s_%s", base, "Func1")
	baseF2 := fmt.Sprintf("%s_%s", base, "Func2")
	baseF3 := fmt.Sprintf("%s_%s", base, "Func3")
	baseMain := fmt.Sprintf("%s_%s", base, "main")
	edge := fmt.Sprintf("%s:%s -> %s:%s", base, baseMain, base, baseF1)
	// Verify the output matches the expected DOT graph
	if strings.Contains(result, baseF1) == false ||
		strings.Contains(result, baseF2) == false ||
		strings.Contains(result, baseF3) == false ||
		strings.Contains(result, baseMain) == false ||
		strings.Contains(result, edge) == false {
		t.Errorf("Dot() returned unexpected output:\nActual: %v", result)
	}
}

// createTestPackage creates a test package in the specified directory
func createTestPackage(dir string) error {
	// create some test files in the package directory
	for i := 1; i <= 3; i++ {
		fileName := fmt.Sprintf("file%d.go", i)
		filePath := filepath.Join(dir, fileName)
		if err := ioutil.WriteFile(filePath, []byte(fmt.Sprintf("package main\n\nfunc Func%d() {}\n", i)), 0644); err != nil {
			return err
		}
	}

	// create the package initialization file
	initFile := filepath.Join(dir, "main.go")
	if err := ioutil.WriteFile(initFile, []byte("package main\n\nfunc main(){Func1()}\n"), 0644); err != nil {
		return err
	}

	return nil
}

type mockRepository struct {
	data map[valueobject.Identifier]entity.DomainObject
}

func (r *mockRepository) Find(id valueobject.Identifier) entity.DomainObject {
	return r.data[id]
}

func (r *mockRepository) Insert(obj entity.DomainObject) error {
	r.data[*obj.Identifier()] = obj
	return nil
}

func (r *mockRepository) Walk(cb func(obj entity.DomainObject) error) {
	for _, obj := range r.data {
		if err := cb(obj); err != nil {
			return
		}
	}
}
