package repository

import (
	"marketplace/internal/domain/enums"
	"time"
)

type JWTRepository interface {
	SaveToken(userID uint64, token string, tokenType enums.Token, expiration time.Duration) error
	GetToken(userID uint64, tokenType enums.Token) (string, error)
	DeleteToken(userID uint64) error
}
