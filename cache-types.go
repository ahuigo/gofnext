package gofnext

import "time"

type CacheMap interface {
	// Goroutine concurrently on **same key**.
	Store(key, value any, err error)
	Load(key any) (value any, existed bool, err error)
	SetTTL(ttl time.Duration) CacheMap
	NeedMarshal() bool
}
