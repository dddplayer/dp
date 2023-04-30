package entity

import (
	"testing"
)

type MockDotElement struct {
	name string
	port string
	attr []any
}

func (e *MockDotElement) Name() string      { return e.name }
func (e *MockDotElement) Port() string      { return e.port }
func (e *MockDotElement) Color() string     { return "" }
func (e *MockDotElement) Attributes() []any { return e.attr }

type MockDotAttribute struct {
	name  string
	port  string
	color string
}

func (a *MockDotAttribute) Name() string  { return a.name }
func (a *MockDotAttribute) Port() string  { return a.port }
func (a *MockDotAttribute) Color() string { return a.color }

func TestDotElements(t *testing.T) {
	// Create some mock DotElements
	attrs1 := []any{
		[]DotAttribute{
			&MockDotAttribute{name: "label", port: "node1", color: "blue"},
			&MockDotAttribute{name: "label", port: "node1", color: "blue"},
			&MockDotAttribute{name: "label", port: "node1", color: "blue"},
		},
		[]DotAttribute{
			&MockDotAttribute{name: "label", port: "node1", color: "blue"},
			&MockDotAttribute{name: "label", port: "node1", color: "blue"},
		},
	}
	elem1 := &MockDotElement{
		name: "node1",
		port: "1234",
		attr: attrs1,
	}

	attrs2 := []any{
		&MockDotAttribute{name: "label", port: "node1", color: "blue"},
		&MockDotAttribute{name: "label", port: "node1", color: "blue"},
		&MockDotAttribute{name: "label", port: "node1", color: "blue"},
		&MockDotAttribute{name: "label", port: "node1", color: "blue"},
	}
	elem2 := &MockDotElement{
		name: "node2",
		port: "5678",
		attr: attrs2,
	}

	dotElements := DotElements{elem1, elem2}

	// Test the First method
	first := dotElements.First()
	if first == nil || first.Name() != "node1" {
		t.Errorf("First() = %v; expected %v", first, elem1)
	}

	// Test the indicators method
	maxLeft, maxRight := dotElements.indicators()
	expectedMaxLeft := 3
	expectedMaxRight := 2
	if maxLeft != expectedMaxLeft {
		t.Errorf("indicators(): maxLeft = %d; expected %d", maxLeft, expectedMaxLeft)
	}
	if maxRight != expectedMaxRight {
		t.Errorf("indicators(): maxRight = %d; expected %d", maxRight, expectedMaxRight)
	}
}
