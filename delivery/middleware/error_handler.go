package middleware

import (
	"marketplace/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorHandlerMiddleware - middleware для обработки ошибок.
func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if appErr, ok := err.(*errors.ErrorResponse); ok {
				return c.JSON(appErr.Code, appErr)
			}
			errors.LogErrorWithTracer(err) // Логируем ошибку со стектрейсом.
			return c.JSON(http.StatusInternalServerError, &errors.ErrorResponse{
				Code:    errors.ErrInternalServerError.Code,
				Details: errors.ErrInternalServerError.Details,
			})
		}
		return nil
	}
}
