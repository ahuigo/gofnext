package decorator

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type cachedObjType struct {
	val       interface{}
	createdAt time.Time
	err       error
}
type Config struct {
	Timeout time.Duration
}

type cachedFn[Ctx any, K any, V any] struct {
	mu             sync.RWMutex
	cacheMap       sync.Map
	routineOnceMap sync.Map
	timeout        time.Duration
	keyLen         int
	getFunc        func(Ctx, K) (V, error)
}

// Cache Function with ctx and 1 parameter
func DecoratorFn2[Ctx any, K any, V any](
	getFunc func(Ctx, K) (V, error),
	config *Config,
) func(Ctx, K) (V, error) {
	ins := &cachedFn[Ctx, K, V]{getFunc: getFunc, keyLen: 2}
	if config != nil {
		ins.timeout = config.Timeout
	}
	return ins.invoke2
}

// Cache Function with 1 parameter
func DecoratorFn1[K any, V any](
	getFunc func(K) (V, error),
	config *Config,
) func(K) (V, error) {
	getFunc0 := func(ctx context.Context, key K) (V, error) {
		return getFunc(key)
	}
	ins := &cachedFn[context.Context, K, V]{getFunc: getFunc0, keyLen: 1}
	if config != nil {
		ins.timeout = config.Timeout
	}
	return ins.invoke1
}

// Cache Function with no parameter
func DecoratorFn0[V any](
	getFunc func() (V, error),
	config *Config,
) func() (V, error) {
	getFunc0 := func(ctx context.Context, i int8) (V, error) {
		return getFunc()
	}
	ins := &cachedFn[context.Context, int8, V]{getFunc: getFunc0, keyLen: 0}
	if config != nil {
		ins.timeout = config.Timeout
	}
	return ins.invoke0
}

// Invoke cached function with no parameter
func (c *cachedFn[any, int, V]) invoke0() (V, error) {
	var ctx any
	var key int
	// key = 0                                    // error: cannot use 0 (untyped int constant) as uint8 value in assignment
	fmt.Printf("cache key: %#v, %T\n", key, key) // cache key: 0, uint8
	return c.invoke2(ctx, key)
}

// Invoke cached function with 1 parameter
func (c *cachedFn[Ctx, K, V]) invoke1(key K) (V, error) {
	var ctx Ctx
	return c.invoke2(ctx, key)
}

func (c *cachedFn[Ctx, K, V]) SetConfig(config Config) *cachedFn[Ctx, K, V] {
	c.mu.Lock()
	c.timeout = config.Timeout
	c.mu.Unlock()
	return c
}

// Invoke cached function with 2 parameter
func (c *cachedFn[Ctx, K, V]) invoke2(key1 Ctx, key2 K) (V, error) {
	// pkey
	var pkey any = key2
	if _, hasCtx := any(key1).(context.Context); hasCtx || c.keyLen <= 1 {
		// ignore context key
		kind := reflect.TypeOf(key2).Kind()
		if kind == reflect.Map || kind == reflect.Slice {
			pkey = fmt.Sprintf("%#v", key2)
		}
	} else {
		pkey = fmt.Sprintf("%#v,%#v", key1, key2)
	}

	// check cache
	needRefresh := false
	value, hasCache := c.cacheMap.Load(pkey)
	if hasCache {
		cachedObj := value.(*cachedObjType)
		if c.timeout > 0 && time.Since(cachedObj.createdAt) > c.timeout {
			needRefresh = true
		}
	}

	if !hasCache || needRefresh {
		var tmpOnce sync.Once
		oncePtr := &tmpOnce
		//1. clean up routineOnceMap key
		if needRefresh {
			c.routineOnceMap.Delete(pkey)
		}
		// 2. load or store routineOnceMap key
		onceInterface, loaded := c.routineOnceMap.LoadOrStore(pkey, oncePtr)
		if loaded {
			oncePtr = onceInterface.(*sync.Once)
		}
		// 3. Execute getFunc(only once)
		oncePtr.Do(func() {
			val, err := c.getFunc(key1, key2)
			createdAt := time.Now()
			c.cacheMap.Store(pkey, &cachedObjType{val: &val, err: err, createdAt: createdAt})
		})
		value, _ = c.cacheMap.Load(pkey)
	}
	cachedObj := value.(*cachedObjType)
	return *(cachedObj.val).(*V), cachedObj.err
}

/*
func (c *cachedFn[string, V]) Get0() (V, error) {
	// var s any
	var s string
	// s = "abc" // error: cannot use "abc" (untyped string constant) as string value in assignment
	fmt.Printf("cache key: %#v, %T\n", s, s)
	return c.Get(s)
}
*/

/*
func (c *cachedFn[int, V]) Get0() (V, error) {
	var s int = 100 //error: cannot use 100 (untyped int constant) as int value in variable declaration
	fmt.Printf("cache key: %#v, %T\n", s, s)
	return c.Get(s)
}
*/
