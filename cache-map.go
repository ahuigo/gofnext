package decorator

import (
	"sync"
	"time"
)
type cachedValue struct{
	val       interface{}
	createdAt time.Time
	err       error
}

type cacheMap struct{
	*sync.Map
	mu sync.RWMutex
	timeout time.Duration
}

func NewCacheMap(timeout time.Duration) *cacheMap{
	return &cacheMap{
		timeout: timeout,
		Map: &sync.Map{},
	}
}

func (m *cacheMap) Store(key, value any, err error) {
	el := cachedValue{
		val: value,
		createdAt: time.Now(),
		err: err,
	}
	m.Map.Store(key, &el)
}

func (m *cacheMap) Load(key any) (value any, existed bool, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	elInter, existed := m.Map.Load(key)
	if existed {
		el := elInter.(*cachedValue)
		if m.timeout>0 && time.Since(el.createdAt) > m.timeout {
			m.Map.Delete(key)
			existed = false
		}else{
			return el.val, existed, el.err
		}
	}
	return
}

func (m *cacheMap) SetTTL(timeout time.Duration) {
	m.timeout = timeout
}