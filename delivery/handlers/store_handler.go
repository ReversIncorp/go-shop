package handlers

import (
	"marketplace/internal/domain/entities"
	storeUsecases "marketplace/internal/domain/usecase/store_usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// StoreHandler обрабатывает HTTP-запросы для магазинов
type StoreHandler struct {
	storeUseCase *storeUsecases.StoreUseCase
}

// NewStoreHandler создает новый экземпляр StoreHandler
func NewStoreHandler(storeUseCase *storeUsecases.StoreUseCase) *StoreHandler {
	return &StoreHandler{storeUseCase: storeUseCase}
}

// CreateStore обрабатывает запрос на создание магазина
func (h *StoreHandler) CreateStore(c echo.Context) error {
	var store entities.Store

	if err := c.Bind(&store); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing user_id from token"})
	}

	if err := h.storeUseCase.CreateStore(store, uint64(uid)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, store)
}

// GetStoreByID обрабатывает запрос на получение магазина по ID
func (h *StoreHandler) GetStoreByID(c echo.Context) error {
	id := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	store, err := h.storeUseCase.GetStoreByID(uint64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, store)
}

// UpdateStore обрабатывает запрос на обновление магазина
func (h *StoreHandler) UpdateStore(c echo.Context) error {
	var store entities.Store
	id := c.Param("store_id")
	storeID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	if err = c.Bind(&store); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	store.ID = storeID

	if err = h.storeUseCase.UpdateStore(store); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, store)
}

// DeleteStore обрабатывает запрос на удаление магазина
func (h *StoreHandler) DeleteStore(c echo.Context) error {
	id := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	if err = h.storeUseCase.DeleteStore(uint64ID); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// GetAllStores обрабатывает запрос на получение всех магазинов
func (h *StoreHandler) GetAllStores(c echo.Context) error {
	stores, err := h.storeUseCase.GetAllStores()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stores)
}

// AddCategoryToStore связывает категорию с магазином
func (h *StoreHandler) AddCategoryToStore(c echo.Context) error {
	var request struct {
		CategoryID uint64 `json:"category_id" validate:"required"`
	}

	storeIDParam := c.Param("store_id")
	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	if err = c.Bind(&request); err != nil || request.CategoryID == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid category ID"})
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing user_id from token"})
	}

	if err = h.storeUseCase.AttachCategoryToStore(storeID, request.CategoryID, uint64(uid)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Category added to store successfully"})
}

// DeleteCategoryFromStore отвязывает категорию от магазина
func (h *StoreHandler) DeleteCategoryFromStore(c echo.Context) error {
	storeIDParam := c.Param("store_id")
	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	categoryIDParam := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing user_id from token"})
	}

	if err = h.storeUseCase.DetachCategoryFromStore(storeID, categoryID, uint64(uid)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
