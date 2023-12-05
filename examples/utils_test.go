package examples

import "sync"

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
