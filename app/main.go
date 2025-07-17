package main

import (
	"fmt"
	_ "marketplace/docs"
	"marketplace/pkg/di"
	errorHandling "marketplace/pkg/error_handling"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// @title Go Shop API
// @version 1.0
// @description REST API для маркетплейса на Go с JWT аутентификацией, управлением магазинами, продуктами и категориями
// @termsOfService http://swagger.io/terms/

// @contact.name Go Shop API Support
// @contact.url https://github.com/ReversIncorp/go-shop
// @contact.email a.savko.developer@goshop.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите "Bearer" за которым следует пробел и JWT токен.

// main.
func main() {
	// Load .env file
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		err := fmt.Sprintf("Error loading .env file.\n%s", err)
		panic(err)
	}
	container := di.Container()
	e := echo.New()

	// Регистрация баз данных
	if err = di.RegisterDatabases(container); err != nil {
		errorHandling.FatalErrorWithTracer("Failed to register databases", err)
	}

	// Регистрация всех зависимостей
	err = di.RegisterDependencies(container)
	if err != nil {
		errorHandling.FatalErrorWithTracer("Failed to register dependencies", err)
	}

	// Регистрация midleware
	if err = di.RegisterMiddleware(container, e); err != nil {
		errorHandling.FatalErrorWithTracer("Failed to register midleware", err)
	}

	// Регистрация маршрутов
	if err = di.RegisterRoutes(container, e); err != nil {
		errorHandling.FatalErrorWithTracer("Failed to register routes", err)
	}
	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
