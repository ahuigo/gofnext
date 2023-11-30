package gofnext

import (
	"testing"
)

func TestIsHashableKey(t *testing.T) {
	// Test case 1: Key is hashable
	key1 := 10
	canHash1 := isHashableKey(key1)
	if !canHash1 {
		t.Errorf("Expected key1 to be hashable, but it is not")
	}

	// Test case 2: Key is not hashable
	key2 := map[string]int{"a": 1, "b": 2}
	canHash2 := isHashableKey(key2)
	if canHash2 {
		t.Errorf("Expected key2 to not be hashable, but it is")
	}

	// Test case 3: Key is nil
	var key3 interface{}
	canHash3 := isHashableKey(key3)
	if !canHash3 {
		t.Errorf("Expected key3 to be hashable, but it is not")
	}

	// Test case 4: Key is a slice
	key4 := []int{1, 2, 3}
	canHash4 := isHashableKey(key4)
	if canHash4 {
		t.Errorf("Expected key4 to not be hashable, but it is")
	}

	// Test case 5: Key is a pointer
	key5 := &key1
	canHash5 := isHashableKey(key5)
	if canHash5 {
		t.Errorf("Expected key5 to not be hashable, but it is")
	}
}