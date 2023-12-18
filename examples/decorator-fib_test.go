package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/gofnext/go18"
)

// https://stackoverflow.com/questions/73379188/how-to-use-cache-decorator-with-a-recursive-function-in-go
func TestFib(t *testing.T) {
	excuteCount := 0
	var fib func(int) int
	fib = func(x int) int {
		excuteCount++
		fmt.Printf("call arg:%d\n", x)
		if x <= 1 {
			return x
		} else {
			return fib(x-1) + fib(x-2)
		}
	}
	fib = gofnext.CacheFn1(fib)

	fmt.Println(fib(5))
	fmt.Println(fib(6))
	if excuteCount != 7 {
		t.Errorf("Expected 7, but got %d", excuteCount)
	}
}
