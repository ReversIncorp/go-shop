package middleware

import (
	"marketplace/config"
	"marketplace/internal/domain/enums"
	userUsecase "marketplace/internal/domain/usecase/user_ucecase"
	"marketplace/pkg/error_handling"

	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware обрабатывает аутентификацию JWT токенов.
func JWTMiddleware(userUseCase *userUsecase.UserUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем токен из заголовков Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return error_handling.ErrMissingToken
			}
			// Проверяем формат токена (Bearer)
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return error_handling.ErrInvalidTokenFormat
			}

			// Получаем сам JWT токен
			tokenString := tokenParts[1]

			// Проверяем валидность Access токена
			token, err := userUseCase.ValidateToken(tokenString, config.GetConfig().JWTKey, enums.Access)
			if err != nil {
				return error_handling.ErrInvalidExpiredToken
			}

			// Извлекаем информацию о пользователе из токена
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return error_handling.ErrInvalidTokenClaims
			}

			// Устанавливаем user_id в контекст для последующего использования в контроллерах.
			userID := claims["user_id"]
			c.Set("user_id", userID)
			// Переходим к следующему обработчику.
			return next(c)
		}
	}
}
