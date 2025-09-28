package main

import (
	"log"

	"github.com/url-shorted2/internal/config"
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
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Khởi tạo Gin router
	router := gin.New()

	// Thêm middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Khởi tạo database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Khởi tạo dependencies theo Clean Architecture
	// 1. Infrastructure layer (repositories)
	urlRepo := repositories.NewURLRepositoryImpl(db)

	// 2. Use case layer
	urlUsecase := usecases.NewURLUsecase(urlRepo, cfg.Server.BaseURL, cfg)

	// 3. Infrastructure layer (handlers)
	urlHandler := handlers.NewURLHandler(urlUsecase)

	// 4. Setup routes
	routes.SetupRoutes(router, urlHandler)

	// Khởi động server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDatabase khởi tạo database và migrate schema
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Sử dụng SQLite cho demo, có thể thay bằng PostgreSQL/MySQL
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{
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

	// Tối ưu connection pool từ config
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

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
