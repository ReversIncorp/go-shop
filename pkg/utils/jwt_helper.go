package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace/internal/domain/entities"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// TokenDetails хранит данные о токенах
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

func (t *TokenDetails) ToTokens() *entities.Tokens {
	return &entities.Tokens{RefreshToken: t.RefreshToken, AccessToken: t.AccessToken}
}

// GenerateTokens создает новые Access и Refresh токены
func GenerateTokens(userID uint64) (*TokenDetails, error) {
	tokenDetails := &TokenDetails{}

	// Генерация Access токена
	tokenDetails.AtExpires = time.Now().Add(time.Hour * 72).Unix() // Срок действия Access токена - 72 часа
	tokenDetails.AccessUUID = uuid.New().String()                  // Генерация нового UUID для Access токена

	accessClaims := jwt.MapClaims{
		"user_id":     userID,
		"exp":         tokenDetails.AtExpires,
		"access_uuid": tokenDetails.AccessUUID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	tokenDetails.AccessToken = tokenString

	// Генерация Refresh токена
	tokenDetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // Срок действия Refresh токена - 7 дней
	tokenDetails.RefreshUUID = uuid.New().String()                     // Генерация нового UUID для Refresh токена

	refreshClaims := jwt.MapClaims{
		"user_id":      userID,
		"exp":          tokenDetails.RtExpires,
		"refresh_uuid": tokenDetails.RefreshUUID,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}
	tokenDetails.RefreshToken = refreshTokenString

	return tokenDetails, nil
}

// ValidateToken проверяет корректность и валидность токена (Access или Refresh)
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		logrus.Infof("User ID from token: %v", claims["user_id"])
		return token, nil // Возвращаем валидный токен
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
