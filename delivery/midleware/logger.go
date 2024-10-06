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
		res := c.Response()
		res.After(func() {
			responseWrapper := NewResponseWriterWrapper(res.Writer)
			l.responseLogger.Infof(responseWrapper.String())
		})
		err := next(c)
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
	rww.body.Write(buf)     // Записываем в тело
	return rww.w.Write(buf) // Записываем в http.ResponseWriter
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
	for k, v := range rww.w.Header() {
		buf.WriteString(fmt.Sprintf("%s: %v\n", k, v))
	}

	buf.WriteString(fmt.Sprintf("Status Code: %d\n", rww.statusCode))
	buf.WriteString("Body:")
	buf.WriteString(rww.body.String())
	return buf.String()
}
