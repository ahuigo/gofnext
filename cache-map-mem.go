package gofnext

import (
	"sync"
	"time"
)

type cachedValue struct {
	val       interface{}
	createdAt time.Time
	err       error
}

type memCacheMap struct {
	*sync.Map
	// mu sync.RWMutex
	ttl      time.Duration
	errTtl   time.Duration
	reuseTtl time.Duration
}

func newCacheMapMem(ttl time.Duration) *memCacheMap {
	return &memCacheMap{
		ttl: ttl,
		Map: &sync.Map{},
	}
}

func (m *memCacheMap) Store(key, value any, err error) {
	el := cachedValue{
		val:       value,
		createdAt: time.Now(),
		err:       err,
	}
	m.Map.Store(key, &el)
}

func (m *memCacheMap) Load(key any) (value any, hasCache bool, alive bool, err error) {
	elInter, hasCache := m.Map.Load(key)
	if hasCache {
		el := elInter.(*cachedValue)
		if (m.ttl > 0 && time.Since(el.createdAt) > m.ttl) ||
			(el.err != nil && m.errTtl >= 0 && time.Since(el.createdAt) > m.errTtl) {
			if m.reuseTtl > 0 && time.Since(el.createdAt) < m.reuseTtl+m.ttl {
				// 1. cache is within reuse ttl
				return el.val, true, false, el.err
			} else {
				// 2. cache is error and within err ttl
				m.Map.Delete(key)
				return el.val, false, false, el.err
			}
		} else {
			// 3. cache is valid
			return el.val, true, true, el.err
		}
	}
	return
}

func (m *memCacheMap) SetTTL(ttl time.Duration) CacheMap {
	m.ttl = ttl
	return m
}

func (m *memCacheMap) SetErrTTL(errTTL time.Duration) CacheMap {
	m.errTtl = errTTL
	return m
}
func (m *memCacheMap) SetReuseTTL(ttl time.Duration) CacheMap {
	m.reuseTtl = ttl
	return m
}

func (m *memCacheMap) NeedMarshal() bool {
	return false
}
