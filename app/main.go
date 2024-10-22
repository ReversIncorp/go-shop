package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"marketplace/pkg/DI"
)

func main() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error loading .env file")
		return
	}
	container := DI.Container()
	e := echo.New()

	// Регистрация базы данных
	if err := DI.RegisterDatabases(container); err != nil {
		fmt.Printf("Failed to register(or connect) db: %v\n", err)
		return
	}

	// Регистрация всех зависимостей
	if err := DI.RegisterDependencies(container); err != nil {
		fmt.Printf("Failed to register dependencies: %v\n", err)
		return
	}

	// Регистрация midleware
	if err := DI.RegisterMiddleware(container, e); err != nil {
		fmt.Printf("Failed to register midleware: %v\n", err)
		return
	}

	if err := DI.RegisterRoutes(container, e); err != nil {
		fmt.Printf("Failed to register routes: %v\n", err)
		return
	}
	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
