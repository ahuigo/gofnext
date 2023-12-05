package serial

import (
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestLoad(t *testing.T) {
	var f float64
	Load([]byte(`-3.14`), &f)
	expectedFloat := -3.14
	if !almostEqual(f, expectedFloat) {
		t.Errorf("got %f, want %f", f, expectedFloat)
	}

	var str string
	data := []byte(`"Hello, World!chars:\"\r\n\t\b"`)
	_ = Load(data, &str)
	expected := "Hello, World!chars:\"\r\n\t\b"
	if str != expected {
		t.Errorf("got %q, want %q", str, expected)
	}
	data = []byte(`"Hello, World!`)
	err := Load(data, &str)
	if err == nil {
		t.Errorf("should be error, but get nil")
	}

	var i int
	Load([]byte(`-42`), &i)
	expectedInt := -42
	if i != expectedInt {
		t.Errorf("got %d, want %d", i, expectedInt)
	}

}
