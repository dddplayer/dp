package valueobject

import (
	"testing"
)

func TestRelation_Identifier(t *testing.T) {
	from := Identifier{Name: "from", Path: "fromPath"}
	to := Identifier{Name: "to", Path: "toPath"}
	r := Relation{
		From: &from,
		To:   &to,
	}

	expected := Identifier{Name: "from->to", Path: "fromPath->toPath"}
	actual := *r.Identifier()

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func TestRelation_Position(t *testing.T) {
	pos := Position{
		Filename: "filename",
		Offset:   10,
		Line:     2,
		Column:   4,
	}
	r := Relation{
		Pos: &pos,
	}

	expected := pos
	actual := *r.Position()

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}
