package entity

import (
	"reflect"
	"testing"
)

func TestCountRow(t *testing.T) {
	testCases := []struct {
		num         int
		rowCapacity int
		expected    int
	}{
		{0, 5, 0},
		{1, 5, 1},
		{5, 5, 1},
		{6, 5, 2},
		{10, 5, 2},
		{11, 5, 3},
	}

	for _, tc := range testCases {
		actual := countRow(tc.num, tc.rowCapacity)
		if actual != tc.expected {
			t.Errorf("countRow(%d, %d): expected %d, but got %d", tc.num, tc.rowCapacity, tc.expected, actual)
		}
	}
}

func TestLineRow(t *testing.T) {
	fs := []DotAttribute{
		&DummyDotAttribute{name: "Attribute1", port: "port1", color: "red"},
		&DummyDotAttribute{name: "Attribute2", port: "port2", color: "green"},
		&DummyDotAttribute{name: "Attribute3", port: "port3", color: "blue"},
	}
	col := 20

	r := lineRow(fs, col)

	// Assert that the row contains the correct number of data items
	expectedNumData := len(fs)*2 - 1 + rowStartEmptyBlockNum + rowEndEmptyBlockNum
	if len(r.Data) != expectedNumData {
		t.Errorf("Unexpected number of data items. Expected %d, but got %d", expectedNumData, len(r.Data))
	}

	// Assert that each data item has the correct column span
	expectedColSpan := (col - len(fs) + 1 - rowStartEmptyBlockNum - rowEndEmptyBlockNum) / len(fs)
	for _, d := range r.Data {
		if d.Text != "" && d.Port != "" && d.ColSpan != expectedColSpan {
			t.Errorf("Unexpected column span for data item. Expected %d, but got %d", expectedColSpan, d.ColSpan)
		}
	}

	// Assert that each data item has the correct name, port, and color
	for i, f := range fs {
		expectedName := f.Name()
		expectedPort := f.Port()
		expectedColor := f.Color()

		d := r.Data[i*2+rowStartEmptyBlockNum]
		name := d.Text
		port := d.Port
		color := d.BgColor

		if name != expectedName {
			t.Errorf("Unexpected name for data item. Expected %s, but got %s", expectedName, name)
		}

		if port != expectedPort {
			t.Errorf("Unexpected port for data item. Expected %s, but got %s", expectedPort, port)
		}

		if color != expectedColor {
			t.Errorf("Unexpected color for data item. Expected %s, but got %s", expectedColor, color)
		}
	}
}

func TestRowFirstPadding(t *testing.T) {
	src := []DotAttribute{
		&DummyDotAttribute{"A", "1", "red"},
		&DummyDotAttribute{"B", "2", "green"},
		&DummyDotAttribute{"C", "3", "blue"},
		&DummyDotAttribute{"D", "4", "yellow"},
		&DummyDotAttribute{"E", "5", "orange"},
	}
	row, column := 2, 3
	expected := [][]DotAttribute{
		{
			&DummyDotAttribute{"A", "1", "red"},
			&DummyDotAttribute{"B", "2", "green"},
			&DummyDotAttribute{"C", "3", "blue"},
		},
		{
			&DummyDotAttribute{"D", "4", "yellow"},
			&DummyDotAttribute{"E", "5", "orange"},
			nil,
		},
	}
	result := rowFirstPadding(src, row, column)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("rowFirstPadding(%v, %d, %d) = %v, expected %v", src, row, column, result, expected)
	}
}

func TestColumnFirstPadding(t *testing.T) {
	src := []DotAttribute{
		&DummyDotAttribute{"A", "1", "red"},
		&DummyDotAttribute{"B", "2", "green"},
		&DummyDotAttribute{"C", "3", "blue"},
		&DummyDotAttribute{"D", "4", "yellow"},
		&DummyDotAttribute{"E", "5", "orange"},
	}
	row, column := 3, 3
	expected := [][]DotAttribute{
		{
			&DummyDotAttribute{"A", "1", "red"},
			&DummyDotAttribute{"B", "2", "green"},
			nil,
		},
		{
			&DummyDotAttribute{"C", "3", "blue"},
			&DummyDotAttribute{"D", "4", "yellow"},
			nil,
		},
		{
			&DummyDotAttribute{"E", "5", "orange"},
			nil,
			nil,
		},
	}
	result := columnFirstPadding(src, row, column)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("columnFirstPadding(%v, %d, %d) = %v, expected %v", src, row, column, result, expected)
	}
}

func TestBuildNodeRow(t *testing.T) {
	// Define some test attributes
	attr1 := &DummyDotAttribute{name: "attr1", port: "port1", color: "color1"}
	attr2 := &DummyDotAttribute{name: "attr2", port: "port2", color: "color2"}

	// Define a test object
	obj := &DummyDotElement{NameVal: "object", PortVal: "port", ColorVal: "color", AttributesVal: []any{attr1}}

	// Define some test input values
	ms := []DotAttribute{attr2}
	as := []DotAttribute{attr1, attr2}
	isFirstRow := true
	right := 2
	row := 3

	// Call the function being tested
	rowResult := buildNodeRow(ms, as, true, obj, right, row)

	objColumn := column("object", "port", "color") // objec
	objColumn.RowSpan = row

	// Define the expected row result
	expectedData := []*Data{
		blankColumn(),                      // ms[1]
		column("attr2", "port2", "color2"), // attr2
		blankColumn(),                      // emptyBlockBeforeDomainObjNum
		objColumn,
		column("attr1", "port1", "color1"), // attr1
		column("attr2", "port2", "color2"), // attr2// t
		blankColumn(),                      // rowEndEmptyBlockNum
	}

	expectedRow := &Row{Data: expectedData}

	// Compare the result with the expected value
	if !rowsEqual(rowResult, expectedRow) {
		t.Errorf("buildNodeRow(%v, %v, %v, %v, %v) = %v; expected %v", ms, as, obj, right, row, rowResult, expectedRow)
	}

	// Test with isFirstRow as false
	isFirstRow = false
	rowResult = buildNodeRow(ms, as, isFirstRow, obj, right, row)

	expectedData = []*Data{
		blankColumn(),                      // ms[1]
		column("attr2", "port2", "color2"), // attr2
		blankColumn(),                      // emptyBlockBeforeDomainObjNum
		column("attr1", "port1", "color1"), // attr1
		column("attr2", "port2", "color2"), // attr2
		blankColumn(),                      // rowEndEmptyBlockNum
	}

	expectedRow = &Row{Data: expectedData}

	if !rowsEqual(rowResult, expectedRow) {
		t.Errorf("buildNodeRow(%v, %v, %v, %v, %v, %v) = %v; expected %v", ms, as, isFirstRow, obj, right, row, rowResult, expectedRow)
	}
}

func TestLeftRightRow(t *testing.T) {
	// Define some test attributes
	attr1 := &DummyDotAttribute{name: "attr1", port: "port1", color: "color1"}
	attr2 := &DummyDotAttribute{name: "attr2", port: "port2", color: "color2"}

	// Define a test object
	obj := &DummyDotElement{NameVal: "object", PortVal: "port", ColorVal: "color",
		AttributesVal: []any{[]DotAttribute{attr2}, []DotAttribute{attr1, attr2}}}

	// Define some test input values
	maxLeft := 2
	maxRight := 2
	e := obj

	n := &Node{}
	leftRightRow(n, e, maxLeft, maxRight)

	// Define the expected result
	expectedData := []*Data{
		blankColumn(), // attr3
		blankColumn(),
		column("attr2", "port2", "color2"), // attr2
		blankColumn(),                      // emptyBlockBeforeDomainObjNum
		column("object", "port", "color"),  // object
		column("attr1", "port1", "color1"), // attr1
		column("attr2", "port2", "color2"), // attr2
		blankColumn(),                      // rowEndEmptyBlockNum
	}

	expectedRow := &Row{Data: expectedData}
	expectedRows := []*Row{expectedRow, blankRow()}
	expectedNode := &Node{Rows: expectedRows}

	// Compare the result with the expected value
	if !nodesEqual(n, expectedNode) {
		t.Errorf("leftRightRow(%v, %v, %v, %v) = %v; expected %v", n, e, maxLeft, maxRight, n, expectedNode)
	}

	// Test with 1 array attribute
	// Define a test object
	obj = &DummyDotElement{NameVal: "object", PortVal: "port", ColorVal: "color",
		AttributesVal: []any{[]DotAttribute{attr1, attr2}}}

	n = &Node{}
	e = obj
	leftRightRow(n, e, maxLeft, maxRight)

	// Define the expected result
	expectedData = []*Data{
		blankColumn(),                      // attr3
		column("attr2", "port2", "color2"), // attr2
		column("attr1", "port1", "color1"), // attr1
		blankColumn(),                      // emptyBlockBeforeDomainObjNum
		column("object", "port", "color"),  // object
		blankColumn(),
		blankColumn(),
		blankColumn(), // rowEndEmptyBlockNum
	}

	expectedRow = &Row{Data: expectedData}
	expectedRows = []*Row{expectedRow, blankRow()}
	expectedNode = &Node{Rows: expectedRows}

	// Compare the result with the expected value
	if !nodesEqual(n, expectedNode) {
		t.Errorf("leftRightRow(%v, %v, %v, %v) = %v; expected %v", n, e, maxLeft, maxRight, n, expectedNode)
	}
}

// nodeEqual is a helper function that compares two Node structs for equality.
func nodesEqual(a, b *Node) bool {
	if len(a.Rows) != len(b.Rows) {
		return false
	}
	for i := range a.Rows {
		if !rowsEqual(a.Rows[i], b.Rows[i]) {
			return false
		}
	}
	return true
}

func rowsEqual(r1, r2 *Row) bool {
	if len(r1.Data) != len(r2.Data) {
		return false
	}

	for i := 0; i < len(r1.Data); i++ {
		if !dataEqual(r1.Data[i], r2.Data[i]) {
			return false
		}
	}

	return true
}

func dataEqual(d1, d2 *Data) bool {
	return d1.Text == d2.Text &&
		d1.Port == d2.Port &&
		d1.BgColor == d2.BgColor &&
		d1.RowSpan == d2.RowSpan &&
		d1.ColSpan == d2.ColSpan
}

type DummyDotElement struct {
	NameVal       string
	PortVal       string
	ColorVal      string
	AttributesVal []any
}

func (e *DummyDotElement) Name() string {
	return e.NameVal
}

func (e *DummyDotElement) Port() string {
	return e.PortVal
}

func (e *DummyDotElement) Color() string {
	return e.ColorVal
}

func (e *DummyDotElement) Attributes() []interface{} {
	return e.AttributesVal
}

type DummyDotAttribute struct {
	name  string
	port  string
	color string
}

func (a DummyDotAttribute) Name() string {
	return a.name
}

func (a DummyDotAttribute) Port() string {
	return a.port
}

func (a DummyDotAttribute) Color() string {
	return a.color
}
