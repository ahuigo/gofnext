package examples

import (
	"testing"

	"github.com/ahuigo/gofnext"
)

func TestRedisCacheMapIntFloatStruct(t *testing.T) {
	count := 0
	type Stu struct {
		Age int
	}
	type M = map[string]any
	type Data struct {
		M M
	}

	getNum := func(i float64) *Data {
		count++
		m := M{
			"i8":  int8(i),
			"u8":  uint8(i),
			"f32": float32(2),
			"stu": Stu{
				Age: 18,
			},
		}
		return &Data{M: m}
	}

	// Cacheable Function
	getNumWithCache := gofnext.CacheFn1(
		getNum,
		&gofnext.Config{
			CacheMap: gofnext.NewCacheRedis("redis-cache-key").ClearAll(),
		},
	)

	// Execute the function multi times in parallel.
	data := getNumWithCache(98)
	m := data.M

	if m["u8"].(uint8) != 98 {
		t.Errorf("u8 should be 98")
	}
	data = getNumWithCache(98)
	m = data.M
	if m["stu"].(map[string]any)["Age"].(int8) != 18 {
		t.Errorf("age should be 18")
	}
	if count != 1 {
		t.Errorf("count should be 1, but get %d", count)
	}
}
