package handlers

import (
	"marketplace/internal/domain/entities"
	productUsecas "marketplace/internal/domain/usecase/product_usecase"
	errorHandling "marketplace/pkg/error_handling"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"
)

// ProductHandler обрабатывает HTTP-запросы для продуктов.
type ProductHandler struct {
	productUseCase *productUsecas.ProductUseCase
	validator      *validator.Validate
}

// NewProductHandler создает новый экземпляр ProductHandler.
func NewProductHandler(
	productUseCase *productUsecas.ProductUseCase,
	validator *validator.Validate,
) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		validator:      validator,
	}
}

// CreateProduct обрабатывает запрос на создание продукта.
// @Summary Создание продукта
// @Description Создает новый продукт в магазине
// @Tags products
// @Consumes application/json
// @Produces application/json
// @Param store_id path int true "ID магазина"
// @Param product body entities.Product true "Данные продукта"
// @Success 201 {object} entities.Product "Продукт создан"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores/{store_id}/products [post].
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product entities.Product
	if err := c.Bind(&product); err != nil {
		return errorHandling.ErrInvalidInput
	}

	id := c.Param("store_id")
	storeID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	product.StoreID = storeID

	if err = h.productUseCase.CreateProduct(&product); err != nil {
		return tracerr.Wrap(err)
	}

	created, err := h.productUseCase.GetProductByID(product.ID)
	if err != nil {
		return errorHandling.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, created)
}

// GetProductByID обрабатывает запрос на получение продукта по ID.
// @Summary Получение продукта по ID
// @Description Возвращает информацию о продукте по его ID
// @Tags products
// @Produces application/json
// @Param id path int true "ID продукта"
// @Success 200 {object} entities.Product "Информация о продукте"
// @Failure 400 {object} errorhandling.ResponseError "Неверный ID"
// @Failure 404 {object} errorhandling.ResponseError "Продукт не найден"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /products/{id} [get].
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}
	product, err := h.productUseCase.GetProductByID(productID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return c.JSON(http.StatusOK, product)
}

// UpdateProduct обрабатывает запрос на обновление продукта.
// @Summary Обновление продукта
// @Description Обновляет информацию о продукте
// @Tags products
// @Consumes application/json
// @Produces application/json
// @Param id path int true "ID продукта"
// @Param store_id path int true "ID магазина"
// @Param product body entities.Product true "Обновленные данные продукта"
// @Success 200 {object} entities.Product "Продукт обновлен"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 404 {object} errorhandling.ResponseError "Продукт не найден"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /stores/{store_id}/products/{id} [put].
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var product entities.Product

	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	sid := c.Param("store_id")
	storeID, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	if err = c.Bind(&product); err != nil {
		return errorHandling.ErrInvalidInput
	}

	product.StoreID = storeID
	product.ID = productID

	if err = h.productUseCase.UpdateProduct(product); err != nil {
		return tracerr.Wrap(err)
	}

	updated, err := h.productUseCase.GetProductByID(productID)
	if err != nil {
		return errorHandling.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, updated)
}

// DeleteProduct обрабатывает запрос на удаление продукта.
// @Summary Удаление продукта
// @Description Удаляет продукт из системы
// @Tags products
// @Produces application/json
// @Param id path int true "ID продукта"
// @Success 204 {string} string "Продукт удален"
// @Failure 400 {object} errorhandling.ResponseError "Неверный ID"
// @Failure 404 {object} errorhandling.ResponseError "Продукт не найден"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /products/{id} [delete].
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errorHandling.ErrInvalidInput
	}

	if err = h.productUseCase.DeleteProduct(productID); err != nil {
		return tracerr.Wrap(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetProductsByFilters обрабатывает запрос на получение продуктов по фильтрам.
// @Summary Получение продуктов по фильтрам
// @Description Возвращает список продуктов с пагинацией и фильтрацией
// @Tags products
// @Consumes application/json
// @Produces application/json
// @Param searchParams body entities.ProductSearchParams true "Параметры поиска"
// @Success 200 {object} map[string]interface{} "Список продуктов"
// @Failure 400 {object} errorhandling.ResponseError "Ошибка валидации"
// @Failure 500 {object} errorhandling.ResponseError "Внутренняя ошибка сервера"
// @Router /products/search [post].
func (h *ProductHandler) GetProductsByFilters(c echo.Context) error {
	var searchParams entities.ProductSearchParams

	if err := c.Bind(&searchParams); err != nil {
		return errorHandling.ErrInvalidInput
	}
	if err := h.validator.Struct(searchParams); err != nil {
		return errorHandling.ErrInvalidInput
	}

	products, nextCursor, err := h.productUseCase.GetProductsByFilters(searchParams)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       products,
		"limit":      searchParams.Limit,
		"nextCursor": nextCursor,
	})
}
