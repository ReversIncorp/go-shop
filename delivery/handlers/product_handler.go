package handlers

import (
	"marketplace/internal/domain/entities"
	productUsecas "marketplace/internal/domain/usecase/product_usecase"
	"marketplace/pkg/error_handling"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product entities.Product
	if err := c.Bind(&product); err != nil {
		return error_handling.ErrInvalidInput
	}

	id := c.Param("store_id")
	storeID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	product.StoreID = storeID

	if err = h.productUseCase.CreateProduct(product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, product)
}

// GetProductByID обрабатывает запрос на получение продукта по ID.
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}
	product, err := h.productUseCase.GetProductByID(productID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

// UpdateProduct обрабатывает запрос на обновление продукта.
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var product entities.Product

	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	sid := c.Param("store_id")
	storeID, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	if err = c.Bind(&product); err != nil {
		return error_handling.ErrInvalidInput
	}

	product.StoreID = storeID
	product.ID = productID

	if err = h.productUseCase.UpdateProduct(product); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

// DeleteProduct обрабатывает запрос на удаление продукта.
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return error_handling.ErrInvalidInput
	}

	if err = h.productUseCase.DeleteProduct(productID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ProductHandler) GetProductsByFilters(c echo.Context) error {
	var searchParams entities.ProductSearchParams

	if err := c.Bind(&searchParams); err != nil {
		return error_handling.ErrInvalidInput
	}
	if err := h.validator.Struct(searchParams); err != nil {
		return error_handling.ErrInvalidInput
	}

	products, nextCursor, err := h.productUseCase.GetProductsByFilters(searchParams)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       products,
		"limit":      searchParams.Limit,
		"nextCursor": nextCursor,
	})
}
