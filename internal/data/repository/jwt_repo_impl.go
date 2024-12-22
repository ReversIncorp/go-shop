package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	"time"

	"github.com/go-redis/redis/v8"
)

// redisJWTRepository - реализация JWTRepository для Redis.
type redisJWTRepository struct {
	redisClient *redis.Client
	context     context.Context
}

// NewRedisJWTRepository - конструктор для создания нового экземпляра redisJWTRepository.
func NewRedisJWTRepository(redisClient *redis.Client) repository.JWTRepository {
	return &redisJWTRepository{
		redisClient: redisClient,
		context:     context.Background(),
	}
}

func (r *redisJWTRepository) SaveSession(userID uint64, sessionID string, session *entities.SessionDetails) error {
	sessionKey := fmt.Sprintf("user:%d:session:%s", userID, sessionID)

	jsonSession, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	ttl := time.Until(time.Unix(session.ExpiresAt, 0))
	if ttl <= 0 {
		return errors.New("session expiration time is invalid or already expired")
	}

	if err = r.redisClient.SetEX(r.context, sessionKey, jsonSession, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save session with expiration: %w", err)
	}

	return nil
}

func (r *redisJWTRepository) GetSession(userID uint64, sessionID string) (*entities.SessionDetails, error) {
	sessionKey := fmt.Sprintf("user:%d:session:%s", userID, sessionID)

	sessionJSON, err := r.redisClient.Get(r.context, sessionKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("session not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session entities.SessionDetails
	if err = json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (r *redisJWTRepository) GetAllSessions(userID uint64) (map[string]*entities.SessionDetails, error) {
	sessionPattern := fmt.Sprintf("user:%d:session:*", userID)

	iter := r.redisClient.Scan(r.context, 0, sessionPattern, 0).Iterator()
	result := make(map[string]*entities.SessionDetails)

	for iter.Next(r.context) {
		sessionKey := iter.Val()
		sessionJSON, err := r.redisClient.Get(r.context, sessionKey).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get session: %w", err)
		}

		var session entities.SessionDetails
		if err = json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			return nil, fmt.Errorf("failed to unmarshal session: %w", err)
		}

		sessionID := sessionKey[len(fmt.Sprintf("user:%d:session:", userID)):]
		result[sessionID] = &session
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sessions: %w", err)
	}

	return result, nil
}

func (r *redisJWTRepository) DeleteSession(userID uint64, sessionID string) error {
	sessionKey := fmt.Sprintf("user:%d:session:%s", userID, sessionID)

	if err := r.redisClient.Del(r.context, sessionKey).Err(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
