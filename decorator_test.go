package decorator

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
	getUserInfoFromDbWithCache := DecoratorFn0(getUserInfoFromDb, &Config{Timeout: 400 * time.Millisecond}) // getFunc can only accept zero parameter
	_ = getUserInfoFromDbWithCache

	// Parallel invocation of multiple functions.
	parallelCall(func() {
		userinfo, err := getUserInfoFromDbWithCache()
		fmt.Println(userinfo, err)
	}, 10)

	// Test timeout
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
		return 98, errors.New("db error")
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := DecoratorFn2(getUserScore, &Config{
		Timeout: time.Hour,
	}) // getFunc can only accept 1 parameter

	// Parallel invocation of multiple functions.
	ctx := context.Background()
	parallelCall(func() {
		score, err := getUserScoreFromDbWithCache(ctx, map[int]int{0: 1})
		fmt.Println(score, err)
		score, err = getUserScoreFromDbWithCache(ctx, map[int]int{0: 2})
		fmt.Println(score, err)
		getUserScoreFromDbWithCache(ctx, map[int]int{0: 3})
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}

}

func TestCacheFuncWithNilContext(t *testing.T) {
	getUserScore := func(c context.Context, arg map[int]int) (int, error) {
		return 98, errors.New("db error")
	}
	getUserScoreFromDbWithCache := DecoratorFn2(getUserScore, nil) // getFunc can only accept 1 parameter
	var ctx context.Context
	getUserScoreFromDbWithCache(ctx, map[int]int{0: 1})
}

func TestCacheFuncWithOneParamLRU(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(more int) (int, error) {
		executeCount++
		return 98+more, errors.New("db error")
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := DecoratorFn1(getUserScore, &Config{
		Timeout: time.Hour,
		CacheMap: NewCacheLru(2, time.Second),
	}) // getFunc can only accept 1 parameter

	// Parallel invocation of multiple functions.
	for i:=0;i<10;i++{
		score, err := getUserScoreFromDbWithCache(1)
		fmt.Println(score, err)
		score, err = getUserScoreFromDbWithCache(2)
		fmt.Println(score, err)
		getUserScoreFromDbWithCache(3)
		getUserScoreFromDbWithCache(3)
	}

	if executeCount != 30 {
		t.Errorf("executeCount should be 30, but get %d", executeCount)
	}

}