package middleware

import (
	"marketplace/pkg/error_handling"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorHandlerMiddleware - middleware для обработки ошибок.
func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}
		if appErr, ok := err.(*error_handling.ErrorResponse); ok {
			return c.JSON(appErr.Code, appErr)
		}
		error_handling.LogErrorWithTracer(err) // Логируем ошибку со стектрейсом.
		return c.JSON(http.StatusInternalServerError, &error_handling.ErrorResponse{
			Code:    error_handling.ErrInternalServerError.Code,
			Details: error_handling.ErrInternalServerError.Details,
		})

	}
}
