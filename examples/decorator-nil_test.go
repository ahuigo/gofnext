package examples

import (
	"fmt"
	"testing"
	"time"

	decorator "github.com/ahuigo/gocache-decorator"
)

func getUserNoError(age int) (UserInfo) {
	time.Sleep(10 * time.Millisecond)
	return UserInfo{Name: "Alex", Age: age}
}

var (
	// Cacheable Function with 1 param and no error
	getUserInfoFromDb= decorator.CacheFn1(getUserNoError, nil) 
)
func TestCacheFuncNil(t *testing.T) {
	// Parallel invocation of multiple functions.
	times := 10
	parallelCall(func() {
		userinfo := getUserInfoFromDb(20)
		fmt.Println(userinfo)
	}, times)
}
