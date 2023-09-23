package valueobject

import (
	"testing"
)

func TestNewStringObj(t *testing.T) {
	v := "testValue"
	obj := NewStringObj(v)

	if obj.id.name != v {
		t.Errorf("Expected id.name to be %s, but got %s", v, obj.id.name)
	}

	if obj.id.pkg != "" {
		t.Errorf("Expected id.pkg to be empty, but got %s", obj.id.pkg)
	}
}
