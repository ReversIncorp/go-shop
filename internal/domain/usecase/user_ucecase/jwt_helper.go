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
func GenerateTokens(userID uint64, key []byte) (*TokenDetails, error) {
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
	tokenString, err := accessToken.SignedString(key)
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
	refreshTokenString, err := refreshToken.SignedString(key)
	if err != nil {
		return nil, err
	}
	tokenDetails.RefreshToken = refreshTokenString

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
