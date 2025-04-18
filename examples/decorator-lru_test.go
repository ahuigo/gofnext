package examples

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func TestCacheFuncWithOneParamLRU(t *testing.T) {
	// Original function
	executeCount := 0
	maxCacheSize := 2
	var getUserScore = func(more int) (int, error) {
		executeCount++
		return 98 + more, nil
	}

	// Cacheable Function
	var getUserScoreFromDbWithLruCache = gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL:      time.Hour,
		CacheMap: gofnext.NewCacheLru(maxCacheSize),
	})

	// Execute the function multi times in parallel.
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

func TestReuseCacheForLRU(t *testing.T) {
	// counter
	var executeCount atomic.Int32
	// Original function
	maxCacheSize := 2
	var getNum = func(more int) (int, error) {
		executeCount.Add(1)
		c := executeCount.Load()
		return int(c) + more, nil
	}

	ttl := time.Millisecond * 50
	reuseTTL := time.Millisecond * 10
	// Cacheable Function
	var getNumWithLruCache = gofnext.CacheFn1Err(getNum, &gofnext.Config{
		TTL:      ttl,
		ReuseTTL: reuseTTL,
		CacheMap: gofnext.NewCacheLru(maxCacheSize),
	})

	score, _ := getNumWithLruCache(1)
	gofnext.AssertEqual(t, score, 2)
	// wait ttl
	time.Sleep(ttl)
	score, _ = getNumWithLruCache(1)
	gofnext.AssertEqual(t, score, 2)

	// wait reuseTTL
	time.Sleep(reuseTTL)
	score, _ = getNumWithLruCache(1)
	gofnext.AssertEqual(t, score, 3)

	// test function call count
	count := executeCount.Load()
	if count != 2 {
		t.Errorf("executeCount should be 2, but get %d", count)
	}

}
