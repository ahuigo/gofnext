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

var (
	// Cacheable Function with 1 param and no error
	getUserInfoFromDb = gofnext.CacheFn0(getUser)
)

func TestCacheFuncWith0Param(t *testing.T) {
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

// Cache Function with more parameter(>2)
func TestCacheFuncWithMoreParams(t *testing.T) {
	executeCount := 0
	type Stu struct {
		name   string
		age    int
		gender int
	}

	// Original function
	fn := func(name string, age, gender int) int {
		executeCount++
		// select score from db where name=name and age=age and gender=gender
		switch name {
		case "Alex":
			return 10
		default:
			return 30
		}
	}

	// Convert to extra parameters to a single parameter(2 prameters is ok)
	fnWrap := func(arg Stu) int {
		return fn(arg.name, arg.age, arg.gender)
	}

	// Cacheable Function
	fnCached := gofnext.CacheFn1(fnWrap)

	// Execute the function multi times in parallel.
	parallelCall(func() {
		score := fnCached(Stu{"Alex", 20, 1})
		if score != 10 {
			t.Errorf("score should be 10, but get %d", score)
		}
		fnCached(Stu{"Jhon", 21, 0})
		fnCached(Stu{"Alex", 20, 1})
	}, 10)

	if executeCount != 2 {
		t.Errorf("executeCount should be 2, but get %d", executeCount)
	}
}
