# Go function extended
This **gofnext** provides the following functions extended. 
- Cache decorators: Similar to Python's `functools.cache` and `functools.lru_cache`. 
    - > Additionally, it supports Redis caching and custom caching.
- Dump: Deep dumper for golang data, it will dump pointer's real value and struct's inner private data.

TOC 
- [Go function extended](#go-function-extended)
  - [Decorator cases](#decorator-cases)
  - [Features](#features)
  - [Decorator examples](#decorator-examples)
    - [Cache fibonacii function](#cache-fibonacii-function)
    - [CachedFunction with zero param](#cachedfunction-with-zero-param)
    - [CachedFunction with 1 param and no error](#cachedfunction-with-1-param-and-no-error)
    - [CachedFunction with 2 param](#cachedfunction-with-2-param)
    - [CachedFunction with lru cache](#cachedfunction-with-lru-cache)
    - [CachedFunction with redis cache](#cachedfunction-with-redis-cache)
  - [Object functions](#object-functions)
  - [Dump](#dump)

## Decorator cases

| function        | decorator             |
|-----------------|-----------------------|
| func f() res    | gofnext.CacheFn0(f,nil) |
| func f(a) res   | gofnext.CacheFn1(f,nil) |
| func f(a,b) res | gofnext.CacheFn2(f,nil) |
| func f() (res,err)    | gofnext.CacheFn0Err(f,nil) |
| func f(a) (res,err)   | gofnext.CacheFn1Err(f,nil)    |
| func f(a,b) (res,err) | gofnext.CacheFn2Err(f,nil)    |
| func f() (res,err) | gofnext.CacheFn0Err(f, &gofnext.Config{TTL: time.Hour})<br/>// memory cache with ttl  |
| func f() (res) | gofnext.CacheFn0(f, &gofnext.Config{CacheMap: gofnext.NewCacheLru(9999)})  <br/>// Maxsize of cache is 9999|
| func f() (res) | gofnext.CacheFn0(f, &gofnext.Config{CacheMap: gofnext.NewCacheRedis("cacheKey", nil)})  <br/>// Warning: redis's marshaling may result in data loss|

## Features
- [x] Cache Decorator (gofnext)
    - [x] Decorator cache for function
    - [x] Goroutine Safe
    - [x] Support memory CacheMap(default)
    - [x] Support memory-lru CacheMap
    - [x] Support redis CacheMap
    - [x] Support customization of the CacheMap(manually)
- [x] Dump (gofnext/dump)
- [x] Object (gofnext/object)
    - [x] `ConvertObjectByte2String(any)`: Convert Object bytes to string

## Decorator examples
Refer to: [examples](https://github.com/ahuigo/gofnext/blob/main/examples)


### Cache fibonacii function
Refer to: [decorator fib example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-fib_test.go)

    package main
    import "fmt"
    import "github.com/ahuigo/gofnext"
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

        fibCached = gofnext.CacheFn1(fib, nil)    

        fmt.Println(fibCached(5))
        fmt.Println(fibCached(6))
    }

### CachedFunction with zero param
Refer to: [decorator example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator_test.go)

    package examples

    import "github.com/ahuigo/gofnext"

    func getUserAnonymouse() (UserInfo, error) {
        fmt.Println("select * from db limit 1", time.Now())
        time.Sleep(10 * time.Millisecond)
        return UserInfo{Name: "Anonymous", Age: 9}, errors.New("db error")
    }

    var (
        // Cacheable Function
        getUserInfoFromDbWithCache = gofnext.CacheFn0Err(getUserAnonymouse, nil) 
    )

    func TestCacheFuncWithNoParam(t *testing.T) {
        // Parallel invocation of multiple functions.
        times := 10
        parallelCall(func() {
            userinfo, err := getUserInfoFromDbWithCache()
            fmt.Println(userinfo, err)
        }, times)
    }


### CachedFunction with 1 param and no error
Refer to: [decorator example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-nil_test.go)

    func getUserNoError(age int) (UserInfo) {
    	time.Sleep(10 * time.Millisecond)
    	return UserInfo{Name: "Alex", Age: age}
    }
    
    var (
    	// Cacheable Function with 1 param and no error
    	getUserInfoFromDbNil= gofnext.CacheFn1(getUserNoError, nil) 
    )

    func TestCacheFuncNil(t *testing.T) {
    	// Parallel invocation of multiple functions.
    	times := 10
    	parallelCall(func() {
    		userinfo := getUserInfoFromDbNil(20)
    		fmt.Println(userinfo)
    	}, times)
    }

### CachedFunction with 2 param
> Refer to: [decorator example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator_test.go)

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

### CachedFunction with lru cache
Refer to: [decorator lru example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-lru_test.go)

### CachedFunction with redis cache
Refer to: [decorator redis example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-redis_test.go)

    var (
        // Cacheable Function
        getUserScoreFromDbWithCache = gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
            TTL:  time.Hour,
            CacheMap: gofnext.NewCacheRedis("redis-cache-key", nil),
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

You could
> Warning: Since redis needs JSON marshaling, this may result in data loss.


## Object functions
Refer to: [object example](https://github.com/ahuigo/gofnext/blob/main/examples/object_test.go)

	import "github.com/ahuigo/gofnext/object"

    func TestConvertMapBytes(t *testing.T) {
        objBytes := map[string][]byte{
            "k1": []byte("v1"),
            "k2": []byte("v2"),
        }
        out, _ := json.Marshal(objBytes)
        fmt.Println(string(out))                 //output: {"k1":"djE=","k2":"djI="}

        objString := object.ConvertObjectByte2String(objBytes)
        out, _ = json.Marshal(objString)
        fmt.Println(string(out))                 //output: {"k1":"v1","k2":"v2"}
    }

## Dump 
Refer to: [dump example](https://github.com/ahuigo/gofnext/blob/main/examples/dump_test.go)

Dump any value to string(include private field)

    type Person struct {
        Name string
        age  int //private
    }
	p := &person
	expectedP := "&Person:{Name:\"John Doe\",age:30}"
	if result := dump.String(p); result != expectedP {
		t.Errorf("Expected %s, but got %s", expectedP, result)
	}