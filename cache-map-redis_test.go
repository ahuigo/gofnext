package decorator

import (
	"testing"
	"time"
)

func TestRedisCacheFuncWithTTL(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(more int) (int, error) {
		executeCount++
		return 98 + more, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := CacheFn1Err(getUserScore, &Config{
		TTL:  time.Hour,
		CacheMap: NewCacheRedis("gocachemap", nil).ClearAll(),
	}) // getFunc can only accept 1 parameter

	// Parallel invocation of multiple functions.
	for i := 0; i < 10; i++ {
		score, err := getUserScoreFromDbWithCache(1)
		if err!=nil ||score != 99 {
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
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}

func TestRedisCacheFuncWithNoTTL(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(more int) (int, error) {
		executeCount++
		return 98 + more, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := CacheFn1Err(
		getUserScore,
		&Config{
			CacheMap: NewCacheRedis("gocachemap", nil).ClearAll(),
		},
	) // getFunc can only accept 1 parameter

	// Parallel invocation of multiple functions.
	parallelCall(func(){
		score, err := getUserScoreFromDbWithCache(1)
		if err!=nil ||score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(2)
		getUserScoreFromDbWithCache(3)
		getUserScoreFromDbWithCache(3)
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}
