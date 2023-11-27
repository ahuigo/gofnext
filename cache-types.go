package decorator

import "time"

type CacheMap interface {
	// Goroutine concurrently on **same key**.
	Store(key, value any, err error)
	Load(key any) (value any, existed bool, err error)
	SetTTL(timeout time.Duration)
	IsMarshalNeeded() bool
}
