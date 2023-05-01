package entity

import (
	"github.com/dddplayer/core/valueobject"
	"testing"
)

func TestIsEntity(t *testing.T) {
	testCases := []struct {
		path     string
		isEntity bool
	}{
		{"/path/to/entity.go", true},
		{"/path/to/other.go", false},
		{"entity", true},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if IsEntity(tc.path) != tc.isEntity {
				t.Errorf("Expected IsEntity(%q) == %v, but got %v", tc.path, tc.isEntity, !tc.isEntity)
			}
		})
	}
}

func TestIsValueObject(t *testing.T) {
	testCases := []struct {
		path          string
		isValueObject bool
	}{
		{"/path/to/valueobject.go", true},
		{"/path/to/other.go", false},
		{"valueobject", true},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if IsValueObject(tc.path) != tc.isValueObject {
				t.Errorf("Expected IsValueObject(%q) == %v, but got %v", tc.path, tc.isValueObject, !tc.isValueObject)
			}
		})
	}
}

func TestIsFactory(t *testing.T) {
	testCases := []struct {
		path      string
		isFactory bool
	}{
		{"/path/to/factory.go", true},
		{"/path/to/other.go", false},
		{"factory", true},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if IsFactory(tc.path) != tc.isFactory {
				t.Errorf("Expected IsFactory(%q) == %v, but got %v", tc.path, tc.isFactory, !tc.isFactory)
			}
		})
	}
}

func TestIsService(t *testing.T) {
	testCases := []struct {
		path      string
		isService bool
	}{
		{"/path/to/service.go", true},
		{"/path/to/other.go", false},
		{"service", true},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if IsService(tc.path) != tc.isService {
				t.Errorf("Expected IsService(%q) == %v, but got %v", tc.path, tc.isService, !tc.isService)
			}
		})
	}
}

func TestNewEntity(t *testing.T) {
	id := &valueobject.Identifier{Name: "entity1", Path: "domain/entities/entity1"}
	pos := &valueobject.Position{Filename: "entity1.go", Offset: 10, Line: 5, Column: 15}

	entity := NewEntity(*id, *pos)

	if entity.id.Path != id.Path && entity.id.Name != id.Name {
		t.Errorf("Expected id to be %v, but got %v", id, entity.keyObj.Class.obj.id)
	}

	if entity.pos.Filename != pos.Filename &&
		entity.pos.Offset != pos.Offset && entity.pos.Line != pos.Line &&
		entity.pos.Column != pos.Column {
		t.Errorf("Expected pos to be %v, but got %v", pos, entity.keyObj.Class.obj.pos)
	}

	if len(entity.keyObj.Class.attrs) != 0 {
		t.Errorf("Expected length of attrs to be 0, but got %v", len(entity.keyObj.Class.attrs))
	}

	if len(entity.keyObj.Class.methods) != 0 {
		t.Errorf("Expected length of methods to be 0, but got %v", len(entity.keyObj.Class.methods))
	}
}

func TestNewValueObject(t *testing.T) {
	id := &valueobject.Identifier{Name: "value1", Path: "domain/valueobjects/value1"}
	pos := &valueobject.Position{Filename: "value1.go", Offset: 10, Line: 5, Column: 15}

	value := NewValueObject(*id, *pos)

	if value.id.Path != id.Path && value.id.Name != id.Name {
		t.Errorf("Expected id to be %v, but got %v", id, value.keyObj.Class.obj.id)
	}

	if value.pos.Filename != pos.Filename &&
		value.pos.Offset != pos.Offset && value.pos.Line != pos.Line &&
		value.pos.Column != pos.Column {
		t.Errorf("Expected pos to be %v, but got %v", pos, value.keyObj.Class.obj.pos)
	}

	if len(value.keyObj.Class.attrs) != 0 {
		t.Errorf("Expected length of attrs to be 0, but got %v", len(value.keyObj.Class.attrs))
	}

	if len(value.keyObj.Class.methods) != 0 {
		t.Errorf("Expected length of methods to be 0, but got %v", len(value.keyObj.Class.methods))
	}
}

func TestNewService(t *testing.T) {
	id := &valueobject.Identifier{Name: "service1", Path: "domain/services/service1"}
	pos := &valueobject.Position{Filename: "service1.go", Offset: 10, Line: 5, Column: 15}

	service := NewService(*id, *pos)

	if service.Function.id.Path != id.Path && service.Function.id.Name != id.Name {
		t.Errorf("Expected id to be %v, but got %v", id, service.Function.id)
	}

	if service.Function.pos.Filename != pos.Filename &&
		service.Function.pos.Offset != pos.Offset && service.Function.pos.Line != pos.Line &&
		service.Function.pos.Column != pos.Column {
		t.Errorf("Expected pos to be %v, but got %v", pos, service.Function.pos)
	}
}

func TestNewFactory(t *testing.T) {
	id := &valueobject.Identifier{Name: "factory1", Path: "domain/factories/factory1"}
	pos := &valueobject.Position{Filename: "factory1.go", Offset: 10, Line: 5, Column: 15}

	factory := NewFactory(*id, *pos)

	if factory.Function.id.Path != id.Path && factory.Function.id.Name != id.Name {
		t.Errorf("Expected id to be %v, but got %v", id, factory.Function.id)
	}

	if factory.Function.pos.Filename != pos.Filename &&
		factory.Function.pos.Offset != pos.Offset && factory.Function.pos.Line != pos.Line &&
		factory.Function.pos.Column != pos.Column {
		t.Errorf("Expected pos to be %v, but got %v", pos, factory.Function.pos)
	}
}

func TestNewInterface(t *testing.T) {
	id := &valueobject.Identifier{Name: "interface1", Path: "domain/interfaces/interface1"}
	pos := &valueobject.Position{Filename: "interface1.go", Offset: 10, Line: 5, Column: 15}

	iface := NewInterface(*id, *pos)

	if iface.id.Path != id.Path && iface.id.Name != id.Name {
		t.Errorf("Expected id to be %v, but got %v", id, iface.id)
	}

	if iface.pos.Filename != pos.Filename &&
		iface.pos.Offset != pos.Offset && iface.pos.Line != pos.Line &&
		iface.pos.Column != pos.Column {
		t.Errorf("Expected pos to be %v, but got %v", pos, iface.pos)
	}

	if len(iface.Methods) != 0 {
		t.Errorf("Expected length of methods to be 0, but got %v", len(iface.Methods))
	}
}

func TestNewInterfaceMethod(t *testing.T) {
	id := &valueobject.Identifier{Name: "method1", Path: "domain/interfaces/interface1/method1"}
	pos := &valueobject.Position{Filename: "interface1.go", Offset: 20, Line: 8, Column: 25}

	method := NewInterfaceMethod(*id, *pos)

	if method.id.Path != id.Path && method.id.Name != id.Name {
		t.Errorf("Expected id to be %v, but got %v", id, method.id)
	}

	if method.pos.Filename != pos.Filename &&
		method.pos.Offset != pos.Offset && method.pos.Line != pos.Line &&
		method.pos.Column != pos.Column {
		t.Errorf("Expected pos to be %v, but got %v", pos, method.pos)
	}
}

func TestNewFunction(t *testing.T) {
	id := &valueobject.Identifier{Name: "function1", Path: "domain/functions/function1"}
	pos := &valueobject.Position{Filename: "function1.go", Offset: 30, Line: 11, Column: 35}

	function := NewFunction(*id, *pos)

	if function.id.Path != id.Path && function.id.Name != id.Name {
		t.Errorf("Expected id to be %v, but got %v", id, function.id)
	}

	if function.pos.Filename != pos.Filename &&
		function.pos.Offset != pos.Offset && function.pos.Line != pos.Line &&
		function.pos.Column != pos.Column {
		t.Errorf("Expected pos to be %v, but got %v", pos, function.pos)
	}
}
