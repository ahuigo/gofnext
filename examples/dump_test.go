// BEGIN: 7f3d8b9g4h5j
package examples

import (
	"testing"

	"github.com/ahuigo/gofnext/dump"
)

type Person struct {
	Name string
	age  int
}

func TestDeepSerial(t *testing.T) {
	// Test case 1: Integer
	num := 42
	expectedNum := "42"
	if result := dump.Dump(num); result != expectedNum {
		t.Errorf("Expected %s, but got %s", expectedNum, result)
	}

	// Test case 2: String
	str := "Hello, World!"
	expectedStr := `"Hello, World!"`
	if result := dump.Dump(str); result != expectedStr {
		t.Errorf("Expected %s, but got %s", expectedStr, result)
	}

	// Test case 3: Struct
	person := Person{Name: "John Doe", age: 30}
	expectedPerson := `{Name:"John Doe",age:30}`
	if result := dump.Dump(person); result != expectedPerson {
		t.Errorf("Expected %s, but got %s", expectedPerson, result)
	}

	// Test case 7: pointer
	p := &person
	expectedP := "&Person:{Name:\"John Doe\",age:30}"
	if result := dump.Dump(p); result != expectedP {
		t.Errorf("Expected %s, but got %s", expectedP, result)
	}

	// Test case 4: Slice
	slice := []int{1, 2, 3, 4, 5}
	expectedSlice := "[1,2,3,4,5]"
	if result := dump.Dump(slice); result != expectedSlice {
		t.Errorf("Expected %s, but got %s", expectedSlice, result)
	}

	// Test case 5: Map(multi)
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	expectedMap := `{"a":1,"b":2,"c":3}`
	if result := dump.Dump(m); result != expectedMap {
		t.Errorf("Expected %s, but got %s", expectedMap, result)
	}

	// Test case 6: interface{}
	var i any = 42
	expectedI := "42"
	if result := dump.Dump(i); result != expectedI {
		t.Errorf("Expected %s, but got %s", expectedI, result)
	}


}
