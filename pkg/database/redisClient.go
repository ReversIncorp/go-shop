package database

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() (*redis.Client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: "",
			DB:       0,
		},
	)
	pong, err := redisClient.Ping(redisClient.Context()).Result()
	if pong != "PONG" || err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return redisClient, nil
}
