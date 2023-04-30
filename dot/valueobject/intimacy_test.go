package valueobject

import (
	"math/rand"
	"sort"
	"testing"
)

func TestIntimacyElementsSort(t *testing.T) {
	// Create test data
	elems := make(IntimacyElements, 10)
	for i := range elems {
		elems[i] = &IntimacyElement{
			Key: i,
			Val: rand.Float64(),
		}
	}

	// Sort the elements
	sort.Sort(elems)

	// Check if the elements are sorted in descending order
	for i := 0; i < len(elems)-1; i++ {
		if elems[i].Val < elems[i+1].Val {
			t.Errorf("IntimacyElements not sorted correctly")
		}
	}
}
