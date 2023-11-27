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

type memCacheMap struct{
	*sync.Map
	// mu sync.RWMutex
	timeout time.Duration
}

func newCacheMapMem(timeout time.Duration) *memCacheMap{
	return &memCacheMap{
		timeout: timeout,
		Map: &sync.Map{},
	}
}

func (m *memCacheMap) Store(key, value any, err error) {
	el := cachedValue{
		val: value,
		createdAt: time.Now(),
		err: err,
	}
	m.Map.Store(key, &el)
}

func (m *memCacheMap) Load(key any) (value any, existed bool, err error) {
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

func (m *memCacheMap) SetTTL(timeout time.Duration) {
	m.timeout = timeout
}

func (m *memCacheMap) IsMarshalNeeded() bool {
	return false
}