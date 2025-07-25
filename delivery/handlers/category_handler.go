package handlers

import (
	"marketplace/internal/domain/entities"
	categoryUsecases "marketplace/internal/domain/usecase/category_usecase"
	errorHandling "marketplace/pkg/error_handling"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"
)

// CategoryHandler обрабатывает HTTP-запросы для категорий.
type CategoryHandler struct {
	categoryUseCase *categoryUsecases.CategoryUseCase
}

// NewCategoryHandler создает новый экземпляр CategoryHandler.
func NewCategoryHandler(categoryUseCase *categoryUsecases.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUseCase: categoryUseCase}
}

// CreateCategory обрабатывает запрос на создание категории.
// @Summary Создание категории
// @Description Создает новую категорию для магазина
// @Tags categories
// @Consumes application/json
// @Produces application/json
// @Param category body entities.Category true "Данные категории"
// @Success 201 {object} entities.Category "Категория создана"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 401 {object} errorhandling.ResponseError "Не авторизован"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /categories [post].
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var category entities.Category

	if err := c.Bind(&category); err != nil {
		return errorHandling.ErrInvalidInput
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return errorHandling.ErrMissingUserFromToken
	}

	if err := h.categoryUseCase.CreateCategory(&category, uint64(uid)); err != nil {
		tracerr.Wrap(err)
	}

	created, err := h.categoryUseCase.GetCategoryByID(category.ID)
	if err != nil {
		return errorHandling.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, created)
}

// DeleteCategory обрабатывает запрос на удаление категории.
// @Summary Удаление категории
// @Description Удаляет категорию из системы
// @Tags categories
// @Produces application/json
// @Param id path int true "ID категории"
// @Success 204 {string} string "Категория удалена"
// @Failure 400 {object} errorhandling.ResponseError "Неверный ID"
// @Failure 401 {object} errorhandling.ResponseError "Не авторизован"
// @Failure 404 {object} errorhandling.ResponseError "Категория не найдена"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /categories/{id} [delete].
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return errorHandling.ErrMissingUserFromToken
	}

	if err = h.categoryUseCase.DeleteCategory(uint64ID, uint64(uid)); err != nil {
		return tracerr.Wrap(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetAllCategoriesByStore обрабатывает запрос на получение всех категорий.
// @Summary Получение категорий магазина
// @Description Возвращает все категории конкретного магазина
// @Tags categories
// @Produces application/json
// @Param store_id path int true "ID магазина"
// @Success 200 {array} entities.Category "Список категорий"
// @Failure 400 {object} errorhandling.ResponseError "Неверный ID магазина"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores/{store_id}/categories [get].
func (h *CategoryHandler) GetAllCategoriesByStore(c echo.Context) error {
	id := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	categories, err := h.categoryUseCase.GetAllCategoriesByStore(uint64ID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return c.JSON(http.StatusOK, categories)
}

// GetCategoryByID обрабатывает запрос на получение категории по айди.
// @Summary Получение категории по ID
// @Description Возвращает информацию о категории по ее ID
// @Tags categories
// @Produces application/json
// @Param id path int true "ID категории"
// @Success 200 {object} entities.Category "Информация о категории"
// @Failure 400 {object} errorhandling.ResponseError "Неверный ID"
// @Failure 404 {object} errorhandling.ResponseError "Категория не найдена"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /categories/{id} [get].
func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	category, err := h.categoryUseCase.GetCategoryByID(uint64ID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return c.JSON(http.StatusOK, category)
}

// GetAllCategories обрабатывает запрос на получение всех категорий.
// @Summary Получение всех категорий
// @Description Возвращает список всех категорий в системе
// @Tags categories
// @Produces application/json
// @Success 200 {array} entities.Category "Список всех категорий"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /categories [get].
func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	category, err := h.categoryUseCase.GetAllCategories()
	if err != nil {
		return tracerr.Wrap(err)
	}

	return c.JSON(http.StatusOK, category)
}
