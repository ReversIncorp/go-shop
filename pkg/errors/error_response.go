package errors

// ErrorResponse - структура для представления ошибок и отправки их в ответе.
type ErrorResponse struct {
	Code    int    `json:"code"`    // HTTP статус-код
	Details string `json:"details"` // Сообщение для пользователя
	Err     error  `json:"-"`       // Внутренняя ошибка, не отправляется клиенту
}

// NewErrorResponse - конструктор для создания ошибок.
func NewErrorResponse(code int, details string, err error) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Details: details,
		Err:     err,
	}
}

// Error - метод, который возвращает строку ошибки.
func (e *ErrorResponse) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Details
}
