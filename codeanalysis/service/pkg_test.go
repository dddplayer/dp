package service

import (
	"fmt"
	"github.com/dddplayer/core/codeanalysis/entity"
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
	if err := createTestPackage(tempDir); err != nil {
		t.Fatalf("failed to create test package: %v", err)
	}

	// define the expected results
	expectedNodeCount := 4
	expectedLinkCount := 1

	// define the callback functions
	var nodes []*entity.Node
	var links []*entity.Link

	nodeCB := func(node *entity.Node) { nodes = append(nodes, node) }
	linkCB := func(link *entity.Link) { links = append(links, link) }

	// call the Visit function with the test package
	if err := Visit(tempDir, "testpkg", nodeCB, linkCB); err != nil {
		t.Fatalf("Visit failed with error: %v", err)
	}

	if len(nodes) != expectedNodeCount {
		t.Errorf("Node callback did not process the expected number of nodes. Got: %d, Expected: %d", len(nodes), expectedNodeCount)
	}

	if len(links) != expectedLinkCount {
		t.Errorf("Link callback did not process the expected number of links. Got: %d, Expected: %d", len(links), expectedLinkCount)
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
