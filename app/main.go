package main

import (
	"fmt"
	"marketplace/pkg/di"
	"marketplace/pkg/error_handling"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

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
