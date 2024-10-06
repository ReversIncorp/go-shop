package midleware

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"marketplace/app/utils"
	"net/http"
	"sync"
	"time"
)

var appLoggers *AppLoggers
var once sync.Once

type AppLoggers struct {
	requestLogger  *logrus.Logger
	responseLogger *logrus.Logger
}

func AppLoggersSingleton() *AppLoggers {
	if appLoggers == nil {
		once.Do(func() {
			//Request
			requestLogger := logrus.New()
			requestLogger.Level = logrus.DebugLevel
			requestLogger.WithTime(time.Now())
			//Response
			responseLogger := logrus.New()
			responseLogger.Level = logrus.DebugLevel
			responseLogger.WithTime(time.Now())

			appLoggers = &AppLoggers{requestLogger: requestLogger, responseLogger: responseLogger}
		})
	}
	return appLoggers
}

func (l *AppLoggers) LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Log Request
		req := c.Request()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			l.requestLogger.Errorf("Error reading request body: %v", err)
			return err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset the request body
		//Formating JSON
		formattedBody, err := utils.AutoFormatJSON(body)
		if err != nil {
			l.requestLogger.Errorf("Error formating request body: %v", err)
			return err
		}
		formattedHeaders, err := utils.AutoFormatJSON(req.Header)
		//Formating JSON
		if err != nil {
			l.requestLogger.Errorf("Error formating request body: %v", err)
			return err
		}
		l.requestLogger.Infof("\nRequest: %s %s\nHeaders:%s \nBody: %s\n",
			req.Method,
			req.URL.String(),
			formattedHeaders,
			formattedBody,
		)

		// Call the next handler
		err = next(c)
		return err
	}
}

func (l *AppLoggers) LoggingResponseMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		response := c.Response()
		originalWriter := response.Writer
		responseWrapper := NewResponseWriterWrapper(originalWriter)
		response.Writer = &responseWrapper

		// Выполняем следующий обработчик
		err := next(c)
		//response.After(func() {
		l.responseLogger.Infof(responseWrapper.String())
		//})
		if err != nil {
			return err
		}
		return err
	}
}

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
		appLoggers.responseLogger.Errorf("Error formating headers: %v", err)
	}
	buf.WriteString(headers + "\n")
	buf.WriteString(fmt.Sprintf("Status Code: %d\n", rww.statusCode))
	body, err := utils.AutoFormatJSON(rww.body.Bytes())
	if err != nil {
		appLoggers.responseLogger.Errorf("Error formating body: %v", err)
	}
	buf.WriteString(fmt.Sprintf("Body: %s\n", body))
	return buf.String()
}
