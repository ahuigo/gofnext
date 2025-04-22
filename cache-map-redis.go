package gofnext

import (
	"errors"
	"hash/fnv"
	"strconv"
	"sync"
	"time"

	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"

	"github.com/ahuigo/gofnext/serial"
	"github.com/go-redis/redis"
)

type redisMap struct {
	mu            sync.Mutex
	redisClient   redis.UniversalClient
	ttl           time.Duration
	errTtl        time.Duration
	reuseTtl      time.Duration
	redisFuncKey  string
	maxHashKeyLen int
}

type redisData struct {
	Data      []byte
	Err       []byte
	CreatedAt time.Time
	// TTL       time.Duration
}

func NewCacheRedis(funcKey string) *redisMap {
	if funcKey == "" {
		panic("NewCacheRedis: funcKey cannot be empty")
	}
	redisAddr := "localhost:6379"
	config := &redis.UniversalOptions{
		Addrs: []string{redisAddr},
		DB:    0,
	}
	redisClient := redis.NewUniversalClient(config)
	return &redisMap{
		redisClient:  redisClient,
		redisFuncKey: "_gofnext:" + funcKey,
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
	m.redisClient.Del(m.redisFuncKey)
	return m
}

func (m *redisMap) HashKeyFunc(key ...any) []byte {
	if len(key) == 0 {
		return nil
	} else if len(key) == 1 {
		return serial.Bytes(key[0], false)
	} else {
		return serial.Bytes(key, false)
	}
}

func (m *redisMap) strkey(key any) string {
	var r string
	switch rt := key.(type) {
	case string:
		r = rt
	default:
		r = serial.String(key, false)
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

func (m *redisMap) Store(key, value any, err0 error) {
	buf, err := marshalMsgpack(value)
	if err != nil {
		slogger.Error("gofnext.redisMap: marshal", "err", err.Error())
		return
	}

	pkey := m.strkey(key)
	// data, _ := json.Marshal(value)
	cacheData := redisData{
		Data: buf,
		// TTL:  m.ttl,
	}
	if err0 != nil && m.errTtl <= 0 {
		// do not cache error
		return
	}
	if m.ttl > 0 || m.errTtl >= 0 {
		cacheData.CreatedAt = time.Now()
	}
	if err0 != nil {
		cacheData.Err = []byte(err0.Error())
	}

	// encode value to bytes
	if buf, err = marshalMsgpack(cacheData); err != nil {
		slogger.Error("gofnext.redisMap", "err", err.Error())
		return
	}
	// buf, _ := json.Marshal(cacheData)
	err = m.redisClient.HSet(m.redisFuncKey, pkey, buf).Err()
	if err != nil {
		slogger.Error("gofnext.redisMap", "err", err.Error())
	}
}

func (m *redisMap) Load(key any) (value any, hasCache, alive bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pkey := m.strkey(key)
	val, err := m.redisClient.HGet(m.redisFuncKey, pkey).Bytes()
	// m.redisClient.TTL()
	if err == redis.Nil {
		hasCache = false
		err = nil
		return
	} else if err != nil {
		return
	}
	cacheData := redisData{}
	err = unmarshalMsgpack(val, &cacheData)
	if err != nil {
		slogger.Error("gofnext.redisMap:decode", "err", err.Error())
		return
	}

	value = cacheData.Data
	if cacheData.Err != nil {
		err = errors.New(string(cacheData.Err))
	}
	if (m.ttl > 0 && time.Since(cacheData.CreatedAt) > m.ttl) ||
		(m.errTtl >= 0 && cacheData.Err != nil && time.Since(cacheData.CreatedAt) > m.errTtl) {
		// 1. cache is within reuse ttl
		if m.reuseTtl > 0 && time.Since(cacheData.CreatedAt) < m.reuseTtl+m.ttl {
			return value, true, false, nil
		} else {
			// 2. cache is not valid
			m.redisClient.HDel(m.redisFuncKey, pkey)
			return value, false, false, nil
		}
	} else {
		// 3. cache is valid
		return value, true, true, nil
	}
}

func (m *redisMap) SetTTL(ttl time.Duration) CacheMap {
	m.ttl = ttl
	return m
}
func (m *redisMap) SetErrTTL(errTTL time.Duration) CacheMap {
	m.errTtl = errTTL
	return m
}
func (m *redisMap) SetReuseTTL(ttl time.Duration) CacheMap {
	m.reuseTtl = ttl
	return m
}

func (m *redisMap) SetMaxHashKeyLen(l int) *redisMap {
	m.maxHashKeyLen = l
	return m
}

func (m *redisMap) NeedMarshal() bool {
	return true
}
