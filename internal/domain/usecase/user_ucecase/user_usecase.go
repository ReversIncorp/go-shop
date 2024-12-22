package userusecase

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

// NewUserUseCase Конструктор для создания новой UserUseCase.
func NewUserUseCase(userRepo repository.UserRepository, tokenRepo repository.JWTRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo, tokenRepo: tokenRepo}
}

// Register Реализация метода Register.
func (u *UserUseCase) Register(user entities.User, ctx echo.Context) (*entities.SessionDetails, error) {
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return nil, errors.New("user already exists")
	}

	userID, err := u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	tokens, err := u.createSession(userID, uuid.New().String(), ctx)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Login Реализация метода Login.
func (u *UserUseCase) Login(email, password string, ctx echo.Context) (*entities.SessionDetails, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user.Password != password { // Здесь должна быть логика хэширования пароля
		return nil, errors.New("invalid credentials")
	}

	tokens, err := u.createSession(user.ID, uuid.New().String(), ctx)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetUserByID Реализация метода GetUserByID.
func (u *UserUseCase) GetUserByID(id uint64) (entities.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *UserUseCase) Logout(accessToken string) error {
	token, err := u.ValidateToken(accessToken, config.GetConfig().JWTKey, enums.Access)
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"]
		if !ok {
			return fmt.Errorf("invalid token: user_id missing or invalid")
		}
		sessionUUID, ok := claims["session_uuid"]
		if !ok {
			return fmt.Errorf("invalid token: session UUID missing or invalid")
		}

		err := u.tokenRepo.DeleteSession(userID.(uint64), sessionUUID.(string))
		if err != nil {
			return errors.New("session deletion failed")
		}

		return nil
	}
	return errors.New("invalid access token")
}

func (u *UserUseCase) UpdateSession(refreshToken string, ctx echo.Context) (*entities.SessionDetails, error) {
	token, err := u.ValidateToken(refreshToken, config.GetConfig().JWTKey, enums.Refresh)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"]
		if !ok {
			return nil, fmt.Errorf("invalid token: user_id missing or invalid")
		}
		sessionUUID, ok := claims["session_uuid"]
		if !ok {
			return nil, fmt.Errorf("invalid token: session UUID missing or invalid")
		}

		session, err := u.createSession(userID.(uint64), sessionUUID.(string), ctx)
		if err != nil {
			return nil, fmt.Errorf("session creation failed")
		}
		err = u.tokenRepo.SaveSession(userID.(uint64), sessionUUID.(string), session)
		if err != nil {
			return nil, fmt.Errorf("session saving failed")
		}

		return session, nil
	}
	return nil, errors.New("invalid refresh token")
}

func (u *UserUseCase) ValidateToken(tokenString string, key []byte, tokenType enums.Token) (*jwt.Token, error) {
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

		if tokenType == enums.Access {
			if session.AccessToken != tokenString {
				return nil, fmt.Errorf("session acess token not equals to user acess token")
			}
		} else if tokenType == enums.Refresh {
			if session.RefreshToken != tokenString {
				return nil, fmt.Errorf("session refresh token not equals to user refresh token")
			}
		} else {
			return nil, fmt.Errorf("invalid token type")
		}

		logrus.Infof("User ID from token: %v, Token session UUID: %v, Token type: %v", userID, sessionUUID, tokenType)

		return token, nil
	}
	return nil, errors.New("invalid token")
}

// UpdateSession Реализация метода UpdateSession
func (u *UserUseCase) createSession(userID uint64, sessionID string, ctx echo.Context) (*entities.SessionDetails, error) {
	session := &entities.SessionDetails{}
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
