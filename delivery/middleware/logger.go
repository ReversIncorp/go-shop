package middleware

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"marketplace/delivery/wrappers"
	"marketplace/pkg/utils"
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

func (l *AppLoggers) LoggingRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
		responseWrapper := wrappers.NewResponseWriterWrapper(originalWriter)
		response.Writer = &responseWrapper

		// Логируем после обработки ответа
		response.After(func() {
			l.responseLogger.Infof(responseWrapper.String())
		})
		// Выполняем следующий обработчик
		err := next(c)
		if err != nil {
			return err
		}
		return err
	}
}
