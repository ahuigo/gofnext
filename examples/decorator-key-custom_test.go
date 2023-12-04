package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func TestCacheFuncKeyCustom(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id int
	}
	// Original function
	executeCount := 0
	getUserScore := func(user *UserInfo, flag bool) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", user.id, time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + user.id, nil
	}

	// hash key function
	hashKeyFunc := func(keys ...any) []byte{
		user := keys[0].(*UserInfo)
		flag := keys[1].(bool)
		// println(fmt.Sprintf("user:%d,flag:%t", user.id, flag))
		return []byte(fmt.Sprintf("user:%d,flag:%t", user.id, flag))
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn2Err(getUserScore, &gofnext.Config{
		HashKeyFunc: hashKeyFunc,
	})

	// Parallel invocation of multiple functions.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(&UserInfo{id: 1}, true)
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(&UserInfo{id: 2}, true)
		getUserScoreFromDbWithCache(&UserInfo{id: 2}, false)
		getUserScoreFromDbWithCache(&UserInfo{id: 2}, false)
	}, 5)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}
