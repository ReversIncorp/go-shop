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

	userID := c.Get("user_id")
	if userID != nil {
		if uid, ok := userID.(float64); ok {
			if err := h.categoryUseCase.CreateCategory(category, int64(uid)); err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
			}
			return c.JSON(http.StatusCreated, category)
		}
	}

	return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing user_id from token"})
}

// GetCategoryByID обрабатывает запрос на получение категории по ID
func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	id := c.Param("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	category, err := h.categoryUseCase.GetCategoryByID(int64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

// UpdateCategory обрабатывает запрос на обновление категории
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	var category entities.Category

	id := c.Param("id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid category ID"})
	}

	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	category.ID = categoryID

	userID := c.Get("user_id")
	if userID != nil {
		if uid, ok := userID.(float64); ok {
			if err := h.categoryUseCase.UpdateCategory(category, int64(uid)); err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
			}
			return c.JSON(http.StatusOK, category)
		}
	}

	return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing user_id from token"})
}

// DeleteCategory обрабатывает запрос на удаление категории
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	userID := c.Get("user_id")
	if userID != nil {
		if uid, ok := userID.(float64); ok {
			if err := h.categoryUseCase.DeleteCategory(int64ID, int64(uid)); err != nil {
				return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
			}
			return c.NoContent(http.StatusNoContent)
		}
	}

	return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing user_id from token"})
}

// GetAllCategories обрабатывает запрос на получение всех категорий
func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	id := c.Param("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	categories, err := h.categoryUseCase.GetAllCategoriesByStore(int64ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, categories)
}
