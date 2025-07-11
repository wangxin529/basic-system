package memory

import (
	cache2 "elevate-hub/db/cache"
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"time"
)

type memoryCache struct {
	cache.Cache
}

func NewMemoryCache() cache2.Cache {
	return &memoryCache{
		*cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (m *memoryCache) Get(key string, result interface{}) bool {
	res, b := m.Cache.Get(key)
	if !b {
		return b
	}
	marshal, _ := json.Marshal(res)
	err := json.Unmarshal(marshal, &result)
	if err != nil {
		return false
	}
	return b
}

func (m *memoryCache) Exists(key string) bool {
	_, b := m.Cache.Get(key)
	return b
}
