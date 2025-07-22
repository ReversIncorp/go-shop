package handlers

import (
	"marketplace/internal/domain/entities"
	categoryUsecases "marketplace/internal/domain/usecase/category_usecase"
	errorHandling "marketplace/pkg/error_handling"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
		return err
	}

	created, err := h.categoryUseCase.GetCategoryByID(category.ID)
	if err != nil {
		return errorHandling.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, created)
}

// DeleteCategory обрабатывает запрос на удаление категории.
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
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// GetAllCategoriesByStore обрабатывает запрос на получение всех категорий.
func (h *CategoryHandler) GetAllCategoriesByStore(c echo.Context) error {
	id := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	categories, err := h.categoryUseCase.GetAllCategoriesByStore(uint64ID)
	if err != nil {
		return errorHandling.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, categories)
}

// GetCategoryByID обрабатывает запрос на получение категории по айди.
func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	category, err := h.categoryUseCase.GetCategoryByID(uint64ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, category)
}

// GetAllCategories обрабатывает запрос на получение всех категорий.
func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	category, err := h.categoryUseCase.GetAllCategories()
	if err != nil {
		return errorHandling.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, category)
}
