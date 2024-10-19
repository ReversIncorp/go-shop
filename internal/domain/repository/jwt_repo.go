package repository

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/enums"
)

type JWTRepository interface {
	SaveToken(userID uint64, token *entities.TokenDetails, tokenType enums.Token) error
	GetToken(userID uint64, tokenType enums.Token) (*entities.TokenDetails, error)
	DeleteToken(userID uint64) error
}
