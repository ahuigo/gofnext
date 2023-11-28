package examples

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func TestCacheFuncWithOneParamLRU(t *testing.T) {
	// Original function
	executeCount := 0
	var getUserScore = func(more int) (int, error) {
		executeCount++
		return 98 + more, errors.New("db error")
	}

	// Cacheable Function
	var getUserScoreFromDbWithLruCache = gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL:  time.Hour,
		CacheMap: gofnext.NewCacheLru(2),
	})

	// Parallel invocation of multiple functions.
	for i := 0; i < 10; i++ {
		score, err := getUserScoreFromDbWithLruCache(1)
		fmt.Println(score, err)
		score, err = getUserScoreFromDbWithLruCache(2)
		fmt.Println(score, err)
		getUserScoreFromDbWithLruCache(3)
		getUserScoreFromDbWithLruCache(3)
	}

	if executeCount != 30 {
		t.Errorf("executeCount should be 30, but get %d", executeCount)
	}

}
