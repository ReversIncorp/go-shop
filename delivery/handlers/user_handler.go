package handlers

import (
	"marketplace/internal/domain/entities"
	userUsecase "marketplace/internal/domain/usecase/user_usecase"
	errorHandling "marketplace/pkg/error_handling"
	"marketplace/pkg/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"
)

// UserHandler обрабатывает HTTP-запросы для пользователей.
type UserHandler struct {
	userUseCase *userUsecase.UserUseCase
	validator   *validator.Validate
}

// NewUserHandler создает новый экземпляр UserHandler.
func NewUserHandler(userUseCase *userUsecase.UserUseCase, validate *validator.Validate) *UserHandler {
	return &UserHandler{userUseCase: userUseCase, validator: validate}
}

// Register обрабатывает запрос на регистрацию пользователя.
func (h *UserHandler) Register(c echo.Context) error {
	var user entities.User

	if err := c.Bind(&user); err != nil {
		return errorHandling.ErrInvalidInput
	}
	if err := h.validator.Struct(user); err != nil {
		return errorHandling.ErrValidationFailed
	}

	// Вызов метода Register и получение токенов
	tokens, err := h.userUseCase.Register(user, c)
	if err != nil {
		return utils.GetHttpErrorOrTracerrError(err)
	}
	// Возвращаем информацию о пользователе и токенах
	return c.JSON(http.StatusCreated, tokens.CleanOutput())
}

// Login обрабатывает запрос на вход пользователя.
func (h *UserHandler) Login(c echo.Context) error {
	var credentials entities.LoginCredentials
	if err := c.Bind(&credentials); err != nil {
		return errorHandling.ErrInvalidInput
	}
	if err := h.validator.Struct(credentials); err != nil {
		return errorHandling.ErrValidationFailed
	}

	// Вызов метода Login и получение токенов
	tokens, err := h.userUseCase.Login(credentials.Email, credentials.Password, c)
	if err != nil {
		return err
	}

	// Возвращаем токены
	return c.JSON(http.StatusOK, tokens.CleanOutput())
}

// GetUserByID обрабатывает запрос на получение информации о пользователе по ID.
func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}
	user, err := h.userUseCase.GetUserByID(uint64ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// RefreshSession обрабатывает рефреш сессии по рефреш токену.
func (h *UserHandler) RefreshSession(c echo.Context) error {
	var request struct {
		Token string `json:"refresh_token" validate:"required"`
	}

	if err := c.Bind(&request); err != nil {
		return errorHandling.ErrMissingToken
	}

	session, err := h.userUseCase.UpdateSession(request.Token, c)
	if err != nil {
		if utils.IsHttpError(err) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid refresh token"})
		} else {
			return tracerr.Wrap(err)
		}
	}
	return c.JSON(http.StatusOK, session.CleanOutput())
}

// Logout обрабатывает логаут сессии по ацесс токену.
func (h *UserHandler) Logout(c echo.Context) error {
	var request struct {
		Token string `json:"access_token" validate:"required"`
	}
	if err := c.Bind(&request); err != nil {
		return errorHandling.ErrMissingToken
	}

	err := h.userUseCase.Logout(request.Token)
	if err != nil {
		if utils.IsHttpError(err) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		} else {
			return tracerr.Wrap(err)
		}
	}
	return c.NoContent(http.StatusOK)
}
