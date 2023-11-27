package decorator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type redisMap struct {
	mu sync.Mutex
	redisClient redis.UniversalClient 
	ttl time.Duration
	redisPreKey string
}

type redisData struct{
	Data []byte
	Err []byte
	CreatedAt time.Time
	TTL time.Duration
}


func NewRedisMap(mapKey string) *redisMap{
	if mapKey == "" {
		panic("mapKey can not be empty")
	}
	redisAddr := "redis:6379"
	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{redisAddr},
		DB:    0,
	})
	return &redisMap{
		redisClient: redisClient,
		redisPreKey: mapKey,
	}
}

func (m *redisMap) ClearAll() *redisMap {
	m.redisClient.Del(m.redisPreKey)
	return m
}

func (m *redisMap) strkey(key any) string {
	r := fmt.Sprintf("%#v", key)
	return r
}

func (m *redisMap) Store(key, value any, err error) {
	pkey := m.strkey(key)
	data, _ := json.Marshal(value)
	cacheData := redisData{
		Data: data,
		TTL: m.ttl,
	}
	if m.ttl>0{
		cacheData.CreatedAt = time.Now()
	}
	if err != nil {
		cacheData.Err = []byte(err.Error())
	}
	val,_ := json.Marshal(cacheData)
	m.redisClient.HSet(m.redisPreKey, pkey, val).Err()
}


func (m *redisMap) Load(key any) (value any, existed bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pkey := m.strkey(key)
	val, err := m.redisClient.HGet(m.redisPreKey,pkey).Bytes()
	// m.redisClient.TTL()
	if err == redis.Nil	{
		existed = false
		err = nil
		return
	}else if err != nil {
		return
	}
	cacheData := redisData{}
	err=json.Unmarshal(val, &cacheData)
	if err != nil {
		return
	}

	value = cacheData.Data
	if cacheData.Err != nil {
		err = errors.New(string(cacheData.Err))
	}
	if cacheData.TTL>0 && time.Since(cacheData.CreatedAt) > cacheData.TTL {
		return value, false, nil //expired
	}
	existed = true
	return
}

func (m *redisMap) SetTTL(timeout time.Duration) {
	m.ttl = timeout
}

func (m *redisMap) IsMarshalNeeded() bool {
	return true
}