package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/enums"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateToken создает новые Access и Refresh токены
func GenerateToken(userID uint64, tokenType enums.Token) (*entities.TokenDetails, error) {
	tokenDetails := &entities.TokenDetails{}

	// Генерация Access токена
	tokenDetails.AtExpires = time.Now().Add(tokenType.Duration()).Unix() // Срок действия Access токена - 72 часа
	tokenDetails.UUID = uuid.New().String()                              // Генерация нового UUID для Access токена

	accessClaims := jwt.MapClaims{
		"user_id":     userID,
		"exp":         tokenDetails.AtExpires,
		"access_uuid": tokenDetails.UUID,
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, err := claims.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	tokenDetails.Token = tokenString

	return tokenDetails, nil
}

// ValidateToken проверяет корректность и валидность токена (Access или Refresh)
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи токена (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
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

		//TODO: Добавить проверку токена на наличие в БД, типо redis для управления сессиями

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

	return nil, fmt.Errorf("invalid token")
}

// ValidateAccessToken проверяет корректность и валидность Access токена
func ValidateAccessToken(accessTokenString string) (*jwt.Token, error) {
	return ValidateToken(accessTokenString)
}

// ValidateRefreshToken проверяет корректность и валидность Refresh токена
func ValidateRefreshToken(refreshTokenString string) (*jwt.Token, error) {
	return ValidateToken(refreshTokenString)
}
