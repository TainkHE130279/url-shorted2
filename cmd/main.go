package main

import (
	"log"
	"os"
	"time"

	"github.com/url-shorted2/internal/domain/entities"
	"github.com/url-shorted2/internal/infrastructure/handlers"
	"github.com/url-shorted2/internal/infrastructure/middleware"
	"github.com/url-shorted2/internal/infrastructure/repositories"
	"github.com/url-shorted2/internal/infrastructure/routes"
	"github.com/url-shorted2/internal/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Khởi tạo Gin router
	router := gin.New()

	// Thêm middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Khởi tạo database
	db, err := initDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Khởi tạo dependencies theo Clean Architecture
	// 1. Infrastructure layer (repositories)
	urlRepo := repositories.NewURLRepositoryImpl(db)

	// 2. Use case layer
	baseURL := getBaseURL()
	urlUsecase := usecases.NewURLUsecase(urlRepo, baseURL)

	// 3. Infrastructure layer (handlers)
	urlHandler := handlers.NewURLHandler(urlUsecase)

	// 4. Setup routes
	routes.SetupRoutes(router, urlHandler)

	// Khởi động server
	port := getPort()
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDatabase khởi tạo database và migrate schema
func initDatabase() (*gorm.DB, error) {
	// Sử dụng SQLite cho demo, có thể thay bằng PostgreSQL/MySQL
	db, err := gorm.Open(sqlite.Open("url_shortener.db"), &gorm.Config{
		// Tối ưu cho high concurrency
		PrepareStmt: true, // Pre-compile statements
	})
	if err != nil {
		return nil, err
	}

	// Cấu hình connection pool cho SQLite
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Tối ưu connection pool
	sqlDB.SetMaxOpenConns(100)          // Tăng số connection tối đa
	sqlDB.SetMaxIdleConns(50)           // Tăng số connection idle
	sqlDB.SetConnMaxLifetime(time.Hour) // Thời gian sống của connection

	// Auto migrate schema
	err = db.AutoMigrate(
		&entities.URL{},
		&entities.Analytics{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getBaseURL lấy base URL từ environment variable
func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return baseURL
}

// getPort lấy port từ environment variable
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
