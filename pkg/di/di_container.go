package di

import (
	"fmt"
	"marketplace/delivery/handlers"
	"marketplace/delivery/middleware"
	"marketplace/internal/data/repository"
	categoryUsecase "marketplace/internal/domain/usecase/category_usecase"
	productUsecase "marketplace/internal/domain/usecase/product_usecase"
	storeUsecase "marketplace/internal/domain/usecase/store_usecase"
	userUsecase "marketplace/internal/domain/usecase/user_ucecase"
	"marketplace/pkg/database"
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
	if err := container.Provide(database.OpenPostgreSQL); err != nil {
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
	if err := container.Provide(repository.NewCategoryRepository); err != nil {
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
	if err := container.Provide(categoryUsecase.NewCategoryUseCase); err != nil {
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
	if err := container.Provide(handlers.NewCategoryHandler); err != nil {
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
	// Инициализация HTTP-хэндлеров и use cases
	var userHandler *handlers.UserHandler
	var productHandler *handlers.ProductHandler
	var storeHandler *handlers.StoreHandler
	var categoryHandler *handlers.CategoryHandler
	var storeUseCase *storeUsecase.StoreUseCase

	// Получаем хэндлеры и use case через контейнер
	if err := container.Invoke(func(
		uh *handlers.UserHandler,
		ph *handlers.ProductHandler,
		sh *handlers.StoreHandler,
		ch *handlers.CategoryHandler,
		su *storeUsecase.StoreUseCase) {
		userHandler = uh
		productHandler = ph
		storeHandler = sh
		categoryHandler = ch
		storeUseCase = su
	}); err != nil {
		fmt.Printf("Failed to invoke handlers or use case: %v\n", err)
		return err
	}

	// Основной скоуп для авторизованных юзеров
	authorizedScope := e.Group("")
	authorizedScope.Use(middleware.JWTMiddleware)

	// Скоуп для админов сторов
	storeAdminScope := authorizedScope.Group("/stores/:store_id")
	storeAdminScope.Use(middleware.StoreAdminMiddleware(storeUseCase))

	// Регистрация маршрутов для пользователей
	e.POST("/users", userHandler.Register)
	e.POST("/users/login", userHandler.Login)

	// Регистрация маршрутов для продуктов
	authorizedScope.GET("/products/:id", productHandler.GetProductByID)
	authorizedScope.GET("/products", productHandler.GetProductsByFilters)
	// Добавление продуктов для админов сторов
	storeAdminScope.POST("/products", productHandler.CreateProduct)
	storeAdminScope.PUT("/products/:id", productHandler.UpdateProduct)
	storeAdminScope.DELETE("/products/:id", productHandler.DeleteProduct)

	// Регистрация маршрутов для магазинов
	authorizedScope.POST("/stores", storeHandler.CreateStore)
	authorizedScope.GET("/stores/:store_id", storeHandler.GetStoreByID)
	authorizedScope.GET("/stores", storeHandler.GetAllStores)
	// Для админов сторов
	storeAdminScope.PUT("", storeHandler.UpdateStore)
	storeAdminScope.DELETE("", storeHandler.DeleteStore)
	authorizedScope.POST("/stores/:store_id/categories", storeHandler.AddCategoryToStore)
	authorizedScope.DELETE("/stores/:store_id/categories/:category_id", storeHandler.DeleteCategoryFromStore)

	// Регистрация маршрутов для категорий
	authorizedScope.POST("/categories", categoryHandler.CreateCategory)
	authorizedScope.GET("/categories/:id", categoryHandler.GetCategoryByID)
	authorizedScope.PUT("/categories/:id", categoryHandler.UpdateCategory)
	authorizedScope.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	authorizedScope.GET("/stores/:store_id/categories", categoryHandler.GetAllCategoriesByStore)
	return nil
}
