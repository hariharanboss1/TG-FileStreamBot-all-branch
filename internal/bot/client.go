package bot

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/commands"
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
)

var Bot *gotgproto.Client

// Assuming the function is defined outside any existing functions
func getUserSubscription(userID int64) (string, error) {
	chatMember, err := bot.GetChatMember(tgbotapi.ChatConfig{
		ChatID:    -1001882519219,
		UserID:    userID,
	})
	if err != nil {
		return "", err // Handle error
	}
	return chatMember.Status, nil
}

func StartClient(log *zap.Logger) (*gotgproto.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	resultChan := make(chan struct {
		client *gotgproto.Client
		err   error
	})
	go func(ctx context.Context) {
		client, err := gotgproto.NewClient(
			int(config.ValueOf.ApiID),
			config.ValueOf.ApiHash,
			gotgproto.ClientTypeBot(config.ValueOf.BotToken),
			&gotgproto.ClientOpts{
				Session: sessionMaker.SqlSession(
					sqlite.Open("fsb.session"),
				),
				DisableCopyright: true,
			},
		)
		resultChan <- struct {
			client *gotgproto.Client
			err   error
		}{client, err}
	}(ctx)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultChan:
		if result.err != nil {
			return nil, result.err
		}

		// Check bot's subscription status (optional)
		userID := client.Self.ID
		subscriptionStatus, err := getUserSubscription(userID)
		if err != nil {
			log.Error("Failed to check bot subscription", zap.Error(err))
			// You can choose to return an error or continue (consider logging a warning)
			// return nil, err  // Uncomment if you want to return an error
		}

		if subscriptionStatus != "member" {
			// Bot is not subscribed to its own channel, handle error (optional)
			log.Warn("Bot is not subscribed to the target channel")
			// You can choose to continue or return an error (consider logging a warning)
			// return nil, errors.New("Bot is not subscribed to the target channel") // Uncomment if you want to return an error
		}

		// Check user subscription status (main logic)
		userID = update.Message.From.ID // Replace with logic to get user ID from update
		subscriptionStatus, err = getUserSubscription(userID)
		if err != nil {
			log.Error("Failed to check user subscription", zap.Error(err))
			// Handle error, log the error and continue (consider logging a warning)
			return nil, err // You can choose to return an error here
		}

		if subscriptionStatus != "member" {
			// User is not a member, send a message and restrict access
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please subscribe to our channel: https://t.me/haris_garage to use this bot.")
			bot.Send(msg)
			return nil, errors.New("User is not subscribed to the target channel") // You can return an error here to prevent further processing
		}

		commands.Load(log, client.Dispatcher)
		log.Info("Client started", zap.String("username", client.Self.Username))
		Bot = client.client
		return result.client, nil
	}
}
