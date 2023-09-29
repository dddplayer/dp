package application

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/dot/valueobject"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGeneralGraph(t *testing.T) {
	// create a temporary directory to hold the test package
	tempDir, err := ioutil.TempDir(".", "testpkg")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}(tempDir)

	// create a test package in the temporary directory
	if err := createGeneralTestPackage(tempDir); err != nil {
		t.Fatalf("failed to create test package: %v", err)
	}

	mockRelRepo := &MockRelationRepository{
		relations: make([]arch.Relation, 0),
	}
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}

	result, err := GeneralGraph(tempDir, path.Join(reflect.TypeOf(MockObjectRepository{}).PkgPath(), path.Base(tempDir)), mockRepo, mockRelRepo)

	if err != nil {
		t.Errorf("GeneralGraph() returned unexpected error:\nActual: %v", err)
	}

	// Verify the output matches the expected DOT directed
	if strings.Contains(result, "Func1") == false ||
		strings.Contains(result, "Func2") == false ||
		strings.Contains(result, "Func3") == false ||
		strings.Contains(result, "main") == false {
		t.Errorf("GeneralGraph() returned unexpected output:\nActual: %v", result)
	}
}

// createTestPackage creates a test package in the specified directory
func createGeneralTestPackage(dir string) error {
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

func createHexagonTestPackage(dir string) error {
	// Create the root directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create subdirectories
	subDirs := []string{
		string(arch.HexagonDirectoryCmd),
		string(arch.HexagonDirectoryPkg),
		string(arch.HexagonDirectoryInternal),
		string(arch.HexagonDirectoryInternal) + "/" + string(arch.HexagonDirectoryDomain),
		string(arch.HexagonDirectoryInternal) + "/" + string(arch.HexagonDirectoryDomain) + "/test",
		string(arch.HexagonDirectoryInternal) + "/" + string(arch.HexagonDirectoryDomain) + "/test/" + string(arch.HexagonDirectoryEntity),
		string(arch.HexagonDirectoryInternal) + "/" + string(arch.HexagonDirectoryDomain) + "/test/" + string(arch.HexagonDirectoryValueObject),
	}

	for _, subDir := range subDirs {
		subDirPath := filepath.Join(dir, subDir)
		if err := os.MkdirAll(subDirPath, 0755); err != nil {
			return err
		}

		switch subDir {
		case string(arch.HexagonDirectoryCmd), string(arch.HexagonDirectoryPkg):
			// Create test files in the package directory
			for i := 1; i <= 3; i++ {
				fileName := fmt.Sprintf("file%d.go", i)
				filePath := filepath.Join(subDirPath, fileName)
				if err := ioutil.WriteFile(filePath, []byte(fmt.Sprintf("package %s\n\nfunc Func%d() {}\n", subDir, i)), 0644); err != nil {
					return err
				}
			}
		}

		if strings.HasSuffix(subDir, string(arch.HexagonDirectoryEntity)) {
			initFile := filepath.Join(subDirPath, "class.go")
			if err := ioutil.WriteFile(initFile, []byte("package entity\n\ntype Test struct {}\n"), 0644); err != nil {
				return err
			}
		}
		if strings.HasSuffix(subDir, string(arch.HexagonDirectoryValueObject)) {
			initFile := filepath.Join(subDirPath, "class.go")
			if err := ioutil.WriteFile(initFile, []byte("package valueobject\n\ntype VO struct {}\n"), 0644); err != nil {
				return err
			}
		}
	}

	// Create the package initialization file
	initFile := filepath.Join(dir, "main.go")
	fileStr := strings.ReplaceAll(main, "module", path.Join(reflect.TypeOf(MockObjectRepository{}).PkgPath(), path.Base(dir)))
	if err := ioutil.WriteFile(initFile, []byte(fileStr), 0644); err != nil {
		return err
	}

	return nil
}

func TestTacticGraph(t *testing.T) {
	// create a temporary directory to hold the test package
	tempDir, err := ioutil.TempDir(".", "testpkg")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}(tempDir)

	// create a test package in the temporary directory
	if err := createHexagonTestPackage(tempDir); err != nil {
		t.Fatalf("failed to create test package: %v", err)
	}

	mockRelRepo := &MockRelationRepository{
		relations: make([]arch.Relation, 0),
	}
	mockRepo := &MockObjectRepository{
		objects: make(map[string]arch.Object),
		idents:  []arch.ObjIdentifier{},
	}

	result, err := TacticGraph(tempDir, path.Join(reflect.TypeOf(MockObjectRepository{}).PkgPath(), path.Base(tempDir)), mockRepo, mockRelRepo)

	// Verify the output matches the expected DOT directed
	if strings.Contains(result, valueobject.GenerateShortURL("test_entity")) == false ||
		strings.Contains(result, valueobject.GenerateShortURL("test_valueobject")) == false {
		t.Errorf("TacticGraph() returned unexpected output:\nActual: %v", result)
	}
}
