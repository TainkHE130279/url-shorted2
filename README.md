# Go Gin Clean Architecture - URL Shortener API

Đây là một bộ khung Go Gin đơn giản với API URL shortener được thiết kế theo cấu trúc Clean Architecture.

## Cấu trúc thư mục

```
url-shorted2/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point của ứng dụng
├── internal/
│   ├── domain/                     # Domain layer (entities, repositories)
│   │   ├── entities/
│   │   │   └── url.go             # URL entities
│   │   └── repositories/
│   │       └── url_repository.go  # URL repository interface
│   ├── usecases/                   # Use case layer (business logic)
│   │   └── url_usecase.go         # URL use case
│   └── infrastructure/             # Infrastructure layer
│       ├── handlers/
│       │   └── url_handler.go     # HTTP handlers
│       ├── repositories/
│       │   └── url_repository_impl.go # Repository implementation
│       ├── routes/
│       │   └── routes.go          # Route definitions
│       └── middleware/
│           ├── cors.go            # CORS middleware
│           └── logger.go          # Logger middleware
├── tests/
│   └── url_test.go                # Unit tests
├── go.mod
└── README.md
```

## Clean Architecture Layers

### 1. Domain Layer (`internal/domain/`)
- **Entities**: Chứa các business objects và data structures
- **Repositories**: Định nghĩa interfaces cho data access

### 2. Use Case Layer (`internal/usecases/`)
- Chứa business logic và use cases
- Không phụ thuộc vào framework hay database
- Chỉ phụ thuộc vào domain layer

### 3. Infrastructure Layer (`internal/infrastructure/`)
- **Handlers**: HTTP request/response handling
- **Repositories**: Implementation của domain repositories
- **Routes**: Route definitions
- **Middleware**: Cross-cutting concerns

## API Endpoints

### POST /api/v1/urls
Tạo short URL mới.

**Request:**
```json
{
  "url": "https://example.com",
  "custom_code": "example",
  "expires_in": 3600
}
```

**Response:**
```json
{
  "short_code": "example",
  "short_url": "http://localhost:8080/example",
  "original_url": "https://example.com",
  "expires_at": "2024-01-01T13:00:00Z",
  "created_at": "2024-01-01T12:00:00Z"
}
```

### GET /:shortCode
Redirect đến original URL.

**Response:** HTTP 301 Redirect

### GET /api/v1/urls/:shortCode
Lấy thông tin URL.

**Response:**
```json
{
  "short_code": "example",
  "original_url": "https://example.com"
}
```

### GET /api/v1/urls/:shortCode/stats
Lấy thống kê URL.

**Response:**
```json
{
  "short_code": "example",
  "original_url": "https://example.com",
  "total_clicks": 42,
  "created_at": "2024-01-01T12:00:00Z",
  "last_clicked": "2024-01-01T14:30:00Z",
  "click_history": [...]
}
```

### DELETE /api/v1/urls/:shortCode
Xóa URL.

**Response:**
```json
{
  "message": "URL deleted successfully"
}
```

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "ok",
  "message": "Service is running"
}
```

## Cách chạy

### 1. Cài đặt dependencies
```bash
go mod tidy
```

### 2. Cấu hình Database (PostgreSQL)
Ứng dụng sử dụng PostgreSQL làm database chính. Cấu hình thông qua environment variables:

```bash
# Database Configuration
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=url_shortener
export DB_SSLMODE=disable

# Server Configuration
export PORT=8080
export BASE_URL=http://localhost:8080
export GIN_MODE=debug
```

Hoặc tạo file `.env` với nội dung tương tự.

### 3. Khởi động PostgreSQL
```bash
# Sử dụng Docker
docker run --name postgres-url-shortener \
  -e POSTGRES_DB=url_shortener \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres:15

# Hoặc cài đặt PostgreSQL locally
```

### 4. Chạy server
```bash
go run cmd/main.go
```

Server sẽ chạy trên port 8080 và tự động migrate database schema.

### 5. Test API
```bash
# Tạo short URL
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Tạo short URL với custom code
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com", "custom_code": "google"}'

# Redirect (thay abc123 bằng short code thực tế)
curl -I http://localhost:8080/abc123

# Lấy thông tin URL
curl http://localhost:8080/api/v1/urls/abc123

# Lấy thống kê URL
curl http://localhost:8080/api/v1/urls/abc123/stats

# Xóa URL
curl -X DELETE http://localhost:8080/api/v1/urls/abc123

# Test health check
curl http://localhost:8080/health
```

### 6. Chạy tests
```bash
go test ./tests/...
```

## Lợi ích của Clean Architecture

1. **Separation of Concerns**: Mỗi layer có trách nhiệm riêng biệt
2. **Dependency Inversion**: High-level modules không phụ thuộc vào low-level modules
3. **Testability**: Dễ dàng viết unit tests cho từng layer
4. **Maintainability**: Code dễ maintain và extend
5. **Flexibility**: Có thể thay đổi implementation mà không ảnh hưởng business logic

## Mở rộng

Để mở rộng ứng dụng, bạn có thể:

1. Thêm entities mới trong `internal/domain/entities/`
2. Thêm repository interfaces trong `internal/domain/repositories/`
3. Thêm use cases mới trong `internal/usecases/`
4. Implement repositories trong `internal/infrastructure/repositories/`
5. Thêm handlers mới trong `internal/infrastructure/handlers/`
6. Thêm routes mới trong `internal/infrastructure/routes/`

## Dependencies

- **Gin**: HTTP web framework
- **GORM**: ORM library
- **PostgreSQL**: Database chính với driver postgres
- **Testify**: Testing toolkit
- **Go 1.21+**: Programming language

## Tính năng

- ✅ Tạo short URL với custom code
- ✅ Redirect đến original URL
- ✅ Analytics và click tracking
- ✅ URL expiration
- ✅ RESTful API
- ✅ Database persistence
- ✅ Clean Architecture
- ✅ Unit tests
- ✅ Error handling
- ✅ CORS support
