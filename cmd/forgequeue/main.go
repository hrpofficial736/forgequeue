package main

import redis "github.com/hrpofficial736/forgequeue/internal/storage/redis"

func main() {
	//cli.Execute()
	redis.LoadRedis()
}
