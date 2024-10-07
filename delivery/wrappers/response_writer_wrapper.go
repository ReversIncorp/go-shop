package wrappers

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"marketplace/app/utils"
	"net/http"
)

type ResponseWriterWrapper struct {
	w          http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

// NewResponseWriterWrapper создает обертку для http.ResponseWriter
func NewResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper {
	var buf bytes.Buffer
	return ResponseWriterWrapper{
		w:          w,
		body:       &buf,
		statusCode: http.StatusOK,
	}
}

func (rww *ResponseWriterWrapper) Write(buf []byte) (int, error) {
	// Сначала записываем данные в локальный буфер
	n, err := rww.body.Write(buf)
	if err != nil {
		return n, err // Возвращаем ошибку, если запись в буфер не удалась
	}

	// Затем записываем данные в оригинальный ResponseWriter
	n, err = rww.w.Write(buf)
	return n, err // Возвращаем количество записанных байт и ошибку (если есть)
}

func (rww *ResponseWriterWrapper) Header() http.Header {
	return rww.w.Header() // Возвращаем заголовки
}

func (rww *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rww.statusCode = statusCode   // Устанавливаем статус код
	rww.w.WriteHeader(statusCode) // Вызываем WriteHeader на http.ResponseWriter
}

func (rww *ResponseWriterWrapper) String() string {
	var buf bytes.Buffer
	buf.WriteString("\nResponse: \n")

	buf.WriteString("Headers:")
	headers, err := utils.AutoFormatJSON(rww.Header())
	if err != nil {
		logrus.Errorf("Error formating headers: %v", err)
	}
	buf.WriteString(headers + "\n")
	buf.WriteString(fmt.Sprintf("Status Code: %d\n", rww.statusCode))
	body, err := utils.AutoFormatJSON(rww.body.Bytes())
	if err != nil {
		logrus.Errorf("Error formating body: %v", err)
	}
	buf.WriteString(fmt.Sprintf("Body: %s\n", body))
	return buf.String()
}
