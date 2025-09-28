package utils

import (
	"testing"
	"time"

	"github.com/url-shorted2/internal/domain/entities"

	"github.com/stretchr/testify/assert"
)

// CreateTestURL tạo URL entity cho testing
func CreateTestURL(id uint, shortCode, originalURL string) *entities.URL {
	return &entities.URL{
		ID:          id,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		IsActive:    true,
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestAnalytics tạo Analytics entity cho testing
func CreateTestAnalytics(id, urlID uint, ipAddress, userAgent, referer string) *entities.Analytics {
	return &entities.Analytics{
		ID:        id,
		URLID:     urlID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Referer:   referer,
		ClickedAt: time.Now(),
	}
}

// CreateTestCreateURLRequest tạo CreateURLRequest cho testing
func CreateTestCreateURLRequest(originalURL string) entities.CreateURLRequest {
	return entities.CreateURLRequest{
		OriginalURL: originalURL,
	}
}

// AssertURLEqual kiểm tra 2 URL entity có giống nhau không
func AssertURLEqual(t *testing.T, expected, actual *entities.URL) {
	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		t.Errorf("Expected: %v, Got: %v", expected, actual)
		return
	}

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.ShortCode, actual.ShortCode)
	assert.Equal(t, expected.OriginalURL, actual.OriginalURL)
	assert.Equal(t, expected.IsActive, actual.IsActive)
	assert.Equal(t, expected.ClickCount, actual.ClickCount)
}

// AssertAnalyticsEqual kiểm tra 2 Analytics entity có giống nhau không
func AssertAnalyticsEqual(t *testing.T, expected, actual *entities.Analytics) {
	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		t.Errorf("Expected: %v, Got: %v", expected, actual)
		return
	}

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.URLID, actual.URLID)
	assert.Equal(t, expected.IPAddress, actual.IPAddress)
	assert.Equal(t, expected.UserAgent, actual.UserAgent)
	assert.Equal(t, expected.Referer, actual.Referer)
}

// AssertCreateURLResponseEqual kiểm tra 2 CreateURLResponse có giống nhau không
func AssertCreateURLResponseEqual(t *testing.T, expected, actual *entities.CreateURLResponse) {
	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		t.Errorf("Expected: %v, Got: %v", expected, actual)
		return
	}

	assert.Equal(t, expected.ShortCode, actual.ShortCode)
	assert.Equal(t, expected.ShortURL, actual.ShortURL)
	assert.Equal(t, expected.OriginalURL, actual.OriginalURL)
}

// AssertURLStatsResponseEqual kiểm tra 2 URLStatsResponse có giống nhau không
func AssertURLStatsResponseEqual(t *testing.T, expected, actual *entities.URLStatsResponse) {
	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		t.Errorf("Expected: %v, Got: %v", expected, actual)
		return
	}

	assert.Equal(t, expected.ShortCode, actual.ShortCode)
	assert.Equal(t, expected.OriginalURL, actual.OriginalURL)
	assert.Equal(t, expected.TotalClicks, actual.TotalClicks)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)

	if expected.LastClicked != nil && actual.LastClicked != nil {
		assert.WithinDuration(t, *expected.LastClicked, *actual.LastClicked, time.Second)
	} else {
		assert.Equal(t, expected.LastClicked, actual.LastClicked)
	}
}

// ValidTestURLs danh sách các URL hợp lệ cho testing
var ValidTestURLs = []string{
	"https://example.com",
	"http://example.com",
	"https://www.example.com",
	"https://example.com/path",
	"https://example.com/path?query=value",
	"https://example.com/path#fragment",
	"example.com", // Sẽ được tự động thêm https://
	"www.example.com",
	"subdomain.example.com",
}

// InvalidTestURLs danh sách các URL không hợp lệ cho testing
var InvalidTestURLs = []string{
	"",
	"not-a-url",
	"https://",
	"http://",
	"ftp://example.com", // Không hỗ trợ FTP
	"://example.com",
	"https:///path",
	"https://example.com:99999", // Port không hợp lệ
}

// TestShortCodes danh sách short codes cho testing
var TestShortCodes = []string{
	"abc123",
	"def456",
	"ghi789",
	"jkl012",
	"mno345",
	"pqr678",
	"stu901",
	"vwx234",
	"yz5678",
	"123abc",
}

// TestIPAddresses danh sách IP addresses cho testing
var TestIPAddresses = []string{
	"192.168.1.1",
	"10.0.0.1",
	"172.16.0.1",
	"127.0.0.1",
	"::1",         // IPv6 localhost
	"2001:db8::1", // IPv6 example
}

// TestUserAgents danh sách User-Agent strings cho testing
var TestUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
	"curl/7.68.0",
	"PostmanRuntime/7.26.8",
	"Go-http-client/1.1",
}

// TestReferers danh sách Referer strings cho testing
var TestReferers = []string{
	"https://google.com",
	"https://bing.com",
	"https://yahoo.com",
	"https://example.com",
	"https://facebook.com",
	"https://twitter.com",
	"", // Direct access
}
