package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func LoadRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	err := client.Set(ctx, "foo", "bar", 0).Err()

	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(val)

	return client
}
