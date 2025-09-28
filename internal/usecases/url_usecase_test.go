package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/url-shorted2/internal/config"
	"github.com/url-shorted2/internal/domain/entities"
	"github.com/url-shorted2/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockURLRepository là mock cho URLRepository interface
type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Create(url *entities.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockURLRepository) GetByShortCode(shortCode string) (*entities.URL, error) {
	args := m.Called(shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.URL), args.Error(1)
}

func (m *MockURLRepository) GetByID(id uint) (*entities.URL, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.URL), args.Error(1)
}

func (m *MockURLRepository) Update(url *entities.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockURLRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockURLRepository) IncrementClickCount(shortCode string) error {
	args := m.Called(shortCode)
	return args.Error(0)
}

func (m *MockURLRepository) GetAnalytics(urlID uint) ([]entities.Analytics, error) {
	args := m.Called(urlID)
	return args.Get(0).([]entities.Analytics), args.Error(1)
}

func (m *MockURLRepository) AddAnalytics(analytics *entities.Analytics) error {
	args := m.Called(analytics)
	return args.Error(0)
}

func (m *MockURLRepository) GetLastID() (uint, error) {
	args := m.Called()
	return args.Get(0).(uint), args.Error(1)
}

// MockIDLock là mock cho IDLock interface
type MockIDLock struct {
	mock.Mock
}

func (m *MockIDLock) Lock(ctx context.Context, key string) (*utils.LockData, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.LockData), args.Error(1)
}

func (m *MockIDLock) Unlock(ctx context.Context, ld *utils.LockData) error {
	args := m.Called(ctx, ld)
	return args.Error(0)
}

func TestURLUsecase_CreateShortURL(t *testing.T) {
	tests := []struct {
		name    string
		req     entities.CreateURLRequest
		setup   func(*MockURLRepository, *MockIDLock)
		want    *entities.CreateURLResponse
		wantErr bool
	}{
		{
			name: "Tạo short URL thành công",
			req: entities.CreateURLRequest{
				OriginalURL: "https://example.com",
			},
			setup: func(mockRepo *MockURLRepository, mockLock *MockIDLock) {
				lockData := &utils.LockData{Key: "lock-create-shorted-link", Value: "lock123"}
				mockLock.On("Lock", mock.Anything, "lock-create-shorted-link").Return(lockData, nil)
				mockLock.On("Unlock", mock.Anything, lockData).Return(nil)
				mockRepo.On("GetLastID").Return(uint(0), gorm.ErrRecordNotFound)
				mockRepo.On("Create", mock.AnythingOfType("*entities.URL")).Return(nil)
			},
			want: &entities.CreateURLResponse{
				ShortCode:   "1",
				ShortURL:    "http://localhost:8080/1",
				OriginalURL: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "Tạo short URL với URL không hợp lệ",
			req: entities.CreateURLRequest{
				OriginalURL: "https://",
			},
			setup: func(mockRepo *MockURLRepository, mockLock *MockIDLock) {
				// Không cần setup mock vì sẽ fail ở validateURL trước khi gọi Lock
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Tạo short URL với URL rỗng",
			req: entities.CreateURLRequest{
				OriginalURL: "",
			},
			setup: func(mockRepo *MockURLRepository, mockLock *MockIDLock) {
				// Không cần setup mock vì sẽ fail ở validateURL trước khi gọi Lock
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Tạo short URL thất bại do lỗi database",
			req: entities.CreateURLRequest{
				OriginalURL: "https://example.com",
			},
			setup: func(mockRepo *MockURLRepository, mockLock *MockIDLock) {
				lockData := &utils.LockData{Key: "lock-create-shorted-link", Value: "lock123"}
				mockLock.On("Lock", mock.Anything, "lock-create-shorted-link").Return(lockData, nil)
				mockLock.On("Unlock", mock.Anything, lockData).Return(nil)
				mockRepo.On("GetLastID").Return(uint(0), gorm.ErrRecordNotFound)
				mockRepo.On("Create", mock.AnythingOfType("*entities.URL")).Return(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Tạo short URL thất bại do không lấy được lock",
			req: entities.CreateURLRequest{
				OriginalURL: "https://example.com",
			},
			setup: func(mockRepo *MockURLRepository, mockLock *MockIDLock) {
				mockLock.On("Lock", mock.Anything, "lock-create-shorted-link").Return(nil, errors.New("lock failed"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			mockLock := &MockIDLock{}
			tt.setup(mockRepo, mockLock)

			usecase := &urlUsecase{
				urlRepo: mockRepo,
				baseURL: "http://localhost:8080",
				locker:  mockLock,
				config:  getTestConfig(),
			}

			got, err := usecase.CreateShortURL(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ShortCode, got.ShortCode)
				assert.Equal(t, tt.want.ShortURL, got.ShortURL)
				assert.Equal(t, tt.want.OriginalURL, got.OriginalURL)
			}

			mockRepo.AssertExpectations(t)
			mockLock.AssertExpectations(t)
		})
	}
}

func TestURLUsecase_GetOriginalURL(t *testing.T) {
	tests := []struct {
		name      string
		shortCode string
		setup     func(*MockURLRepository)
		want      string
		wantErr   bool
	}{
		{
			name:      "Lấy original URL thành công",
			shortCode: "abc123",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "abc123").Return(&entities.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					IsActive:    true,
				}, nil)
			},
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name:      "URL không tồn tại",
			shortCode: "notfound",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "notfound").Return(nil, gorm.ErrRecordNotFound)
			},
			want:    "",
			wantErr: true,
		},
		{
			name:      "URL không active",
			shortCode: "inactive",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "inactive").Return(&entities.URL{
					ShortCode:   "inactive",
					OriginalURL: "https://example.com",
					IsActive:    false,
				}, nil)
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			usecase := &urlUsecase{
				urlRepo: mockRepo,
				baseURL: "http://localhost:8080",
			}

			got, err := usecase.GetOriginalURL(tt.shortCode)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestURLUsecase_Redirect(t *testing.T) {
	tests := []struct {
		name      string
		shortCode string
		ipAddress string
		userAgent string
		referer   string
		setup     func(*MockURLRepository)
		want      string
		wantErr   bool
	}{
		{
			name:      "Redirect thành công",
			shortCode: "abc123",
			ipAddress: "192.168.1.1",
			userAgent: "Mozilla/5.0",
			referer:   "https://google.com",
			setup: func(mockRepo *MockURLRepository) {
				// GetOriginalURL call
				mockRepo.On("GetByShortCode", "abc123").Return(&entities.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					IsActive:    true,
					ID:          1,
				}, nil)
				// IncrementClickCount call
				mockRepo.On("IncrementClickCount", "abc123").Return(nil)
				// GetByShortCode call for analytics
				mockRepo.On("GetByShortCode", "abc123").Return(&entities.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					IsActive:    true,
					ID:          1,
				}, nil)
				// AddAnalytics call
				mockRepo.On("AddAnalytics", mock.AnythingOfType("*entities.Analytics")).Return(nil)
			},
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name:      "Redirect thất bại do URL không tồn tại",
			shortCode: "notfound",
			ipAddress: "192.168.1.1",
			userAgent: "Mozilla/5.0",
			referer:   "https://google.com",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "notfound").Return(nil, gorm.ErrRecordNotFound)
			},
			want:    "",
			wantErr: true,
		},
		{
			name:      "Redirect thành công nhưng increment click count thất bại",
			shortCode: "abc123",
			ipAddress: "192.168.1.1",
			userAgent: "Mozilla/5.0",
			referer:   "https://google.com",
			setup: func(mockRepo *MockURLRepository) {
				// GetOriginalURL call
				mockRepo.On("GetByShortCode", "abc123").Return(&entities.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					IsActive:    true,
					ID:          1,
				}, nil)
				// IncrementClickCount call fails
				mockRepo.On("IncrementClickCount", "abc123").Return(errors.New("increment failed"))
				// GetByShortCode call for analytics
				mockRepo.On("GetByShortCode", "abc123").Return(&entities.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					IsActive:    true,
					ID:          1,
				}, nil)
				// AddAnalytics call
				mockRepo.On("AddAnalytics", mock.AnythingOfType("*entities.Analytics")).Return(nil)
			},
			want:    "https://example.com",
			wantErr: false, // Redirect vẫn thành công dù increment thất bại
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			usecase := &urlUsecase{
				urlRepo: mockRepo,
				baseURL: "http://localhost:8080",
			}

			got, err := usecase.Redirect(tt.shortCode, tt.ipAddress, tt.userAgent, tt.referer)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestURLUsecase_GetURLStats(t *testing.T) {
	tests := []struct {
		name      string
		shortCode string
		setup     func(*MockURLRepository)
		want      *entities.URLStatsResponse
		wantErr   bool
	}{
		{
			name:      "Lấy stats thành công",
			shortCode: "abc123",
			setup: func(mockRepo *MockURLRepository) {
				urlEntity := &entities.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					ClickCount:  5,
					CreatedAt:   time.Now(),
					ID:          1,
				}
				mockRepo.On("GetByShortCode", "abc123").Return(urlEntity, nil)

				analytics := []entities.Analytics{
					{
						URLID:     1,
						IPAddress: "192.168.1.1",
						UserAgent: "Mozilla/5.0",
						ClickedAt: time.Now(),
					},
				}
				mockRepo.On("GetAnalytics", uint(1)).Return(analytics, nil)
			},
			want: &entities.URLStatsResponse{
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
				TotalClicks: 5,
			},
			wantErr: false,
		},
		{
			name:      "URL không tồn tại",
			shortCode: "notfound",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "notfound").Return(nil, gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:      "Lấy stats thành công nhưng không có analytics",
			shortCode: "abc123",
			setup: func(mockRepo *MockURLRepository) {
				urlEntity := &entities.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					ClickCount:  0,
					CreatedAt:   time.Now(),
					ID:          1,
				}
				mockRepo.On("GetByShortCode", "abc123").Return(urlEntity, nil)
				mockRepo.On("GetAnalytics", uint(1)).Return([]entities.Analytics{}, gorm.ErrRecordNotFound)
			},
			want: &entities.URLStatsResponse{
				ShortCode:    "abc123",
				OriginalURL:  "https://example.com",
				TotalClicks:  0,
				ClickHistory: []entities.Analytics{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			usecase := &urlUsecase{
				urlRepo: mockRepo,
				baseURL: "http://localhost:8080",
			}

			got, err := usecase.GetURLStats(tt.shortCode)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ShortCode, got.ShortCode)
				assert.Equal(t, tt.want.OriginalURL, got.OriginalURL)
				assert.Equal(t, tt.want.TotalClicks, got.TotalClicks)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestURLUsecase_DeleteURL(t *testing.T) {
	tests := []struct {
		name      string
		shortCode string
		setup     func(*MockURLRepository)
		wantErr   bool
	}{
		{
			name:      "Xóa URL thành công",
			shortCode: "abc123",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "abc123").Return(&entities.URL{
					ShortCode: "abc123",
					ID:        1,
				}, nil)
				mockRepo.On("Delete", uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "URL không tồn tại",
			shortCode: "notfound",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "notfound").Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:      "Xóa URL thất bại",
			shortCode: "abc123",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "abc123").Return(&entities.URL{
					ShortCode: "abc123",
					ID:        1,
				}, nil)
				mockRepo.On("Delete", uint(1)).Return(errors.New("delete failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			usecase := &urlUsecase{
				urlRepo: mockRepo,
				baseURL: "http://localhost:8080",
			}

			err := usecase.DeleteURL(tt.shortCode)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestURLUsecase_validateURL(t *testing.T) {
	usecase := &urlUsecase{}

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "URL hợp lệ với https",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "URL hợp lệ với http",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "URL không có protocol - tự động thêm https",
			url:     "example.com",
			wantErr: false,
		},
		{
			name:    "URL rỗng",
			url:     "",
			wantErr: true,
		},
		// {
		// 	name:    "URL không hợp lệ",
		// 	url:     "http://[::1]:99999",
		// 	wantErr: true,
		// },
		{
			name:    "URL chỉ có protocol",
			url:     "https://",
			wantErr: true,
		},
		{
			name:    "URL có host rỗng",
			url:     "https:///path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.validateURL(tt.url)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestURLUsecase_generateShortCode(t *testing.T) {
	mockRepo := &MockURLRepository{}
	usecase := &urlUsecase{
		urlRepo: mockRepo,
	}

	tests := []struct {
		name    string
		setup   func(*MockURLRepository)
		wantErr bool
	}{
		{
			name: "Tạo short code thành công",
			setup: func(mockRepo *MockURLRepository) {
				// Mock để short code chưa tồn tại
				mockRepo.On("GetByShortCode", mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: false,
		},
		// {
		// 	name: "Tạo short code thất bại sau nhiều lần thử",
		// 	setup: func(mockRepo *MockURLRepository) {
		// 		// Mock để tất cả short code đều đã tồn tại
		// 		mockRepo.On("GetByShortCode", mock.MatchedBy(func(s string) bool {
		// 			return len(s) == 6 // Short code length is 6
		// 		})).Return(&entities.URL{}, nil).Maybe()
		// 	},
		// 	wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockRepo)

			got, err := usecase.generateShortCode()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.Len(t, got, 6) // Short code length should be 6
			}

			mockRepo.AssertExpectations(t)
		})
	}
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
			Type:             "sqlite",
			Path:             ":memory:",
			MaxOpenConns:     10,
			MaxIdleConns:     5,
			ConnMaxLifetime:  time.Hour,
		},
		Redis: config.RedisConfig{
			URL:            "redis://localhost:6379",
			PoolSize:       10,
			MinIdleConns:   5,
			MaxRetries:     3,
			DialTimeout:    5 * time.Second,
			ReadTimeout:    3 * time.Second,
			WriteTimeout:   3 * time.Second,
			PoolTimeout:    4 * time.Second,
		},
		Lock: config.LockConfig{
			MaxTime:    30 * time.Second,
			MaxTryTime: 10 * time.Second,
		},
	}
}
