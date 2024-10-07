package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"marketplace/delivery/handlers"
	"marketplace/delivery/midleware"
	"marketplace/internal/data/repository"
	"marketplace/internal/domain/usecase"
	"os"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	e := echo.New()

	httpLogger := midleware.AppLoggersSingleton()
	if os.Getenv("APP_ENV") == "Dev" {
		e.Use(httpLogger.LoggingRequestMiddleware)
		e.Use(httpLogger.LoggingResponseMiddleware)
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository()       // Репозиторий пользователей
	productRepo := repository.NewProductRepository() // Репозиторий продуктов
	storeRepo := repository.NewStoreRepository()     // Репозиторий магазинов

	// Инициализация use cases
	userUseCase := usecase.NewUserUseCase(userRepo)          // Use case для пользователей
	productUseCase := usecase.NewProductUseCase(productRepo) // Use case для продуктов
	storeUseCase := usecase.NewStoreUseCase(storeRepo)       // Use case для магазинов

	// Инициализация HTTP-хэндлеров
	userHandler := handlers.NewUserHandler(userUseCase)          // Обработчик для пользователей
	productHandler := handlers.NewProductHandler(productUseCase) // Обработчик для продуктов
	storeHandler := handlers.NewStoreHandler(storeUseCase)       // Обработчик для магазинов

	// Регистрация маршрутов для пользователей
	e.POST("/users", userHandler.Register)
	e.POST("/users/login", userHandler.Login)

	// Регистрация маршрутов для продуктов
	e.POST("/products", productHandler.CreateProduct)
	e.GET("/products/:id", productHandler.GetProductByID)
	e.PUT("/products/:id", productHandler.UpdateProduct)
	e.DELETE("/products/:id", productHandler.DeleteProduct)
	e.GET("/stores/:store_id/products", productHandler.GetProductsByStore)

	// Регистрация маршрутов для магазинов
	e.POST("/stores", storeHandler.CreateStore)
	e.GET("/stores/:id", storeHandler.GetStoreByID)
	e.PUT("/stores/:id", storeHandler.UpdateStore)
	e.DELETE("/stores/:id", storeHandler.DeleteStore)
	e.GET("/stores", storeHandler.GetAllStores)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
