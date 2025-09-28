package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/url-shorted2/internal/config"
	"github.com/url-shorted2/internal/domain/entities"
	"github.com/url-shorted2/internal/infrastructure/handlers"
	"github.com/url-shorted2/internal/infrastructure/repositories"
	"github.com/url-shorted2/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Auto migrate
	db.AutoMigrate(&entities.URL{}, &entities.Analytics{})
	return db
}

// getTestConfig tạo config cho test
func getTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port:    "8080",
			BaseURL: "http://localhost:8080",
			GinMode: "test",
		},
		Database: config.DatabaseConfig{
			Type:            "sqlite",
			Path:            ":memory:",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 3600,
		},
		Redis: config.RedisConfig{
			URL:          "redis://localhost:6379",
			PoolSize:     10,
			MinIdleConns: 5,
			MaxRetries:   3,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolTimeout:  4,
		},
		Lock: config.LockConfig{
			MaxTime:    30,
			MaxTryTime: 10,
		},
	}
}

func TestCreateShortURL(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	// Khởi tạo dependencies
	urlRepo := repositories.NewURLRepositoryImpl(db)
	urlUsecase := usecases.NewURLUsecase(urlRepo, "http://localhost:8080", getTestConfig())
	urlHandler := handlers.NewURLHandler(urlUsecase)

	// Tạo router
	router := gin.New()
	router.POST("/api/v1/urls", urlHandler.CreateShortURL)

	// Test case 1: Tạo short URL thành công
	t.Run("Create short URL successfully", func(t *testing.T) {
		requestBody := map[string]string{
			"url": "https://example.com",
		}
		jsonData, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.CreateURLResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.ShortCode)
		assert.Equal(t, "https://example.com", response.OriginalURL)
		assert.Contains(t, response.ShortURL, "http://localhost:8080/")
	})

	// Test case 2: Tạo short URL với URL khác
	t.Run("Create short URL with different URL", func(t *testing.T) {
		requestBody := map[string]string{
			"url": "https://google.com",
		}
		jsonData, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.CreateURLResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.ShortCode) // Short code được generate tự động
		assert.Equal(t, "https://google.com", response.OriginalURL)
	})

	// Test case 3: Invalid URL
	t.Run("Create short URL with invalid URL", func(t *testing.T) {
		requestBody := map[string]string{
			"url": "",
		}
		jsonData, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRedirect(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	// Khởi tạo dependencies
	urlRepo := repositories.NewURLRepositoryImpl(db)
	urlUsecase := usecases.NewURLUsecase(urlRepo, "http://localhost:8080", getTestConfig())
	urlHandler := handlers.NewURLHandler(urlUsecase)

	// Tạo router
	router := gin.New()
	router.GET("/:shortCode", urlHandler.Redirect)

	// Tạo URL test trước
	urlEntity := &entities.URL{
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		IsActive:    true,
		ClickCount:  0,
	}
	db.Create(urlEntity)

	// Test case 1: Redirect thành công
	t.Run("Redirect successfully", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Location"))
	})

	// Test case 2: Short code không tồn tại
	t.Run("Redirect with non-existent short code", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetURLStats(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	// Khởi tạo dependencies
	urlRepo := repositories.NewURLRepositoryImpl(db)
	urlUsecase := usecases.NewURLUsecase(urlRepo, "http://localhost:8080", getTestConfig())
	urlHandler := handlers.NewURLHandler(urlUsecase)

	// Tạo router
	router := gin.New()
	router.GET("/api/v1/urls/:shortCode/stats", urlHandler.GetURLStats)

	// Tạo URL test trước
	urlEntity := &entities.URL{
		ShortCode:   "stats123",
		OriginalURL: "https://example.com",
		IsActive:    true,
		ClickCount:  5,
	}
	db.Create(urlEntity)

	// Test case 1: Lấy stats thành công
	t.Run("Get URL stats successfully", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/urls/stats123/stats", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.URLStatsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "stats123", response.ShortCode)
		assert.Equal(t, "https://example.com", response.OriginalURL)
		assert.Equal(t, int64(5), response.TotalClicks)
	})

	// Test case 2: Short code không tồn tại
	t.Run("Get stats for non-existent short code", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/urls/nonexistent/stats", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
