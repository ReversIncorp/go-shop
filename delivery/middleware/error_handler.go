package middleware

import (
	"errors"
	errorHandling "marketplace/pkg/error_handling"
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

		var appErr *errorHandling.ResponseError
		if errors.As(err, &appErr) {
			return c.JSON(appErr.Code, appErr)
		}

		errorHandling.LogErrorWithTracer(err) // Логируем ошибку со стектрейсом.
		return c.JSON(http.StatusInternalServerError, &errorHandling.ResponseError{
			Code:    errorHandling.ErrInternalServerError.Code,
			Details: errorHandling.ErrInternalServerError.Details,
		})
	}
}
