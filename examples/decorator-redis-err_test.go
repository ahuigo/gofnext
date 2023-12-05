package examples

import (
	"errors"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func TestRedisCacheClientPanic(t *testing.T) {
	defer func() {
		r := recover() //r.(string)
		if r==nil{
			t.Error("should panic")
		}
	}()
	gofnext.NewCacheRedis("")

}

func TestRedisCacheFuncErr(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(more int) (int, error) {
		executeCount++
		return 98 + more, errors.New("db error")
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL:      500*time.Second,
		CacheMap: gofnext.NewCacheRedis("redis-cache-key").ClearAll(),
	})

	// Parallel invocation of multiple functions.
	for i := 0; i < 10; i++ {
		score, err := getUserScoreFromDbWithCache(1)
		if err == nil {
			t.Errorf("should be error, but get nil")
		}
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		score, _ = getUserScoreFromDbWithCache(2)
		if score != 100 {
			t.Fatalf("score should be 100, but get %d", score)
		}
		getUserScoreFromDbWithCache(3)
		getUserScoreFromDbWithCache(3)
	}

	if executeCount != 3 {
		t.Errorf("executeCount should be 1, but get %d", executeCount)
	}
}
