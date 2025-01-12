package userusecase

import (
	"marketplace/config"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/enums"
	"marketplace/internal/domain/repository"
	errorHandling "marketplace/pkg/error_handling"
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
		return nil, errorHandling.ErrUserExists
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
		return nil, errorHandling.ErrInvalidCredentials
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
		userID, ok := claims["user_id"].(uint64)
		if !ok {
			return errorHandling.ErrInvalidTokenFormat
		}
		sessionUUID, ok := claims["session_uuid"].(string)
		if !ok {
			return errorHandling.ErrInvalidTokenFormat
		}

		err := u.tokenRepo.DeleteSession(userID, sessionUUID)
		if err != nil {
			return err
		}

		return nil
	}
	return errorHandling.ErrInvalidExpiredToken
}

func (u *UserUseCase) UpdateSession(refreshToken string, ctx echo.Context) (*entities.SessionDetails, error) {
	token, err := u.ValidateToken(refreshToken, config.GetConfig().JWTKey, enums.Refresh)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errorHandling.ErrInvalidTokenFormat
		}
		sessionUUID, ok := claims["session_uuid"].(string)
		if !ok {
			return nil, errorHandling.ErrInvalidTokenFormat
		}

		session, err := u.createSession(uint64(userID), sessionUUID, ctx)
		if err != nil {
			return nil, err
		}
		err = u.tokenRepo.SaveSession(uint64(userID), sessionUUID, session)
		if err != nil {
			return nil, err
		}

		return session, nil
	}
	return nil, errorHandling.ErrInvalidExpiredToken
}

func (u *UserUseCase) ValidateToken(
	tokenString string,
	key []byte,
	tokenType enums.Token,
) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorHandling.ErrInvalidExpiredToken
		}
		return key, nil
	})
	if err != nil {
		return nil, errorHandling.ErrInvalidExpiredToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errorHandling.ErrInvalidTokenClaims
	}
	if !token.Valid {
		return nil, errorHandling.ErrInvalidTokenClaims
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errorHandling.ErrInvalidTokenClaims
	}
	sessionUUID, ok := claims["session_uuid"].(string)
	if !ok {
		return nil, errorHandling.ErrInvalidTokenClaims
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errorHandling.ErrInvalidExpiredToken
	}

	currentTime := time.Now().Unix()
	if int64(exp) < currentTime {
		return nil, errorHandling.ErrInvalidExpiredToken
	}

	session, err := u.tokenRepo.GetSession(uint64(userID), sessionUUID)
	if err != nil {
		return nil, errorHandling.ErrInvalidExpiredToken
	}
	switch tokenType {
	case enums.Access:
		if session.AccessToken != tokenString {
			return nil, errorHandling.ErrInvalidTokenType
		}
	case enums.Refresh:
		if session.RefreshToken != tokenString {
			return nil, errorHandling.ErrInvalidTokenType
		}
	default:
		return nil, errorHandling.ErrInvalidTokenType
	}

	logrus.Infof("User ID from token: %v, Token session UUID: %v, Token type: %v",
		userID,
		sessionUUID,
		tokenType,
	)

	return token, nil
}

// UpdateSession Реализация метода UpdateSession.
func (u *UserUseCase) createSession(
	userID uint64,
	sessionID string,
	ctx echo.Context,
) (*entities.SessionDetails, error) {
	session := &entities.SessionDetails{}
	session.DeviceInfo = ctx.Request().Header.Get("User-Agent")
	session.IPAddress = ctx.RealIP()

	accessToken, err := GenerateToken(userID, sessionID, enums.Access, config.GetConfig().JWTKey)
	if err != nil {
		return nil, err
	}
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
