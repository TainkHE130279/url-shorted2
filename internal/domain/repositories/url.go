package repositories

import "github.com/url-shorted2/internal/domain/entities"

// URLRepository định nghĩa interface cho URL repository
type IURLRepository interface {
	Create(url *entities.URL) error
	GetByShortCode(shortCode string) (*entities.URL, error)
	GetByID(id uint) (*entities.URL, error)
	Update(url *entities.URL) error
	Delete(id uint) error
	IncrementClickCount(shortCode string) error
	GetAnalytics(urlID uint) ([]entities.Analytics, error)
	AddAnalytics(analytics *entities.Analytics) error
	GetLastID() (uint, error)
}
