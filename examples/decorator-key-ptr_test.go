package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func TestCacheFuncKeyPointerAddr(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id int
	}
	// Original function
	executeCount := 0
	getUserScore := func(user *UserInfo) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", user.id, time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + user.id, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		NeedHashKeyPointerAddr: true,
	})

	// Parallel invocation of multiple functions.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(&UserInfo{id: 1})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(&UserInfo{id: 2})
	}, 5)

	if executeCount != 10 {
		t.Errorf("executeCount should be 10, but get %d", executeCount)
	}
}
