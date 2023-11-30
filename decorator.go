package gofnext

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/ahuigo/gofnext/dump"
)

type Config struct {
	TTL                   time.Duration
	CacheMap              CacheMap
	NeedDumpKey           bool
	NeedCmpKeyPointerAddr bool
}

type cachedFn[K1 any, K2 any, V any] struct {
	needDumpKey           bool
	needCmpKeyPointerAddr bool
	cacheMap              CacheMap
	pkeyLockMap           sync.Map
	keyLen                int
	getFunc               func(K1, K2) (V, error)
}

func (c *cachedFn[K1, K2, V]) setConfig(config *Config) *cachedFn[K1, K2, V] {
	// default value
	if config == nil {
		config = &Config{}
	}
	if config.CacheMap == nil {
		config.CacheMap = newCacheMapMem(config.TTL)
	}

	// init value
	c.cacheMap = config.CacheMap
	c.needCmpKeyPointerAddr = config.NeedCmpKeyPointerAddr
	c.needDumpKey = config.NeedDumpKey
	if config.TTL > 0 {
		c.cacheMap.SetTTL(config.TTL)
	}
	return c
}

// Cache Function with 2 parameter
func CacheFn2Err[K1 any, K2 any, V any](
	getFunc func(K1, K2) (V, error),
	config *Config,
) func(K1, K2) (V, error) {
	ins := &cachedFn[K1, K2, V]{getFunc: getFunc, keyLen: 2}
	ins.setConfig(config)
	return ins.invoke2err
}

// Cache Function with 2 parameter
func CacheFn2[K1 any, K2 any, V any](
	getFunc func(K1, K2) V,
	config *Config,
) func(K1, K2) V {
	getFunc0 := func(ctx K1, key K2) (V, error) {
		return getFunc(ctx, key), nil
	}
	ins := &cachedFn[K1, K2, V]{getFunc: getFunc0, keyLen: 2}
	ins.setConfig(config)
	return ins.invoke2
}

// Cache Function with 1 parameter
func CacheFn1Err[K any, V any](
	getFunc func(K) (V, error),
	config *Config,
) func(K) (V, error) {
	getFunc0 := func(ctx context.Context, key K) (V, error) {
		return getFunc(key)
	}
	ins := &cachedFn[context.Context, K, V]{getFunc: getFunc0, keyLen: 1}
	ins.setConfig(config)
	return ins.invoke1
}

func CacheFn1[K any, V any](
	getFunc func(K) V,
	config *Config,
) func(K) V {
	getFunc0 := func(ctx context.Context, key K) (V, error) {
		return getFunc(key), nil
	}
	ins := &cachedFn[context.Context, K, V]{getFunc: getFunc0, keyLen: 1}
	ins.setConfig(config)
	return ins.invoke1err
}

// Cache Function with 0 parameter
func CacheFn0Err[V any](
	getFunc func() (V, error),
	config *Config,
) func() (V, error) {
	getFunc0 := func(ctx context.Context, i int8) (V, error) {
		return getFunc()
	}
	ins := &cachedFn[context.Context, int8, V]{getFunc: getFunc0, keyLen: 0}
	ins.setConfig(config)
	return ins.invoke0err
}

// Cache Function with 0 parameter
func CacheFn0[V any](
	getFunc func() V,
	config *Config,
) func() V {
	getFunc0 := func(ctx context.Context, i int8) (V, error) {
		return getFunc(), nil
	}
	ins := &cachedFn[context.Context, int8, V]{getFunc: getFunc0, keyLen: 0}
	ins.setConfig(config)
	return ins.invoke0
}

// Invoke cached function with no parameter
func (c *cachedFn[any, int, V]) invoke0err() (V, error) {
	var ctx any
	var key int
	// key = 0                                    // error: cannot use 0 (untyped int constant) as uint8 value in assignment
	return c.invoke2err(ctx, key)
}
func (c *cachedFn[any, int, V]) invoke0() V {
	var ctx any
	var key int
	retv, _ := c.invoke2err(ctx, key)
	return retv
}

// Invoke cached function with 1 parameter
func (c *cachedFn[Ctx, K, V]) invoke1(key K) (V, error) {
	var ctx Ctx
	return c.invoke2err(ctx, key)
}

func (c *cachedFn[Ctx, K, V]) invoke1err(key K) V {
	var ctx Ctx
	val, _ := c.invoke2err(ctx, key)
	return val
}

func (c *cachedFn[K1, K2, V]) invoke2(key1 K1, key2 K2) (retv V) {
	retv, _ = c.invoke2err(key1, key2)
	return retv
}

var _isHashKey map[any]int

func isHashableKey(key any, cmpPtr bool) (canHash bool) {
	defer func() {
		if err := recover(); err != nil {
			canHash = false
		}
	}()
	_ = _isHashKey[key]
	if cmpPtr {
		return true
	}
	return reflect.ValueOf(key).Kind() != reflect.Pointer
}

// Invoke cached function with 2 parameter
func (c *cachedFn[Ctx, K, V]) invoke2err(key1 Ctx, key2 K) (retv V, err error) {
	// 1. generate pkey
	var pkey any
	cmpPtr := c.needCmpKeyPointerAddr
	if c.keyLen == 2 {
		if _, hasCtx := any(key1).(context.Context); hasCtx {
			pkey = key2
			if !c.needDumpKey {
				c.needDumpKey = !isHashableKey(key2, cmpPtr)
			}
		} else {
			pkey = [2]any{key1, key2}
			if !c.needDumpKey {
				c.needDumpKey = !isHashableKey(key1, cmpPtr) || !isHashableKey(key2, cmpPtr)
			}
		}
	} else if c.keyLen == 1 {
		if _, hasCtx := any(key2).(context.Context); hasCtx {
			pkey = 0
		} else {
			pkey = key2
			if !c.needDumpKey {
				c.needDumpKey = !isHashableKey(key2, cmpPtr)
			}
		}
	} else {
		pkey = 0
	}
	if c.needDumpKey {
		pkey = dump.String(pkey, c.needCmpKeyPointerAddr)
	}

	// 2. require lock for each pkey(go routine safe)
	var tmpOnce sync.RWMutex
	pkeyLock := &tmpOnce
	pkeyLockInter, loaded := c.pkeyLockMap.LoadOrStore(pkey, pkeyLock)
	if loaded {
		pkeyLock = pkeyLockInter.(*sync.RWMutex)
	}

	// 3. check cache
	checkCacheCount := 0
checkCache:
	checkCacheCount++
	pkeyLock.RLock()
	value, hasCache, err := c.cacheMap.Load(pkey)
	pkeyLock.RUnlock()

	// 3.1 check if marshal needed
	if hasCache && c.cacheMap.NeedMarshal() {
		err = json.Unmarshal(value.([]byte), &retv)
		return retv, err
	}

	// 4. Execute getFunc(only once)
	if !hasCache {
		// 4.1 try lock
		// If 100 goroutines call the same function at the same time,
		// only one goroutine can execute the getFunc.
		isLocked := pkeyLock.TryLock()

		if !isLocked {
			if checkCacheCount < 3 {
				// wait for other goroutine to finish
				time.Sleep(time.Millisecond * 10)

				// Avoid all goroutines calling TryCount simultaneously, which may lead to failure.
				sleepRandom(0, 50*time.Millisecond)

				// try lock again
				goto checkCache
			}
			pkeyLock.Lock()
		}
		defer pkeyLock.Unlock()

		// 4.2 check cache again
		val, err := c.getFunc(key1, key2)
		c.cacheMap.Store(pkey, &val, err)
		return val, err
	}
	return *(value).(*V), err
}
