package valueobject

import "testing"

func TestIdentifierString(t *testing.T) {
	id := &identifier{
		path: "path/to/something",
		name: "Identifier",
	}
	expected := "path/to/something.Identifier"
	if id.String() != expected {
		t.Errorf("Unexpected result: %s, expected: %s", id.String(), expected)
	}
}

func TestIdentifierPath(t *testing.T) {
	id := &identifier{
		path: "path/to/something",
		name: "Identifier",
	}
	if id.Path() != "path/to/something" {
		t.Errorf("Unexpected path: %s", id.Path())
	}
}

func TestIdentifierName(t *testing.T) {
	id := &identifier{
		path: "path/to/something",
		name: "Identifier",
	}
	if id.Name() != "Identifier" {
		t.Errorf("Unexpected name: %s", id.Name())
	}
}
