# gocache-decorator
## Features
- [x] Decorator for function cache
- [x] Goroutine Safe
- [x] Support customization of the CacheMap
- [x] Support redis CacheMap

# Examples
> Refer to: [examples](https://github.com/ahuigo/gocache-decorator/blob/main/examples)

## CachedFunction with zero param
Refer to: [decorator example](https://github.com/ahuigo/gocache-decorator/blob/main/examples/decorator_test.go)

    package examples

    import "github.com/ahuigo/gocache-decorator"

    func getUserAnonymouse() (UserInfo, error) {
        fmt.Println("select * from db limit 1", time.Now())
        time.Sleep(10 * time.Millisecond)
        return UserInfo{Name: "Anonymous", Age: 9}, errors.New("db error")
    }

    var (
        // Cacheable Function
        getUserInfoFromDbWithCache = decorator.DecoratorFn0(getUserAnonymouse, nil) 
    )

    func TestCacheFuncWithNoParam(t *testing.T) {
        // Parallel invocation of multiple functions.
        times := 10
        parallelCall(func() {
            userinfo, err := getUserInfoFromDbWithCache()
            fmt.Println(userinfo, err)
        }, times)
    }


## CachedFunction with 2 param
> Refer to: [decorator example](https://github.com/ahuigo/gocache-decorator/blob/main/examples/decorator_test.go)

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

## CachedFunction with redis
Refer to: [decorator redis example](https://github.com/ahuigo/gocache-decorator/blob/main/examples/decorator-redis_test.go)

    var (
        // Cacheable Function
        getUserScoreFromDbWithCache = decorator.DecoratorFn1(getUserScore, &decorator.Config{
            Timeout:  time.Hour,
            CacheMap: decorator.NewRedisMap("redis-cache-key"),
        }) 
    )

    func TestRedisCacheFuncWithTTL(t *testing.T) {

        // Parallel invocation of multiple functions.
        for i := 0; i < 10; i++ {
            score, err := getUserScoreFromDbWithCache(1)
            if err != nil || score != 99 {
                t.Errorf("score should be 99, but get %d", score)
            }
        }

    }

> Warning: Since JSON marshaling cannot serialize private data, Redis will lose private data.


