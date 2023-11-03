package application

import (
	"errors"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/dot/valueobject"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

func TestStrategicGraph(t *testing.T) {
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

	result, err := StrategicGraph(tempDir,
		path.Join(reflect.TypeOf(MockObjectRepository{}).PkgPath(), path.Base(tempDir)),
		false,
		mockRepo, mockRelRepo)

	if err != nil {
		t.Errorf("StrategicGraph() returned unexpected error:\nActual: %v", err)
	}

	// Verify the output matches the expected DOT directed
	if strings.Contains(result, valueobject.GenerateShortURL("test_entity_Test")) == false ||
		strings.Contains(result, valueobject.GenerateShortURL("test_valueobject_VO")) == false {
		t.Errorf("StrategicGraph() returned unexpected output:\nActual: %v", result)
	}
}

func TestStrategicGraph_ArchFactoryError(t *testing.T) {
	_, err := StrategicGraph("", "", false, nil, nil)

	if err == nil || err.Error() != "objRepo cannot be nil" {
		t.Errorf("Expected error 'objRepo cannot be nil', but got: %v", err)
	}
}

func TestStrategicGraph_EntityNewCodeError(t *testing.T) {
	// 创建一个模拟的对象仓库和关系仓库
	mockObjRepo := &MockObjectRepository{}
	mockRelRepo := &MockRelationRepository{}

	// 模拟 entity.NewCode 函数返回错误
	expectedError := errors.New("packages contain errors")

	_, err := StrategicGraph("non-exist", "dummy", false, mockObjRepo, mockRelRepo)

	// 验证返回的错误是否符合预期
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}
