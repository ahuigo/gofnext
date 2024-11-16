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
	list    *list.List
	listMap *sync.Map
	maxSize int
	mu      sync.RWMutex
	ttl     time.Duration
	errTtl  time.Duration
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

func (m *cacheLru) Load(key any) (value any, existed bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	elInter, existed := m.listMap.Load(key)
	if existed {
		el := elInter.(*cachedNode)
		if (m.ttl > 0 && time.Since(el.createdAt) > m.ttl) ||
			(el.err != nil && m.errTtl >= 0 && time.Since(el.createdAt) > m.errTtl) {
			m.listMap.Delete(key)
			m.list.Remove(el.element)
			existed = false
		} else {
			// move to front
			m.list.MoveToFront(el.element)
			return el.val, existed, el.err
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

func (m *cacheLru) NeedMarshal() bool {
	return false
}
