# 🛠️ Go function extended
[![标签](https://img.shields.io/github/tag/ahuigo/gofnext.svg)](https://github.com/ahuigo/gofnext/tags)
![Go 版本](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/ahuigo/gofnext?status.svg)](https://pkg.go.dev/github.com/ahuigo/gofnext)
![构建状态](https://github.com/ahuigo/gofnext/actions/workflows/test.yml/badge.svg)
[![Go 报告](https://goreportcard.com/badge/github.com/ahuigo/gofnext)](https://goreportcard.com/report/github.com/ahuigo/gofnext)
[![覆盖率](https://img.shields.io/codecov/c/github/ahuigo/gofnext)](https://codecov.io/gh/ahuigo/gofnext)
[![贡献者](https://img.shields.io/github/contributors/ahuigo/gofnext)](https://github.com/ahuigo/gofnext/graphs/contributors)
[![许可证](https://img.shields.io/github/license/ahuigo/gofnext)](./LICENSE)

这个 **gofnext** 提供以下函数扩展（go>=1.21）。

缓存装饰器（并发安全）：类似于 Python 的 `functools.cache` 和 `functools.lru_cache`。除了内存缓存，它也支持 Redis 缓存和自定义缓存。

- [🛠️ Go function extended](#️-go-function-extended)
  - [装饰器cases](#装饰器cases)
  - [特性](#特性)
  - [装饰器示例](#装饰器示例)
    - [缓存斐波那契函数](#缓存斐波那契函数)
    - [带有0个参数的缓存函数](#带有0个参数的缓存函数)
    - [带有1个参数的缓存函数](#带有1个参数的缓存函数)
    - [带有2个参数的缓存函数](#带有2个参数的缓存函数)
    - [带有2个以上参数的缓存函数](#带有2个以上参数的缓存函数)
    - [带LRU 缓存的函数](#带lru-缓存的函数)
    - [带redis缓存的函数(unstable)](#带redis缓存的函数unstable)
    - [定制缓存函数](#定制缓存函数)
  - [装饰器配置](#装饰器配置)
    - [配置项清单(`gofnext.Config`)](#配置项清单gofnextconfig)
    - [缓存时间](#缓存时间)
    - [如果有error就不缓存](#如果有error就不缓存)
    - [哈希指针地址还是值？](#哈希指针地址还是值)
    - [自定义哈希键函数](#自定义哈希键函数)
  - [Roadmap](#roadmap)

[Egnlish](/)/[中文]()

## 装饰器cases

| 函数             | 装饰器                      |
|------------------|-----------------------------|
| func f() R     | gofnext.CacheFn0(f)         |
| func f(K1) R    | gofnext.CacheFn1(f)         |
| func f(K1, K2) R  | gofnext.CacheFn2(f)         |
| func f() (R, error)     | gofnext.CacheFn0Err(f)      |
| func f(K1) (R, error)    | gofnext.CacheFn1Err(f)      |
| func f(K1, K2) (R, error)  | gofnext.CacheFn2Err(f)      |
| func f() (R,error) | gofnext.CacheFn0Err(f, &gofnext.Config{TTL: time.Hour})<br/>// 带有 ttl 的内存缓存  |
| func f() (R) | gofnext.CacheFn0(f, &gofnext.Config{CacheMap: gofnext.NewCacheLru(9999)})  <br/>// 缓存的最大大小为 9999|
| func f() (R) | gofnext.CacheFn0(f, &gofnext.Config{CacheMap: gofnext.NewCacheRedis("cacheKey")})  <br/>// 警告：redis 的序列化可能会导致数据丢失|

## 特性
- [x] 缓存装饰器 (gofnext)
    - [x] 函数的装饰器缓存
    - [x] 并发协程安全
    - [x] 支持内存 CacheMap（默认）
    - [x] 支持内存-LRU CacheMap
    - [x] 支持 redis CacheMap
    - [x] 手动支持自定义 CacheMap

## 装饰器示例
参考：[示例](https://github.com/ahuigo/gofnext/blob/main/examples)

### 缓存斐波那契函数
参考：[装饰器斐波那契示例](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-fib_test.go)

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

### 带有0个参数的缓存函数
参考: [decorator example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator_test.go)

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

### 带有1个参数的缓存函数
参考: [decorator example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-nil_test.go)

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

### 带有2个参数的缓存函数
> 参考: [decorator example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator_test.go)

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

        // Execute the function multi times in parallel.
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

### 带有2个以上参数的缓存函数
参考: [decorator example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator_test.go)

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

### 缓存method
### Cache method
Original method:

    type User Struct{}
    func NewUser() (*User){
        u := &User{}
        return u
    }

    func (u *User) getUserInfo(uid int) (*User, error){  
        ...  
    }

Cached method(singleton mode):

    type User Struct{}
    var getUserInfoCached func(uid int) (&User, error)
    func NewUser() *User{
        u := &User{}
       	getUserInfoCached = gofnext.CacheFn1Err(u.getUserInfo) 
        return u
    }

    // proxy method with cached
    func (u *User) GetUserInfo(uid int) (*User, error){  
        return getUserInfoCached(uid)
    }

    func (u *User) getUserInfo(uid int) (*User, error){  
        ...  
    }

Cached method(factory mode):

    type User Struct{
        getUserInfoCached func(uid int) (&User, error)
    }
    func NewUser() *User{
        u := &User{}
       	u.getUserInfoCached = gofnext.CacheFn1Err(u.getUserInfo) 
        return u
    }
    ....

### 带LRU 缓存的函数
参考: [decorator lru example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-lru_test.go)

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

### 带redis缓存的函数(unstable)
> 警告: 目前使用json序列化,可能会有私有属性丢失
> 后续序列化方法可能会有变化, 请不要用于生产

参考: [decorator redis example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-redis_test.go)

    var (
        // Cacheable Function
        getUserScoreFromDbWithCache = gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
            TTL:  time.Hour,
            CacheMap: gofnext.NewCacheRedis("redis-cache-key"),
        }) 
    )

    func TestRedisCacheFuncWithTTL(t *testing.T) {
        // Execute the function multi times in parallel.
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

### 定制缓存函数
参考: https://github.com/ahuigo/gofnext/blob/main/cache-map-mem.go

## 装饰器配置
### 配置项清单(`gofnext.Config`)
gofnext.Config 清单:

| 键 | 描述             |默认                |
|-----|------------------|
| TTL    | 缓存时间 | 0(不过期)|
| ErrTTL| 控制error返回的缓存时间|0(0:不缓存error; >0:缓存error有ErrTTL限制；-1: 只依赖TTL)  |
| CacheMap| 自定义缓存map |默认内存Map|
| HashKeyPointerAddr | 哈希key时，使用指针本身地址(&p)，而不是实际的值 |默认使用pointer指向实际值(*p)|
| HashKeyFunc| 自定义哈希键函数 |内置hashFunc|

### 缓存时间
e.g.

    gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
        TTL:  time.Hour,
    }) 

### 如果有error就不缓存
> 默认有函数返回error时, 就不会用缓存.

如果存在error时, 也需要缓存的话。 参考: https://github.com/ahuigo/gofnext/blob/main/examples/decorator-err_test.go

    gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
        ErrTTL<=0, // 不会缓存error
        ErrTTL: time.Seconds * 60, // err缓存的errTTL 是60秒
    }) 

### 哈希指针地址还是值？
> 装饰器将函数的所有参数哈希成哈希键（hashkey）。 默认情况下，如果参数是指针，装饰器将哈希其实际值而不是指针地址。

如果您想要哈希指针地址，您应该打开 `HashKeyPointerAddr` 选项：

	getUserScoreFromDbWithCache := gofnext.CacheFn1Err(getUserScore, &gofnext.Config{
		HashKeyPointerAddr: true,
	})

### 自定义哈希键函数
> 这种情况下，您需要保证不会有生成重复的key。

参考: [example](https://github.com/ahuigo/gofnext/blob/main/examples/decorator-key-custom_test.go)

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

## Roadmap
- [] Redis CacheMap 支持序列化所有私有属性
