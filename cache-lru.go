package decorator

import (
	"container/list"
	"fmt"
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
	timeout time.Duration
}

func NewCacheLru(size int, timeout time.Duration) *cacheLru {
	if size <= 0 {
		size = 100
	}
	return &cacheLru{
		timeout: timeout,
		maxSize: size,
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
	fmt.Println(m.maxSize, m.list.Len(), key)
	if m.maxSize > 0 && m.list.Len() >= m.maxSize {
		elInter := m.list.Back()
		m.list.Remove(elInter)
		m.listMap.Delete(elInter.Value)
	}
	el.element = m.list.PushFront(key)
	m.listMap.Store(key, &el)
}

func (m *cacheLru) Load(key any) (value any, existed bool, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	elInter, existed := m.listMap.Load(key)
	if existed {
		el := elInter.(*cachedNode)
		if m.timeout > 0 && time.Since(el.createdAt) > m.timeout {
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

func (m *cacheLru) SetTTL(timeout time.Duration) {
	m.timeout = timeout
}
