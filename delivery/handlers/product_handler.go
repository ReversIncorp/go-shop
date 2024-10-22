package handlers

import (
	"github.com/labstack/echo/v4"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/usecase"
	"net/http"
	"strconv"
)

// ProductHandler обрабатывает HTTP-запросы для продуктов
type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

// NewProductHandler создает новый экземпляр ProductHandler
func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

// CreateProduct обрабатывает запрос на создание продукта
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product entities.Product

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing user_id from token"})
	}

	if err := h.productUseCase.CreateProduct(product, int64(uid)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, product)
}

// GetProductByID обрабатывает запрос на получение продукта по ID
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	product, err := h.productUseCase.GetProductByID(int64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

// UpdateProduct обрабатывает запрос на обновление продукта
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var product entities.Product
	product.ID = int64ID

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing user_id from token"})
	}

	if err := h.productUseCase.UpdateProduct(product, int64(uid)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, product)
}

// DeleteProduct обрабатывает запрос на удаление продукта
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userID := c.Get("user_id")
	uid, ok := userID.(float64)
	if !ok || userID == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing user_id from token"})
	}

	if err := h.productUseCase.DeleteProduct(int64ID, int64(uid)); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// GetProductsByStore обрабатывает запрос на получение всех продуктов по ID магазина
func (h *ProductHandler) GetProductsByStore(c echo.Context) error {
	storeID := c.Param("store_id")
	int64ID, err := strconv.ParseInt(storeID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	products, err := h.productUseCase.GetProductsByStore(int64ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}

// GetProductsByStoreAndCategory обрабатывает запрос на получение всех продуктов по ID магазина и ID категории
func (h *ProductHandler) GetProductsByStoreAndCategory(c echo.Context) error {
	storeID := c.Param("store_id")
	int64StoreID, err := strconv.ParseInt(storeID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	categoryID := c.Param("category_id")
	int64CategoryID, err := strconv.ParseInt(categoryID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	products, err := h.productUseCase.GetProductsByStoreAndCategory(int64StoreID, int64CategoryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}
