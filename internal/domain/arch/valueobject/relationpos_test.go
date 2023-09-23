package valueobject

import "testing"

func TestRelationPosMethods(t *testing.T) {
	fromPosition := &MockPosition{
		filename: "from_file.txt",
		offset:   100,
		line:     5,
		column:   10,
	}

	toPosition := &MockPosition{
		filename: "to_file.txt",
		offset:   200,
		line:     8,
		column:   15,
	}

	rPos := &relationPos{
		from: fromPosition,
		to:   toPosition,
	}

	expectedFrom := fromPosition
	expectedTo := toPosition

	actualFrom := rPos.From()
	actualTo := rPos.To()

	if actualFrom != expectedFrom {
		t.Errorf("For relationPos.From():\nExpected: %v\nGot: %v", expectedFrom, actualFrom)
	}

	if actualTo != expectedTo {
		t.Errorf("For relationPos.To():\nExpected: %v\nGot: %v", expectedTo, actualTo)
	}
}

func TestNewRelationPos(t *testing.T) {
	fromPosition := &MockPosition{
		filename: "from_file.txt",
		offset:   100,
		line:     5,
		column:   10,
	}

	toPosition := &MockPosition{
		filename: "to_file.txt",
		offset:   200,
		line:     8,
		column:   15,
	}

	relationPosInstance := NewRelationPos(fromPosition, toPosition)

	expectedFrom := fromPosition
	expectedTo := toPosition

	actualFrom := relationPosInstance.From()
	actualTo := relationPosInstance.To()

	if actualFrom != expectedFrom {
		t.Errorf("For NewRelationPos:\nExpected 'From' position: %v\nGot: %v", expectedFrom, actualFrom)
	}

	if actualTo != expectedTo {
		t.Errorf("For NewRelationPos:\nExpected 'To' position: %v\nGot: %v", expectedTo, actualTo)
	}
}

func TestNewEmptyRelationPos(t *testing.T) {
	emptyRelationPos := NewEmptyRelationPos()

	if emptyRelationPos.From() != nil {
		t.Errorf("For NewEmptyRelationPos:\nExpected empty 'From' position, but got: %v", emptyRelationPos.From())
	}

	if emptyRelationPos.To() != nil {
		t.Errorf("For NewEmptyRelationPos:\nExpected empty 'To' position, but got: %v", emptyRelationPos.To())
	}
}
