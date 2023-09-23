package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

func TestRelationMetaMethods(t *testing.T) {
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

	relationType := arch.RelationTypeAssociation

	relationMetaInstance := NewRelationMeta(relationType, fromPosition, toPosition)

	expectedType := relationType
	expectedFrom := fromPosition
	expectedTo := toPosition

	actualType := relationMetaInstance.Type()
	actualFrom := relationMetaInstance.Position().From()
	actualTo := relationMetaInstance.Position().To()

	if actualType != expectedType {
		t.Errorf("For RelationMeta.Type():\nExpected: %d\nGot: %d", expectedType, actualType)
	}

	if actualFrom != expectedFrom {
		t.Errorf("For RelationMeta.Position().From():\nExpected: %v\nGot: %v", expectedFrom, actualFrom)
	}

	if actualTo != expectedTo {
		t.Errorf("For RelationMeta.Position().To():\nExpected: %v\nGot: %v", expectedTo, actualTo)
	}
}

func TestNewRelationMeta(t *testing.T) {
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

	relationType := arch.RelationTypeAssociation

	relationMetaInstance := NewRelationMeta(relationType, fromPosition, toPosition)

	expectedType := relationType
	expectedFrom := fromPosition
	expectedTo := toPosition

	actualType := relationMetaInstance.Type()
	actualFrom := relationMetaInstance.Position().From()
	actualTo := relationMetaInstance.Position().To()

	if actualType != expectedType {
		t.Errorf("For NewRelationMeta:\nExpected Type: %d\nGot Type: %d", expectedType, actualType)
	}

	if actualFrom != expectedFrom {
		t.Errorf("For NewRelationMeta:\nExpected From: %v\nGot From: %v", expectedFrom, actualFrom)
	}

	if actualTo != expectedTo {
		t.Errorf("For NewRelationMeta:\nExpected To: %v\nGot To: %v", expectedTo, actualTo)
	}
}
