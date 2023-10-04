package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/code"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestVisit(t *testing.T) {
	// create a temporary directory to hold the test package
	tempDir, err := ioutil.TempDir(".", "testpkg")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// create a test package in the temporary directory
	if err := createPkgTestPackage(tempDir); err != nil {
		t.Fatalf("failed to create test package: %v", err)
	}

	// define the expected results
	expectedNodeCount := 4
	expectedLinkCount := 1

	// define the callback functions
	ch := &MockCodeHandler{}

	// call the VisitFast function with the test package
	c, err := NewCode(tempDir, "testpkg")
	if err != nil {
		t.Fatalf("NewCode failed with error: %v", err)
	}
	if err := c.VisitFast(ch); err != nil {
		t.Fatalf("VisitFast failed with error: %v", err)
	}

	if len(ch.nodes) != expectedNodeCount {
		t.Errorf("Node callback did not process the expected number of nodes. Got: %d, Expected: %d", len(ch.nodes), expectedNodeCount)
	}

	if len(ch.links) != expectedLinkCount {
		t.Errorf("Link callback did not process the expected number of links. Got: %d, Expected: %d", len(ch.links), expectedLinkCount)
	}
}

func TestVisitDeep(t *testing.T) {
	// create a temporary directory to hold the test package
	tempDir, err := ioutil.TempDir(".", "testpkg")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// create a test package in the temporary directory
	if err := createPkgTestPackage(tempDir); err != nil {
		t.Fatalf("failed to create test package: %v", err)
	}

	// define the expected results
	expectedNodeCount := 4
	expectedLinkCount := 1

	// define the callback functions
	ch := &MockCodeHandler{}

	// call the VisitFast function with the test package
	c, err := NewCode(tempDir, "testpkg")
	if err != nil {
		t.Fatalf("NewCode failed with error: %v", err)
	}
	if err := c.VisitDeep(ch); err != nil {
		t.Fatalf("VisitFast failed with error: %v", err)
	}

	if len(ch.nodes) != expectedNodeCount {
		t.Errorf("Node callback did not process the expected number of nodes. Got: %d, Expected: %d", len(ch.nodes), expectedNodeCount)
	}

	if len(ch.links) != expectedLinkCount {
		t.Errorf("Link callback did not process the expected number of links. Got: %d, Expected: %d", len(ch.links), expectedLinkCount)
	}
}

// createTestPackage creates a test package in the specified directory
func createPkgTestPackage(dir string) error {
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

type MockCodeHandler struct {
	nodes []*code.Node
	links []*code.Link
}

func (ch *MockCodeHandler) NodeHandler(node *code.Node) {
	ch.nodes = append(ch.nodes, node)
}

func (ch *MockCodeHandler) LinkHandler(link *code.Link) {
	ch.links = append(ch.links, link)
}

func TestNewGo_Error(t *testing.T) {
	path := "github.com/example/mypackage"
	domain := "example.com"

	_, err := newGo(path, domain)

	if err == nil {
		t.Errorf("Expected an error, got no error")
	}
}

func TestNewCode_Error(t *testing.T) {
	mainPkgPath := "github.com/example/mypackage"
	domain := "example.com"

	_, err := NewCode(mainPkgPath, domain)

	if err == nil {
		t.Errorf("Expected an error, got no error")
	}
}
