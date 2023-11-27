# gocache-decorator
## Features
- [x] Decorator for function cache
- [x] Goroutine Safe
- [x] Support customization of the CacheMap
- [x] Support redis CacheMap

# Examples
> Refer to: [examples](https://github.com/ahuigo/gocache-decorator/blob/main/examples)

## CachedFunction with zero param
Refer to: [cache-function](https://github.com/ahuigo/gocache-decorator/blob/main/examples/cache-function_test.go)

    func getUserAnonymouse() (UserInfo, error) {
        fmt.Println("select * from db limit 1", time.Now())
        time.Sleep(10 * time.Millisecond)
        return UserInfo{Name: "Anonymous", Age: 9}, errors.New("db error")
    }

    func TestCacheFuncWithNoParam(t *testing.T) {
        // Cacheable Function
        getUserInfoFromDbWithCache := decorator.DecoratorFn0(getUserAnonymouse, nil) 

        // Parallel invocation of multiple functions.
        times := 10
        parallelCall(func() {
            userinfo, err := getUserInfoFromDbWithCache()
            fmt.Println(userinfo, err)
        }, times)
    }


## CachedFunction with 2 param
Refer to: [cache-function](https://github.com/ahuigo/gocache-decorator/blob/main/examples/cache-function_test.go)

    func TestCacheFuncWith2Param(t *testing.T) {
        // Original function
        executeCount := 0
        getUserScore := func(c context.Context, id int) (int, error) {
            executeCount++
            fmt.Println("select score from db where id=", id, time.Now())
            time.Sleep(10 * time.Millisecond)
            return 98 + id, errors.New("db error")
        }

        // Cacheable Function
        getUserScoreFromDbWithCache := decorator.DecoratorFn2(getUserScore, &decorator.Config{
            Timeout: time.Hour,
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
