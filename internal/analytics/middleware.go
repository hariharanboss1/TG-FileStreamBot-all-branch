package analytics

import (
	"time"

	"github.com/gin-gonic/gin"
)

// AnalyticsMiddleware tracks request performance metrics
func AnalyticsMiddleware(service *AnalyticsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate metrics
		duration := time.Since(start)
		status := c.Writer.Status()

		// Record performance metrics
		service.mu.Lock()
		service.performance = append(service.performance, &PerformanceMetrics{
			Timestamp:       start,
			ResponseTime:    duration.Seconds(),
			ErrorRate:       calculateErrorRate(status),
			ConcurrentUsers: len(service.userStats),
		})
		if len(service.performance) > service.maxMetrics {
			service.performance = service.performance[1:]
		}
		service.mu.Unlock()
	}
}

// calculateErrorRate determines if the status code indicates an error
func calculateErrorRate(status int) float64 {
	if status >= 400 {
		return 1.0
	}
	return 0.0
}

// TrackFileUploadMiddleware tracks file upload metrics
func TrackFileUploadMiddleware(service *AnalyticsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get file information from context
		fileID := c.GetString("file_id")
		fileName := c.GetString("file_name")
		fileSize := c.GetInt64("file_size")
		fileType := c.GetString("file_type")
		userID := c.GetInt64("user_id")

		// Process request
		c.Next()

		// Record file upload metrics
		if fileID != "" && fileName != "" && fileSize > 0 {
			service.RecordFileUpload(userID, fileID, fileName, fileSize, fileType)
		}
	}
}

// TrackFileDownloadMiddleware tracks file download metrics
func TrackFileDownloadMiddleware(service *AnalyticsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get file information from context
		fileID := c.GetString("file_id")
		userID := c.GetInt64("user_id")

		// Process request
		c.Next()

		// Record file download metrics
		if fileID != "" {
			service.RecordFileDownload(userID, fileID)
		}
	}
} 