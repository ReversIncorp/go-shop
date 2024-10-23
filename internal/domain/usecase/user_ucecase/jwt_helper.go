package userUsecase

import (
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GenerateToken создает новые Access и Refresh токены
func GenerateToken(userID uint64, tokenType enums.Token, key []byte) (*entities.TokenDetails, error) {
	tokenDetails := &entities.TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(tokenType.Duration()).Unix()
	tokenDetails.UUID = uuid.New().String()

	claims := jwt.MapClaims{
		"user_id":     userID,
		"exp":         tokenDetails.AtExpires,
		"access_uuid": tokenDetails.UUID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	tokenDetails.Token = tokenString

	return tokenDetails, nil
}

// validateToken проверяет корректность и валидность токена (Access или Refresh)
func validateToken(tokenString string, key []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи токена (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем наличие user_id в токене
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid token: user_id missing or invalid")
		}
		// Проверяем наличие UUID токена
		tokenUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token: UUID missing or invalid")
		}

		// Проверяем время истечения токена
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid token: expiration time missing")
		}

		// Текущее время в Unix формате
		currentTime := time.Now().Unix()

		// Если время истекло, возвращаем ошибку
		if int64(exp) < currentTime {
			return nil, fmt.Errorf("token has expired")
		}

		// Логируем полезные данные из токена
		logrus.Infof("User ID from token: %v, Token UUID: %v", userID, tokenUUID)

		// Если все проверки пройдены, возвращаем валидный токен
		return token, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateAccessToken проверяет корректность и валидность Access токена
func ValidateAccessToken(accessTokenString string, key []byte) (*jwt.Token, error) {
	return validateToken(accessTokenString, key)
}

// ValidateRefreshToken проверяет корректность и валидность Refresh токена
//
//goland:noinspection GoUnusedExportedFunction
func ValidateRefreshToken(refreshTokenString string, key []byte) (*jwt.Token, error) {
	return validateToken(refreshTokenString, key)
}
