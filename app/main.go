package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func main() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error loading .env file")
		return
	}

	container := dig.New()

	// Регистрация всех зависимостей
	if err := registerDependencies(container); err != nil {
		fmt.Printf("Failed to register dependencies: %v\n", err)
		return
	}

	// Регистрация midleware
	if err := registerMiddleware(container); err != nil {
		fmt.Printf("Failed to register midleware: %v\n", err)
		return
	}

	if _, err := registerRoutes(container); err != nil {
		fmt.Printf("Failed to register routes: %v\n", err)
		return
	}
	if err := container.Invoke(func(e *echo.Echo) {
		// Запуск сервера
		e.Logger.Fatal(e.Start(":8080"))
	}); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
}
