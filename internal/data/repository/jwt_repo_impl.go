package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"marketplace/internal/domain/enums"
	"marketplace/internal/domain/repository"
	"strconv"
	"time"
)

// RedisJWTRepository - реализация JWTRepository для Redis
type RedisJWTRepository struct {
	redisClient *redis.Client
	context     context.Context
}

// NewRedisJWTRepository - конструктор для создания нового экземпляра RedisJWTRepository
func NewRedisJWTRepository(redisClient *redis.Client) repository.JWTRepository {
	return &RedisJWTRepository{
		redisClient: redisClient,
		context:     context.Background(),
	}
}

// SaveToken saves TokenDetails in Redis with a specified expiration time
func (r *RedisJWTRepository) SaveToken(
	userID uint64,
	token string,
	tokenType enums.Token,
	expiration time.Duration,
) error {
	err := r.redisClient.SetEX(r.context, fmt.Sprintf("%s:%d", tokenType.String(), userID), token, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to save token into Redis: %v", err)
	}
	return nil
}

// GetToken retrieves TokenDetails by userID
func (r *RedisJWTRepository) GetToken(
	userID uint64,
	tokenType enums.Token,
) (string, error) {
	// Retrieve serialized JSON from Redis
	token, err := r.redisClient.Get(r.context, fmt.Sprintf("%s:%d", tokenType.String(), userID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("tokens not found")
	} else if err != nil {
		return "", fmt.Errorf("error retrieving tokens: %v", err)
	}

	return token, nil
}

// DeleteToken removes tokens by userID
func (r *RedisJWTRepository) DeleteToken(userID uint64) error {
	err := r.redisClient.Del(r.context, strconv.FormatUint(userID, 10)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete tokens: %v", err)
	}
	return nil
}
