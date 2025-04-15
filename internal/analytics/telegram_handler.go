package analytics

import (
	"fmt"
	"time"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/tg"
)

// RegisterTelegramCommands registers the stats command with the bot
func RegisterTelegramCommands(dp dispatcher.Dispatcher, service *AnalyticsService) {
	dp.AddHandler(handlers.NewCommand("stats", func(ctx *ext.Context, update *types.UpdateNewMessage) error {
		// Get analytics summary
		summary := service.GetAnalyticsSummary()

		// Format the message
		message := fmt.Sprintf("📊 *Bot Statistics*\n\n"+
			"👥 *Users*\n"+
			"Total Users: %d\n"+
			"Active Users (24h): %d\n\n"+
			"📁 *Files*\n"+
			"Total Files: %d\n"+
			"Total Size: %s\n"+
			"Uploads Today: %d\n"+
			"Downloads Today: %d\n\n"+
			"⚡ *Performance*\n"+
			"Average Latency: %.2fms\n"+
			"Error Rate: %.2f%%\n\n"+
			"Last Updated: %s",
			summary.TotalUsers,
			summary.ActiveUsers,
			summary.TotalFiles,
			formatFileSize(summary.TotalSize),
			summary.UploadsToday,
			summary.DownloadsToday,
			summary.AverageLatency,
			summary.ErrorRate*100,
			summary.LastUpdated.Format("2006-01-02 15:04:05"))

		// Send the message
		_, err := ctx.Reply(update, message, &tg.MessagesSendMessageRequest{
			ParseMode: "Markdown",
		})
		return err
	}))
}

// formatFileSize converts bytes to human readable format
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
} 