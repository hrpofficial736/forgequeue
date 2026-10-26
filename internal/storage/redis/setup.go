package redis

import (
	"github.com/redis/go-redis/v9"
)

func LoadRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return client
}
