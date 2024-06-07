package db

import (
	"context"

	"github.com/redis/go-redis/v9"
	cfg "github.com/seedlings-calm/prst/config"
)

var Redis redis.UniversalClient

func RedisInit() {
	var client redis.UniversalClient
	// 使用集群模式
	if cfg.Config.Redis.UseCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    cfg.Config.Redis.ClusterAddrs,
			Password: cfg.Config.Redis.Password,
		})
	} else {
		// 使用单例模式
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Config.Redis.Addr,
			Password: cfg.Config.Redis.Password,
			DB:       cfg.Config.Redis.DB,
		})
	}
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	} else {
		Redis = client
	}
}
