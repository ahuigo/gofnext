package examples

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/ahuigo/gofnext/go18"
)

type UserInfo struct {
	Name string
	Age  int
}

func getUserAndErr() (UserInfo, error) {
	// fmt.Println("select * from db limit 1", time.Now())
	return UserInfo{Name: "Anonymous", Age: 9}, errors.New("db error")
}

var (
	// Cacheable Function
	getUserAndErrCached = gofnext.CacheFn0Err(getUserAndErr, nil)
)

func TestCacheFunc0WithErr(t *testing.T) {
	times := 10
	parallelCall(func() {
		userinfo, err := getUserAndErrCached()
		if err == nil {
			t.Error("should be error, but get nil")
		}
		fmt.Println(userinfo, err)
	}, times)
}

func TestCacheFunc0CacheErr(t *testing.T) {
	count := atomic.Uint32{}
	getUserAndErr := func(age int) (UserInfo, error) {
		count.Add(1)
		if age <= 0 {
			return UserInfo{}, errors.New("invalid age")
		}
		return UserInfo{Name: "Anonymous", Age: 9}, nil
	}
	// Cacheable Function
	getUserAndErrCached := gofnext.CacheFn1Err(getUserAndErr, &gofnext.Config{
		NeedCacheIfErr: true,
	})

	times := 5
	parallelCall(func() {
		_, err := getUserAndErrCached(0) //1 times
		if err == nil {
			t.Error("should be error, but get nil")
		}
	}, times)
	if count.Load() != 1 {
		t.Errorf("Execute count should be 1, but get %d", count.Load())
	}
}
