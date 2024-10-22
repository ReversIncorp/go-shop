package handlers

import (
	"marketplace/internal/domain/entities"
	productUsecas "marketplace/internal/domain/usecase/product_usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ProductHandler обрабатывает HTTP-запросы для продуктов.
type ProductHandler struct {
	productUseCase *productUsecas.ProductUseCase
}

// NewProductHandler создает новый экземпляр ProductHandler.
func NewProductHandler(productUseCase *productUsecas.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

// CreateProduct обрабатывает запрос на создание продукта.
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

// GetProductByID обрабатывает запрос на получение продукта по ID.
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	product, err := h.productUseCase.GetProductByID(uint64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

// UpdateProduct обрабатывает запрос на обновление продукта.
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

// DeleteProduct обрабатывает запрос на удаление продукта.
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := h.productUseCase.DeleteProduct(uint64ID); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetProductsByStore обрабатывает запрос на получение всех продуктов по ID магазина.
func (h *ProductHandler) GetProductsByStore(c echo.Context) error {
	storeID := c.Param("store_id")
	uint64ID, err := strconv.ParseUint(storeID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	products, err := h.productUseCase.GetProductsByStore(uint64ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}
