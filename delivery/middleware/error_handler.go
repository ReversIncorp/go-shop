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
		if appErr, ok := err.(*errorHandling.ErrorResponse); ok {
			return c.JSON(appErr.Code, appErr)
		}
		errorHandling.LogErrorWithTracer(err) // Логируем ошибку со стектрейсом.
		return c.JSON(http.StatusInternalServerError, &errorHandling.ErrorResponse{
			Code:    errorHandling.ErrInternalServerError.Code,
			Details: errorHandling.ErrInternalServerError.Details,
		})

	}
}
