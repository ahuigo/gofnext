package gofnext

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/ahuigo/gofnext/serial"
)

type Config struct {
	/* Set cache's TTL time:
	if TTL==0, use permanent cache;
	if TTL>0, cache's live time is TTL
	*/
	TTL time.Duration
	/* Set error cache's TTL time:
	if ErrTTL=0, do not cache error;
	if ErrTTL>0, error cache's live time is ErrTTL;
	if ErrTTL=-1, error cache's live time is controlled by TTL
	*/
	ErrTTL             time.Duration
	CacheMap           CacheMap
	NeedDumpKey        bool
	HashKeyPointerAddr bool
	HashKeyFunc        func(args ...any) []byte
	/* ReuseTTl controls how to handle expired cache:
	if ReuseTTl>0: When cache is expired but within ReuseTTl duration, return the expired cache and update it asynchronously
	if ReuseTTl=0: When cache is expired, wait for the cache to be updated
	*/
	ReuseTTL time.Duration
}

type cachedFn[K1, K2, K3 any, V any] struct {
	needDumpKey        bool
	hashKeyPointerAddr bool
	hashKeyFunc        func(args ...any) []byte
	cacheMap           CacheMap
	pkeyLockMap        sync.Map
	keyLen             int
	getFunc            func(K1, K2, K3) (V, error)
}

func (c *cachedFn[K1, K2, K3, V]) setConfigs(configs ...*Config) *cachedFn[K1, K2, K3, V] {
	if len(configs) > 0 {
		return c.setConfig(configs[0])
	} else {
		return c.setConfig(nil)
	}
}

func (c *cachedFn[K1, K2, K3, V]) setConfig(config *Config) *cachedFn[K1, K2, K3, V] {
	// default value
	if config == nil {
		config = &Config{}
	}
	if config.CacheMap == nil {
		config.CacheMap = newCacheMapMem(config.TTL)
	}

	// init value
	c.hashKeyPointerAddr = config.HashKeyPointerAddr
	c.needDumpKey = config.NeedDumpKey
	c.cacheMap = config.CacheMap
	if config.ErrTTL < -1 {
		panic("ErrTTL should not be less than -1")
	}
	if config.TTL < 0 {
		panic("TTL should not be less than 0")
	}
	c.cacheMap.SetErrTTL(config.ErrTTL)
	c.cacheMap.SetReuseTTL(config.ReuseTTL)
	if config.TTL > 0 {
		c.cacheMap.SetTTL(config.TTL)
	}
	// init hashKeyFuncMethod
	if config.HashKeyFunc != nil {
		c.hashKeyFunc = config.HashKeyFunc
	} else {
		cacheMapRefV := reflect.ValueOf(c.cacheMap)
		methodValue := cacheMapRefV.MethodByName("HashKeyFunc")

		if methodValue.IsValid() {
			c.hashKeyFunc = methodValue.Interface().(func(...any) []byte)
			// c.hashKeyFunc = func(keys ...any) []byte {
			// 	reflectKeys := make([]reflect.Value, len(keys))
			// 	for i, key := range keys {
			// 		reflectKeys[i] = reflect.ValueOf(key)
			// 	}
			// 	result := methodValue.Call(reflectKeys)
			// 	if len(result) > 0 {
			// 		if bytes, ok := result[0].Interface().([]byte); ok {
			// 			return bytes
			// 		}
			// 	}
			// 	return nil
			// }
		}
	}
	return c
}

// Cache Function with 3 parameter(with error)
func CacheFn3Err[K1 any, K2 any, K3 any, V any](
	getFunc func(K1, K2, K3) (V, error),
	configs ...*Config,
) func(K1, K2, K3) (V, error) {
	getFunc0 := func(k1 K1, k2 K2, k3 K3) (V, error) {
		return getFunc(k1, k2, k3)
	}
	ins := &cachedFn[K1, K2, K3, V]{getFunc: getFunc0, keyLen: 3}
	ins.setConfigs(configs...)
	return ins.invoke3err
}

// Cache Function with 3 parameter
func CacheFn3[K1 any, K2 any, K3 any, V any](
	getFunc func(K1, K2, K3) V,
	configs ...*Config,
) func(K1, K2, K3) V {
	getFunc0 := func(k1 K1, k2 K2, k3 K3) (V, error) {
		return getFunc(k1, k2, k3), nil
	}
	ins := &cachedFn[K1, K2, K3, V]{getFunc: getFunc0, keyLen: 3}
	ins.setConfigs(configs...)
	return ins.invoke3
}

// Cache Function with 2 parameter(with error)
func CacheFn2Err[K1 any, K2 any, V any](
	getFunc func(K1, K2) (V, error),
	configs ...*Config,
) func(K1, K2) (V, error) {
	getFunc0 := func(k1 K1, k2 K2, k3 int8) (V, error) {
		return getFunc(k1, k2)
	}
	ins := &cachedFn[K1, K2, int8, V]{getFunc: getFunc0, keyLen: 2}
	ins.setConfigs(configs...)
	return ins.invoke2err
}

// Cache Function with 2 parameter
func CacheFn2[K1 any, K2 any, V any](
	getFunc func(K1, K2) V,
	configs ...*Config,
) func(K1, K2) V {
	getFunc0 := func(k1 K1, k2 K2, k3 any) (V, error) {
		return getFunc(k1, k2), nil
	}
	ins := &cachedFn[K1, K2, any, V]{getFunc: getFunc0, keyLen: 2}
	ins.setConfigs(configs...)
	return ins.invoke2
}

// Cache Function with 1 parameter(with error)
func CacheFn1Err[K any, V any](
	getFunc func(K) (V, error),
	configs ...*Config,
) func(K) (V, error) {
	getFunc0 := func(key K, k2 context.Context, k3 any) (V, error) {
		return getFunc(key)
	}
	ins := &cachedFn[K, context.Context, any, V]{getFunc: getFunc0, keyLen: 1}
	ins.setConfigs(configs...)
	x := ins.invoke1err
	return x
}

// Cache Function with 1 parameter
func CacheFn1[K any, V any](
	getFunc func(K) V,
	configs ...*Config,
) func(K) V {
	getFunc0 := func(k1 K, k2 context.Context, k3 int8) (V, error) {
		return getFunc(k1), nil
	}
	ins := &cachedFn[K, context.Context, int8, V]{getFunc: getFunc0, keyLen: 1}
	ins.setConfigs(configs...)
	return ins.invoke1
}

// Cache Function with 0 parameter(with error)
func CacheFn0Err[V any](
	getFunc func() (V, error),
	configs ...*Config,
) func() (V, error) {
	getFunc0 := func(ctx context.Context, i int8, a byte) (V, error) {
		return getFunc()
	}
	ins := &cachedFn[context.Context, int8, byte, V]{getFunc: getFunc0, keyLen: 0}
	ins.setConfigs(configs...)
	return ins.invoke0err
}

// Cache Function with 0 parameter
func CacheFn0[V any](
	getFunc func() V,
	configs ...*Config,
) func() V {
	getFunc0 := func(ctx context.Context, i int8, a byte) (V, error) {
		return getFunc(), nil
	}
	ins := &cachedFn[context.Context, int8, byte, V]{getFunc: getFunc0, keyLen: 0}
	ins.setConfigs(configs...)
	return ins.invoke0
}

// Invoke cached function with no parameter
func (c *cachedFn[any, int, A, V]) invoke0() V {
	var k1 any
	var k2 int
	var k3 A
	retv, _ := c.invoke3err(k1, k2, k3)
	return retv
}

// Invoke cached function with no parameter(with error)
func (c *cachedFn[any, int8, A, V]) invoke0err() (V, error) {
	var k1 any
	var k2 int8
	// k2 = 0                                    // error: cannot use 0 (untyped int constant) as uint8 value in assignment
	var k3 A
	return c.invoke3err(k1, k2, k3)
}

// Invoke cached function with 1 parameter
func (c *cachedFn[K1, K2, K3, V]) invoke1(k1 K1) V {
	var k2 K2
	var k3 K3
	val, _ := c.invoke3err(k1, k2, k3)
	return val
}
func (c *cachedFn[K1, K2, K3, V]) invoke1err(k1 K1) (V, error) {
	var k2 K2
	var k3 K3
	return c.invoke3err(k1, k2, k3)
}

// Invoke cached function with 2 parameter
func (c *cachedFn[K1, K2, A, V]) invoke2(key1 K1, key2 K2) (retv V) {
	var a A
	retv, _ = c.invoke3err(key1, key2, a)
	return retv
}
func (c *cachedFn[K1, K2, K3, V]) invoke2err(k1 K1, k2 K2) (V, error) {
	var k3 K3
	return c.invoke3err(k1, k2, k3)
}

func (c *cachedFn[K1, K2, K3, V]) invoke3(key1 K1, key2 K2, key3 K3) (retv V) {
	retv, _ = c.invoke3err(key1, key2, key3)
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

func (c *cachedFn[K1, K2, K3, V]) hashKeyFuncWrap(key1 K1, key2 K2, key3 K3) (pkey any) {
	// outer hash key func
	if c.hashKeyFunc != nil {
		if c.keyLen == 3 {
			pkey = string(c.hashKeyFunc(key1, key2, key3))
		} else if c.keyLen == 2 {
			pkey = string(c.hashKeyFunc(key1, key2))
		} else if c.keyLen == 1 {
			pkey = string(c.hashKeyFunc(key1))
		} else {
			pkey = 0
		}
		return pkey
	}

	// inner hash key func
	needHashPtrAddr := c.hashKeyPointerAddr
	needDumpKey := c.needDumpKey
	if c.keyLen == 3 {
		if _, hasCtx := any(key1).(context.Context); hasCtx {
			pkey = [2]any{key2, key3}
			if !needDumpKey {
				needDumpKey = !isHashableKey(key2, needHashPtrAddr) || !isHashableKey(key3, needHashPtrAddr)
			}
		} else {
			pkey = [3]any{key1, key2, key3}
			if !needDumpKey {
				needDumpKey = !isHashableKey(key1, needHashPtrAddr) || !isHashableKey(key2, needHashPtrAddr) || !isHashableKey(key3, needHashPtrAddr)
			}
		}
	} else if c.keyLen == 2 {
		if _, hasCtx := any(key1).(context.Context); hasCtx {
			pkey = key2
			if !needDumpKey {
				needDumpKey = !isHashableKey(key2, needHashPtrAddr)
			}
		} else {
			pkey = [2]any{key1, key2}
			if !needDumpKey {
				needDumpKey = !isHashableKey(key1, needHashPtrAddr) || !isHashableKey(key2, needHashPtrAddr)
			}
		}
	} else if c.keyLen == 1 {
		if _, hasCtx := any(key1).(context.Context); hasCtx {
			pkey = 0
		} else {
			pkey = key1
			if !needDumpKey {
				needDumpKey = !isHashableKey(key1, needHashPtrAddr)
			}
		}
	} else {
		pkey = 0
	}
	if needDumpKey {
		pkey = serial.String(pkey, needHashPtrAddr)
	}
	return pkey
}

// Invoke cached function with 2 parameter
func (c *cachedFn[K1, K2, K3, V]) invoke3err(key1 K1, key2 K2, key3 K3) (retv V, err error) {
	// 1. generate pkey
	var pkey any = c.hashKeyFuncWrap(key1, key2, key3)

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
	value, hasCache, alive, err := c.cacheMap.Load(pkey)
	pkeyLock.RUnlock()

	// 3.1 check if marshal needed
	if hasCache && c.cacheMap.NeedMarshal() {
		err2 := json.Unmarshal(value.([]byte), &retv)
		if err == nil {
			err = err2
		}
		value = &retv
		// return retv, err
	}

	// 4. Execute getFunc(only once)
	if !hasCache {
		// 4.1 try lock
		// If multiple goroutines call the same function at the same time,
		// only one goroutine can execute the getFunc.
		isLocked := pkeyLock.TryLock()

		if !isLocked {
			if checkCacheCount < 3 {
				// wait for other goroutine to finish
				time.Sleep(time.Millisecond * 10)

				// Avoid all goroutines calling TryCount simultaneously, which may lead to failure.
				sleepRandom(0, 50*time.Millisecond)

				// try check again
				goto checkCache
			}
			// if checkCacheCount >= 3, run lock and getFunc
			pkeyLock.Lock()
		}
		defer pkeyLock.Unlock()

		// 4.2 check cache again
		val, err2 := c.getFunc(key1, key2, key3)
		c.cacheMap.Store(pkey, &val, err2)
		return val, err2
	} else if hasCache && !alive {
		// If the cache is not alive,  it will return the expired cache (and update the cache asynchronously)
		// 5.1 try lock
		// If 100 goroutines call the same function at the same time,
		// only one goroutine can execute the getFunc.
		go func() {
			isLocked := pkeyLock.TryLock()
			if !isLocked {
				return
			}
			defer pkeyLock.Unlock()
			// 5.2 check cache again
			val, err2 := c.getFunc(key1, key2, key3)
			c.cacheMap.Store(pkey, &val, err2)
		}()

	}
	return *(value).(*V), err
}
