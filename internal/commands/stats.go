package commands

import (
	"fmt"
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/utils"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
)

func (m *command) LoadStats(dispatcher dispatcher.Dispatcher) {
	log := m.log.Named("stats")
	defer log.Sugar().Info("Loaded")
	dispatcher.AddHandler(handlers.NewCommand("stats", stats))
}

func stats(ctx *ext.Context, u *ext.Update) error {
	chatId := u.EffectiveChat().GetID()
	peerChatId := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatId.Type != int(storage.TypeUser) {
		return dispatcher.EndGroups
	}
	if len(config.ValueOf.AllowedUsers) != 0 && !utils.Contains(config.ValueOf.AllowedUsers, chatId) {
		ctx.Reply(u, "You are not allowed to use this command.", nil)
		return dispatcher.EndGroups
	}

	// Get analytics summary
	summary := analyticsService.GetAnalyticsSummary()

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
	_, err := ctx.Reply(u, message, &tg.MessagesSendMessageRequest{
		ParseMode: "Markdown",
	})
	return err
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
