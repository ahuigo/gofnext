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
	ttl    time.Duration
	errTtl time.Duration
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

func (m *memCacheMap) Load(key any) (value any, existed bool, err error) {
	elInter, existed := m.Map.Load(key)
	if existed {
		el := elInter.(*cachedValue)
		if ( m.ttl > 0 && time.Since(el.createdAt) > m.ttl) ||
		(el.err != nil && m.errTtl >= 0 && time.Since(el.createdAt) > m.errTtl) {
			m.Map.Delete(key)
			existed = false
		} else {
			return el.val, existed, el.err
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

func (m *memCacheMap) NeedMarshal() bool {
	return false
}
