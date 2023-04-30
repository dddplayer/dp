package entity

import "testing"

func TestIsLeftRightStructure(t *testing.T) {
	// Define some test attributes
	attr1 := &DummyDotAttribute{name: "attr1", port: "port1", color: "color1"}
	attr2 := &DummyDotAttribute{name: "attr2", port: "port2", color: "color2"}

	// Define some test objects
	obj1 := &DummyDotElement{NameVal: "obj1", PortVal: "port1", ColorVal: "color1", AttributesVal: []interface{}{[]DotAttribute{attr1}, []DotAttribute{attr2}}}
	obj2 := &DummyDotElement{NameVal: "obj2", PortVal: "port2", ColorVal: "color2", AttributesVal: []interface{}{[]DotAttribute{attr1}}}
	obj3 := &DummyDotElement{NameVal: "obj3", PortVal: "port3", ColorVal: "color3", AttributesVal: []interface{}{attr1}}

	// Define the expected results
	expectedResults := []bool{true, true, false}

	// Test the function on each object and compare the result with the expected value
	for i, obj := range []*DummyDotElement{obj1, obj2, obj3} {
		result := isLeftRightStructure(obj)
		expected := expectedResults[i]
		if result != expected {
			t.Errorf("isLeftRightStructure(%v) = %v; expected %v", obj, result, expected)
		}
	}
}

func TestValidateElements(t *testing.T) {
	dummyAttr := &DummyDotAttribute{name: "dummy", port: "port", color: "color"}
	dummyAttrs := []DotAttribute{dummyAttr, dummyAttr}
	dummyAttrsSlice := []interface{}{dummyAttrs}

	tests := []struct {
		name  string
		els   DotElements
		valid bool
	}{
		{
			name: "Valid Elements",
			els: DotElements{
				&DummyDotElement{NameVal: "dummy", PortVal: "port", ColorVal: "color", AttributesVal: []interface{}{dummyAttr}},
				&DummyDotElement{NameVal: "dummy", PortVal: "port", ColorVal: "color", AttributesVal: dummyAttrsSlice},
			},
			valid: true,
		},
		{
			name: "Invalid Elements",
			els: DotElements{
				&DummyDotElement{NameVal: "dummy", PortVal: "port", ColorVal: "color", AttributesVal: []interface{}{1, "dummy"}},
				&DummyDotElement{NameVal: "dummy", PortVal: "port", ColorVal: "color", AttributesVal: []interface{}{dummyAttr, 2}},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateElements(tt.els); got != tt.valid {
				t.Errorf("ValidateElements() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestNode_Build(t *testing.T) {
	// Valid elements
	els1 := DotElements{
		&DummyDotElement{
			NameVal: "Test1",
			AttributesVal: []interface{}{
				DummyDotAttribute{name: "attr1", port: "p1", color: "red"},
			},
		},
		&DummyDotElement{
			NameVal: "Test2",
			AttributesVal: []interface{}{
				[]DotAttribute{
					DummyDotAttribute{name: "attr2", port: "p2", color: "blue"},
					DummyDotAttribute{name: "attr3", port: "p3", color: "green"},
				},
			},
		},
	}

	n1 := &Node{Name: "TestNode1"}
	err := n1.Build(els1)
	if err != nil {
		t.Errorf("Unexpected error with valid elements: %v", err)
	}

	// Invalid elements
	els2 := DotElements{
		&DummyDotElement{
			NameVal: "Test3",
			AttributesVal: []interface{}{
				DummyDotAttribute{name: "attr4", port: "p4", color: "red"},
				123,
			},
		},
	}

	n2 := &Node{Name: "TestNode2"}
	err = n2.Build(els2)
	if err == nil {
		t.Error("Expected an error with invalid elements, but got none")
	}

	// Valid elements with left-right structure
	els3 := DotElements{
		&DummyDotElement{
			NameVal: "Test4",
			AttributesVal: []interface{}{
				[]DotAttribute{
					DummyDotAttribute{name: "attr5", port: "p5", color: "blue"},
				},
			},
		},
		&DummyDotElement{
			NameVal: "Test5",
			AttributesVal: []interface{}{
				[]DotAttribute{
					DummyDotAttribute{name: "attr6", port: "p6", color: "green"},
				},
			},
		},
	}

	n3 := &Node{Name: "TestNode3"}
	err = n3.Build(els3)
	if err != nil {
		t.Errorf("Unexpected error with valid elements with left-right structure: %v", err)
	}

	// Invalid elements with left-right structure
	els4 := DotElements{
		&DummyDotElement{
			NameVal: "Test6",
			AttributesVal: []interface{}{
				DummyDotAttribute{name: "attr7", port: "p7", color: "red"},
			},
		},
		&DummyDotElement{
			NameVal: "Test7",
			AttributesVal: []interface{}{
				DummyDotAttribute{name: "attr8", port: "p8", color: "green"},
			},
		},
		&DummyDotElement{
			NameVal: "Test8",
			AttributesVal: []interface{}{
				123,
			},
		},
	}

	n4 := &Node{Name: "TestNode4"}
	err = n4.Build(els4)
	if err == nil {
		t.Error("Expected an error with invalid elements with left-right structure, but got none")
	}
}
