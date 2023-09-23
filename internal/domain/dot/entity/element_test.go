package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

func TestElementsIndicators(t *testing.T) {

	// Define some test attributes
	attr1 := &DummyDotAttribute{name: "attr1", port: "port1", color: "color1"}
	attr2 := &DummyDotAttribute{name: "attr2", port: "port2", color: "color2"}
	attr3 := &DummyDotAttribute{name: "attr3", port: "port3", color: "color3"}
	attr4 := &DummyDotAttribute{name: "attr4", port: "port4", color: "color4"}

	// Define a test object
	obj1 := &DummyDotElement{NameVal: "object", PortVal: "port", ColorVal: "color",
		AttributesVal: []arch.Nodes{[]arch.Node{attr2}, []arch.Node{attr1, attr2}}}
	obj2 := &DummyDotElement{NameVal: "object", PortVal: "port", ColorVal: "color",
		AttributesVal: []arch.Nodes{[]arch.Node{attr1, attr2, attr3, attr4}, []arch.Node{attr1, attr2, attr3, attr4}}}

	testCases := []struct {
		name                              string
		elements                          []arch.Element
		expectedMaxLeft, expectedMaxRight int
	}{
		{
			name:             "Single Object with Attributes",
			elements:         []arch.Element{obj1},
			expectedMaxLeft:  1, // Since the first child has 1 attribute and the second child has 2 attributes.
			expectedMaxRight: 2, // The second child has 2 attributes.
		},
		{
			name:             "Multiple Objects",
			elements:         []arch.Element{obj1, obj2}, // Two identical objects.
			expectedMaxLeft:  3,                          // The first and second objects have the same attributes counts.
			expectedMaxRight: 3,                          // The second child has 2 attributes.
		},
		{
			name:             "Empty Elements",
			elements:         []arch.Element{}, // Empty input.
			expectedMaxLeft:  1,
			expectedMaxRight: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			maxLeft, maxRight := ElementsIndicators(tc.elements)
			if maxLeft != tc.expectedMaxLeft || maxRight != tc.expectedMaxRight {
				t.Errorf("For test case %s: expected (%d, %d), got (%d, %d)",
					tc.name, tc.expectedMaxLeft, tc.expectedMaxRight, maxLeft, maxRight)
			}
		})
	}
}
