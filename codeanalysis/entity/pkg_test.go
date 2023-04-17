package entity

import (
	"fmt"
	"testing"
)

func TestPkg_Load(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedError error
	}{
		{
			name:          "Invalid Package Path",
			path:          "github.com/your-username/non-existent-package",
			expectedError: fmt.Errorf("packages contain errors"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkg := &Pkg{Path: test.path}
			err := pkg.Load()

			if test.expectedError == nil {
				if err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}
				if len(pkg.Initial) == 0 {
					t.Error("Expected non-empty package, but got empty package")
				}
			} else {
				if err == nil {
					t.Error("Expected error, but got nil")
				} else if err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error '%v', but got '%v'", test.expectedError, err)
				}
				if len(pkg.Initial) != 0 {
					t.Error("Expected empty package, but got non-empty package")
				}
			}
		})
	}
}
