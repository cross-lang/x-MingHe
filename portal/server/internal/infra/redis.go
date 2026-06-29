package infra

import (
	"context"
	"fmt"
	"portal/internal/config"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

var RedisDB redis.UniversalClient
var Limiter *redis_rate.Limiter

func InitRedis(conf config.RedisConfig) (redis.UniversalClient, error) {

	switch conf.Mode {
	case "single":
		RedisDB = redis.NewClient(&redis.Options{
			Addr:         conf.Addr,
			Password:     conf.Password,
			DB:           conf.DB,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConns,
		})
	case "cluster":
		RedisDB = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        conf.Addrs,
			Password:     conf.Password,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConns,
		})
	case "sentinel":
		RedisDB = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.MasterName,
			SentinelAddrs: conf.Addrs,
			Password:      conf.Password,
			DB:            conf.DB,
			PoolSize:      conf.PoolSize,
			MinIdleConns:  conf.MinIdleConns,
		})
	default:
		return nil, fmt.Errorf("未知的redis连接模式: %s", conf.Mode)
	}
	err := RedisDB.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("redis连接失败: %w", err)
	}
	return RedisDB, nil
}

func InitRedisLimiter(redisDB redis.UniversalClient) *redis_rate.Limiter {
	return redis_rate.NewLimiter(redisDB)
}
