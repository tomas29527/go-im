package redislib

import (
	"github.com/go-redis/redis"
	"go-im/config"
	"go-im/logger"
	"go.uber.org/zap"
)

var (
	client *redis.Client
)

func ExampleNewClient() {

	client = redis.NewClient(&redis.Options{
		Addr:         config.GlobalConfig.Redis.Addr,
		Password:     config.GlobalConfig.Redis.Password,
		DB:           config.GlobalConfig.Redis.DB,
		PoolSize:     config.GlobalConfig.Redis.PoolSize,
		MinIdleConns: config.GlobalConfig.Redis.MinIdleConns,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		logger.Logger.Error("初始化redis错误:", zap.Error(err))
	}
	logger.Logger.Info("初始化redis:", zap.String("pong", pong))
	// Output: PONG <nil>
}

func GetClient() (c *redis.Client) {

	return client
}
