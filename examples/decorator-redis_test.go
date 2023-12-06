package examples

import (
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
	"github.com/go-redis/redis"
)

func TestRedisCacheClient(t *testing.T) {
	// method 1: by default: localhost:6379
	cache := gofnext.NewCacheRedis("redis-cache-key")

	// method 2: set redis addr
	cache.SetRedisAddr("192.168.1.1:6379")

	// method 3: set redis options
	cache.SetRedisOpts(&redis.Options{
		Addr: "localhost:6379",
	})

	// method 4: set redis universal options
	cache.SetRedisUniversalOpts(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})
	cache.SetMaxHashKeyLen(0)
	cache.SetMaxHashKeyLen(100)
}

func TestRedisCacheFuncWithTTL(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(more int) (int, error) {
		executeCount++
		return 98 + more, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL:      time.Hour,
		CacheMap: gofnext.NewCacheRedis("redis-cache-key").ClearAll(),
	})

	// Parallel invocation of multiple functions.
	for i := 0; i < 10; i++ {
		score, err := getUserScoreFromDbWithCache(1)
		if err != nil || score != 99 {
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

func TestRedisCacheFuncWithNoTTL(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(more int, flag bool) (int, error) {
		executeCount++
		return 98 + more, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn2Err(
		getUserScore,
		&gofnext.Config{
			CacheMap: gofnext.NewCacheRedis("redis-cache-key").ClearAll(),
		},
	) // getFunc can only accept 1 parameter

	// Parallel invocation of multiple functions.
	parallelCall(func() {
		score, err := getUserScoreFromDbWithCache(1, true)
		if err != nil || score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(2, true)
		getUserScoreFromDbWithCache(3, true)
		getUserScoreFromDbWithCache(3, true)
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}

func TestRedisCacheFuncWithTTLTimeout(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(more int) (int, error) {
		executeCount++
		return 98 + more, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL:      time.Millisecond * 200,
		CacheMap: gofnext.NewCacheRedis("redis-cache-key").ClearAll(),
	})

	// Parallel invocation of multiple functions.
	for i := 0; i < 5; i++ {//2+4=6 times
		getUserScoreFromDbWithCache(1)
		time.Sleep(time.Millisecond * 500)
		getUserScoreFromDbWithCache(1)
	}

	if executeCount != 6 {
		t.Errorf("executeCount should be 6, but get %d", executeCount)
	}
}