package routes

import (
	"github.com/url-shorted2/internal/infrastructure/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes thiết lập tất cả routes cho ứng dụng
func SetupRoutes(router *gin.Engine, urlHandler *handlers.URLHandler) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// URL routes
		v1.POST("/urls", urlHandler.CreateShortURL)
		v1.GET("/urls/:shortCode", urlHandler.GetURLInfo)
		v1.GET("/urls/:shortCode/stats", urlHandler.GetURLStats)
		v1.DELETE("/urls/:shortCode", urlHandler.DeleteURL)
	}

	// Redirect route (short code without prefix)
	router.GET("/:shortCode", urlHandler.Redirect)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})
}
