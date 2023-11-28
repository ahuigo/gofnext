package examples

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

type UserInfo struct {
	Name string
	Age  int
}

func getUserAnonymouse() (UserInfo, error) {
	fmt.Println("select * from db limit 1", time.Now())
	time.Sleep(10 * time.Millisecond)
	return UserInfo{Name: "Anonymous", Age: 9}, errors.New("db error")
}
var (
	// Cacheable Function
	getUserInfoFromDbWithCache = gofnext.CacheFn0Err(getUserAnonymouse, nil) 
)

func TestCacheFunc0Err(t *testing.T) {
	// Parallel invocation of multiple functions.
	times := 10
	parallelCall(func() {
		userinfo, err := getUserInfoFromDbWithCache()
		fmt.Println(userinfo, err)
	}, times)
}