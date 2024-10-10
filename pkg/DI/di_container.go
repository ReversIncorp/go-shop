package DI

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
	"marketplace/delivery/handlers"
	"marketplace/delivery/midleware"
	"marketplace/internal/data/repository"
	"marketplace/internal/domain/usecase"
	"marketplace/pkg/utils"
	"os"
)

var container = dig.New()

func Container() *dig.Container {
	return container
}

func RegisterDependencies(container *dig.Container) error {
	if err := container.Provide(utils.AppValidate); err != nil {
		return err
	}
	// Регистрация логгера
	if err := container.Provide(midleware.AppLoggersSingleton); err != nil {
		return err
	}
	// Регистрация репозиториев
	if err := container.Provide(repository.NewUserRepository); err != nil {
		return err
	}
	if err := container.Provide(repository.NewProductRepository); err != nil {
		return err
	}
	if err := container.Provide(repository.NewStoreRepository); err != nil {
		return err
	}

	// Регистрация use cases
	if err := container.Provide(usecase.NewUserUseCase); err != nil {
		return err
	}
	if err := container.Provide(usecase.NewProductUseCase); err != nil {
		return err
	}
	if err := container.Provide(usecase.NewStoreUseCase); err != nil {
		return err
	}

	// Регистрация обработчиков
	if err := container.Provide(handlers.NewUserHandler); err != nil {
		return err
	}
	if err := container.Provide(handlers.NewProductHandler); err != nil {
		return err
	}
	if err := container.Provide(handlers.NewStoreHandler); err != nil {
		return err
	}

	return nil
}

func RegisterMiddleware(container *dig.Container, e *echo.Echo) error {
	// Используем логгер из контейнера
	var httpLogger *midleware.AppLoggers
	if err := container.Invoke(func(logger *midleware.AppLoggers) {
		httpLogger = logger
	}); err != nil {
		return fmt.Errorf("failed to invoke logger: %w", err)
	}

	// Добавляем midleware для логирования в зависимости от окружения
	if os.Getenv("APP_ENV") == "Dev" {
		e.Use(httpLogger.LoggingRequestMiddleware)
		e.Use(httpLogger.LoggingResponseMiddleware)
	}

	return nil
}

func RegisterRoutes(container *dig.Container, e *echo.Echo) error {
	// Инициализация HTTP-хэндлеров
	var userHandler *handlers.UserHandler
	var productHandler *handlers.ProductHandler
	var storeHandler *handlers.StoreHandler

	// Получаем хэндлеры через контейнер
	if err := container.Invoke(func(uh *handlers.UserHandler, ph *handlers.ProductHandler, sh *handlers.StoreHandler) {
		userHandler = uh
		productHandler = ph
		storeHandler = sh
	}); err != nil {
		fmt.Printf("Failed to invoke handlers: %v\n", err)
		return err
	}

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
	return nil
}
