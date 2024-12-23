package errorhandling

// ResponseError - структура для представления ошибок и отправки их в ответе.
type ResponseError struct {
	Code    int    `json:"code"`    // HTTP статус-код
	Details string `json:"details"` // Сообщение для пользователя
	Err     error  `json:"-"`       // Внутренняя ошибка, не отправляется клиенту
}

// NewResponseError - конструктор для создания ошибок.
func NewResponseError(code int, details string, err error) *ResponseError {
	return &ResponseError{
		Code:    code,
		Details: details,
		Err:     err,
	}
}

// Error - метод, который возвращает строку ошибки.
func (e *ResponseError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Details
}
