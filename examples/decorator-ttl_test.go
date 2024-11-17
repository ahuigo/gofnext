package examples

import (
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func TestCacheFuncWithTTLTimeout(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(more int) (int, error) {
		executeCount++
		return 98 + more, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL: time.Millisecond * 10,
	})

	// Execute the function multi times
	getUserScoreFromDbWithCache(1)
	getUserScoreFromDbWithCache(1)
	time.Sleep(time.Millisecond * 11) // wait for ttl timeout
	getUserScoreFromDbWithCache(1)
	getUserScoreFromDbWithCache(1)

	if executeCount != 2 {
		t.Errorf("executeCount should be 2, but get %d", executeCount)
	}
}
