package bench

import (
	"testing"
	"time"

	"github.com/ahuigo/gofnext"
)
type UserInfo struct {
	Name string
	Age  int
	Desc string
}
func getUser() UserInfo {
	desc := ""
	for i := 0; i < 100; i++ {
		desc += letterBytes
	}
	time.Sleep(10 * time.Millisecond)
	return UserInfo{Name: "Alex", Age: 20, Desc: desc}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var (
	getUserWithMemCache = gofnext.CacheFn0(getUser)
	getUserWithLruCache = gofnext.CacheFn0(getUser, &gofnext.Config{
		CacheMap: gofnext.NewCacheLru(100),
	})
	getUserWithRedisCache = gofnext.CacheFn0(getUser, &gofnext.Config{
		CacheMap: gofnext.NewCacheRedis("gofnext-test-key"),
	})
)

func benchmark(b *testing.B, f func() UserInfo) {
	for i := 0; i < b.N; i++ {
		f()
	}
}

// go test -bench="Cache$" -benchmem .
func BenchmarkGetDataWithNoCache(b *testing.B) { benchmark(b, getUser) }
func BenchmarkGetDataWithMemCache(b *testing.B) { benchmark(b, getUserWithMemCache) }
func BenchmarkGetDataWithLruCache(b *testing.B) { benchmark(b, getUserWithLruCache) }
func BenchmarkGetDataWithRedisCache(b *testing.B) { benchmark(b, getUserWithRedisCache) }
