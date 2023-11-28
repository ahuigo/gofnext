# gocache-decorator
Cache Decorator for Go functions, similar to Python's `functools.cache` and `functools.lru_cache`. \
Additionally, it supports Redis caching and custom caching.

## Features
- [x] Decorator cache for function
- [x] Goroutine Safe
- [x] Support memory CacheMap(default)
- [x] Support memory-lru CacheMap
- [x] Support redis CacheMap
- [x] Support customization of the CacheMap(manually)

# Examples
> Refer to: [examples](https://github.com/ahuigo/gocache-decorator/blob/main/examples)

| function        | decorator             |
|-----------------|-----------------------|
| func f() res    | decorator.CacheFn0(f,nil) |
| func f(a) res   | decorator.CacheFn1(f,nil) |
| func f(a,b) res | decorator.CacheFn2(f,nil) |
| func f() (res,err)    | decorator.CacheFn0Err(f,nil) |
| func f(a) (res,err)   | decorator.CacheFn1Err(f,nil)    |
| func f(a,b) (res,err) | decorator.CacheFn2Err(f,nil)    |
| func f() (res,err) | decorator.CacheFn0Err(f, &decorator.Config{TTL: time.Hour})<br/>// memory cache with ttl  |
| func f() (res) | decorator.CacheFn0(f, &decorator.Config{CacheMap: decorator.NewCacheLru(9999)})  <br/>// Maxsize of cache is 9999|
| func f() (res) | decorator.CacheFn0(f, &decorator.Config{CacheMap: decorator.NewCacheRedis("cacheKey", nil)})  <br/>// Warning: redis's marshaling may result in data loss|

## Cache fibonacii function
Refer to: [decorator fib example](https://github.com/ahuigo/gocache-decorator/blob/main/examples/fib_test.go)

    package main
    import "fmt"
    import "github.com/ahuigo/gocache-decorator"
    func main() {
        var fib func(int) int
        var fibCached func(int) int
        fib = func(x int) int {
            fmt.Printf("call arg:%d\n", x)
            if x <= 1 {
                return x
            } else {
                return fibCached(x-1) + fibCached(x-2)
            }
        }

        fibCached = decorator.CacheFn1(fib, nil)    

        fmt.Println(fibCached(5))
        fmt.Println(fibCached(6))
    }

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
        getUserInfoFromDbWithCache = decorator.CacheFn0Err(getUserAnonymouse, nil) 
    )

    func TestCacheFuncWithNoParam(t *testing.T) {
        // Parallel invocation of multiple functions.
        times := 10
        parallelCall(func() {
            userinfo, err := getUserInfoFromDbWithCache()
            fmt.Println(userinfo, err)
        }, times)
    }


## CachedFunction with 1 param and no error
Refer to: [decorator example](https://github.com/ahuigo/gocache-decorator/blob/main/examples/decorator-nil_test.go)

    func getUserNoError(age int) (UserInfo) {
    	time.Sleep(10 * time.Millisecond)
    	return UserInfo{Name: "Alex", Age: age}
    }
    
    var (
    	// Cacheable Function with 1 param and no error
    	getUserInfoFromDbNil= decorator.CacheFn1(getUserNoError, nil) 
    )

    func TestCacheFuncNil(t *testing.T) {
    	// Parallel invocation of multiple functions.
    	times := 10
    	parallelCall(func() {
    		userinfo := getUserInfoFromDbNil(20)
    		fmt.Println(userinfo)
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
        getUserScoreFromDbWithCache := decorator.CacheFn2Err(getUserScore, &decorator.Config{
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

## CachedFunction with lru cache
Refer to: [decorator lru example](https://github.com/ahuigo/gocache-decorator/blob/main/examples/decorator-lru_test.go)

## CachedFunction with redis cache
Refer to: [decorator redis example](https://github.com/ahuigo/gocache-decorator/blob/main/examples/decorator-redis_test.go)

    var (
        // Cacheable Function
        getUserScoreFromDbWithCache = decorator.CacheFn1Err(getUserScore, &decorator.Config{
            TTL:  time.Hour,
            CacheMap: decorator.NewCacheRedis("redis-cache-key", nil),
        }) 
    )

    func TestRedisCacheFuncWithTTL(t *testing.T) {

        // Parallel invocation of multiple functions.
        for i := 0; i < 10; i++ {
            score, _ := getUserScoreFromDbWithCache(1)
            if score != 99 {
                t.Errorf("score should be 99, but get %d", score)
            }
        }

    }

> Warning: Since redis needs JSON marshaling, this may result in data loss.


