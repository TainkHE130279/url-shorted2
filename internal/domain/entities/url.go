package entities

import "time"

// URL represents a shortened URL in the database
type URL struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ShortCode   string    `json:"short_code" gorm:"uniqueIndex;not null"`
	OriginalURL string    `json:"original_url" gorm:"not null;size:2048"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	ClickCount  int64     `json:"click_count" gorm:"default:0"`

	// Analytics relationship
	Analytics []Analytics `json:"analytics,omitempty" gorm:"foreignKey:URLID"`
}

// Analytics represents click analytics for a URL
type Analytics struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	URLID     uint      `json:"url_id" gorm:"not null;index"`
	IPAddress string    `json:"ip_address" gorm:"size:45"` // IPv6 compatible
	UserAgent string    `json:"user_agent" gorm:"size:500"`
	Referer   string    `json:"referer" gorm:"size:500"`
	Country   string    `json:"country" gorm:"size:2"`
	City      string    `json:"city" gorm:"size:100"`
	ClickedAt time.Time `json:"clicked_at"`

	// Relationship
	URL URL `json:"url,omitempty" gorm:"foreignKey:URLID"`
}

type CreateURLRequest struct {
	OriginalURL string `json:"url"`
}

// CreateURLResponse represents the response after creating a new URL
type CreateURLResponse struct {
	ShortCode   string    `json:"short_code"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// URLStatsResponse represents analytics data for a URL
type URLStatsResponse struct {
	ShortCode    string      `json:"short_code"`
	OriginalURL  string      `json:"original_url"`
	TotalClicks  int64       `json:"total_clicks"`
	CreatedAt    time.Time   `json:"created_at"`
	LastClicked  *time.Time  `json:"last_clicked,omitempty"`
	ClickHistory []Analytics `json:"click_history,omitempty"`
}
