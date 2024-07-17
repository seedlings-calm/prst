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
}

var (
	sets *RedisSets
	once sync.Once
)

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
					Addr:       fmt.Sprintf("%s:%s", r.Addr, r.Port),
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
