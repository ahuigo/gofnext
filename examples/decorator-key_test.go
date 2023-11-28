package examples

import (
	"fmt"
	"testing"
	"time"

	decorator "github.com/ahuigo/gocache-decorator"
)

func TestCacheFuncKeyStruct(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id int
	}
	// Original function
	executeCount := 0
	getUserScore := func(user UserInfo) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", user.id, time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + user.id, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := decorator.CacheFn1Err(getUserScore, &decorator.Config{
		TTL: time.Hour,
	})

	// Parallel invocation of multiple functions.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(UserInfo{id: 1})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(UserInfo{id: 2})
		getUserScoreFromDbWithCache(UserInfo{id: 3})
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}

func TestCacheFuncKeyMap(t *testing.T) {
	// Original function
	type usermap = map[string]int
	executeCount := 0
	getUserScore := func(user usermap) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", user["id"], time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + user["id"], nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := decorator.CacheFn1Err(getUserScore, &decorator.Config{
		TTL: time.Hour,
	})

	// Parallel invocation of multiple functions.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(usermap{"id": 1})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(usermap{"id": 2})
		getUserScoreFromDbWithCache(usermap{"id": 3})
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}
func TestCacheFuncKeySlice(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id int
	}
	// Original function
	executeCount := 0
	getUserScore := func(users []UserInfo) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", users[0].id, time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + users[0].id, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := decorator.CacheFn1Err(getUserScore, &decorator.Config{
		TTL: time.Hour,
	})

	// Parallel invocation of multiple functions.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache([]UserInfo{{id: 1}})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache([]UserInfo{{id: 2}})
		getUserScoreFromDbWithCache([]UserInfo{{id: 3}})
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}

func TestCacheFuncKeyPointer(t *testing.T) {
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
	getUserScoreFromDbWithCache := decorator.CacheFn1Err(getUserScore, &decorator.Config{
		TTL: time.Hour,
	})

	// Parallel invocation of multiple functions.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(&UserInfo{id: 1})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(&UserInfo{id: 2})
		getUserScoreFromDbWithCache(&UserInfo{id: 3})
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}
