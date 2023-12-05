# Go function extended(go>=1.21)
[![tag](https://img.shields.io/github/tag/ahuigo/gofnext.svg)](https://github.com/ahuigo/gofnext/tags)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/ahuigo/gofnext?status.svg)](https://pkg.go.dev/github.com/ahuigo/gofnext)
![Build Status](https://github.com/ahuigo/gofnext/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/ahuigo/gofnext)](https://goreportcard.com/report/github.com/ahuigo/gofnext)
[![Coverage](https://img.shields.io/codecov/c/github/ahuigo/gofnext)](https://codecov.io/gh/ahuigo/gofnext)
[![Contributors](https://img.shields.io/github/contributors/ahuigo/gofnext)](https://github.com/ahuigo/gofnext/graphs/contributors)
[![License](https://img.shields.io/github/license/ahuigo/gofnext)](./LICENSE)

This **gofnext** provides the following functions extended. 
- Cache decorators: Similar to Python's `functools.cache` and `functools.lru_cache`. 
    - Additionally, it supports Redis caching and custom caching.

TOC 
- [Go function extended(go\>=1.21)](#go-function-extendedgo121)
  - [Decorator cases](#decorator-cases)
  - [Features](#features)
  - [Decorator examples](#decorator-examples)
    - [Cache fibonacii function](#cache-fibonacii-function)
    - [Cache function with 0 param](#cache-function-with-0-param)
    - [Cache function with 1 param](#cache-function-with-1-param)
    - [Cache function with 2 params](#cache-function-with-2-params)
    - [Cache function with more params(\>2)](#cache-function-with-more-params2)
    - [Cache function with lru cache](#cache-function-with-lru-cache)
    - [Cache function with redis cache](#cache-function-with-redis-cache)
    - [Hash Pointer address or value?](#hash-pointer-address-or-value)
    - [Custom hash key function](#custom-hash-key-function)

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

### Cache function with 0 param
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


### Cache function with 1 param
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

### Cache function with 2 params
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

### Cache function with more params(>2)
Refer to: [decorator example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator_test.go)

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

### Cache function with lru cache
Refer to: [decorator lru example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-lru_test.go)

### Cache function with redis cache
> Warning: Since redis needs JSON marshaling, this may result in data loss.

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

To avoid keys being too long, you can limit the length of Redis key:

    cacheMap := gofnext.NewCacheRedis("redis-cache-key").SetMaxHashKeyLen(256);

Set redis config:

	// method 1: by default: localhost:6379
	cache := gofnext.NewCacheRedis("redis-cache-key") 

	// method 2: set redis addr
	cache.SetRedisAddr("192.168.1.1:6379")

	// method 3: set redis options
	cache.SetRedisOpts(&redis.Options{
		Addr: "localhost:6379",
	})

	// method 4: set redis universal options
	cache.SetRedisUniversalOpts(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})


### Hash Pointer address or value?
Decorator will hash all function's parameters into hashkey.
By default, if parameter is pointer, decorator will hash its real value instead of pointer address.

If you wanna hash pointer address, you should turn on `HashKeyPointerAddr`:

	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		HashKeyPointerAddr: true,
	})

### Custom hash key function
> You need to resolve key collisions yourself.
Refer to: [example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-key-custom_test.go)

	// hash key function
	hashKeyFunc := func(keys ...any) []byte{
		user := keys[0].(*UserInfo)
		flag := keys[1].(bool)
		return []byte(fmt.Sprintf("user:%d,flag:%t", user.id, flag))
	}

	// Cacheable Function
	getUserScoreFromDbWithCache := gofnext.CacheFn2Err(getUserScore, &gofnext.Config{
		HashKeyFunc: hashKeyFunc,
	})
