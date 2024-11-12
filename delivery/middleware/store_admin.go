package middleware

import (
	usecase "marketplace/internal/domain/usecase/store_usecase"
	"marketplace/pkg/errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

// StoreAdminMiddleware проверяет, является ли пользователь администратором или владельцем магазина.
func StoreAdminMiddleware(storeUsecase *usecase.StoreUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get("user_id").(float64)
			if !ok {
				return errors.ErrUnauthorizedAccess
			}
			uid := uint64(userID)

			storeIDParam := c.Param("store_id")
			storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
			if err != nil {
				return errors.ErrInvalidInput
			}

			isAdmin, err := storeUsecase.IsUserStoreAdmin(storeID, uid)
			if err != nil {
				return err
			}
			if !isAdmin {
				return errors.ErrUserNotAdminStore
			}

			return next(c)
		}
	}
}
