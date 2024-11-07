package middleware

import (
	usecase "marketplace/internal/domain/usecase/store_usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// StoreAdminMiddleware проверяет, является ли пользователь администратором или владельцем магазина.
func StoreAdminMiddleware(storeUsecase *usecase.StoreUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get("user_id").(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized access: invalid user"})
			}
			uid := uint64(userID)

			storeIDParam := c.Param("store_id")
			storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid store ID"})
			}

			isAdmin, err := storeUsecase.IsUserStoreAdmin(storeID, uid)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check store admin status"})
			}
			if !isAdmin {
				return c.JSON(http.StatusForbidden, echo.Map{"error": "Forbidden: user is not an admin of this store"})
			}

			return next(c)
		}
	}
}
