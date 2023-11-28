package gofnext

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"crypto/sha512"
	"encoding/hex"

	"github.com/go-redis/redis"
)

type redisMap struct {
	mu          sync.Mutex
	redisClient redis.UniversalClient
	ttl         time.Duration
	redisPreKey string
	maxHashKeyLen int
}

type redisData struct {
	Data      []byte
	Err       []byte
	CreatedAt time.Time
	TTL       time.Duration
}

func NewCacheRedis(mapKey string, config *redis.UniversalOptions) *redisMap {
	if mapKey == "" {
		panic("mapKey can not be empty")
	}
	if config == nil {
		redisAddr := "redis:6379"
		config = &redis.UniversalOptions{
			Addrs: []string{redisAddr},
			DB:    0,
		}
	}
	redisClient := redis.NewUniversalClient(config)
	return &redisMap{
		redisClient: redisClient,
		redisPreKey: mapKey,
		maxHashKeyLen: 2000,
	}
}

func (m *redisMap) ClearAll() *redisMap {
	m.redisClient.Del(m.redisPreKey)
	return m
}


func (m *redisMap) strkey(key any) string {
	var r string
	switch rt := key.(type) {
	case string:
		r = rt	
	default:
		r = fmt.Sprintf("%#v", key)
	}
	if len(r) > m.maxHashKeyLen{
		hash := sha512.Sum512([]byte(r))
		r = hex.EncodeToString(hash[:])
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
