package di

import (
	"fmt"
	"marketplace/delivery/handlers"
	"marketplace/delivery/middleware"
	"marketplace/internal/data/repository"
	productUsecase "marketplace/internal/domain/usecase/product_usecase"
	storeUsecase "marketplace/internal/domain/usecase/store_usecase"
	userUsecase "marketplace/internal/domain/usecase/user_ucecase"
	"marketplace/pkg/utils"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

var container = dig.New() //nolint:gochecknoglobals

func Container() *dig.Container {
	return container
}

func RegisterDatabases(container *dig.Container) error {
	if err := container.Provide(registerRedisClient); err != nil {
		return err
	}
	return nil
}

func registerRedisClient() (*redis.Client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: "",
			DB:       0,
		},
	)
	pong, err := redisClient.Ping(redisClient.Context()).Result()
	if pong != "PONG" || err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return redisClient, nil
}

func RegisterDependencies(container *dig.Container) error {
	if err := container.Provide(utils.AppValidate); err != nil {
		return err
	}
	// Регистрация логгера
	if err := container.Provide(middleware.AppLoggersSingleton); err != nil {
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
	if err := container.Provide(repository.NewRedisJWTRepository); err != nil {
		return err
	}

	// Регистрация use cases
	if err := container.Provide(userUsecase.NewUserUseCase); err != nil {
		return err
	}
	if err := container.Provide(productUsecase.NewProductUseCase); err != nil {
		return err
	}
	if err := container.Provide(storeUsecase.NewStoreUseCase); err != nil {
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
	var httpLogger *middleware.AppLoggers
	if err := container.Invoke(func(logger *middleware.AppLoggers) {
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

	authorizedScope := e.Group("")
	authorizedScope.Use(middleware.JWTMiddleware)

	// Регистрация маршрутов для пользователей
	e.POST("/users", userHandler.Register)
	e.POST("/users/login", userHandler.Login)

	// Регистрация маршрутов для продуктов
	authorizedScope.POST("/products", productHandler.CreateProduct)
	authorizedScope.GET("/products/:id", productHandler.GetProductByID)
	authorizedScope.PUT("/products/:id", productHandler.UpdateProduct)
	authorizedScope.DELETE("/products/:id", productHandler.DeleteProduct)
	authorizedScope.GET("/stores/:store_id/products", productHandler.GetProductsByStore)

	// Регистрация маршрутов для магазинов
	authorizedScope.POST("/stores", storeHandler.CreateStore)
	authorizedScope.GET("/stores/:id", storeHandler.GetStoreByID)
	authorizedScope.PUT("/stores/:id", storeHandler.UpdateStore)
	authorizedScope.DELETE("/stores/:id", storeHandler.DeleteStore)
	authorizedScope.GET("/stores", storeHandler.GetAllStores)
	return nil
}
