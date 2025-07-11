package cache

import (
	"time"
)

type Cache interface {
	Set(key string, value interface{}, expireTime time.Duration)
	Get(key string, result interface{}) bool
	Delete(key string)
	Increment(key string, count int64) error
	Exists(key string) bool
}
