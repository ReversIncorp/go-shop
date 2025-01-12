package errorhandling

import (
	"net/http"
)

// Предопределённые ошибки.
var (
	ErrStoreNotFound = NewResponseError(http.StatusNotFound, "Store not found", nil)

	ErrProductNotFound          = NewResponseError(http.StatusNotFound, "Product not found", nil)
	ErrProductNotBelongsToStore = NewResponseError(http.StatusBadRequest, "Product not belongs to this store", nil)

	ErrCategoryNotFound    = NewResponseError(http.StatusNotFound, "Category not found", nil)
	ErrCategoryAttached    = NewResponseError(http.StatusBadRequest, "Category already attached", nil)
	ErrCategoryNotAttached = NewResponseError(http.StatusBadRequest, "Category not attached", nil)

	ErrUserNotFound       = NewResponseError(http.StatusNotFound, "User not found", nil)
	ErrUserNotSeller      = NewResponseError(http.StatusBadRequest, "User is not a seller", nil)
	ErrUserExists         = NewResponseError(http.StatusBadRequest, "User already exists", nil)
	ErrInvalidCredentials = NewResponseError(http.StatusUnauthorized, "Invalid credentials", nil)

	ErrInvalidInput        = NewResponseError(http.StatusBadRequest, "Invalid input", nil)
	ErrValidationFailed    = NewResponseError(http.StatusNotFound, "Input validation failed", nil)
	ErrUnexpectedError     = NewResponseError(http.StatusInternalServerError, "Unexpected error occupied", nil)
	ErrInternalServerError = NewResponseError(http.StatusInternalServerError, "Internal server error", nil)

	ErrUnauthorizedAccess = NewResponseError(http.StatusUnauthorized, "Unauthorized access: invalid user", nil)
	ErrUserNotAdminStore  = NewResponseError(http.StatusForbidden, "Forbidden: user is not admin of this store", nil)

	ErrMissingUserFromToken = NewResponseError(http.StatusBadRequest, "Invalid or missing user from token", nil)
	ErrInvalidTokenType     = NewResponseError(http.StatusUnauthorized, "Invalid token type", nil)
	ErrMissingToken         = NewResponseError(http.StatusUnauthorized, "Missing token", nil)
	ErrInvalidTokenFormat   = NewResponseError(http.StatusUnauthorized, "Invalid token format", nil)
	ErrInvalidExpiredToken  = NewResponseError(http.StatusUnauthorized, "Invalid or expired token", nil)
	ErrInvalidTokenClaims   = NewResponseError(http.StatusUnauthorized, "Invalid token claims", nil)
)
