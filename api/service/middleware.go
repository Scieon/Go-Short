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

func InitializeRedisClient(addr string) {
	redisClient = &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		id: 100,
	}
}

func GetBaseURL() string {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")

	host := viper.GetString("server.host")
	port := viper.GetString("server.port")

	url := fmt.Sprintf("%s:%s/", host, port)
	return url
}

func GetRedisClient() *RedisClient {
	return redisClient
}

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		l, _ := zap.NewDevelopment()
		logger = l.Sugar()
	}

	return logger
}
