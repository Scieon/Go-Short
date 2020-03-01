package service

import (
	"github.com/go-redis/redis/v7"
)

type RedisClient struct {
	client *redis.Client
	id int
}

var redisClient *RedisClient

func GetRedisClient() *RedisClient {
	if redisClient == nil {
		redisClient = &RedisClient{
			client: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
			id: 1,
		}
	}
	return redisClient
}
