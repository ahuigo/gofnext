package examples

import (
	"strings"
	"sync"
	"testing"
)

// Parallel caller via goroutines
func parallelCall(fn func(), times int) {
	var wg sync.WaitGroup
	for k := 0; k < times; k++ {
		wg.Add(1)
		go func() {
			fn()
			wg.Done()
		}()
	}
	wg.Wait()
}

func assertContains(t *testing.T, s string, search string) {
	if !strings.Contains(s, search) {
		t.Fatalf("%s should contains %s", s, search)
	}
}
