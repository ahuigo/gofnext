package gofnext

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheFuncWithNoParam(t *testing.T) {
	executeCount := 0
	// Original function
	getNumFromDb := func() (int, error) {
		executeCount++
		time.Sleep(10 * time.Millisecond)
		return executeCount, errors.New("db error")
	}

	// Cacheable Function
	getNumWithCache := CacheFn0Err(getNumFromDb, &Config{
		TTL:    400 * time.Millisecond,
		ErrTTL: time.Hour,
	})

	// Execute the function multi times in parallel.
	parallelCall(func() {
		userinfo, err := getNumWithCache()
		fmt.Println(userinfo, err)
	}, 10)

	// Test ttl
	num, _ := getNumWithCache()
	AssertEqual(t, num, 1)

	// Test expired ttl
	time.Sleep(600 * time.Millisecond)
	num, _ = getNumWithCache()
	AssertEqual(t, num, 2)

	if executeCount != 2 {
		t.Error("executeCount should be 2", ", but get ", executeCount)
	}
}

func TestCacheFuncReuseCache(t *testing.T) {
	var executeCount atomic.Int32
	// Original function
	getNum := func() (int32, error) {
		executeCount.Add(1)
		return executeCount.Load(), nil
	}

	// Cacheable Function
	getNumWithCache := CacheFn0Err(getNum, &Config{
		TTL:      200 * time.Millisecond,
		ErrTTL:   time.Hour,
		ReuseTTL: 200 * time.Millisecond,
	})

	// Test ttl
	num, _ := getNumWithCache()
	AssertEqual(t, num, 1)

	// Test expired ttl(reuse ttl)
	time.Sleep(200 * time.Millisecond)
	num, _ = getNumWithCache()
	AssertEqual(t, num, 1)

	// Wait for async goroutine to finish function call
	time.Sleep(100 * time.Millisecond)
	num, _ = getNumWithCache()
	AssertEqual(t, num, 2)

	time.Sleep(100 * time.Millisecond)
	count := executeCount.Load()
	if count != 2 {
		t.Error("executeCount should be 2", ", but get ", count)
	}
}

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

func TestCacheFuncWith2Param(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(c context.Context, arg map[int]int) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", arg[0], time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + arg[0], nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := CacheFn2Err(getUserScore, &Config{
		TTL: time.Hour,
	}) // getFunc can only accept 2 parameter

	// Execute the function multi times in parallel.
	ctx := context.Background()
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(ctx, map[int]int{0: 1})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(ctx, map[int]int{0: 2})
		getUserScoreFromDbWithCache(ctx, map[int]int{0: 3})
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}

}

func TestCacheFunc2WithErr(t *testing.T) {
	getUserScore := func(c context.Context, arg map[int]int) (int, error) {
		return 98, errors.New("db error")
	}
	getUserScoreFromDbWithCache := CacheFn2Err(getUserScore, nil) // getFunc can only accept 2 parameter
	var ctx context.Context
	getUserScoreFromDbWithCache(ctx, map[int]int{0: 1})
}
