package main

import (
	"fmt"
	"marketplace/pkg/di"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env file
	var err error
	err = godotenv.Load("../.env")
	if err != nil {
		err := fmt.Sprintf("Error loading .env file.\n%s", err)
		panic(err)
	}
	container := di.Container()
	e := echo.New()

	// Регистрация всех зависимостей
	err = di.RegisterDependencies(container)
	if err != nil {
		panic(fmt.Sprintf("Failed to register dependencies: %v\n", err))
	}

	// Регистрация midleware
	if err = di.RegisterMiddleware(container, e); err != nil {
		err := fmt.Sprintf("Failed to register midleware: %v\n", err)
		panic(err)
	}

	if err = di.RegisterRoutes(container, e); err != nil {
		panic(fmt.Sprintf("Failed to register routes: %v\n", err))
	}
	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
