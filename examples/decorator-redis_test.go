package examples

import (
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
	"github.com/go-redis/redis"
)

func TestRedisCacheClient(t *testing.T) {
	// method 1: by default: localhost:6379
	cache := gofnext.NewCacheRedis("redis-cache-key") // you can list value `HGETALL _gofnext:redis-cache-key`

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
	add98Origin := func(more int) (int, error) {
		executeCount++
		return 98 + more, nil
	}
	redisCache := gofnext.NewCacheRedis("redis-cache-key")
	redisCache.ClearAll() // redis> del _gofnext:redis-cache-key
	// redisCache.SetRedisOpts(&redis.Options{
	// 	Addr: "localhost:6379",
	// })

	// Cacheable Function
	add98 := gofnext.CacheFn1Err(add98Origin, &gofnext.Config{
		TTL:      time.Hour,
		CacheMap: redisCache,
	})

	// Execute the function 10 times
	for i := 0; i < 10; i++ {
		score, err := add98(1)
		if err != nil || score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		score, _ = add98(2)
		if score != 100 {
			t.Fatalf("score should be 100, but get %d", score)
		}
		add98(3)
		add98(3)
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

	// Execute the function multi times in parallel.
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

	// Execute the function multi times in parallel.
	for i := 0; i < 5; i++ { //5 times
		getUserScoreFromDbWithCache(1)
		getUserScoreFromDbWithCache(1) // cache hit: read from redis
		time.Sleep(time.Millisecond * 200)
	}

	if executeCount != 5 {
		t.Errorf("executeCount should be 5, but get %d", executeCount)
	}
}
