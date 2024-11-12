package errors

import (
	"net/http"
)

// Предопределённые ошибки.
var (
	ErrStoreNotFound = NewErrorResponse(http.StatusNotFound, "Store not found", nil)

	ErrProductNotFound          = NewErrorResponse(http.StatusNotFound, "Product not found", nil)
	ErrProductNotBelongsToStore = NewErrorResponse(http.StatusBadRequest, "Product not belongs to this store", nil)

	ErrCategoryNotFound    = NewErrorResponse(http.StatusNotFound, "Category not found", nil)
	ErrCategoryAttached    = NewErrorResponse(http.StatusBadRequest, "Category already attached", nil)
	ErrCategoryNotAttached = NewErrorResponse(http.StatusBadRequest, "Category not attached", nil)

	ErrUserNotFound       = NewErrorResponse(http.StatusNotFound, "User not found", nil)
	ErrUserNotSeller      = NewErrorResponse(http.StatusBadRequest, "User is not a seller", nil)
	ErrUserExists         = NewErrorResponse(http.StatusBadRequest, "User already exists", nil)
	ErrInvalidCredentials = NewErrorResponse(http.StatusUnauthorized, "Invalid credentials", nil)

	ErrInvalidInput        = NewErrorResponse(http.StatusBadRequest, "Invalid input", nil)
	ErrValidationFailed    = NewErrorResponse(http.StatusNotFound, "Input validation failed", nil)
	ErrUnexpectedError     = NewErrorResponse(http.StatusInternalServerError, "Unexpected error occupied", nil)
	ErrInternalServerError = NewErrorResponse(http.StatusInternalServerError, "Internal server error", nil)

	ErrUnauthorizedAccess = NewErrorResponse(http.StatusUnauthorized, "Unauthorized access: invalid user", nil)
	ErrUserNotAdminStore  = NewErrorResponse(http.StatusForbidden, "Forbidden: user is not admin of this store", nil)

	ErrMissingUserFromToken = NewErrorResponse(http.StatusBadRequest, "Invalid or missing user from token", nil)
	ErrMissingToken         = NewErrorResponse(http.StatusUnauthorized, "Missing token", nil)
	ErrInvalidTokenFormat   = NewErrorResponse(http.StatusUnauthorized, "Invalid token format", nil)
	ErrInvalidExpiredToken  = NewErrorResponse(http.StatusUnauthorized, "Invalid or expired token", nil)
	ErrInvalidTokenClaims   = NewErrorResponse(http.StatusUnauthorized, "Invalid token claims", nil)
)
