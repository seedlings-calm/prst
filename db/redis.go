package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	cfg "github.com/seedlings-calm/prst/config"
	"github.com/spf13/cast"
)

type RedisSets struct {
	redis map[string]*redis.Client
	l     sync.RWMutex
}

var (
	sets *RedisSets
	once sync.Once
)

const (
	USER_INSTANCE_NAME = "user" //用户redis
)

func (r *RedisSets) getRedis(key ...string) *redis.Client {
	r.l.RLock()
	defer r.l.RUnlock()
	name := "default"
	if len(key) > 0 {
		name = key[0]
	}
	if client, ok := r.redis[name]; ok {
		return client
	}
	return nil
}

func NewRedis() *RedisSets {
	once.Do(func() {
		conf := cfg.GetGlobalConf()
		redisSets := map[string]*redis.Client{}
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, r := range conf.Redis {
			wg.Add(1)
			go func(r cfg.Redis) {
				defer wg.Done()
				client := redis.NewClient(&redis.Options{
					Addr:       fmt.Sprintf("%s:%d", r.Addr, r.Port),
					Password:   r.Password,
					DB:         cast.ToInt(r.DB),
					MaxRetries: 3, //重试次数
				})
				_, err := client.Ping(context.Background()).Result()
				if err != nil {
					// 记录错误日志而不是 panic
					log.Printf("redis初始化失败: %s", err.Error())
					return
				}
				mu.Lock()
				redisSets[r.Name] = client
				mu.Unlock()
			}(r)
		}
		wg.Wait()
		sets = &RedisSets{
			redis: redisSets,
		}
	})
	return sets
}

func getRedis(keys ...string) *redis.Client {
	rds := NewRedis()
	key := ""
	if len(keys) > 0 {
		key = keys[0]
	}
	return rds.getRedis(key)
}

func (r *RedisSets) GetUserRedis() *redis.Client {
	return getRedis(USER_INSTANCE_NAME)
}
