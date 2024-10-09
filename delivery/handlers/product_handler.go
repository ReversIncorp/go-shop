package handlers

import (
	"github.com/labstack/echo/v4"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/usecase"
	"net/http"
)

// ProductHandler обрабатывает HTTP-запросы для продуктов
type ProductHandler struct {
	productUseCase usecase.ProductUseCase
}

// NewProductHandler создает новый экземпляр ProductHandler
func NewProductHandler(productUseCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

// CreateProduct обрабатывает запрос на создание продукта
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product entities.Product

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := h.productUseCase.CreateProduct(product); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, product)
}

// GetProductByID обрабатывает запрос на получение продукта по ID
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")

	product, err := h.productUseCase.GetProductByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

// UpdateProduct обрабатывает запрос на обновление продукта
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var product entities.Product

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := h.productUseCase.UpdateProduct(product); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

// DeleteProduct обрабатывает запрос на удаление продукта
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	if err := h.productUseCase.DeleteProduct(id); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetProductsByStore обрабатывает запрос на получение всех продуктов по ID магазина
func (h *ProductHandler) GetProductsByStore(c echo.Context) error {
	storeID := c.Param("store_id")

	products, err := h.productUseCase.GetProductsByStore(storeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}
