package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/enums"
	"marketplace/internal/domain/repository"
	"strconv"
)

// redisJWTRepository - реализация JWTRepository для Redis
type redisJWTRepository struct {
	redisClient *redis.Client
	context     context.Context
}

// NewRedisJWTRepository - конструктор для создания нового экземпляра redisJWTRepository
func NewRedisJWTRepository(redisClient *redis.Client) repository.JWTRepository {
	return &redisJWTRepository{
		redisClient: redisClient,
		context:     context.Background(),
	}
}

// SaveToken saves TokenDetails in Redis with a specified expiration time
func (r *redisJWTRepository) SaveToken(
	userID uint64,
	token *entities.TokenDetails,
	tokenType enums.Token,
) error {
	jsonToken, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}
	err = r.redisClient.SetEX(r.context, fmt.Sprintf("%s:%d", tokenType.String(), userID), jsonToken, tokenType.Duration()).Err()
	if err != nil {
		return fmt.Errorf("failed to save token into Redis: %v", err)
	}
	return nil
}

// GetToken retrieves TokenDetails by userID
func (r *redisJWTRepository) GetToken(
	userID uint64,
	tokenType enums.Token,
) (*entities.TokenDetails, error) {
	var token *entities.TokenDetails
	tokenJson, err := r.redisClient.Get(r.context, fmt.Sprintf("%s:%d", tokenType.String(), userID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("tokens not found")
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving tokens: %v", err)
	}

	if err = json.Unmarshal([]byte(tokenJson), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %v", err)
	}
	return token, nil
}

// DeleteToken removes tokens by userID
func (r *redisJWTRepository) DeleteToken(userID uint64) error {
	err := r.redisClient.Del(r.context, strconv.FormatUint(userID, 10)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete tokens: %v", err)
	}
	return nil
}
