package handlers

import (
	"github.com/labstack/echo/v4"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/usecase"
	"net/http"
	"strconv"
)

// StoreHandler обрабатывает HTTP-запросы для магазинов
type StoreHandler struct {
	storeUseCase *usecase.StoreUseCase
}

// NewStoreHandler создает новый экземпляр StoreHandler
func NewStoreHandler(storeUseCase *usecase.StoreUseCase) *StoreHandler {
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

	store.OwnerID = int64(uid)

	if err := h.storeUseCase.CreateStore(store, int64(uid)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, store)
}

// GetStoreByID обрабатывает запрос на получение магазина по ID
func (h *StoreHandler) GetStoreByID(c echo.Context) error {
	id := c.Param("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	store, err := h.storeUseCase.GetStoreByID(int64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, store)
}

// UpdateStore обрабатывает запрос на обновление магазина
func (h *StoreHandler) UpdateStore(c echo.Context) error {
	var store entities.Store

	id := c.Param("id")
	storeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	if err := c.Bind(&store); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	store.ID = storeID

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing user_id from token"})
	}

	if err := h.storeUseCase.UpdateStore(store, int64(uid)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, store)
}

// DeleteStore обрабатывает запрос на удаление магазина
func (h *StoreHandler) DeleteStore(c echo.Context) error {
	id := c.Param("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing user_id from token"})
	}

	if err := h.storeUseCase.DeleteStore(int64ID, int64(uid)); err != nil {
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
