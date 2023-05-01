package valueobject

import (
	"testing"
)

func TestIdentifier_String(t *testing.T) {
	id := Identifier{
		Name: "testName",
		Path: "/path/to/identifier",
	}
	expected := "/path/to/identifier.testName"
	if id.String() != expected {
		t.Errorf("Expected %s, but got %s", expected, id.String())
	}
}

func TestIdentifier_DomainName(t *testing.T) {
	id := Identifier{
		Name: "testName",
		Path: "/path/to/identifier",
	}
	expected := "to"
	if id.DomainName() != expected {
		t.Errorf("Expected %s, but got %s", expected, id.DomainName())
	}
}

func TestIdentifier_Base(t *testing.T) {
	id := Identifier{
		Name: "testName",
		Path: "/path/to/identifier",
	}
	expected := "identifier.testName"
	if id.Base() != expected {
		t.Errorf("Expected %s, but got %s", expected, id.Base())
	}
}

func TestNewIdentifier(t *testing.T) {
	// create a mock valueobject.Identifier
	mockID := &MockIdentifier{
		mockName: "testName",
		mockPath: "/test/path",
	}
	// create an Identifier instance using NewIdentifier
	identifier := NewIdentifier(mockID)

	// check if the identifier has the correct name and path values
	if identifier.Name != mockID.Name() {
		t.Errorf("Expected identifier name %q, but got %q", mockID.Name(), identifier.Name)
	}
	if identifier.Path != mockID.Path() {
		t.Errorf("Expected identifier path %q, but got %q", mockID.Path(), identifier.Path)
	}
}

type MockIdentifier struct {
	mockPath string
	mockName string
}

func (mi *MockIdentifier) Path() string {
	return mi.mockPath
}

func (mi *MockIdentifier) Name() string {
	return mi.mockName
}

func (mi *MockIdentifier) String() string {
	return mi.mockPath + mi.mockName
}
