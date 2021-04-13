package config

import (
	"os"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func InitRedis() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     dsn, //redis port
		Password: "123",
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}
