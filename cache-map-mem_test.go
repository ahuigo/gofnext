package gofnext

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheFuncWithNoParam(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
	}

	executeCount := 0
	// Original function
	getUserInfoFromDb := func() (UserInfo, error) {
		executeCount++
		fmt.Println("select * from db limit 1", time.Now())
		time.Sleep(10 * time.Millisecond)
		return UserInfo{Name: "Anonymous", Age: 9}, errors.New("db error")
	}

	// Cacheable Function
	getUserInfoFromDbWithCache := CacheFn0Err(getUserInfoFromDb, &Config{
		TTL:    400 * time.Millisecond,
		ErrTTL: time.Hour,
	})
	_ = getUserInfoFromDbWithCache

	// Execute the function multi times in parallel.
	parallelCall(func() {
		userinfo, err := getUserInfoFromDbWithCache()
		fmt.Println(userinfo, err)
	}, 10)

	// Test ttl
	_, _ = getUserInfoFromDbWithCache()
	time.Sleep(600 * time.Millisecond)
	_, _ = getUserInfoFromDbWithCache()

	if executeCount != 2 {
		t.Error("executeCount should be 2", ", but get ", executeCount)
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
