package analytics

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service *AnalyticsService
}

func NewAnalyticsHandler(service *AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

// RegisterRoutes sets up the analytics endpoints
func (h *AnalyticsHandler) RegisterRoutes(router *gin.Engine) {
	analytics := router.Group("/analytics")
	{
		analytics.GET("/summary", h.getAnalyticsSummary)
		analytics.GET("/users", h.getUserStats)
		analytics.GET("/files", h.getFileStats)
		analytics.GET("/system", h.getSystemMetrics)
		analytics.GET("/performance", h.getPerformanceMetrics)
	}
}

// getAnalyticsSummary returns a summary of all analytics data
func (h *AnalyticsHandler) getAnalyticsSummary(c *gin.Context) {
	summary := h.service.GetAnalyticsSummary()
	c.JSON(http.StatusOK, summary)
}

// getUserStats returns statistics for all users
func (h *AnalyticsHandler) getUserStats(c *gin.Context) {
	h.service.mu.RLock()
	defer h.service.mu.RUnlock()

	// Convert map to slice for JSON serialization
	userStats := make([]*UserStats, 0, len(h.service.userStats))
	for _, stat := range h.service.userStats {
		userStats = append(userStats, stat)
	}

	c.JSON(http.StatusOK, userStats)
}

// getFileStats returns statistics for all files
func (h *AnalyticsHandler) getFileStats(c *gin.Context) {
	h.service.mu.RLock()
	defer h.service.mu.RUnlock()

	// Convert map to slice for JSON serialization
	fileStats := make([]*FileStats, 0, len(h.service.fileStats))
	for _, stat := range h.service.fileStats {
		fileStats = append(fileStats, stat)
	}

	c.JSON(http.StatusOK, fileStats)
}

// getSystemMetrics returns recent system metrics
func (h *AnalyticsHandler) getSystemMetrics(c *gin.Context) {
	h.service.mu.RLock()
	defer h.service.mu.RUnlock()

	// Get time range from query parameters
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var start, end time.Time
	var err error

	if startTime != "" {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	}

	if endTime != "" {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	}

	// Filter metrics by time range
	var filteredMetrics []*SystemMetrics
	for _, metric := range h.service.systemMetrics {
		if (start.IsZero() || metric.Timestamp.After(start)) &&
			(end.IsZero() || metric.Timestamp.Before(end)) {
			filteredMetrics = append(filteredMetrics, metric)
		}
	}

	c.JSON(http.StatusOK, filteredMetrics)
}

// getPerformanceMetrics returns recent performance metrics
func (h *AnalyticsHandler) getPerformanceMetrics(c *gin.Context) {
	h.service.mu.RLock()
	defer h.service.mu.RUnlock()

	// Get time range from query parameters
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var start, end time.Time
	var err error

	if startTime != "" {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	}

	if endTime != "" {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	}

	// Filter metrics by time range
	var filteredMetrics []*PerformanceMetrics
	for _, metric := range h.service.performance {
		if (start.IsZero() || metric.Timestamp.After(start)) &&
			(end.IsZero() || metric.Timestamp.Before(end)) {
			filteredMetrics = append(filteredMetrics, metric)
		}
	}

	c.JSON(http.StatusOK, filteredMetrics)
}

// ExportMetrics exports all metrics to a JSON file
func (h *AnalyticsHandler) ExportMetrics(c *gin.Context) {
	h.service.mu.RLock()
	defer h.service.mu.RUnlock()

	export := struct {
		UserStats     []*UserStats     `json:"user_stats"`
		FileStats     []*FileStats     `json:"file_stats"`
		SystemMetrics []*SystemMetrics `json:"system_metrics"`
		Performance   []*PerformanceMetrics `json:"performance_metrics"`
	}{
		UserStats:     make([]*UserStats, 0, len(h.service.userStats)),
		FileStats:     make([]*FileStats, 0, len(h.service.fileStats)),
		SystemMetrics: h.service.systemMetrics,
		Performance:   h.service.performance,
	}

	// Convert maps to slices
	for _, stat := range h.service.userStats {
		export.UserStats = append(export.UserStats, stat)
	}
	for _, stat := range h.service.fileStats {
		export.FileStats = append(export.FileStats, stat)
	}

	// Set headers for file download
	c.Header("Content-Disposition", "attachment; filename=analytics_export.json")
	c.Header("Content-Type", "application/json")

	// Stream the JSON response
	encoder := json.NewEncoder(c.Writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(export); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export metrics"})
		return
	}
} 