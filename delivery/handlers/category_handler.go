package handlers

import (
	"github.com/labstack/echo/v4"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/usecase"
	"net/http"
	"strconv"
)

// CategoryHandler обрабатывает HTTP-запросы для категорий
type CategoryHandler struct {
	categoryUseCase *usecase.CategoryUseCase
}

// NewCategoryHandler создает новый экземпляр CategoryHandler
func NewCategoryHandler(categoryUseCase *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUseCase: categoryUseCase}
}

// CreateCategory обрабатывает запрос на создание категории
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var category entities.Category

	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := h.categoryUseCase.CreateCategory(category); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, category)
}

// GetCategoryByID обрабатывает запрос на получение категории по ID
func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	category, err := h.categoryUseCase.GetCategoryByID(uint64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

// UpdateCategory обрабатывает запрос на обновление категории
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	var category entities.Category

	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := h.categoryUseCase.UpdateCategory(category); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

// DeleteCategory обрабатывает запрос на удаление категории
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := h.categoryUseCase.DeleteCategory(uint64ID); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetAllCategories обрабатывает запрос на получение всех категорий
func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	categories, err := h.categoryUseCase.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, categories)
}
