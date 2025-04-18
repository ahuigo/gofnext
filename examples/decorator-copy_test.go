package examples

import (
	"testing"

	"github.com/ahuigo/gofnext"
)

func TestCacheFuncCopy(t *testing.T) {
	count := 0
	getUser := func(age int) UserInfo {
		count += 1
		return UserInfo{Name: "Alex", Age: age}
	}

	// cacheable function
	getUserCached1 := gofnext.CacheFn1(getUser)
	getUserCached2 := getUserCached1
	getUserCached1(20)
	getUserCached2(20)
	if count != 1 {
		t.Error("count should be 1")
	}
}
