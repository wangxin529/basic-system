package redis

import (
	"context"
	"elevate-hub/db/cache"
	"elevate-hub/db/conf"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisCache struct {
	*redis.Client
}

func NewRedisCache(conf *conf.Redis) cache.Cache {
	if conf == nil {
		panic("redis config is nil")
	}
	var client *redis.Client
	if conf.MasterName == "" && conf.SentinelPassword == "" {
		options := &redis.Options{
			Password: conf.Password,
			Addr:     conf.Addr[0],
			DB:       conf.DB,
		}
		client = redis.NewClient(options)
	} else {
		options := &redis.FailoverOptions{
			MasterName:       conf.MasterName,
			SentinelPassword: conf.SentinelPassword,
			SentinelAddrs:    conf.Addr,
			Password:         conf.Password,
			DB:               conf.DB,
		}
		client = redis.NewFailoverClient(options)
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(errors.Wrap(err, "redis init failed"))
	}
	fmt.Println("redis int success!")
	return &redisCache{
		client,
	}
}

func (r *redisCache) Set(key string, value interface{}, expireTime time.Duration) {
	marshal, _ := json.Marshal(value)
	err := r.Client.Set(context.Background(), key, marshal, expireTime).Err()
	if err != nil {
	}

}

func (r *redisCache) Get(key string, result interface{}) bool {
	bytes, err := r.Client.Get(context.Background(), key).Bytes()
	if err != nil {
		return false
	}
	err = json.Unmarshal(bytes, &result)
	return err == nil
}

func (r *redisCache) Delete(key string) {
	r.Client.Del(context.Background(), key)
}
func (r *redisCache) Increment(key string, count int64) error {
	return r.Client.IncrBy(context.Background(), key, count).Err()
}

func (r *redisCache) Exists(key string) bool {
	val := r.Client.Exists(context.Background(), key).Val()
	return val == 1
}
