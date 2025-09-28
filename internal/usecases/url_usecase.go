package usecases

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/url-shorted2/internal/config"
	"github.com/url-shorted2/internal/domain/entities"
	"github.com/url-shorted2/internal/domain/repositories"
	"github.com/url-shorted2/internal/utils"

	"gorm.io/gorm"
)

type IURLUsecase interface {
	CreateShortURL(req entities.CreateURLRequest) (*entities.CreateURLResponse, error)
	GetOriginalURL(shortCode string) (string, error)
	Redirect(shortCode string, ipAddress, userAgent, referer string) (string, error)
	GetURLStats(shortCode string) (*entities.URLStatsResponse, error)
	DeleteURL(shortCode string) error
}

type urlUsecase struct {
	urlRepo repositories.IURLRepository
	baseURL string
	locker  utils.IDLock
	config  *config.Config
}

func NewURLUsecase(urlRepo repositories.IURLRepository, baseURL string, cfg *config.Config) IURLUsecase {
	var locker utils.IDLock

	// Sử dụng mock lock trong test environment
	if cfg.Server.GinMode == "test" {
		locker = utils.NewMockLock()
	} else {
		locker = utils.NewRedisLock(utils.GetRedisWithConfig(
			cfg.Redis.URL,
			cfg.Redis.PoolSize,
			cfg.Redis.MinIdleConns,
			cfg.Redis.MaxRetries,
			cfg.Redis.DialTimeout,
			cfg.Redis.ReadTimeout,
			cfg.Redis.WriteTimeout,
			cfg.Redis.PoolTimeout,
		), &utils.Config{
			MaxLockTime: cfg.Lock.MaxTime,
			MaxTryTime:  cfg.Lock.MaxTryTime,
		})
	}

	return &urlUsecase{
		urlRepo: urlRepo,
		baseURL: baseURL,
		locker:  locker,
		config:  cfg,
	}
}

// CreateShortURL tạo short URL
func (u *urlUsecase) CreateShortURL(req entities.CreateURLRequest) (*entities.CreateURLResponse, error) {
	// Validate URL
	if err := u.validateURL(req.OriginalURL); err != nil {
		return nil, err
	}
	var key = "lock-create-shorted-link"
	//lock key
	lrs, err := u.locker.Lock(context.TODO(), key)
	if nil != err {
		return nil, err
	}
	// unlock key
	defer func() {
		if err := u.locker.Unlock(context.TODO(), lrs); err != nil {
			fmt.Printf("Failed to unlock: %v\n", err)
		}
	}()
	//get lastID for create short link
	id, err := u.urlRepo.GetLastID()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil
	}

	//defer unlock key

	next := id + 1
	//Create URL entity
	urlEntity := &entities.URL{
		ShortCode:   fmt.Sprintf("%v", next),
		OriginalURL: req.OriginalURL,
		IsActive:    true,
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// // Save to database
	if err := u.urlRepo.Create(urlEntity); err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	// // Return response
	response := &entities.CreateURLResponse{
		ShortCode:   fmt.Sprintf("%v", next),
		ShortURL:    fmt.Sprintf("%s/%v", u.baseURL, next),
		OriginalURL: req.OriginalURL,
		CreatedAt:   urlEntity.CreatedAt,
	}

	return response, nil
}

// GetOriginalURL lấy original URL từ short code
func (u *urlUsecase) GetOriginalURL(shortCode string) (string, error) {
	urlEntity, err := u.urlRepo.GetByShortCode(shortCode)
	if err != nil {
		return "", fmt.Errorf("URL not found: %w", err)
	}

	// Check if URL is active
	if !urlEntity.IsActive {
		return "", errors.New("URL is inactive")
	}

	return urlEntity.OriginalURL, nil
}

// Redirect thực hiện redirect và ghi analytics
func (u *urlUsecase) Redirect(shortCode string, ipAddress, userAgent, referer string) (string, error) {
	// Get original URL
	originalURL, err := u.GetOriginalURL(shortCode)
	if err != nil {
		return "", err
	}

	// Increment click count
	if err := u.urlRepo.IncrementClickCount(shortCode); err != nil {
		// Log error but don't fail the redirect
		fmt.Printf("Failed to increment click count: %v\n", err)
	}

	// Add analytics
	analytics := &entities.Analytics{
		URLID:     0, // Will be set by repository
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Referer:   referer,
		ClickedAt: time.Now(),
	}

	// Get URL ID for analytics
	urlEntity, err := u.urlRepo.GetByShortCode(shortCode)
	if err == nil {
		analytics.URLID = urlEntity.ID
		if err := u.urlRepo.AddAnalytics(analytics); err != nil {
			// Log error but don't fail the redirect
			fmt.Printf("Failed to add analytics: %v\n", err)
		}
	}

	return originalURL, nil
}

// GetURLStats lấy thống kê URL
func (u *urlUsecase) GetURLStats(shortCode string) (*entities.URLStatsResponse, error) {
	urlEntity, err := u.urlRepo.GetByShortCode(shortCode)
	if err != nil {
		return nil, fmt.Errorf("URL not found: %w", err)
	}

	// Get analytics
	analytics, err := u.urlRepo.GetAnalytics(urlEntity.ID)
	if err != nil {
		analytics = []entities.Analytics{} // Return empty if no analytics
	}

	// Find last clicked time
	var lastClicked *time.Time
	if len(analytics) > 0 {
		lastClicked = &analytics[len(analytics)-1].ClickedAt
	}

	response := &entities.URLStatsResponse{
		ShortCode:    urlEntity.ShortCode,
		OriginalURL:  urlEntity.OriginalURL,
		TotalClicks:  urlEntity.ClickCount,
		CreatedAt:    urlEntity.CreatedAt,
		LastClicked:  lastClicked,
		ClickHistory: analytics,
	}

	return response, nil
}

// DeleteURL xóa URL
func (u *urlUsecase) DeleteURL(shortCode string) error {
	urlEntity, err := u.urlRepo.GetByShortCode(shortCode)
	if err != nil {
		return fmt.Errorf("URL not found: %w", err)
	}

	return u.urlRepo.Delete(urlEntity.ID)
}

// validateURL kiểm tra URL có hợp lệ không
func (u *urlUsecase) validateURL(rawURL string) error {
	if rawURL == "" {
		return errors.New("URL cannot be empty")
	}

	// Add protocol if missing
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Host == "" {
		return errors.New("URL must have a valid host")
	}

	return nil
}

// generateShortCode tạo short code ngẫu nhiên
func (u *urlUsecase) generateShortCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	for attempts := 0; attempts < 10; attempts++ {
		code := make([]byte, codeLength)
		for i := range code {
			// Use crypto/rand for secure random number generation
			n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", fmt.Errorf("failed to generate random number: %w", err)
			}
			code[i] = charset[n.Int64()]
		}
		shortCode := string(code)

		// Check if code already exists
		_, err := u.urlRepo.GetByShortCode(shortCode)
		if err != nil {
			return shortCode, nil
		}
	}

	return "", errors.New("failed to generate unique short code")
}
