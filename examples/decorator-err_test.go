package examples

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

type UserInfo struct {
	Name string
	Age  int
}

var (
	count = atomic.Uint32{}
)
func getUserWithErr() (UserInfo, error) {
	count.Add(1)
	// fmt.Println("select * from db limit 1", time.Now())
	return UserInfo{Name: "Anonymous", Age: 9}, errors.New("db error")
}


// Test: if there is any error, do not cache error(default)
func TestNoCacheErr(t *testing.T) {
	times := 10
	count.Store(0)
	// 0. Cacheable Function
	getUserAndErrCached := gofnext.CacheFn0Err(getUserWithErr, nil)

	// 1. run 10 times
	parallelCall(func() {
		userinfo, err := getUserAndErrCached()
		if err == nil {
			t.Error("should be error, but get nil")
		}
		fmt.Println(userinfo, err)
	}, times)

	// 2. check count
	if count.Load() !=  uint32(times){
		t.Fatalf("Execute count should be %d, but get %d", times, count.Load())
	}
}

// Test: if there is an error, cache the error
func TestNeedCacheErrWithTTL(t *testing.T) {
	count := atomic.Uint32{}
	getUserAndErr := func(age int) (UserInfo, error) {
		count.Add(1)
		if age <= 0 {
			return UserInfo{}, errors.New("invalid age")
		}
		return UserInfo{Name: "Anonymous", Age: 9}, nil
	}
	// 1. Cacheable Function
	getUserAndErrCached := gofnext.CacheFn1Err(getUserAndErr, &gofnext.Config{
		ErrTTL: time.Hour,
	})

	times := 5
	// 2. run 5 times
	parallelCall(func() {
		_, err := getUserAndErrCached(0) //1 times
		if err == nil {
			t.Error("should be error, but get nil")
		}
	}, times)

	// 3. check count
	if count.Load() != 1 {
		t.Errorf("Execute count should be 1, but get %d", count.Load())
	}
}

// Test: if there is an error, cache the error
func TestNeedCacheErrWithTtlTimeout(t *testing.T) {
	count := atomic.Uint32{}
	getUserWithErr := func(age int) (UserInfo, error) {
		count.Add(1)
		if age <= 0 {
			return UserInfo{}, errors.New("invalid age")
		}
		return UserInfo{Name: "Anonymous", Age: 9}, nil
	}
	interval := time.Millisecond * 3
	// 1. Cacheable Function
	getUserWithErrCached := gofnext.CacheFn1Err(getUserWithErr, &gofnext.Config{
		ErrTTL: interval * 1,
	})

	times := 5
	// 2. run 5 times
	for i := 0; i < times; i++ {
		_, err := getUserWithErrCached(0) // cost 1us-50us(0.05ms)
		if err == nil {
			t.Error("should be error, but get nil")
		}
		time.Sleep(time.Millisecond * 1)
	}

	// 3. check count
	if count.Load() != 2 {
		t.Errorf("Execute count should be 1, but get %d", count.Load())
	}
}

func TestNeedLruCacheErrWithTtlTimeout(t *testing.T) {
	count := atomic.Uint32{}
	getUserWithErr := func(age int) (UserInfo, error) {
		count.Add(1)
		if age <= 0 {
			return UserInfo{}, errors.New("invalid age")
		}
		return UserInfo{Name: "Anonymous", Age: 9}, nil
	}
	interval := time.Millisecond * 3
	// 1. Cacheable Function
	getUserWithErrCached := gofnext.CacheFn1Err(getUserWithErr, &gofnext.Config{
		ErrTTL: interval * 1,
		CacheMap: gofnext.NewCacheLru(100),
	})

	times := 5
	// 2. run 5 times
	for i := 0; i < times; i++ {
		_, err := getUserWithErrCached(0) // cost 1us-50us(0.05ms)
		if err == nil {
			t.Error("should be error, but get nil")
		}
		time.Sleep(time.Millisecond * 1)
	}

	// 3. check count
	if count.Load() != 2 {
		t.Errorf("Execute count should be 1, but get %d", count.Load())
	}
}
