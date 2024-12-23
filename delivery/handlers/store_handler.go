package handlers

import (
	"marketplace/internal/domain/entities"
	storeUsecases "marketplace/internal/domain/usecase/store_usecase"
	"marketplace/pkg/error_handling"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// StoreHandler обрабатывает HTTP-запросы для магазинов.
type StoreHandler struct {
	storeUseCase *storeUsecases.StoreUseCase
	validator    *validator.Validate
}

// NewStoreHandler создает новый экземпляр StoreHandler.
func NewStoreHandler(
	storeUseCase *storeUsecases.StoreUseCase,
	validator *validator.Validate,
) *StoreHandler {
	return &StoreHandler{
		storeUseCase: storeUseCase,
		validator:    validator,
	}
}

// CreateStore обрабатывает запрос на создание магазина.
func (h *StoreHandler) CreateStore(c echo.Context) error {
	var store entities.Store

	if err := c.Bind(&store); err != nil {
		return error_handling.ErrInvalidInput
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return error_handling.ErrMissingUserFromToken
	}

	if err := h.storeUseCase.CreateStore(store, uint64(uid)); err != nil {
		return error_handling.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, store)
}

// GetStoreByID обрабатывает запрос на получение магазина по ID.
func (h *StoreHandler) GetStoreByID(c echo.Context) error {
	id := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}
	store, err := h.storeUseCase.GetStoreByID(uint64ID)
	if err != nil {
		return error_handling.ErrStoreNotFound
	}

	return c.JSON(http.StatusOK, store)
}

// UpdateStore обрабатывает запрос на обновление магазина.
func (h *StoreHandler) UpdateStore(c echo.Context) error {
	var store entities.Store
	id := c.Param("store_id")
	storeID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	if err = c.Bind(&store); err != nil {
		return error_handling.ErrInvalidInput
	}

	store.ID = storeID

	if err = h.storeUseCase.UpdateStore(store); err != nil {
		return error_handling.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, store)
}

// DeleteStore обрабатывает запрос на удаление магазина.
func (h *StoreHandler) DeleteStore(c echo.Context) error {
	id := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	if err = h.storeUseCase.DeleteStore(uint64ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// GetStoresByFilters обрабатывает запрос на получение магазинов по фильтрам.
func (h *StoreHandler) GetStoresByFilters(c echo.Context) error {
	var searchParams entities.StoreSearchParams

	if err := c.Bind(&searchParams); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	if err := h.validator.Struct(searchParams); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Input validation failed"})
	}

	products, nextCursor, err := h.storeUseCase.GetStoresByFilters(searchParams)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       products,
		"limit":      searchParams.Limit,
		"nextCursor": nextCursor,
	})
}

// AttachCategoryToStore связывает категорию с магазином.
func (h *StoreHandler) AttachCategoryToStore(c echo.Context) error {
	var request struct {
		CategoryID uint64 `json:"category_id" validate:"required"`
	}

	storeIDParam := c.Param("store_id")
	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	if err = c.Bind(&request); err != nil || request.CategoryID == 0 {
		return error_handling.ErrInvalidInput
	}

	if err = h.storeUseCase.AttachCategoryToStore(storeID, request.CategoryID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Category attached to store successfully"})
}

// DetachCategoryFromStore отвязывает категорию от магазина.
func (h *StoreHandler) DetachCategoryFromStore(c echo.Context) error {
	storeIDParam := c.Param("store_id")
	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	categoryIDParam := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	if err = h.storeUseCase.DetachCategoryFromStore(storeID, categoryID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
