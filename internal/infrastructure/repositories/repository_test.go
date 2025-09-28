package repositories

import (
	"testing"
	"time"

	"github.com/url-shorted2/internal/domain/entities"

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

// TestURLRepositoryImpl_Create tests Create method
func TestURLRepositoryImpl_Create(t *testing.T) {
	tests := []struct {
		name    string
		url     *entities.URL
		setup   func(*MockURLRepository)
		wantErr bool
	}{
		{
			name: "Tạo URL thành công",
			url: &entities.URL{
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
				IsActive:    true,
				ClickCount:  0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*entities.URL")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Tạo URL thất bại do lỗi database",
			url: &entities.URL{
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
			},
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("Create", mock.AnythingOfType("*entities.URL")).Return(gorm.ErrInvalidData)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			err := mockRepo.Create(tt.url)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestURLRepositoryImpl_GetByShortCode tests GetByShortCode method
func TestURLRepositoryImpl_GetByShortCode(t *testing.T) {
	tests := []struct {
		name      string
		shortCode string
		setup     func(*MockURLRepository)
		want      *entities.URL
		wantErr   bool
	}{
		{
			name:      "Lấy URL thành công",
			shortCode: "abc123",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByShortCode", "abc123").Return(&entities.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					IsActive:    true,
					ClickCount:  0,
				}, nil)
			},
			want: &entities.URL{
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
				IsActive:    true,
				ClickCount:  0,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			got, err := mockRepo.GetByShortCode(tt.shortCode)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ShortCode, got.ShortCode)
				assert.Equal(t, tt.want.OriginalURL, got.OriginalURL)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestURLRepositoryImpl_GetByID tests GetByID method
func TestURLRepositoryImpl_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		setup   func(*MockURLRepository)
		want    *entities.URL
		wantErr bool
	}{
		{
			name: "Lấy URL theo ID thành công",
			id:   1,
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByID", uint(1)).Return(&entities.URL{
					ID:          1,
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					IsActive:    true,
				}, nil)
			},
			want: &entities.URL{
				ID:          1,
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
				IsActive:    true,
			},
			wantErr: false,
		},
		{
			name: "URL không tồn tại",
			id:   999,
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			got, err := mockRepo.GetByID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestURLRepositoryImpl_Update tests Update method
func TestURLRepositoryImpl_Update(t *testing.T) {
	tests := []struct {
		name    string
		url     *entities.URL
		setup   func(*MockURLRepository)
		wantErr bool
	}{
		{
			name: "Cập nhật URL thành công",
			url: &entities.URL{
				ID:          1,
				ShortCode:   "abc123",
				OriginalURL: "https://updated.com",
				IsActive:    true,
			},
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("Update", mock.AnythingOfType("*entities.URL")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Cập nhật URL thất bại",
			url: &entities.URL{
				ID:          1,
				ShortCode:   "abc123",
				OriginalURL: "https://updated.com",
			},
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("Update", mock.AnythingOfType("*entities.URL")).Return(gorm.ErrInvalidData)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			err := mockRepo.Update(tt.url)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestURLRepositoryImpl_Delete tests Delete method
func TestURLRepositoryImpl_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		setup   func(*MockURLRepository)
		wantErr bool
	}{
		{
			name: "Xóa URL thành công",
			id:   1,
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("Delete", uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Xóa URL thất bại",
			id:   999,
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("Delete", uint(999)).Return(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			err := mockRepo.Delete(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestURLRepositoryImpl_IncrementClickCount tests IncrementClickCount method
func TestURLRepositoryImpl_IncrementClickCount(t *testing.T) {
	tests := []struct {
		name      string
		shortCode string
		setup     func(*MockURLRepository)
		wantErr   bool
	}{
		{
			name:      "Tăng click count thành công",
			shortCode: "abc123",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("IncrementClickCount", "abc123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Tăng click count thất bại",
			shortCode: "notfound",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("IncrementClickCount", "notfound").Return(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			err := mockRepo.IncrementClickCount(tt.shortCode)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestURLRepositoryImpl_GetAnalytics tests GetAnalytics method
func TestURLRepositoryImpl_GetAnalytics(t *testing.T) {
	tests := []struct {
		name    string
		urlID   uint
		setup   func(*MockURLRepository)
		want    []entities.Analytics
		wantErr bool
	}{
		{
			name:  "Lấy analytics thành công",
			urlID: 1,
			setup: func(mockRepo *MockURLRepository) {
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
			want: []entities.Analytics{
				{
					URLID:     1,
					IPAddress: "192.168.1.1",
					UserAgent: "Mozilla/5.0",
				},
			},
			wantErr: false,
		},
		{
			name:  "Không có analytics",
			urlID: 1,
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetAnalytics", uint(1)).Return([]entities.Analytics{}, gorm.ErrRecordNotFound)
			},
			want:    []entities.Analytics{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			got, err := mockRepo.GetAnalytics(tt.urlID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestURLRepositoryImpl_AddAnalytics tests AddAnalytics method
func TestURLRepositoryImpl_AddAnalytics(t *testing.T) {
	tests := []struct {
		name      string
		analytics *entities.Analytics
		setup     func(*MockURLRepository)
		wantErr   bool
	}{
		{
			name: "Thêm analytics thành công",
			analytics: &entities.Analytics{
				URLID:     1,
				IPAddress: "192.168.1.1",
				UserAgent: "Mozilla/5.0",
				ClickedAt: time.Now(),
			},
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("AddAnalytics", mock.AnythingOfType("*entities.Analytics")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Thêm analytics thất bại",
			analytics: &entities.Analytics{
				URLID:     1,
				IPAddress: "192.168.1.1",
			},
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("AddAnalytics", mock.AnythingOfType("*entities.Analytics")).Return(gorm.ErrInvalidData)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			err := mockRepo.AddAnalytics(tt.analytics)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestURLRepositoryImpl_GetLastID tests GetLastID method
func TestURLRepositoryImpl_GetLastID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockURLRepository)
		want    uint
		wantErr bool
	}{
		{
			name: "Lấy last ID thành công",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetLastID").Return(uint(5), nil)
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "Không có record nào",
			setup: func(mockRepo *MockURLRepository) {
				mockRepo.On("GetLastID").Return(uint(0), gorm.ErrRecordNotFound)
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			tt.setup(mockRepo)

			got, err := mockRepo.GetLastID()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, uint(0), got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
