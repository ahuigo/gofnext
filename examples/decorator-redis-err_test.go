package examples

import (
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func TestRedisCacheClientPanic(t *testing.T) {
	defer func() {
		r := recover() //r.(string)
		if r == nil {
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
		return 1 + more, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL:      500 * time.Second,
		CacheMap: gofnext.NewCacheRedis("redis-cache-key").ClearAll(),
	})

	// Execute the function multi times in parallel.
	for i := 0; i < 10; i++ {
		score, err := getUserScoreFromDbWithCache(1)
		if err!=nil{
			t.Fatal(err)
		}
		if score != 2 {
			t.Errorf("score should be 2, but get %d", score)
		}
		score, _ = getUserScoreFromDbWithCache(2)
		if score != 3 {
			t.Fatalf("score should be 3, but get %d", score)
		}
		getUserScoreFromDbWithCache(3)
		getUserScoreFromDbWithCache(3)
	}

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}
