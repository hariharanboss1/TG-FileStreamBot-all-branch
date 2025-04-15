package analytics

import (
	"context"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/disk"
)

type AnalyticsService struct {
	mu sync.RWMutex

	// In-memory storage for quick access
	userStats     map[int64]*UserStats
	fileStats     map[string]*FileStats
	systemMetrics []*SystemMetrics
	performance   []*PerformanceMetrics

	// Configuration
	metricsInterval time.Duration
	maxMetrics      int
}

func NewAnalyticsService(metricsInterval time.Duration, maxMetrics int) *AnalyticsService {
	return &AnalyticsService{
		userStats:       make(map[int64]*UserStats),
		fileStats:       make(map[string]*FileStats),
		systemMetrics:   make([]*SystemMetrics, 0),
		performance:     make([]*PerformanceMetrics, 0),
		metricsInterval: metricsInterval,
		maxMetrics:      maxMetrics,
	}
}

// StartMonitoring begins collecting system metrics
func (s *AnalyticsService) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(s.metricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.collectSystemMetrics()
		}
	}
}

// collectSystemMetrics gathers system performance data
func (s *AnalyticsService) collectSystemMetrics() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get CPU usage
	cpuPercent, _ := cpu.Percent(time.Second, false)
	
	// Get memory usage
	memInfo, _ := mem.VirtualMemory()
	
	// Get disk usage
	diskUsage, _ := disk.Usage("/")

	metrics := &SystemMetrics{
		Timestamp:      time.Now(),
		CPUUsage:       cpuPercent[0],
		MemoryUsage:    memInfo.UsedPercent,
		DiskUsage:      diskUsage.UsedPercent,
		ActiveUsers:    len(s.userStats),
		TotalRequests:  s.getTotalRequests(),
		ErrorRate:      s.calculateErrorRate(),
		AverageLatency: s.calculateAverageLatency(),
	}

	s.systemMetrics = append(s.systemMetrics, metrics)
	if len(s.systemMetrics) > s.maxMetrics {
		s.systemMetrics = s.systemMetrics[1:]
	}
}

// RecordFileUpload tracks file upload statistics
func (s *AnalyticsService) RecordFileUpload(userID int64, fileID, fileName string, fileSize int64, fileType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update user stats
	userStat, exists := s.userStats[userID]
	if !exists {
		userStat = &UserStats{
			UserID:      userID,
			LastActive:  time.Now(),
			UploadCount: 0,
		}
		s.userStats[userID] = userStat
	}
	userStat.TotalFiles++
	userStat.TotalSize += fileSize
	userStat.UploadCount++
	userStat.LastActive = time.Now()

	// Update file stats
	s.fileStats[fileID] = &FileStats{
		FileID:       fileID,
		FileName:     fileName,
		FileSize:     fileSize,
		FileType:     fileType,
		UploadTime:   time.Now(),
		UploaderID:   userID,
		DownloadCount: 0,
	}
}

// RecordFileDownload tracks file download statistics
func (s *AnalyticsService) RecordFileDownload(userID int64, fileID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update user stats
	if userStat, exists := s.userStats[userID]; exists {
		userStat.DownloadCount++
		userStat.LastActive = time.Now()
	}

	// Update file stats
	if fileStat, exists := s.fileStats[fileID]; exists {
		fileStat.DownloadCount++
		fileStat.LastAccessed = time.Now()
	}
}

// GetAnalyticsSummary returns a summary of all analytics data
func (s *AnalyticsService) GetAnalyticsSummary() *AnalyticsSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var uploadsToday, downloadsToday int
	var totalSize int64

	for _, file := range s.fileStats {
		totalSize += file.FileSize
		if file.UploadTime.After(today) {
			uploadsToday++
		}
		if file.LastAccessed.After(today) {
			downloadsToday++
		}
	}

	return &AnalyticsSummary{
		TotalUsers:     len(s.userStats),
		ActiveUsers:    s.getActiveUsersCount(),
		TotalFiles:     len(s.fileStats),
		TotalSize:      totalSize,
		UploadsToday:   uploadsToday,
		DownloadsToday: downloadsToday,
		AverageLatency: s.calculateAverageLatency(),
		ErrorRate:      s.calculateErrorRate(),
		LastUpdated:    now,
	}
}

// Helper methods
func (s *AnalyticsService) getTotalRequests() int {
	total := 0
	for _, user := range s.userStats {
		total += user.UploadCount + user.DownloadCount
	}
	return total
}

func (s *AnalyticsService) calculateErrorRate() float64 {
	// Implement error rate calculation based on your error tracking
	return 0.0
}

func (s *AnalyticsService) calculateAverageLatency() float64 {
	// Implement latency calculation based on your performance metrics
	return 0.0
}

func (s *AnalyticsService) getActiveUsersCount() int {
	activeCount := 0
	activeThreshold := time.Now().Add(-24 * time.Hour)
	
	for _, user := range s.userStats {
		if user.LastActive.After(activeThreshold) {
			activeCount++
		}
	}
	return activeCount
} 