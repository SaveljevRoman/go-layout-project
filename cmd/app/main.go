package main

import (
	"context"
	"github.com/SaveljevRoman/go-layout-project/internal/api"
	"github.com/SaveljevRoman/go-layout-project/internal/config"
	"github.com/SaveljevRoman/go-layout-project/internal/repository/mysql"
	"github.com/SaveljevRoman/go-layout-project/internal/repository/redis"
	"github.com/SaveljevRoman/go-layout-project/internal/service"
	mysqlpkg "github.com/SaveljevRoman/go-layout-project/pkg/mysql"
	redispkg "github.com/SaveljevRoman/go-layout-project/pkg/redis"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Инициализация подключений к БД
	mysqlDB, err := mysqlpkg.NewConnection(cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	redisClient, err := redispkg.NewConnection(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Инициализация репозиториев
	userRepo := mysql.NewUserRepository(mysqlDB)
	userCache := redis.NewUserCache(redisClient)
	productRepo := mysql.NewProductRepository(mysqlDB)
	productCache := redis.NewProductCache(redisClient)
	purchaseRepo := mysql.NewPurchaseRepository(mysqlDB)
	purchaseCache := redis.NewPurchaseCache(redisClient)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo, userCache)
	productService := service.NewProductService(productRepo, productCache)
	purchaseService := service.NewPurchaseService(purchaseRepo, purchaseCache, userService, productService)

	// Запуск фоновых горутин для обновления кеша
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go userService.StartCacheUpdater(ctx, time.Duration(cfg.CacheUpdateInterval)*time.Second)
	go productService.StartCacheUpdater(ctx, time.Duration(cfg.CacheUpdateInterval)*time.Second)
	go purchaseService.StartCacheUpdater(ctx, time.Duration(cfg.CacheUpdateInterval)*time.Second)

	// Инициализация роутера и хендлеров
	router := api.NewRouter(userService, productService, purchaseService)

	// Запуск HTTP сервера
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s", cfg.ServerAddress)

	// Обработка сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Завершение контекста для остановки фоновых задач
	cancel()

	// Остановка HTTP сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited properly")
}
