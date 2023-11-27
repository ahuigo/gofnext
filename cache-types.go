package decorator

import "time"

type CacheMap interface {
	// It's safe to call Store() concurrently on **same key**.
	// It need lock if operate on **different keys** concurrently
	Store(key, value any, err error)

	// Need lock always
	Load(key any) (value any, existed bool, err error)
	SetTTL(timeout time.Duration)
}
