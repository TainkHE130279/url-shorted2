package repositories

import (
	"github.com/url-shorted2/internal/domain/entities"
	"github.com/url-shorted2/internal/domain/repositories"

	"gorm.io/gorm"
)

// urlRepositoryImpl implement URLRepository interface
type urlRepositoryImpl struct {
	db *gorm.DB
}

// NewURLRepositoryImpl tạo instance mới của URLRepository
func NewURLRepositoryImpl(db *gorm.DB) repositories.IURLRepository {
	return &urlRepositoryImpl{
		db: db,
	}
}

// Create tạo URL mới
func (r *urlRepositoryImpl) Create(url *entities.URL) error {
	return r.db.Create(url).Error
}

// GetByShortCode lấy URL theo short code
func (r *urlRepositoryImpl) GetByShortCode(shortCode string) (*entities.URL, error) {
	var url entities.URL
	err := r.db.Where("short_code = ?", shortCode).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// GetByID lấy URL theo ID
func (r *urlRepositoryImpl) GetByID(id uint) (*entities.URL, error) {
	var url entities.URL
	err := r.db.First(&url, id).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// Update cập nhật URL
func (r *urlRepositoryImpl) Update(url *entities.URL) error {
	return r.db.Save(url).Error
}

// Delete xóa URL
func (r *urlRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.URL{}, id).Error
}

// IncrementClickCount tăng số lần click
func (r *urlRepositoryImpl) IncrementClickCount(shortCode string) error {
	return r.db.Model(&entities.URL{}).
		Where("short_code = ?", shortCode).
		Update("click_count", gorm.Expr("click_count + 1")).Error
}

// GetAnalytics lấy analytics cho URL
func (r *urlRepositoryImpl) GetAnalytics(urlID uint) ([]entities.Analytics, error) {
	var analytics []entities.Analytics
	err := r.db.Where("url_id = ?", urlID).
		Order("clicked_at DESC").
		Find(&analytics).Error
	return analytics, err
}

// AddAnalytics thêm analytics record
func (r *urlRepositoryImpl) AddAnalytics(analytics *entities.Analytics) error {
	return r.db.Create(analytics).Error
}

// GetLastID lấy ID cuối cùng (cao nhất) trong table
func (r *urlRepositoryImpl) GetLastID() (uint, error) {
	var url entities.URL
	err := r.db.Order("id DESC").First(&url).Error
	if err != nil {
		return 0, err
	}
	return url.ID, nil
}
