package examples

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func getUser() UserInfo {
	time.Sleep(10 * time.Millisecond)
	return UserInfo{Name: "Alex", Age: 20}
}

func TestCacheFuncWith0Param(t *testing.T) {
	getUserInfoFromDb := gofnext.CacheFn0(getUser)
	// Execute the function 10 times in parallel.
	parallelCall(func() {
		userinfo := getUserInfoFromDb()
		fmt.Println(userinfo)
	}, 10) //10 times
}

func TestCacheFuncWith1Param(t *testing.T) {
	getUser := func(age int) UserInfo {
		time.Sleep(10 * time.Millisecond)
		return UserInfo{Name: "Alex", Age: age}
	}

	// cacheable function
	getUserInfoFromDb := gofnext.CacheFn1(getUser)

	parallelCall(func() {
		userinfo := getUserInfoFromDb(20)
		fmt.Println(userinfo)
	}, 10)
}

func TestCacheFuncWith2Params(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(c context.Context, id int) int {
		executeCount++
		fmt.Println("select score from db where id=", id, time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + id
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn2(getUserScore, &gofnext.Config{
		TTL: time.Hour,
	}) // getFunc can only accept 2 parameter

	// Execute the function multi times in parallel.
	ctx := context.Background()
	parallelCall(func() {
		score := getUserScoreFromDbWithCache(ctx, 1)
		if score != 99 {
			t.Errorf("score should be 99, but get %d", score)
		}
		getUserScoreFromDbWithCache(ctx, 2)
		getUserScoreFromDbWithCache(ctx, 3)
	}, 10)

	if executeCount != 3 {
		t.Errorf("executeCount should be 3, but get %d", executeCount)
	}
}

func TestCacheFuncWith3Params(t *testing.T) {
	// Original function
	executeCount := 0
	sum := func(a, b, c int) int {
		executeCount++
		time.Sleep(10 * time.Millisecond)
		return a + b + c
	}

	// Cacheable Function
	sumCache := gofnext.CacheFn3(sum, &gofnext.Config{
		TTL: time.Hour,
	}) // getFunc can only accept 2 parameter

	// Execute the function multi times in parallel.
	// ctx := context.Background()
	parallelCall(func() {
		score := sumCache(1, 3, 5)
		if score != 9 {
			t.Errorf("score should be 99, but get %d", score)
		}
		sumCache(1, 3, 5)
		sumCache(1, 3, 6)
	}, 5)

	if executeCount != 2 {
		t.Errorf("executeCount should be 2, but get %d", executeCount)
	}
}

func TestCacheCtxFuncWith3Params(t *testing.T) {
	// Original function
	executeCount := 0
	sum := func(ctx context.Context, b, c int) int {
		executeCount++
		time.Sleep(10 * time.Millisecond)
		return b + c
	}

	// Cacheable Function
	sumCache := gofnext.CacheFn3(sum, &gofnext.Config{
		TTL: time.Hour,
	}) // accept 3 parameter

	// Execute the function multi times in parallel.
	parallelCall(func() {
		ctx := context.Background()
		score := sumCache(ctx, 3, 5)
		if score != 8 {
			t.Errorf("score should be 99, but get %d", score)
		}
		ctx = context.Background()
		sumCache(ctx, 3, 5)
		ctx = context.Background()
		sumCache(ctx, 3, 6)
	}, 5)

	if executeCount != 2 {
		t.Errorf("executeCount should be 2, but get %d", executeCount)
	}
}

// Cache Function with more parameter(>3)
func TestCacheFuncWithMoreParams(t *testing.T) {
	executeCount := 0
	type Stu struct {
		name   string
		age    int
		gender int
		height int
	}

	// Original function
	getUserScoreOrigin := func(name string, age, gender, height int) int {
		_ = age + gender + height
		executeCount++
		// select score from db where name=name and age=age and gender=gender
		switch name {
		case "Alex":
			return 10
		default:
			return 30
		}
	}

	// Convert to extra parameters to a 1 parameter(2 or 3 prameters)
	fnWrap := func(arg Stu) int {
		return getUserScoreOrigin(arg.name, arg.age, arg.gender, arg.height)
	}

	// Cacheable Function
	fnCachedInner := gofnext.CacheFn1(fnWrap)
	getUserScore := func(name string, age, gender, height int) int {
		return fnCachedInner(Stu{name, age, gender, height})
	}

	// Execute the function multi times in parallel.
	parallelCall(func() {
		score := getUserScore("Alex", 20, 1, 160)
		if score != 10 {
			t.Errorf("score should be 10, but get %d", score)
		}
		getUserScore("Jhon", 21, 0, 160)
		getUserScore("Alex", 20, 1, 160)
	}, 10)

	// Test count
	if executeCount != 2 {
		t.Errorf("executeCount should be 2, but get %d", executeCount)
	}
}
