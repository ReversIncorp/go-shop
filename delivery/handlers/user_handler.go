package handlers

import (
	"marketplace/internal/domain/entities"
	userUsecase "marketplace/internal/domain/usecase/user_ucecase"
	"marketplace/pkg/error_handling"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// UserHandler обрабатывает HTTP-запросы для пользователей
type UserHandler struct {
	userUseCase *userUsecase.UserUseCase
	validator   *validator.Validate
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userUseCase *userUsecase.UserUseCase, validate *validator.Validate) *UserHandler {
	return &UserHandler{userUseCase: userUseCase, validator: validate}
}

// Register обрабатывает запрос на регистрацию пользователя
func (h *UserHandler) Register(c echo.Context) error {
	var user entities.User

	if err := c.Bind(&user); err != nil {
		return error_handling.ErrInvalidInput
	}
	if err := h.validator.Struct(user); err != nil {
		return error_handling.ErrValidationFailed
	}

	// Вызов метода Register и получение токенов
	tokens, err := h.userUseCase.Register(user, c)
	if err != nil {
		return err
	}

	// Возвращаем информацию о пользователе и токенах
	return c.JSON(http.StatusCreated, tokens.CleanOutput())
}

// Login обрабатывает запрос на вход пользователя
func (h *UserHandler) Login(c echo.Context) error {
	var credentials entities.LoginCredentials
	if err := c.Bind(&credentials); err != nil {
		return error_handling.ErrInvalidInput
	}
	if err := h.validator.Struct(credentials); err != nil {
		return error_handling.ErrValidationFailed
	}

	// Вызов метода Login и получение токенов
	tokens, err := h.userUseCase.Login(credentials.Email, credentials.Password, c)
	if err != nil {
		return err
	}

	// Возвращаем токены
	return c.JSON(http.StatusOK, tokens.CleanOutput())
}

// GetUserByID обрабатывает запрос на получение информации о пользователе по ID
func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}
	user, err := h.userUseCase.GetUserByID(uint64ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateToken обрабатывает запрос на получение информации о пользователе по ID
func (h *UserHandler) UpdateToken(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}
	user, err := h.userUseCase.GetUserByID(uint64ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
