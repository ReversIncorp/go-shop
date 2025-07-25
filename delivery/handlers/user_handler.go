package handlers

import (
	"marketplace/internal/domain/entities"
	userUsecase "marketplace/internal/domain/usecase/user_usecase"
	errorHandling "marketplace/pkg/error_handling"
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
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags users
// @Consumes application/json
// @Produces application/json
// @Param user body entities.User true "Данные пользователя"
// @Success 201 {object} map[string]interface{} "Пользователь успешно зарегистрирован"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /users/register [post].
func (h *UserHandler) Register(c echo.Context) error {
	var user entities.User

	if err := c.Bind(&user); err != nil {
		return errorHandling.ErrInvalidInput
	}
	if err := h.validator.Struct(user); err != nil {
		return errorHandling.ErrValidationFailed
	}

	// Вызов метода Register и получение токенов
	tokens, err := h.userUseCase.Register(&user, c)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// Возвращаем информацию о пользователе и токенах
	return c.JSON(http.StatusCreated, tokens.CleanOutput())
}

// Login обрабатывает запрос на вход пользователя.
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает токены доступа
// @Tags users
// @Consumes application/json
// @Produces application/json
// @Param credentials body entities.LoginCredentials true "Данные для входа"
// @Success 200 {object} map[string]interface{} "Успешный вход"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 401 {object} errorhandling.ResponseError "Неверные учетные данные"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /users/login [post].
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
		return tracerr.Wrap(err)
	}

	// Возвращаем токены
	return c.JSON(http.StatusOK, tokens.CleanOutput())
}

// GetUserByID обрабатывает запрос на получение информации о пользователе по ID.
// @Summary Получение пользователя по ID
// @Description Возвращает информацию о пользователе по его ID
// @Tags users
// @Produces application/json
// @Param id path int true "ID пользователя"
// @Success 200 {object} entities.User "Информация о пользователе"
// @Failure 400 {object} errorhandling.ResponseError "Неверный ID"
// @Failure 404 {object} errorhandling.ResponseError "Пользователь не найден"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /users/{id} [get].
func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}
	user, err := h.userUseCase.GetUserByID(uint64ID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return c.JSON(http.StatusOK, user)
}

// RefreshSession обрабатывает рефреш сессии по рефреш токену.
// @Summary Обновление сессии
// @Description Обновляет токены доступа используя refresh token
// @Tags users
// @Consumes application/json
// @Produces application/json
// @Param request body object true "Refresh token" {"refresh_token": "string"}
// @Success 200 {object} map[string]interface{} "Токены обновлены"
// @Failure 400 {object} errorhandling.ResponseError "Неверный refresh token"
// @Router /users/refresh [post].
func (h *UserHandler) RefreshSession(c echo.Context) error {
	var request struct {
		Token string `json:"refresh_token" validate:"required"`
	}

	if err := c.Bind(&request); err != nil {
		return errorHandling.ErrMissingToken
	}

	session, err := h.userUseCase.UpdateSession(request.Token, c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid refresh token"})
	}
	return c.JSON(http.StatusOK, session.CleanOutput())
}

// Logout обрабатывает логаут сессии по ацесс токену.
// @Summary Выход из системы
// @Description Завершает сессию пользователя
// @Tags users
// @Consumes application/json
// @Produces application/json
// @Param request body object true "Access token" {"access_token": "string"}
// @Success 200 {string} string "Успешный выход"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка при выходе"
// @Router /users/logout [post].
func (h *UserHandler) Logout(c echo.Context) error {
	var request struct {
		Token string `json:"access_token" validate:"required"`
	}
	if err := c.Bind(&request); err != nil {
		return errorHandling.ErrMissingToken
	}

	err := h.userUseCase.Logout(request.Token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
