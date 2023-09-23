package valueobject

import "testing"

func TestMetaMethods(t *testing.T) {
	testCases := []struct {
		name               string
		meta               *meta
		expectedPkg        string
		expectedName       string
		expectedParentName string
		expectedHasParent  bool
	}{
		{
			name: "Meta with Parent",
			meta: &meta{
				pkg:        "package1",
				name:       "name1",
				parentName: "parent1",
			},
			expectedPkg:        "package1",
			expectedName:       "name1",
			expectedParentName: "parent1",
			expectedHasParent:  true,
		},
		{
			name: "Meta without Parent",
			meta: &meta{
				pkg:        "package2",
				name:       "name2",
				parentName: "",
			},
			expectedPkg:        "package2",
			expectedName:       "name2",
			expectedParentName: "",
			expectedHasParent:  false,
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualPkg := tc.meta.Pkg()
			actualName := tc.meta.Name()
			actualParentName := tc.meta.Parent()
			actualHasParent := tc.meta.HasParent()

			if actualPkg != tc.expectedPkg || actualName != tc.expectedName ||
				actualParentName != tc.expectedParentName || actualHasParent != tc.expectedHasParent {
				t.Errorf("For test case %s:\nExpected: (%s, %s, %s, %v)\nGot: (%s, %s, %s, %v)",
					tc.name, tc.expectedPkg, tc.expectedName, tc.expectedParentName, tc.expectedHasParent,
					actualPkg, actualName, actualParentName, actualHasParent)
			}
		})
	}
}

func TestNewMeta(t *testing.T) {
	testCases := []struct {
		name               string
		pkg                string
		nameVal            string
		expectedPkg        string
		expectedName       string
		expectedParentName string
		expectedHasParent  bool
	}{
		{
			name:               "Meta with Package and Name",
			pkg:                "package1",
			nameVal:            "name1",
			expectedPkg:        "package1",
			expectedName:       "name1",
			expectedParentName: "",
			expectedHasParent:  false,
		},
		{
			name:               "Meta with Another Package and Name",
			pkg:                "package2",
			nameVal:            "name2",
			expectedPkg:        "package2",
			expectedName:       "name2",
			expectedParentName: "",
			expectedHasParent:  false,
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metaInfo := NewMeta(tc.pkg, tc.nameVal)

			actualPkg := metaInfo.Pkg()
			actualName := metaInfo.Name()
			actualParentName := metaInfo.Parent()
			actualHasParent := metaInfo.HasParent()

			if actualPkg != tc.expectedPkg || actualName != tc.expectedName ||
				actualParentName != tc.expectedParentName || actualHasParent != tc.expectedHasParent {
				t.Errorf("For test case %s:\nExpected: (%s, %s, %s, %v)\nGot: (%s, %s, %s, %v)",
					tc.name, tc.expectedPkg, tc.expectedName, tc.expectedParentName, tc.expectedHasParent,
					actualPkg, actualName, actualParentName, actualHasParent)
			}
		})
	}
}

func TestNewMetaWithParent(t *testing.T) {
	testCases := []struct {
		name               string
		pkg                string
		nameVal            string
		parentNameVal      string
		expectedPkg        string
		expectedName       string
		expectedParentName string
		expectedHasParent  bool
	}{
		{
			name:               "Meta with Package, Name, and Parent Name",
			pkg:                "package1",
			nameVal:            "name1",
			parentNameVal:      "parent1",
			expectedPkg:        "package1",
			expectedName:       "name1",
			expectedParentName: "parent1",
			expectedHasParent:  true,
		},
		{
			name:               "Meta with Another Package, Name, and Parent Name",
			pkg:                "package2",
			nameVal:            "name2",
			parentNameVal:      "parent2",
			expectedPkg:        "package2",
			expectedName:       "name2",
			expectedParentName: "parent2",
			expectedHasParent:  true,
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metaInfo := NewMetaWithParent(tc.pkg, tc.nameVal, tc.parentNameVal)

			actualPkg := metaInfo.Pkg()
			actualName := metaInfo.Name()
			actualParentName := metaInfo.Parent()
			actualHasParent := metaInfo.HasParent()

			if actualPkg != tc.expectedPkg || actualName != tc.expectedName ||
				actualParentName != tc.expectedParentName || actualHasParent != tc.expectedHasParent {
				t.Errorf("For test case %s:\nExpected: (%s, %s, %s, %v)\nGot: (%s, %s, %s, %v)",
					tc.name, tc.expectedPkg, tc.expectedName, tc.expectedParentName, tc.expectedHasParent,
					actualPkg, actualName, actualParentName, actualHasParent)
			}
		})
	}
}
