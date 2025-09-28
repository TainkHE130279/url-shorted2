package handlers

import (
	"net/http"

	"github.com/url-shorted2/internal/domain/entities"
	"github.com/url-shorted2/internal/usecases"

	"github.com/gin-gonic/gin"
)

// URLHandler xử lý các request liên quan đến URL
type URLHandler struct {
	urlUsecase usecases.IURLUsecase
}

// NewURLHandler tạo instance mới của URLHandler
func NewURLHandler(urlUsecase usecases.IURLUsecase) *URLHandler {
	return &URLHandler{
		urlUsecase: urlUsecase,
	}
}

// CreateShortURL xử lý POST /api/v1/urls
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var request entities.CreateURLRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	// Create short URL
	response, err := h.urlUsecase.CreateShortURL(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create short URL",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Redirect xử lý GET /:shortCode
func (h *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Short code is required",
		})
		return
	}

	// Get client info
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referer := c.GetHeader("Referer")

	// Redirect
	originalURL, err := h.urlUsecase.Redirect(shortCode, ipAddress, userAgent, referer)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "URL not found or expired",
			"details": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// GetURLStats xử lý GET /api/v1/urls/:shortCode/stats
func (h *URLHandler) GetURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Short code is required",
		})
		return
	}

	stats, err := h.urlUsecase.GetURLStats(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "URL not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DeleteURL xử lý DELETE /api/v1/urls/:shortCode
func (h *URLHandler) DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Short code is required",
		})
		return
	}

	err := h.urlUsecase.DeleteURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "URL not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "URL deleted successfully",
	})
}

// GetURLInfo xử lý GET /api/v1/urls/:shortCode
func (h *URLHandler) GetURLInfo(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Short code is required",
		})
		return
	}

	originalURL, err := h.urlUsecase.GetOriginalURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "URL not found or expired",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code":   shortCode,
		"original_url": originalURL,
	})
}
