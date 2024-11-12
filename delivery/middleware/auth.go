package middleware

import (
	"marketplace/config"
	userUsecase "marketplace/internal/domain/usecase/user_ucecase"
	"marketplace/pkg/errors"

	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware обрабатывает аутентификацию JWT токенов.
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Извлекаем токен из заголовков Authorization
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return errors.ErrMissingToken
		}
		// Проверяем формат токена (Bearer)
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return errors.ErrInvalidTokenFormat
		}

		// Получаем сам JWT токен
		tokenString := tokenParts[1]

		// Проверяем валидность Access токена
		token, err := userUsecase.ValidateAccessToken(tokenString, config.GetConfig().JWTKey)
		if err != nil {
			return errors.ErrInvalidExpiredToken
		}

		// Извлекаем информацию о пользователе из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return errors.ErrInvalidTokenClaims
		}

		// Устанавливаем user_id в контекст для последующего использования в контроллерах.
		userID := claims["user_id"]
		c.Set("user_id", userID)
		// Переходим к следующему обработчику.
		return next(c)
	}
}
