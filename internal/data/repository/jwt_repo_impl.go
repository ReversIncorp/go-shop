package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/enums"
	"marketplace/internal/domain/repository"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"
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
	ctx echo.Context,
) error {
	jsonToken, err := json.Marshal(token)
	if err != nil {
		return tracerr.Wrap(fmt.Errorf("failed to marshal token: %v", err))
	}
	if err = r.redisClient.SetEX(
		r.context,
		fmt.Sprintf(
			"%s_%s:%d",
			tokenType.String(),
			ctx.Request().Header.Get("User-Agent"),
			userID,
		),
		jsonToken, tokenType.Duration(),
	).Err(); err != nil {
		return tracerr.Wrap(fmt.Errorf("failed to save token into Redis: %v", err))
	}
	return nil
}

// GetToken retrieves TokenDetails by userID
func (r *redisJWTRepository) GetToken(
	userID uint64,
	tokenType enums.Token,
	ctx echo.Context,
) (*entities.TokenDetails, error) {
	var token *entities.TokenDetails
	tokenJson, err := r.redisClient.Get(
		r.context,
		fmt.Sprintf(
			"%s_%s:%d",
			tokenType.String(),
			ctx.Request().Header.Get("User-Agent()"),
			userID,
		),
	).Result()
	if errors.Is(err, redis.Nil) {
		return nil, tracerr.Wrap(fmt.Errorf("tokens not found"))
	} else if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("error retrieving tokens: %v", err))
	}

	if err = json.Unmarshal([]byte(tokenJson), &token); err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("failed to unmarshal token: %v", err))
	}
	return token, nil
}

// DeleteToken removes tokens by userID
func (r *redisJWTRepository) DeleteToken(userID uint64, tokenType enums.Token, ctx echo.Context) error {
	err := r.redisClient.Del(r.context,
		fmt.Sprintf(
			"%s_%s:%d",
			tokenType.String(),
			ctx.Request().Header.Get("User-Agent"),
			userID,
		),
	).Err()
	if err != nil {
		return fmt.Errorf("failed to delete tokens: %v", err)
	}
	return nil
}
