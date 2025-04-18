package gofnext

import "time"

type CacheMap interface {
	// Goroutine concurrently on **same key**.
	Store(key, value any, err error)
	/*
	   hasCache && alive: valid cache
	   hasCache && !alive: reuse dead cache, call function
	   !hasCache: no cache, call function
	*/
	Load(key any) (value any, hasCache, alive bool, err error)
	SetTTL(ttl time.Duration) CacheMap
	SetErrTTL(ttl time.Duration) CacheMap
	SetReuseTTL(ttl time.Duration) CacheMap
	NeedMarshal() bool
}
