package userUsecase

import (
	"errors"
	"fmt"
	"marketplace/config"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/enums"
	"marketplace/internal/domain/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserUseCase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.JWTRepository
}

// NewUserUseCase Конструктор для создания новой UserUseCase
func NewUserUseCase(userRepo repository.UserRepository, tokenRepo repository.JWTRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo, tokenRepo: tokenRepo}
}

// Register Реализация метода Register
func (u *UserUseCase) Register(user entities.User, ctx echo.Context) (*entities.SessionDetails, error) {
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return nil, errors.New("user already exists")
	}

	userID, err := u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	tokens, err := u.createSession(userID, ctx)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Login Реализация метода Login
func (u *UserUseCase) Login(email, password string, ctx echo.Context) (*entities.SessionDetails, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user.Password != password { // Здесь должна быть логика хэширования пароля
		return nil, errors.New("invalid credentials")
	}

	tokens, err := u.createSession(user.ID, ctx)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetUserByID Реализация метода GetUserByID
func (u *UserUseCase) GetUserByID(id uint64) (entities.User, error) {
	return u.userRepo.FindByID(id)
}

// UpdateToken Реализация метода GetUserByID
func (u *UserUseCase) UpdateToken(id uint64) (entities.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *UserUseCase) ValidateToken(tokenString string, key []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid token: user_id missing or invalid")
		}
		sessionUUID, ok := claims["session_uuid"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token: session UUID missing or invalid")
		}
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid token: expiration time missing")
		}

		currentTime := time.Now().Unix()
		if int64(exp) < currentTime {
			return nil, fmt.Errorf("token has expired")
		}

		session, err := u.tokenRepo.GetSession(uint64(userID), sessionUUID)
		if err != nil {
			return nil, fmt.Errorf("session not found")
		}
		if session.AccessToken != tokenString {
			return nil, fmt.Errorf("session acess token not equals to user acess token")
		}

		logrus.Infof("User ID from token: %v, Token session UUID: %v", userID, sessionUUID)

		return token, nil
	}
	return nil, errors.New("invalid token")
}

// RefreshToken Реализация метода RefreshToken
func (u *UserUseCase) createSession(userID uint64, ctx echo.Context) (*entities.SessionDetails, error) {
	session := &entities.SessionDetails{}
	sessionID := uuid.New().String()
	session.DeviceInfo = ctx.Request().Header.Get("User-Agent")
	session.IPAddress = ctx.RealIP()

	accessToken, err := GenerateToken(userID, sessionID, enums.Access, config.GetConfig().JWTKey)
	refreshToken, err := GenerateToken(userID, sessionID, enums.Refresh, config.GetConfig().JWTKey)
	if err != nil {
		return nil, err
	}

	session.ExpiresAt = time.Now().Add(enums.Refresh.Duration()).Unix()
	session.AccessToken = accessToken
	session.RefreshToken = refreshToken

	if err = u.tokenRepo.SaveSession(
		userID,
		sessionID,
		session,
	); err != nil {
		return nil, err
	}

	return session, nil
}
