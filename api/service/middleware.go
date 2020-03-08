package service

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

type RedisClient struct {
	Client *redis.Client
	id     int
}

var redisClient *RedisClient
var logger *zap.SugaredLogger

func GetRedisClient() *RedisClient {
	if redisClient == nil {
		redisPort := viper.GetInt64("redis.port")
		redisAddr := fmt.Sprintf("localhost:%d", redisPort)

		redisClient = &RedisClient{
			Client: redis.NewClient(&redis.Options{
				Addr:     redisAddr,
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
			id: 100,
		}
	}
	return redisClient
}

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		l, _ := zap.NewDevelopment()
		logger = l.Sugar()
	}

	return logger
}
