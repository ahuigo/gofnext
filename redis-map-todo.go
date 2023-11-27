package decorator

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RedisMap struct {
	mu sync.Mutex
	redisClient redis.UniversalClient 
	ttl time.Duration
	redisPreKey string
}


func NewRedisMap() *RedisMap{
	redisAddr := "redis:6379"
	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{redisAddr},
		DB:    0,
	})
	return &RedisMap{
		redisClient: redisClient,
		redisPreKey: "cachemap",
	}
}

func (m *RedisMap) Store(key, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pkey := m.strkey(key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = m.redisClient.Set(pkey, data, m.ttl).Err()
	return err
}

func (m *RedisMap) strkey(key any) string {
	return fmt.Sprintf("%s:%#v", m.redisPreKey, key)
}

func (m *RedisMap) Load(key any) (value any, existed bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pkey := m.strkey(key)
	value, err = m.redisClient.Get(pkey).Bytes()
	// m.redisClient.TTL()
	if err == redis.Nil	{
		existed = false
		err = nil
		return
	}else if err != nil {
		return
	}
	existed = true
	return
}
