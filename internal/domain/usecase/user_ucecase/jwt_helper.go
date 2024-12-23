package userusecase

import (
	"marketplace/internal/domain/enums"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateToken создает новые Access и Refresh токены.
func GenerateToken(userID uint64, sessionID string, tokenType enums.Token, key []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      userID,
		"exp":          time.Now().Add(tokenType.Duration()).Unix(),
		"session_uuid": sessionID, // Пока не знаю где это можно заюзать, фактически бесполезно
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
