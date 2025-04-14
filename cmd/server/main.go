package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/volchok96/todoapp/docs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	pg "github.com/volchok96/todoapp/internal/db"
	"github.com/volchok96/todoapp/internal/domain"
	httpDelivery "github.com/volchok96/todoapp/internal/http"
	"github.com/volchok96/todoapp/internal/usecase"
	"go.uber.org/zap"
)

// @title TodoApp API
// @version 1.0
// @description This is a simple todo application
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "local" || appEnv == "" {
		fmt.Println("Running in LOCAL mode")
		_ = godotenv.Load(".env.local")
	} else {
		fmt.Println("Running in DOCKER/DEFAULT mode")
		_ = godotenv.Load(".env")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Получение переменных окружения с fallback значениями
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "todoapp")

	// Формирование DSN строки
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}

	db.AutoMigrate(&domain.Task{})

	repo := pg.NewTaskRepo(db, logger)
	uc := usecase.NewTaskUsecase(repo)

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	httpDelivery.RegisterRoutes(app, uc, logger)

	go func() {
		logger.Info("Starting server", zap.String("port", "8080"))
		if err := app.Listen(":8080"); err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server exiting")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
