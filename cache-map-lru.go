package gofnext

import (
	"container/list"
	"sync"
	"time"
)

type cachedNode struct {
	val       interface{}
	createdAt time.Time
	err       error
	element   *list.Element
}

type cacheLru struct {
	list     *list.List
	listMap  *sync.Map
	maxSize  int
	mu       sync.RWMutex
	ttl      time.Duration
	errTtl   time.Duration
	resueTtl time.Duration
}

func NewCacheLru(maxSize int) *cacheLru {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &cacheLru{
		maxSize: maxSize,
		list:    list.New(),
		listMap: &sync.Map{},
	}
}

func (m *cacheLru) Store(key, value any, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	el := cachedNode{
		val:       value,
		createdAt: time.Now(),
		err:       err,
	}
	if m.maxSize > 0 && m.list.Len() >= m.maxSize {
		elInter := m.list.Back()
		m.list.Remove(elInter)
		m.listMap.Delete(elInter.Value)
	}
	el.element = m.list.PushFront(key)
	m.listMap.Store(key, &el)
}

func (m *cacheLru) Load(key any) (value any, hasCache, alive bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	elInter, hasCache := m.listMap.Load(key)
	if hasCache {
		el := elInter.(*cachedNode)
		if (m.ttl > 0 && time.Since(el.createdAt) > m.ttl) ||
			(el.err != nil && m.errTtl >= 0 && time.Since(el.createdAt) > m.errTtl) {
			// 1. cache is within reuse ttl
			if m.resueTtl > 0 && time.Since(el.createdAt) < m.resueTtl+m.ttl {
				return el.val, true, false, el.err
			} else {
				// 2. cache is not valid
				m.listMap.Delete(key)
				m.list.Remove(el.element)
				return el.val, false, false, el.err
			}
		} else {
			// 3. cache is valid: move to front
			m.list.MoveToFront(el.element)
			return el.val, true, true, el.err
		}
	}
	return
}

func (m *cacheLru) SetTTL(ttl time.Duration) CacheMap {
	m.ttl = ttl
	return m
}
func (m *cacheLru) SetErrTTL(errTTL time.Duration) CacheMap {
	m.errTtl = errTTL
	return m
}
func (m *cacheLru) SetReuseTTL(errTTL time.Duration) CacheMap {
	m.resueTtl = errTTL
	return m
}

func (m *cacheLru) NeedMarshal() bool {
	return false
}
