package examples

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/ahuigo/gofnext"
)

type UserInfo struct {
	Name string
	Age  int
}

var (
	count = atomic.Uint32{}
)
func getUserAndErr() (UserInfo, error) {
	count.Add(1)
	// fmt.Println("select * from db limit 1", time.Now())
	return UserInfo{Name: "Anonymous", Age: 9}, errors.New("db error")
}

var (
	// Cacheable Function
	getUserAndErrCached = gofnext.CacheFn0Err(getUserAndErr, nil)
)

func TestNoCacheIfErr(t *testing.T) {
	times := 10
	count.Store(0)

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

func TestNeedCacheIfErr(t *testing.T) {
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
		NeedCacheIfErr: true,
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
