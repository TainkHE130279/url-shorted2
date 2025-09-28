# URL Shortener API

Một URL Shortener service được xây dựng với Go và Clean Architecture, hỗ trợ tạo short URL, redirect và analytics.

## 🛠️ Công nghệ sử dụng

### **Backend**
- **Go 1.23** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **SQLite** - Database (có thể chuyển sang PostgreSQL/MySQL)
- **Redis** - Distributed locking và caching

### **Infrastructure**
- **Docker** - Containerization với multi-stage build
- **Docker Compose** - Orchestration

### **Architecture**
- **Clean Architecture** - Separation of concerns
- **Repository Pattern** - Data access abstraction
- **Dependency Injection** - Loose coupling

## ⚙️ Cấu hình ứng dụng

### **Environment Variables**

```bash
# Server Configuration
PORT=8080
BASE_URL=http://localhost:8080
GIN_MODE=release

# Database Configuration
DB_TYPE=sqlite
DB_PATH=/app/data/url_shortener.db
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=50
DB_CONN_MAX_LIFETIME=3600

# Redis Configuration
REDIS_URL=redis://localhost:6379
REDIS_POOL_SIZE=100
REDIS_MIN_IDLE_CONNS=20
REDIS_MAX_RETRIES=3
REDIS_DIAL_TIMEOUT=5
REDIS_READ_TIMEOUT=3
REDIS_WRITE_TIMEOUT=3
REDIS_POOL_TIMEOUT=4

# Lock Configuration
LOCK_MAX_TIME=30
LOCK_MAX_TRY_TIME=10
```

### **Cách chạy**

#### **1. Development Mode**
```bash
# Clone repository
git clone <repository-url>
cd url-shorted2

# Install dependencies
go mod tidy

# Copy environment template
cp env.example .env

# Run application
go run cmd/main.go
```

#### **2. Docker Mode**
```bash
# Build base image
docker build -f Dockerfile.golang -t url-shortener:golang .

# Build application image
docker build -f Dockerfile -t url-shortener:latest .

# Run with docker-compose (includes Redis)
docker-compose up -d
```

**Services sẽ chạy:**
- URL Shortener: http://localhost:8080
- Redis: localhost:6379

#### **3. Production Mode**
```bash
# Build optimized binary
CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o main cmd/main.go

# Run with production config
GIN_MODE=release ./main
```

## 📡 API Documentation

### **Base URL**
```
http://localhost:8080
```

### **1. Tạo Short URL**
**POST** `/api/v1/urls`

Tạo một short URL mới từ URL gốc. (Đã thêm distribute lock sử dụng redis để tránh trường hợp tạo short-link tạo ra link id trùng lặp)

**Request Body:**
```json
{
  "url": "https://example.com"
}
```

**Response:**
```json
{
  "short_code": "abc123",
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://example.com",
  "created_at": "2024-01-01T12:00:00Z"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

### **2. Redirect đến Original URL**
**GET** `/{shortCode}`

Redirect người dùng đến URL gốc và tăng click count.

**Response:** HTTP 301 Redirect

**Example:**
```bash
curl -I http://localhost:8080/abc123
# Response: HTTP/1.1 301 Moved Permanently
# Location: https://example.com
```

### **3. Lấy thông tin URL**
**GET** `/api/v1/urls/{shortCode}`

Lấy thông tin chi tiết của một short URL.

**Response:**
```json
{
  "short_code": "abc123",
  "original_url": "https://example.com",
  "is_active": true,
  "click_count": 42,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T14:30:00Z"
}
```

**Example:**
```bash
curl http://localhost:8080/api/v1/urls/abc123
```

### **4. Lấy thống kê URL**
**GET** `/api/v1/urls/{shortCode}/stats`

Lấy thống kê chi tiết và lịch sử click của URL.

**Response:**
```json
{
  "short_code": "abc123",
  "original_url": "https://example.com",
  "total_clicks": 42,
  "created_at": "2024-01-01T12:00:00Z",
  "last_clicked": "2024-01-01T14:30:00Z",
  "click_history": [
    {
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "referer": "https://google.com",
      "country": "VN",
      "city": "Ho Chi Minh",
      "clicked_at": "2024-01-01T14:30:00Z"
    }
  ]
}
```

**Example:**
```bash
curl http://localhost:8080/api/v1/urls/abc123/stats
```

### **5. Xóa URL**
**DELETE** `/api/v1/urls/{shortCode}`

Xóa một short URL (soft delete).

**Response:**
```json
{
  "message": "URL deleted successfully"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/urls/abc123
```

### **6. Health Check**
**GET** `/health`

Kiểm tra trạng thái sức khỏe của service.

**Response:**
```json
{
  "status": "ok",
  "message": "Service is running",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Example:**
```bash
curl http://localhost:8080/health
```

## 🔧 Error Handling

### **Error Response Format**
```json
{
  "error": "Error message",
  "details": "Detailed error description",
  "code": "ERROR_CODE",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### **Common Error Codes**
- `INVALID_URL` - URL không hợp lệ
- `URL_NOT_FOUND` - Short URL không tồn tại
- `URL_INACTIVE` - URL đã bị vô hiệu hóa
- `DATABASE_ERROR` - Lỗi database
- `REDIS_ERROR` - Lỗi Redis connection

## 📊 Monitoring & Observability

### **Health Checks**
- Application health endpoint (`/health`)
- Database connectivity check
- Redis connectivity check

### **Logging**
- Structured logging với JSON format
- Request/response logging
- Error logging với stack trace
- Performance metrics logging

### **Metrics**
- Request count và response time
- Database connection pool status
- Redis connection status
- Click count per URL
- Error rate

## 🚀 Performance Features

- **Connection Pooling** - Database và Redis connection pooling
- **Distributed Locking** - Redis-based locking cho concurrent access
- **Static Binary** - Optimized Go binary với stripped symbols
- **Health Checks** - Container health monitoring
- **Caching** - Redis caching cho frequently accessed data
- **Environment Configuration** - Flexible config management

## 🧪 Testing

```bash
# Run unit tests
go test ./internal/...

# Run integration tests
go test ./tests/...

# Run tests with coverage
go test ./... -cover

# Run tests with race detection
go test ./... -race
```

## 📁 Project Structure

```
url-shorted2/
├── cmd/main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── domain/                 # Domain entities và interfaces
│   ├── usecases/               # Business logic
│   ├── infrastructure/         # External concerns
│   └── utils/                  # Utilities và helpers
├── tests/                      # Integration tests
├── scripts/                    # Build và deployment scripts
├── Dockerfile.golang          # Base Docker image
├── Dockerfile                 # Application Docker image
├── docker-compose.yml         # Development environment
└── env.example               # Environment template
```
