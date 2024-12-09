# 🛠️ Go function extended
[![tag](https://img.shields.io/github/tag/ahuigo/gofnext.svg)](https://github.com/ahuigo/gofnext/tags)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/ahuigo/gofnext?status.svg)](https://pkg.go.dev/github.com/ahuigo/gofnext)
![Build Status](https://github.com/ahuigo/gofnext/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/ahuigo/gofnext)](https://goreportcard.com/report/github.com/ahuigo/gofnext)
[![Coverage](https://img.shields.io/codecov/c/github/ahuigo/gofnext)](https://codecov.io/gh/ahuigo/gofnext)
[![Contributors](https://img.shields.io/github/contributors/ahuigo/gofnext)](https://github.com/ahuigo/gofnext/graphs/contributors)
[![License](https://img.shields.io/github/license/ahuigo/gofnext)](./LICENSE)

This **gofnext** provides the following functions extended(go>=1.21).

Cache decorators(concurrent safe): Similar to Python's `functools.cache` and `functools.lru_cache`. 

In addition to memory caching, it also supports Redis caching and custom caching.

[简体中文](/readme.zh.md)

- [🛠️ Go function extended](#️-go-function-extended)
  - [Decorator cases](#decorator-cases)
  - [Features](#features)
  - [Decorator examples](#decorator-examples)
    - [Cache fibonacii function](#cache-fibonacii-function)
    - [Cache function with 0 param](#cache-function-with-0-param)
    - [Cache function with 1 param](#cache-function-with-1-param)
    - [Cache function with 2 params](#cache-function-with-2-params)
    - [Cache function with more params(\>2)](#cache-function-with-more-params2)
    - [Cache function with lru cache](#cache-function-with-lru-cache)
    - [Cache function with redis cache(unstable)](#cache-function-with-redis-cacheunstable)
    - [Custom cache map](#custom-cache-map)
    - [Extension(pg)](#extensionpg)
  - [Decorator config](#decorator-config)
    - [Config item(`gofnext.Config`)](#config-itemgofnextconfig)
    - [Cache's Live Time(TTL)](#caches-live-timettl)
    - [Error Cache's Live Time(ErrTTl)](#error-caches-live-timeerrttl)
    - [Hash Pointer address or value?](#hash-pointer-address-or-value)
    - [Custom hash key function](#custom-hash-key-function)
  - [Roadmap](#roadmap)

## Decorator cases

| function        | decorator             |
|-----------------|-----------------------|
| func f() R    | gofnext.CacheFn0(f) |
| func f(K) R   | gofnext.CacheFn1(f) |
| func f(K1, K2) R | gofnext.CacheFn2(f) |
| func f() (R,error)    | gofnext.CacheFn0Err(f) |
| func f(T) (R,error)   | gofnext.CacheFn1Err(f)    |
| func f(T,P) (R,error) | gofnext.CacheFn2Err(f)    |
| func f() (R,error) | gofnext.CacheFn0Err(f, &gofnext.Config{TTL: time.Hour})<br/>// memory cache with ttl  |
| func f() R | gofnext.CacheFn0(f, &gofnext.Config{CacheMap: gofnext.NewCacheLru(9999)})  <br/>// Maxsize of cache is 9999|
| func f() R | gofnext.CacheFn0(f, &gofnext.Config{CacheMap: gofnext.NewCacheRedis("cacheKey")})  <br/>// Warning: redis's marshaling may result in data loss|

**Benchmark**
Benchmark case: https://github.com/ahuigo/gofnext/blob/main/bench/
```
# golang1.22
pkg: github.com/ahuigo/gofnext/bench
BenchmarkGetDataWithNoCache-10               100          11179015 ns/op          281220 B/op         99 allocs/op
BenchmarkGetDataWithMemCache-10         11036955                95.49 ns/op           72 B/op          2 allocs/op
BenchmarkGetDataWithLruCache-10         11362039               104.8 ns/op            72 B/op          2 allocs/op
BenchmarkGetDataWithRedisCache-10          15850             74653 ns/op           28072 B/op         29 allocs/op
```

## Features
- Cache Decorator (gofnext)
    - [x] Decorator cache for function
    - [x] Concurrent goroutine Safe
    - [x] Support memory CacheMap(default)
    - [x] Support memory-lru CacheMap
    - [x] Support redis CacheMap
    - [x] Support [postgres CacheMap](https://github.com/ahuigo/gofnext_pg)
    - [x] Support customization of the CacheMap(manually)
- Common functions
    - I recommend [gox](https://github.com/icza/gox), it provides ternary operator(If/IfFunc),Ptr,Pie,....


## Decorator examples
Refer to: [examples](https://github.com/ahuigo/gofnext/blob/main/examples)

### Cache fibonacii function
Refer to: [decorator fib example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-fib_test.go)
> Play: https://go.dev/play/p/7BCINKENJzA

```go
package main
import "fmt"
import "github.com/ahuigo/gofnext"
func main() {
    var fib func(int) int
    fib = func(x int) int {
        fmt.Printf("call arg:%d\n", x)
        if x <= 1 {
            return x
        } else {
            return fib(x-1) + fib(x-2)
        }
    }
    fib = gofnext.CacheFn1(fib)

    fmt.Println(fib(5))
    fmt.Println(fib(6))
}
```

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
        getUserInfoFromDbWithCache = gofnext.CacheFn0Err(getUserAnonymouse) 
    )

    func TestCacheFuncWithNoParam(t *testing.T) {
        // Execute the function multi times in parallel.
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
    	getUserInfoFromDbNil= gofnext.CacheFn1(getUserNoError) 
    )

    func TestCacheFuncNil(t *testing.T) {
    	// Execute the function multi times in parallel.
    	times := 10
    	parallelCall(func() {
    		userinfo := getUserInfoFromDbNil(20)
    		fmt.Println(userinfo)
    	}, times)
    }

### Cache function with 2 or 3 params
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
        getUserScoreWithCache := gofnext.CacheFn2Err(getUserScore, &gofnext.Config{
            TTL: time.Hour,
        }) // getFunc can only accept 2 parameter

        // Execute the function multi times in parallel.
        ctx := context.Background()
        parallelCall(func() {
            score, _ := getUserScoreWithCache(ctx, 1)
            if score != 99 {
                t.Errorf("score should be 99, but get %d", score)
            }
            getUserScoreWithCache(ctx, 2)
            getUserScoreWithCache(ctx, 3)
        }, 10)

        if executeCount != 3 {
            t.Errorf("executeCount should be 3, but get %d", executeCount)
        }
    }

### Cache function with more params(>3)
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
	fnCachedInner := gofnext.CacheFn1(fnWrap)
	fnCached := func(name string, age, gender int) int {
		return fnCachedInner(Stu{name, age, gender})
	}

	// Execute the function multi times in parallel.
	parallelCall(func() {
		score := fnCached("Alex", 20, 1)
		if score != 10 {
			t.Errorf("score should be 10, but get %d", score)
		}
		fnCached("Jhon", 21, 0)
		fnCached("Alex", 20, 1)
	}, 10)

    // Test Count
    if executeCount != 2 {
		t.Errorf("executeCount should be 2, but get %d", executeCount)
	}

### Cache function with lru cache
Refer to: [decorator lru example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-lru_test.go)

	executeCount := 0
	maxCacheSize := 2
	var getUserScore = func(more int) (int, error) {
		executeCount++
		return 98 + more, errors.New("db error")
	}

	// Cacheable Function
	var getUserScoreFromDbWithLruCache = gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		TTL:      time.Hour,
		CacheMap: gofnext.NewCacheLru(maxCacheSize),
	})

### Cache function with redis cache(unstable)
> Warning: Since redis needs JSON marshaling, this may result in data loss.

Refer to: [decorator redis example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-redis_test.go)

    var (
        // Cacheable Function
        getUserScoreWithCache = gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
            TTL:  time.Hour,
            CacheMap: gofnext.NewCacheRedis("redis-cache-key"),
        }) 
    )

    func TestRedisCacheFuncWithTTL(t *testing.T) {
        // Execute the function multi times in parallel.
        for i := 0; i < 10; i++ {
            score, _ := getUserScoreWithCache(1)
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

### Custom cache map
Refer to: https://github.com/ahuigo/gofnext/blob/main/cache-map-mem.go

### Extension(pg)
- Postgres cache extension: https://github.com/ahuigo/gofnext_pg

## Decorator config
### Config item(`gofnext.Config`)
gofnext.Config item list:

| Key | Description      |      Default       |
|-----|------------------|--------------------|
| TTL    | Cache Time to Live |0(if TTL=0:use permanent cache; if TTL>0:set cache with TTL)  |
| ErrTTL | cache TTL for error return if there is an error | 0(0:Donot cache error;  >0:Cache error with TTL; -1:rely on TTL only; )  |
| CacheMap|Custom own cache   | Inner Memory  |
| HashKeyPointerAddr | Use Pointer Addr(&p) as key instead of its value when hashing key |false(Use real value`*p` as key) |
| HashKeyFunc| Custom hash key function | Inner hash func|

### Cache's Live Time(TTL)
For example: set cache's live time to 1hour.

    gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
        /* Set cache's TTL time: 
            if TTL==0, use permanent cache; 
            if TTL>0, cache's live time is TTL
        */
        TTL:  time.Hour, 
    }) 

### Error Cache's Live Time(ErrTTl)
> By default, gofnext won't cache error when there is an error.

If there is an **error**, and you wanna control the error cache's TTL, simply add `ErrTTL: time.Duration`.
Refer to: https://github.com/ahuigo/gofnext/blob/main/examples/decorator-err_test.go

    gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
        /* Set error cache's TTL time:
            if ErrTTL=0, do not cache error;
            if ErrTTL>0, error cache's live time is ErrTTL;
            if ErrTTL=-1, error cache's live time is controlled by TTL
        */
        ErrTTL: 0, // Do not cache error(default:0)
        ErrTTL: time.Seconds * 60, // error cache's live time is 60s
        ErrTTL: -1, // rely on TTL only
    }) 

### Hash Pointer address or value?
Decorator will hash function's all parameters into hashkey.
By default, if parameter is pointer, decorator will hash its real value instead of pointer address.

If you wanna hash pointer address, you should turn on `HashKeyPointerAddr`:

	getUserScoreWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		HashKeyPointerAddr: true,
	})

### Custom hash key function
> In this case, you need to ensure that duplicate keys are not generated.
Refer to: [example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-key-custom_test.go)

	// hash key function
	hashKeyFunc := func(keys ...any) []byte{
		user := keys[0].(*UserInfo)
		flag := keys[1].(bool)
		return []byte(fmt.Sprintf("user:%d,flag:%t", user.id, flag))
	}

	// Cacheable Function
	getUserScoreWithCache := gofnext.CacheFn2Err(getUserScore, &gofnext.Config{
		HashKeyFunc: hashKeyFunc,
	})

## Roadmap
- [] Include private property when serializating for redis(#spec/reflect/unexported)
