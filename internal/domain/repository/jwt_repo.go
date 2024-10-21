package repository

import (
	"github.com/labstack/echo/v4"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/enums"
)

type JWTRepository interface {
	SaveToken(userID uint64, token *entities.TokenDetails, tokenType enums.Token, ctx echo.Context) error
	GetToken(userID uint64, tokenType enums.Token, ctx echo.Context) (*entities.TokenDetails, error)
	DeleteToken(userID uint64, tokenType enums.Token, ctx echo.Context) error
}
