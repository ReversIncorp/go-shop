package handlers

import (
	"github.com/labstack/echo/v4"
	"marketplace/internal/domain/models"
	"marketplace/internal/domain/usecase"
	"net/http"
)

// UserHandler обрабатывает HTTP-запросы для пользователей
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// Register обрабатывает запрос на регистрацию пользователя
func (h *UserHandler) Register(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := h.userUseCase.Register(user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// Login обрабатывает запрос на вход пользователя
func (h *UserHandler) Login(c echo.Context) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	token, err := h.userUseCase.Login(credentials.Email, credentials.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

// GetUserByID обрабатывает запрос на получение информации о пользователе по ID
func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")

	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
