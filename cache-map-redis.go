package gofnext

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"strconv"
	"sync"
	"time"

	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"

	"github.com/ahuigo/gofnext/dump"
	"github.com/go-redis/redis"
)

type redisMap struct {
	mu            sync.Mutex
	redisClient   redis.UniversalClient
	ttl           time.Duration
	redisPreKey   string
	maxHashKeyLen int
}

type redisData struct {
	Data      []byte
	Err       []byte
	CreatedAt time.Time
	TTL       time.Duration
}

func NewCacheRedis(mapKey string) *redisMap {
	if mapKey == "" {
		panic("mapKey can not be empty")
	}
	redisAddr := "localhost:6379"
	config := &redis.UniversalOptions{
		Addrs: []string{redisAddr},
		DB:    0,
	}
	redisClient := redis.NewUniversalClient(config)
	return &redisMap{
		redisClient: redisClient,
		redisPreKey: mapKey,
	}
}

func (m *redisMap) SetRedisAddr(addr string) *redisMap {
	m.redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return m
}

func (m *redisMap) SetRedisOpts(opts *redis.Options) *redisMap {
	m.redisClient = redis.NewClient(opts)
	return m
}

func (m *redisMap) SetRedisUniversalOpts(opts *redis.UniversalOptions) *redisMap {
	m.redisClient = redis.NewUniversalClient(opts)
	return m
}

func (m *redisMap) ClearAll() *redisMap {
	m.redisClient.Del(m.redisPreKey)
	return m
}

func (m *redisMap) HashKeyFunc(key ...any) []byte {
	if len(key) == 0 {
		return nil
	} else if len(key) == 1 {
		return dump.Bytes(key[0], false)
	} else {
		return dump.Bytes(key, false)
	}
}

func (m *redisMap) strkey(key any) string {
	var r string
	switch rt := key.(type) {
	case string:
		r = rt
	default:
		r = dump.String(key, false)
	}
	if m.maxHashKeyLen > 0 && len(r) > m.maxHashKeyLen {
		if m.maxHashKeyLen <= 8 {
			h := fnv.New64a()
			_, _ = h.Write([]byte(r))
			return strconv.FormatUint(h.Sum64(), 16)
		} else if m.maxHashKeyLen <= 32 {
			hash := md5.Sum([]byte(r))
			r = hex.EncodeToString(hash[:])
		} else if m.maxHashKeyLen <= 64 {
			hash := sha512.Sum512_256([]byte(r))
			r = hex.EncodeToString(hash[:])
		} else {
			hash := sha512.Sum512([]byte(r))
			r = hex.EncodeToString(hash[:])
		}
	}
	return r
}

func (m *redisMap) Store(key, value any, err error) {
	pkey := m.strkey(key)
	data, _ := json.Marshal(value)
	cacheData := redisData{
		Data: data,
		TTL:  m.ttl,
	}
	if m.ttl > 0 {
		cacheData.CreatedAt = time.Now()
	}
	if err != nil {
		cacheData.Err = []byte(err.Error())
	}
	val, _ := json.Marshal(cacheData)
	m.redisClient.HSet(m.redisPreKey, pkey, val).Err()
}

func (m *redisMap) Load(key any) (value any, existed bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pkey := m.strkey(key)
	val, err := m.redisClient.HGet(m.redisPreKey, pkey).Bytes()
	// m.redisClient.TTL()
	if err == redis.Nil {
		existed = false
		err = nil
		return
	} else if err != nil {
		return
	}
	cacheData := redisData{}
	err = json.Unmarshal(val, &cacheData)
	if err != nil {
		return
	}

	value = cacheData.Data
	if cacheData.Err != nil {
		err = errors.New(string(cacheData.Err))
	}
	if cacheData.TTL > 0 && time.Since(cacheData.CreatedAt) > cacheData.TTL {
		return value, false, nil //expired
	}
	existed = true
	return
}

func (m *redisMap) SetTTL(ttl time.Duration) CacheMap {
	m.ttl = ttl
	return m
}

func (m *redisMap) SetMaxHashKeyLen(l int) *redisMap {
	m.maxHashKeyLen = l
	return m
}

func (m *redisMap) NeedMarshal() bool {
	return true
}
