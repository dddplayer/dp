package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

func TestRelationMethods(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj := &obj{
		id:  &ident{name: "toObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	testCases := []struct {
		name         string
		relationType arch.RelationType
		from         *obj
		expectedType arch.RelationType
		expectedFrom *obj
	}{
		{
			name:         "Relation with Association Type",
			relationType: arch.RelationTypeAssociation,
			from:         fromObj,
			expectedType: arch.RelationTypeAssociation,
			expectedFrom: fromObj,
		},
		{
			name:         "Relation with Composition Type",
			relationType: arch.RelationTypeComposition,
			from:         fromObj,
			expectedType: arch.RelationTypeComposition,
			expectedFrom: fromObj,
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &relation{from: tc.from, to: toObj, relType: tc.relationType}

			actualType := r.Type()
			actualFrom := r.From()

			if actualType != tc.expectedType || actualFrom != tc.expectedFrom {
				t.Errorf("For test case %s:\nExpected: (%d, %v)\nGot: (%d, %v)",
					tc.name, tc.expectedType, tc.expectedFrom, actualType, actualFrom)
			}
		})
	}
}

func TestDependenceMethods(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj := &obj{
		id:  &ident{name: "toObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	dependenceType := arch.RelationTypeAssociation

	dependence := &Dependence{
		relation: &relation{
			from:    fromObj,
			to:      toObj,
			relType: dependenceType,
		},
	}

	expectedDependsOn := toObj

	actualDependsOn := dependence.DependsOn()

	if actualDependsOn != expectedDependsOn {
		t.Errorf("For Dependence:\nExpected: %v\nGot: %v", expectedDependsOn, actualDependsOn)
	}
}

func TestNewDependence(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj := &obj{
		id:  &ident{name: "toObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	dependence := NewDependence(fromObj, toObj).(*Dependence)

	expectedType := arch.RelationTypeDependency
	expectedFrom := fromObj
	expectedTo := toObj

	actualType := dependence.Type()
	actualFrom := dependence.From()
	actualDependsOn := dependence.DependsOn()

	if actualType != expectedType || actualFrom != expectedFrom || actualDependsOn != expectedTo {
		t.Errorf("For NewDependence:\nExpected: (%d, %v, %v)\nGot: (%d, %v, %v)",
			expectedType, expectedFrom, expectedTo, actualType, actualFrom, actualDependsOn)
	}
}

func TestCompositionMethods(t *testing.T) {
	parentObj := &obj{
		id:  &ident{name: "parentObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	childObj := &obj{
		id:  &ident{name: "childObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	compositionType := arch.RelationTypeComposition

	composition := &Composition{
		relation: &relation{
			from:    parentObj,
			to:      childObj,
			relType: compositionType,
		},
	}

	expectedChild := childObj

	actualChild := composition.Child()

	if actualChild != expectedChild {
		t.Errorf("For Composition:\nExpected: %v\nGot: %v", expectedChild, actualChild)
	}
}

func TestNewComposition(t *testing.T) {
	parentObj := &obj{
		id:  &ident{name: "parentObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	childObj := &obj{
		id:  &ident{name: "childObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	composition := NewComposition(parentObj, childObj).(*Composition)

	expectedType := arch.RelationTypeComposition
	expectedParent := parentObj
	expectedChild := childObj

	actualType := composition.Type()
	actualFrom := composition.From()
	actualChild := composition.Child()

	if actualType != expectedType || actualFrom != expectedParent || actualChild != expectedChild {
		t.Errorf("For NewComposition:\nExpected: (%d, %v, %v)\nGot: (%d, %v, %v)",
			expectedType, expectedParent, expectedChild, actualType, actualFrom, actualChild)
	}
}

func TestEmbeddingMethods(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj := &obj{
		id:  &ident{name: "toObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	embeddingType := arch.RelationTypeEmbedding

	embedding := &Embedding{
		relation: &relation{
			from:    fromObj,
			to:      toObj,
			relType: embeddingType,
		},
	}

	expectedEmbedded := toObj

	actualEmbedded := embedding.Embedded()

	if actualEmbedded != expectedEmbedded {
		t.Errorf("For Embedding:\nExpected: %v\nGot: %v", expectedEmbedded, actualEmbedded)
	}
}

func TestNewEmbedding(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj := &obj{
		id:  &ident{name: "toObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	embedding := NewEmbedding(fromObj, toObj).(*Embedding)

	expectedType := arch.RelationTypeEmbedding
	expectedFrom := fromObj
	expectedEmbedded := toObj

	actualType := embedding.Type()
	actualFrom := embedding.From()
	actualEmbedded := embedding.Embedded()

	if actualType != expectedType || actualFrom != expectedFrom || actualEmbedded != expectedEmbedded {
		t.Errorf("For NewEmbedding:\nExpected: (%d, %v, %v)\nGot: (%d, %v, %v)",
			expectedType, expectedFrom, expectedEmbedded, actualType, actualFrom, actualEmbedded)
	}
}

func TestImplementationMethods(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj1 := &obj{
		id:  &ident{name: "toObj1", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}
	toObj2 := &obj{
		id:  &ident{name: "toObj2", pkg: "package3"},
		pos: &pos{filename: "file3.txt", offset: 300, line: 12, column: 20},
	}

	implementationType := arch.RelationTypeImplementation

	implementation := &Implementation{
		relation: &relation{
			from:    fromObj,
			to:      nil,
			relType: implementationType,
		},
		to: make([]arch.Object, 0),
	}

	implementation.Implemented(toObj1)
	implementation.Implemented(toObj2)

	expectedImplemented := []arch.Object{toObj1, toObj2}

	actualImplemented := implementation.Implements()

	if len(actualImplemented) != len(expectedImplemented) {
		t.Errorf("For Implementation:\nExpected: %v\nGot: %v", expectedImplemented, actualImplemented)
	} else {
		for i := range expectedImplemented {
			if actualImplemented[i] != expectedImplemented[i] {
				t.Errorf("For Implementation:\nExpected: %v\nGot: %v", expectedImplemented, actualImplemented)
				break
			}
		}
	}
}

func TestNewImplementation(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj := &obj{
		id:  &ident{name: "toObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	implementation := NewImplementation(fromObj, toObj).(*Implementation)

	expectedType := arch.RelationTypeImplementation
	expectedFrom := fromObj
	expectedImplemented := []arch.Object{toObj}

	actualType := implementation.Type()
	actualFrom := implementation.From()
	actualImplemented := implementation.Implements()

	if actualType != expectedType || actualFrom != expectedFrom || len(actualImplemented) != len(expectedImplemented) {
		t.Errorf("For NewImplementation:\nExpected: (%d, %v, %v)\nGot: (%d, %v, %v)",
			expectedType, expectedFrom, expectedImplemented, actualType, actualFrom, actualImplemented)
	} else {
		for i := range expectedImplemented {
			if actualImplemented[i] != expectedImplemented[i] {
				t.Errorf("For NewImplementation:\nExpected: (%d, %v, %v)\nGot: (%d, %v, %v)",
					expectedType, expectedFrom, expectedImplemented, actualType, actualFrom, actualImplemented)
				break
			}
		}
	}
}

func TestAssociationMethods(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj := &obj{
		id:  &ident{name: "toObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	associationType := arch.RelationTypeAssociation

	association := &Association{
		relation: &relation{
			from:    fromObj,
			to:      toObj,
			relType: associationType,
		},
		ship: associationType,
	}

	expectedRefer := toObj
	expectedAssociationType := associationType

	actualRefer := association.Refer()
	actualAssociationType := association.AssociationType()

	if actualRefer != expectedRefer {
		t.Errorf("For Association.Refer():\nExpected: %v\nGot: %v", expectedRefer, actualRefer)
	}

	if actualAssociationType != expectedAssociationType {
		t.Errorf("For Association.AssociationType():\nExpected: %v\nGot: %v", expectedAssociationType, actualAssociationType)
	}
}

func TestNewAssociation(t *testing.T) {
	fromObj := &obj{
		id:  &ident{name: "fromObj", pkg: "package1"},
		pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
	}
	toObj := &obj{
		id:  &ident{name: "toObj", pkg: "package2"},
		pos: &pos{filename: "file2.txt", offset: 200, line: 8, column: 15},
	}

	associationType := arch.RelationTypeAssociation

	association := NewAssociation(fromObj, toObj, associationType).(*Association)

	expectedType := arch.RelationTypeAssociation
	expectedFrom := fromObj
	expectedAssociationType := associationType

	actualType := association.Type()
	actualFrom := association.From()
	actualAssociationType := association.AssociationType()

	if actualType != expectedType || actualFrom != expectedFrom || actualAssociationType != expectedAssociationType {
		t.Errorf("For NewAssociation:\nExpected: (%d, %v, %v)\nGot: (%d, %v, %v)",
			expectedType, expectedFrom, expectedAssociationType, actualType, actualFrom, actualAssociationType)
	}
}
