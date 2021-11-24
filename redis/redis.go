package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type MapRedisCache struct {
	instances map[string]*redis.Client
}

func NewMapRedisCache() *MapRedisCache {
	return &MapRedisCache{
		instances: make(map[string]*redis.Client),
	}
}

func (c MapRedisCache) Close() {
	for _, r := range c.instances {
		_ = r.Close()
	}
}

// AddRedisInstance add a redis instance into poll
// note that, all instance should be added in init status of Application
func (c MapRedisCache) AddRedisInstance(key string, addr string, port string, pwd string, dbNum int) error {
	if _, ok := c.instances[key]; !ok {
		redisDB := redis.NewClient(&redis.Options{
			Addr:       addr + ":" + port,
			Password:   pwd,
			DB:         dbNum,
			MaxRetries: 2, // retry 3 times (<=MaxRetries)
			PoolSize:   1024,
		})

		if _, err := redisDB.Ping(context.Background()).Result(); err == nil {
			c.instances[key] = redisDB
		} else {
			return err
		}

	} else {
		return errors.New("repeated key")
	}

	return nil
}

func (c MapRedisCache) GetRedisInstance(key string) (*redis.Client, bool) {
	r, ok := c.instances[key]
	return r, ok
}

// -----------------------------------------------------------------------------

var (
	defaultMapCache *MapRedisCache
)

func init() {
	defaultMapCache = NewMapRedisCache()
}

func AddRedisInstance(key string, addr string, port string, pwd string, dbNum int) error {
	return defaultMapCache.AddRedisInstance(key, addr, port, pwd, dbNum)
}

func GetRedisInstance(key string) (*redis.Client, bool) {
	return defaultMapCache.GetRedisInstance(key)
}

func MustGetRedisInstance(instance ...string) *redis.Client {
	name := ""
	if len(instance) == 1 {
		name = instance[0]
	}
	redisInstance, ok := GetRedisInstance(name)
	if !ok {
		panic("redis instance [" + name + "] not exists")
	}

	return redisInstance
}

func ReleaseRedisPool() {
	defaultMapCache.Close()
}
