package main

import (
	"context"
	_ "effective_mobile/docs"
	"effective_mobile/internal/api"
	"effective_mobile/internal/config"
	"effective_mobile/internal/repository"
	"effective_mobile/internal/service"
	"effective_mobile/pkg/logger_module"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/pressly/goose"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func main() {
	// 1. Инициализация логгера записывает в файл и в stdout
	log_file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed open log file", err)
	}
	defer log_file.Close()

	multi_writer := io.MultiWriter(log_file, os.Stdout)

	logger := logger_module.New(multi_writer, "[APP]", log.LstdFlags|log.Lshortfile)

	logger.Info("Logger starter")
	// 2. Загрузка конфигурации
	conf, err := config.Load_Config_PG()
	if err != nil {
		logger.Fatal("Failed to load config", "error", err)
	}

	// 3. Подключение к PostgreSQL
	db, err := repository.NewConnectPostgresDB(conf)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}

	// 4. Применение миграций
	if err := applyMigrations(db); err != nil {
		logger.Fatal("Failed to apply migrations", "error", err)
	}

	// 5. Инициализация слоёв приложения
	gorm_repo := repository.NewGormRepo(db)
	subService := service.NewSubciptionService(gorm_repo)
	subHandler := api.NewSubciptionHandler(subService, logger)

	// 6. Настройка роутера
	router := mux.NewRouter()
	CreateRoutes(router, subHandler)

	// 7. Настройка HTTP-сервера
	server := &http.Server{
		Addr:         ":" + conf.Http_Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 8. Запускаем сервер в отдельной горутине чтобы не заблочить основной поток
	go func() {
		logger.Info("starting server", "port", conf.Http_Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", "error", err)
		}
	}()

	// Ожидание сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Завершение работы с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown error", "error", err)
	}

	logger.Info("Server stopped gracefully")
}

// applyMigrations применяет миграции к БД
func applyMigrations(db *gorm.DB) error {
	logger := logger_module.Get()
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get sql.DB: %w", err)
	}
	// Указываем dialect и папку с миграциями
	goose.SetDialect("postgres")
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		logger.Fatal("failed to apply migrations: %w", err)
	}

	return nil

}

// Регистрируем все HTTP-роуты
func CreateRoutes(router *mux.Router, handler *api.SubscriptionHandler) {
	// Добавляем Swagger UI к роутеру
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Добавляем префикс для работы с endpoints
	api := router.PathPrefix("/api/").Subrouter()
	handler.RegisterRouter(api)
}
