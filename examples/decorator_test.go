package examples

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)

func getUser(age int) UserInfo {
	time.Sleep(10 * time.Millisecond)
	return UserInfo{Name: "Alex", Age: age}
}

var (
	// Cacheable Function with 1 param and no error
	getUserInfoFromDb = gofnext.CacheFn1(getUser, nil)
)

func TestCacheFuncWith0Param(t *testing.T) {
	// Parallel invocation of multiple functions.
	times := 10
	parallelCall(func() {
		userinfo := getUserInfoFromDb(20)
		fmt.Println(userinfo)
	}, times)
}

func TestCacheFuncWith2Params(t *testing.T) {
	// Original function
	executeCount := 0
	getUserScore := func(c context.Context, id int) (int, error) {
		executeCount++
		fmt.Println("select score from db where id=", id, time.Now())
		time.Sleep(10 * time.Millisecond)
		return 98 + id, errors.New("db error")
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn2Err(getUserScore, &gofnext.Config{
		TTL: time.Hour,
	}) // getFunc can only accept 2 parameter

	// Parallel invocation of multiple functions.
	ctx := context.Background()
	parallelCall(func() {
		score, _ := getUserScoreFromDbWithCache(ctx, 1)
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
	fnCached := gofnext.CacheFn1(fnWrap, nil)

	// Parallel invocation of multiple functions.
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
