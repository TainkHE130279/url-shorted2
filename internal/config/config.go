package config

import (
	"os"
	"strconv"
	"time"
)

// Config chứa tất cả cấu hình của ứng dụng
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Lock     LockConfig
}

// ServerConfig cấu hình server
type ServerConfig struct {
	Port    string
	BaseURL string
	GinMode string
}

// DatabaseConfig cấu hình database
type DatabaseConfig struct {
	Type            string
	Path            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig cấu hình Redis
type RedisConfig struct {
	URL          string
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
}

// LockConfig cấu hình lock
type LockConfig struct {
	MaxTime    time.Duration
	MaxTryTime time.Duration
}

// LoadConfig load cấu hình từ environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Type:            getEnv("DB_TYPE", "sqlite"),
			Path:            getEnv("DB_PATH", "url_shortener.db"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 50),
			ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 3600)) * time.Second,
		},
		Redis: RedisConfig{
			URL:          getEnv("REDIS_URL", "redis://localhost:6379"),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 100),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 20),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
			DialTimeout:  time.Duration(getEnvAsInt("REDIS_DIAL_TIMEOUT", 5)) * time.Second,
			ReadTimeout:  time.Duration(getEnvAsInt("REDIS_READ_TIMEOUT", 3)) * time.Second,
			WriteTimeout: time.Duration(getEnvAsInt("REDIS_WRITE_TIMEOUT", 3)) * time.Second,
			PoolTimeout:  time.Duration(getEnvAsInt("REDIS_POOL_TIMEOUT", 4)) * time.Second,
		},
		Lock: LockConfig{
			MaxTime:    time.Duration(getEnvAsInt("LOCK_MAX_TIME", 30)) * time.Second,
			MaxTryTime: time.Duration(getEnvAsInt("LOCK_MAX_TRY_TIME", 10)) * time.Second,
		},
	}
}

// getEnv lấy environment variable với default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt lấy environment variable dưới dạng int với default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
