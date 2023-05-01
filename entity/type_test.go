package entity

import (
	"testing"

	"github.com/dddplayer/core/valueobject"
)

func TestObj_Identifier(t *testing.T) {
	id := valueobject.NewIdentifier(&MockIdentifier{mockName: "testName", mockPath: "/test/path"})
	pos := valueobject.NewPosition(&MockPosition{filename: "test.go", offset: 1, line: 2, column: 3})

	o := &obj{id: &id, pos: &pos}

	if o.Identifier() != &id {
		t.Errorf("Expected id %v, but got %v", id, o.Identifier())
	}
}

func TestObj_Position(t *testing.T) {
	id := valueobject.NewIdentifier(&MockIdentifier{mockName: "testName", mockPath: "/test/path"})
	pos := valueobject.NewPosition(&MockPosition{filename: "test.go", offset: 1, line: 2, column: 3})

	o := &obj{id: &id, pos: &pos}

	if o.Position() != &pos {
		t.Errorf("Expected position %v, but got %v", pos, o.Position())
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

type MockPosition struct {
	filename string
	offset   int
	line     int
	column   int
}

func (p *MockPosition) Filename() string { return p.filename }
func (p *MockPosition) Offset() int      { return p.offset }
func (p *MockPosition) Line() int        { return p.line }
func (p *MockPosition) Column() int      { return p.column }
