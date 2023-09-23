package valueobject

import (
	"testing"
)

func TestNewPosition(t *testing.T) {
	filename := "test.go"
	offset := 100
	line := 5
	column := 10

	vp := &MockPosition{
		filename: filename,
		offset:   offset,
		line:     line,
		column:   column,
	}

	p := newPosition(vp)

	if p.filename != filename {
		t.Errorf("Expected filename '%s', but got '%s'", filename, p.filename)
	}

	if p.offset != offset {
		t.Errorf("Expected offset '%d', but got '%d'", offset, p.offset)
	}

	if p.line != line {
		t.Errorf("Expected line '%d', but got '%d'", line, p.line)
	}

	if p.column != column {
		t.Errorf("Expected column '%d', but got '%d'", column, p.column)
	}
}

func TestPosMethods(t *testing.T) {
	testCases := []struct {
		name             string
		position         *pos
		expectedValid    bool
		expectedLine     int
		expectedColumn   int
		expectedOffset   int
		expectedFilename string
	}{
		{
			name: "Valid Position",
			position: &pos{
				filename: "file1.txt",
				offset:   100,
				line:     5,
				column:   10,
			},
			expectedValid:    true,
			expectedLine:     5,
			expectedColumn:   10,
			expectedOffset:   100,
			expectedFilename: "file1.txt",
		},
		{
			name: "Invalid Position",
			position: &pos{
				filename: "file2.txt",
				offset:   200,
				line:     -1,
				column:   -1,
			},
			expectedValid:    false,
			expectedLine:     -1,
			expectedColumn:   -1,
			expectedOffset:   200,
			expectedFilename: "file2.txt",
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualValid := tc.position.Valid()
			actualLine := tc.position.Line()
			actualColumn := tc.position.Column()
			actualOffset := tc.position.Offset()
			actualFilename := tc.position.Filename()

			if actualValid != tc.expectedValid || actualLine != tc.expectedLine ||
				actualColumn != tc.expectedColumn || actualOffset != tc.expectedOffset ||
				actualFilename != tc.expectedFilename {
				t.Errorf("For test case %s:\nExpected: (%v, %d, %d, %d, %s)\nGot: (%v, %d, %d, %d, %s)",
					tc.name, tc.expectedValid, tc.expectedLine, tc.expectedColumn, tc.expectedOffset,
					tc.expectedFilename, actualValid, actualLine, actualColumn, actualOffset, actualFilename)
			}
		})
	}
}

func TestEmptyPosition(t *testing.T) {
	emptyPos := emptyPosition()

	if emptyPos.filename != "" || emptyPos.offset != 0 ||
		emptyPos.line != -1 || emptyPos.column != 0 {
		t.Errorf("Empty position is not as expected:\nGot: %v", emptyPos)
	}
}
