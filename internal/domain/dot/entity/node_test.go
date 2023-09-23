package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

func TestIsLeftRightStructure(t *testing.T) {
	// Define some test attributes
	attr1 := &DummyDotAttribute{name: "attr1", port: "port1", color: "color1"}
	attr2 := &DummyDotAttribute{name: "attr2", port: "port2", color: "color2"}

	// Define some test objects
	obj1 := &DummyDotElement{NameVal: "obj1", PortVal: "port1", ColorVal: "color1", AttributesVal: []arch.Nodes{[]arch.Node{attr1}, []arch.Node{attr2}}}
	obj2 := &DummyDotElement{NameVal: "obj2", PortVal: "port2", ColorVal: "color2", AttributesVal: []arch.Nodes{[]arch.Node{attr1}}}
	obj3 := &DummyDotElement{NameVal: "obj3", PortVal: "port3", ColorVal: "color3", AttributesVal: []arch.Nodes{[]arch.Node{attr1}}}

	// Define the expected results
	expectedResults := []bool{true, false, false}

	// Test the function on each object and compare the result with the expected value
	for i, obj := range []*DummyDotElement{obj1, obj2, obj3} {
		result := isLeftRightStructure(obj)
		expected := expectedResults[i]
		if result != expected {
			t.Errorf("isLeftRightStructure(%v) = %v; expected %v", obj, result, expected)
		}
	}
}

func TestNode_Build(t *testing.T) {
	// Valid elements
	els1 := []arch.Element{
		&DummyDotElement{
			NameVal: "Test1",
			AttributesVal: []arch.Nodes{
				[]arch.Node{DummyDotAttribute{name: "attr1", port: "p1", color: "red"}},
			},
		},
		&DummyDotElement{
			NameVal: "Test2",
			AttributesVal: []arch.Nodes{
				[]arch.Node{
					DummyDotAttribute{name: "attr2", port: "p2", color: "blue"},
					DummyDotAttribute{name: "attr3", port: "p3", color: "green"},
				},
			},
		},
	}

	n1 := &Node{Name: "TestNode1", Table: &Table{Rows: []*Row{}}}
	err := n1.Build(els1)
	if err != nil {
		t.Errorf("Unexpected error with valid elements: %v", err)
	}

	// Valid elements with left-right structure
	els2 := []arch.Element{
		&DummyDotElement{
			NameVal: "Test4",
			AttributesVal: []arch.Nodes{
				[]arch.Node{
					DummyDotAttribute{name: "attr5", port: "p5", color: "blue"},
				},
			},
		},
		&DummyDotElement{
			NameVal: "Test5",
			AttributesVal: []arch.Nodes{
				[]arch.Node{
					DummyDotAttribute{name: "attr6", port: "p6", color: "green"},
				},
			},
		},
	}

	n3 := &Node{Name: "TestNode3", Table: &Table{Rows: []*Row{}}}
	err = n3.Build(els2)
	if err != nil {
		t.Errorf("Unexpected error with valid elements with left-right structure: %v", err)
	}
}
