package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func TestCacheFuncKeyStruct(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id   int
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
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL: time.Hour,
	})

	// Execute the function multi times in parallel.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(UserInfo{id: 1})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(UserInfo{id: 2})
		getUserScoreFromDbWithCache(UserInfo{id: 3})
	}, 5)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}
func TestCacheFuncKeyStructUnexportedKey(t *testing.T) {
    type info struct{
        name string // unexported field
    }

	type UserInfo struct {
        info struct{
            name string // unexported field
        }
	}
	// Original function
	getUserName := func(user UserInfo, flag *string) string {
		return user.info.name
	}

	// Cacheable Function
	getUserName2 := gofnext.CacheFn2(getUserName, )

	// Execute the function
    flag := "flag"
	name := getUserName2(UserInfo{info{name: "Alex"}}, &flag)
	if name != "Alex" {
		t.Errorf("name should be 'Alex', but get %s", name)
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
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL: time.Hour,
	})

	// Execute the function multi times in parallel.
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

func TestCacheFuncKeyDeepMap(t *testing.T) {
	// Original function
	type Params struct {
		m map[string]int
	}
	executeCount := 0
	getUserScore := func(params Params) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", params.m["id"], time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + params.m["id"], nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL: time.Hour,
	})

	// Execute the function multi times in parallel.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(Params{m: map[string]int{"id": 1}})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(Params{m: map[string]int{"id": 2}})
		getUserScoreFromDbWithCache(Params{m: map[string]int{"id": 3}})
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}
func TestCacheFuncKeySlice(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id   int
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
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL: time.Hour,
	})

	// Execute the function multi times in parallel.
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
		id   int
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
		TTL: time.Hour,
	})

	// Execute the function multi times in parallel.
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

func TestCacheFuncKeyInnerPointer(t *testing.T) {
	type UserInfo struct {
		Name string
		Age  int
		id   int
	}
	type ExtUserInfo struct {
		u *UserInfo
	}
	// Original function
	executeCount := 0
	getUserScore := func(user ExtUserInfo) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", user.u.id, time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + user.u.id, nil
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL:         time.Hour,
		NeedDumpKey: true,
	})

	// Execute the function multi times in parallel.
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(ExtUserInfo{u: &UserInfo{id: 1}})
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(ExtUserInfo{u: &UserInfo{id: 2}})
		getUserScoreFromDbWithCache(ExtUserInfo{u: &UserInfo{id: 3}})
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}
