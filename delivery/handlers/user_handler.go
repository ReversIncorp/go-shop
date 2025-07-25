package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/usecase"
	"net/http"
	"strconv"
)

// UserHandler обрабатывает HTTP-запросы для пользователей
type UserHandler struct {
	userUseCase *usecase.UserUseCase
	validator   *validator.Validate
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase, validate *validator.Validate) *UserHandler {
	return &UserHandler{userUseCase: userUseCase, validator: validate}
}

// Register обрабатывает запрос на регистрацию пользователя
func (h *UserHandler) Register(c echo.Context) error {
	var user entities.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	if err := h.validator.Struct(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Вызов метода Register и получение токенов
	tokens, err := h.userUseCase.Register(user, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Возвращаем информацию о пользователе и токенах
	return c.JSON(http.StatusCreated, tokens.CleanOutput())
}

// Login обрабатывает запрос на вход пользователя
func (h *UserHandler) Login(c echo.Context) error {
	var credentials entities.LoginCredentials
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	if err := h.validator.Struct(credentials); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Вызов метода Login и получение токенов
	tokens, err := h.userUseCase.Login(credentials.Email, credentials.Password, c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	// Возвращаем токены
	return c.JSON(http.StatusOK, tokens.CleanOutput())
}

// GetUserByID обрабатывает запрос на получение информации о пользователе по ID
func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user, err := h.userUseCase.GetUserByID(uint64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateToken обрабатывает запрос на получение информации о пользователе по ID
func (h *UserHandler) UpdateToken(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user, err := h.userUseCase.GetUserByID(uint64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
