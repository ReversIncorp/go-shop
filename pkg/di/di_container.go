package di

import (
	"database/sql"
	"marketplace/delivery/handlers"
	"marketplace/delivery/middleware"
	"marketplace/internal/data/repository"
	categoryUsecase "marketplace/internal/domain/usecase/category_usecase"
	productUsecase "marketplace/internal/domain/usecase/product_usecase"
	storeUsecase "marketplace/internal/domain/usecase/store_usecase"
	userUsecase "marketplace/internal/domain/usecase/user_usecase"
	"marketplace/pkg/database"
	"marketplace/pkg/utils"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"
	"go.uber.org/dig"
)

var container = dig.New()

func Container() *dig.Container {
	return container
}

func RegisterDatabases(container *dig.Container) error {
	// Создаем экземпляр Redis и проверяем соединение
	redisClient, err := database.NewRedisClient()
	if err != nil {
		return tracerr.Errorf("failed to connect to Redis: %w", err)
	}
	if err = container.Provide(func() *redis.Client { return redisClient }); err != nil {
		return tracerr.Errorf("failed to register Redis client: %w", err)
	}

	// Создаем экземпляр PostgreSQL и проверяем соединение
	postgresDB, err := database.OpenPostgreSQL()
	if err != nil {
		return tracerr.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	if err = container.Provide(func() *sql.DB { return postgresDB }); err != nil {
		return tracerr.Errorf("failed to register PostgreSQL client: %w", err)
	}

	return nil
}

func RegisterDependencies(container *dig.Container) error {
	if err := container.Provide(utils.AppValidate); err != nil {
		return tracerr.Wrap(err)
	}
	// Регистрация логгера
	if err := container.Provide(middleware.AppLoggersSingleton); err != nil {
		return tracerr.Wrap(err)
	}
	// Регистрация репозиториев
	if err := container.Provide(repository.NewUserRepository); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(repository.NewProductRepository); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(repository.NewStoreRepository); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(repository.NewCategoryRepository); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(repository.NewRedisJWTRepository); err != nil {
		return tracerr.Wrap(err)
	}

	// Регистрация use cases
	if err := container.Provide(userUsecase.NewUserUseCase); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(productUsecase.NewProductUseCase); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(storeUsecase.NewStoreUseCase); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(categoryUsecase.NewCategoryUseCase); err != nil {
		return tracerr.Wrap(err)
	}

	// Регистрация обработчиков
	if err := container.Provide(handlers.NewUserHandler); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(handlers.NewProductHandler); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(handlers.NewStoreHandler); err != nil {
		return tracerr.Wrap(err)
	}
	if err := container.Provide(handlers.NewCategoryHandler); err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func RegisterMiddleware(container *dig.Container, e *echo.Echo) error {
	// Используем логгер из контейнера
	var httpLogger *middleware.AppLoggers
	if err := container.Invoke(func(logger *middleware.AppLoggers) {
		httpLogger = logger
	}); err != nil {
		return tracerr.Errorf("failed to invoke logger: %w", err)
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
	var userUseCase *userUsecase.UserUseCase

	// Получаем хэндлеры и use case через контейнер
	if err := container.Invoke(func(
		uh *handlers.UserHandler,
		ph *handlers.ProductHandler,
		sh *handlers.StoreHandler,
		ch *handlers.CategoryHandler,
		su *storeUsecase.StoreUseCase,
		uu *userUsecase.UserUseCase) {
		userHandler = uh
		productHandler = ph
		storeHandler = sh
		categoryHandler = ch
		storeUseCase = su
		userUseCase = uu
	}); err != nil {
		return tracerr.Errorf("Failed to invoke handlers or use case: %v\n", err)
	}

	// Первичный скоуп, на который вешаем миддлвейр для хендла ошибок и возможно другие общие вещи.
	mainScope := e.Group("")
	mainScope.Use(middleware.ErrorHandlerMiddleware)

	// Основной скоуп для авторизованных юзеров
	authorizedScope := mainScope.Group("")
	authorizedScope.Use(middleware.JWTMiddleware(userUseCase))

	// Скоуп для админов сторов
	storeAdminScope := authorizedScope.Group("/stores/:store_id")
	storeAdminScope.Use(middleware.StoreAdminMiddleware(storeUseCase))

	// Регистрация маршрутов для пользователей
	mainScope.POST("/users", userHandler.Register)
	mainScope.POST("/users/login", userHandler.Login)
	mainScope.POST("/users/refresh-session", userHandler.RefreshSession)
	mainScope.POST("/users/logout", userHandler.Logout)

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
	authorizedScope.GET("/stores", storeHandler.GetStoresByFilters)
	authorizedScope.GET("/stores/:store_id/categories", categoryHandler.GetAllCategoriesByStore)
	// Для админов сторов
	storeAdminScope.PUT("", storeHandler.UpdateStore)
	storeAdminScope.DELETE("", storeHandler.DeleteStore)
	storeAdminScope.POST("/categories", storeHandler.AttachCategoryToStore)
	storeAdminScope.DELETE("/categories/:category_id", storeHandler.DetachCategoryFromStore)

	// Регистрация маршрутов для категорий
	authorizedScope.POST("/categories", categoryHandler.CreateCategory)
	authorizedScope.GET("/categories/:id", categoryHandler.GetCategoryByID)
	authorizedScope.GET("/categories", categoryHandler.GetAllCategories)
	authorizedScope.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	return nil
}
