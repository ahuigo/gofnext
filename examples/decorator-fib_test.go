package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/gofnext"
)

// https://stackoverflow.com/questions/73379188/how-to-use-cache-decorator-with-a-recursive-function-in-go
func TestFib(t *testing.T) {
	var fib func(int) int
	var fibCached func(int) int
	fib = func(x int) int {
		fmt.Printf("call arg:%d\n", x)
		if x <= 1 {
			return x
		} else {
			return fibCached(x-1) + fibCached(x-2)
		}
	}

	fibCached = gofnext.CacheFn1(fib, nil)

	fmt.Println(fibCached(5))
	fmt.Println(fibCached(6))
}
