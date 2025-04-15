package analytics

import (
	"time"
)

// Analytics models for monitoring and tracking

type UserStats struct {
	UserID        int64     `json:"user_id"`
	TotalFiles    int       `json:"total_files"`
	TotalSize     int64     `json:"total_size"`
	LastActive    time.Time `json:"last_active"`
	UploadCount   int       `json:"upload_count"`
	DownloadCount int       `json:"download_count"`
}

type FileStats struct {
	FileID        string    `json:"file_id"`
	FileName      string    `json:"file_name"`
	FileSize      int64     `json:"file_size"`
	FileType      string    `json:"file_type"`
	UploadTime    time.Time `json:"upload_time"`
	DownloadCount int       `json:"download_count"`
	LastAccessed  time.Time `json:"last_accessed"`
	UploaderID    int64     `json:"uploader_id"`
}

type SystemMetrics struct {
	Timestamp      time.Time `json:"timestamp"`
	CPUUsage       float64   `json:"cpu_usage"`
	MemoryUsage    float64   `json:"memory_usage"`
	DiskUsage      float64   `json:"disk_usage"`
	ActiveUsers    int       `json:"active_users"`
	TotalRequests  int       `json:"total_requests"`
	ErrorRate      float64   `json:"error_rate"`
	AverageLatency float64   `json:"average_latency"`
}

type PerformanceMetrics struct {
	UploadSpeed     float64   `json:"upload_speed"`
	DownloadSpeed   float64   `json:"download_speed"`
	ResponseTime    float64   `json:"response_time"`
	ErrorRate       float64   `json:"error_rate"`
	ConcurrentUsers int       `json:"concurrent_users"`
	Timestamp       time.Time `json:"timestamp"`
}

type AnalyticsSummary struct {
	TotalUsers      int       `json:"total_users"`
	ActiveUsers     int       `json:"active_users"`
	TotalFiles      int       `json:"total_files"`
	TotalSize       int64     `json:"total_size"`
	UploadsToday    int       `json:"uploads_today"`
	DownloadsToday  int       `json:"downloads_today"`
	AverageLatency  float64   `json:"average_latency"`
	ErrorRate       float64   `json:"error_rate"`
	LastUpdated     time.Time `json:"last_updated"`
} 