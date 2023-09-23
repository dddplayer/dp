package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/code"
	"testing"
)

func TestIdentifier_String(t *testing.T) {
	id := ident{
		name: "testName",
		pkg:  "/pkg/to/identifier",
	}
	expected := "/pkg/to/identifier/testName"
	if id.ID() != expected {
		t.Errorf("Expected %s, but got %s", expected, id.ID())
	}
}

func TestNewIdentifier(t *testing.T) {
	testCases := []struct {
		name         string
		meta         code.MetaInfo
		expectedName string
		expectedPkg  string
	}{
		{
			name: "Identifier with Parent",
			meta: &DummyMetaInfo{
				pkg:        "package1",
				name:       "name1",
				parentName: "parent1",
			},
			expectedName: "parent1" + DotJoiner + "name1",
			expectedPkg:  "package1",
		},
		{
			name: "Identifier without Parent",
			meta: &DummyMetaInfo{
				pkg:        "package2",
				name:       "name2",
				parentName: "",
			},
			expectedName: "name2",
			expectedPkg:  "package2",
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualIdentifier := newIdentifier(tc.meta)

			if actualIdentifier.name != tc.expectedName || actualIdentifier.pkg != tc.expectedPkg {
				t.Errorf("For test case %s:\nExpected: (%s, %s)\nGot: (%s, %s)",
					tc.name, tc.expectedName, tc.expectedPkg, actualIdentifier.name, actualIdentifier.pkg)
			}
		})
	}
}
