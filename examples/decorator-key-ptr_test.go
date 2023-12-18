//go:build !race

package examples

import (
	"testing"

	"github.com/ahuigo/gofnext/go18"
)

func TestCacheFuncKeyPointerAddr(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id   int
	}
	// Original function
	executeCount := 0
	getUserScore := func(user *UserInfo) (int, error) {
		executeCount++
		// fmt.Println("select score from db where id=", user.id, time.Now())
		return 98 + user.id, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		HashKeyPointerAddr: true,
	})

	// Execute the function multi times in parallel.
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

func TestCacheFuncKeyPointerCycle(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id   int
		User *UserInfo
	}
	// Original function
	getUserScore := func(user *UserInfo) (int, error) {
		return 98 + user.id, nil
	}

	// Cacheable Function
	getUserScoreCached := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{})

	// Execute the function multi times in parallel.
	u := &UserInfo{id: 1}
	u.User = u
	score, _ := getUserScoreCached(u)
	if score != 99 {
		t.Errorf("score should be 99, but get %d", score)
	}
}
