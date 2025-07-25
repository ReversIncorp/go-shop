package handlers

import (
	"marketplace/internal/domain/entities"
	storeUsecases "marketplace/internal/domain/usecase/store_usecase"
	errorHandling "marketplace/pkg/error_handling"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"
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
// @Summary Создание магазина
// @Description Создает новый магазин для пользователя
// @Tags stores
// @Consumes application/json
// @Produces application/json
// @Param store body entities.Store true "Данные магазина"
// @Success 201 {object} entities.Store "Магазин создан"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 401 {object} errorhandling.ResponseError "Не авторизован"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores [post].
func (h *StoreHandler) CreateStore(c echo.Context) error {
	var store entities.Store

	if err := c.Bind(&store); err != nil {
		return errorHandling.ErrInvalidInput
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return errorHandling.ErrMissingUserFromToken
	}

	if err := h.storeUseCase.CreateStore(&store, uint64(uid)); err != nil {
		return errorHandling.ErrInternalServerError
	}

	created, err := h.storeUseCase.GetStoreByID(store.ID)
	if err != nil {
		return errorHandling.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, created)
}

// GetStoreByID обрабатывает запрос на получение магазина по ID.
// @Summary Получение магазина по ID
// @Description Возвращает информацию о магазине по его ID
// @Tags stores
// @Produces application/json
// @Param store_id path int true "ID магазина"
// @Success 200 {object} entities.Store "Информация о магазине"
// @Failure 400 {object} errorhandling.ResponseError "Неверный ID"
// @Failure 404 {object} errorhandling.ResponseError "Магазин не найден"
// @Router /stores/{store_id} [get].
func (h *StoreHandler) GetStoreByID(c echo.Context) error {
	id := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}
	store, err := h.storeUseCase.GetStoreByID(uint64ID)
	if err != nil {
		return errorHandling.ErrStoreNotFound
	}

	return c.JSON(http.StatusOK, store)
}

// UpdateStore обрабатывает запрос на обновление магазина.
// @Summary Обновление магазина
// @Description Обновляет информацию о магазине
// @Tags stores
// @Consumes application/json
// @Produces application/json
// @Param store_id path int true "ID магазина"
// @Param store body entities.Store true "Обновленные данные магазина"
// @Success 200 {object} entities.Store "Магазин обновлен"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 404 {object} errorhandling.ResponseError "Магазин не найден"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores/{store_id} [put].
func (h *StoreHandler) UpdateStore(c echo.Context) error {
	var store entities.Store
	id := c.Param("store_id")
	storeID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	if err = c.Bind(&store); err != nil {
		return errorHandling.ErrInvalidInput
	}

	store.ID = storeID

	if err = h.storeUseCase.UpdateStore(store); err != nil {
		return errorHandling.ErrInternalServerError
	}

	updated, err := h.storeUseCase.GetStoreByID(storeID)
	if err != nil {
		return errorHandling.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, updated)
}

// DeleteStore обрабатывает запрос на удаление магазина.
// @Summary Удаление магазина
// @Description Удаляет магазин из системы
// @Tags stores
// @Produces application/json
// @Param store_id path int true "ID магазина"
// @Success 204 {string} string "Магазин удален"
// @Failure 400 {object} errorhandling.ResponseError "Неверный ID"
// @Failure 404 {object} errorhandling.ResponseError "Магазин не найден"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores/{store_id} [delete].
func (h *StoreHandler) DeleteStore(c echo.Context) error {
	id := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	if err = h.storeUseCase.DeleteStore(uint64ID); err != nil {
		return tracerr.Wrap(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetStoresByFilters обрабатывает запрос на получение магазинов по фильтрам.
// @Summary Получение магазинов по фильтрам
// @Description Возвращает список магазинов с пагинацией и фильтрацией
// @Tags stores
// @Consumes application/json
// @Produces application/json
// @Param searchParams body entities.StoreSearchParams true "Параметры поиска"
// @Success 200 {object} map[string]interface{} "Список магазинов"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores/search [post].
func (h *StoreHandler) GetStoresByFilters(c echo.Context) error {
	var searchParams entities.StoreSearchParams

	if err := c.Bind(&searchParams); err != nil {
		return errorHandling.ErrInvalidInput
	}
	if err := h.validator.Struct(searchParams); err != nil {
		return errorHandling.ErrInvalidInput
	}

	products, nextCursor, err := h.storeUseCase.GetStoresByFilters(searchParams)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       products,
		"limit":      searchParams.Limit,
		"nextCursor": nextCursor,
	})
}

// AttachCategoryToStore связывает категорию с магазином.
// @Summary Привязка категории к магазину
// @Description Связывает категорию с магазином
// @Tags stores
// @Consumes application/json
// @Produces application/json
// @Param store_id path int true "ID магазина"
// @Param request body object true "ID категории" {"category_id": "int"}
// @Success 200 {object} map[string]interface{} "Категория привязана"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores/{store_id}/categories [post].
func (h *StoreHandler) AttachCategoryToStore(c echo.Context) error {
	var request struct {
		CategoryID uint64 `json:"category_id" validate:"required"`
	}

	storeIDParam := c.Param("store_id")
	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	if err = c.Bind(&request); err != nil || request.CategoryID == 0 {
		return errorHandling.ErrInvalidInput
	}

	if err = h.storeUseCase.AttachCategoryToStore(storeID, request.CategoryID); err != nil {
		return tracerr.Wrap(err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Category attached to store successfully"})
}

// DetachCategoryFromStore отвязывает категорию от магазина.
// @Summary Отвязка категории от магазина
// @Description Отвязывает категорию от магазина
// @Tags stores
// @Produces application/json
// @Param store_id path int true "ID магазина"
// @Param category_id path int true "ID категории"
// @Success 204 {string} string "Категория отвязана"
// @Failure 400 {object} errorhandling.ResponseError "Неверные ID"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores/{store_id}/categories/{category_id} [delete].
func (h *StoreHandler) DetachCategoryFromStore(c echo.Context) error {
	storeIDParam := c.Param("store_id")
	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	categoryIDParam := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	if err = h.storeUseCase.DetachCategoryFromStore(storeID, categoryID); err != nil {
		return tracerr.Wrap(err)
	}

	return c.NoContent(http.StatusNoContent)
}
