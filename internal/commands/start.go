package commands

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/utils"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
	"fmt"
)

func (m *command) LoadStart(dispatcher dispatcher.Dispatcher) {
	log := m.log.Named("start")
	defer log.Sugar().Info("Loaded")
	dispatcher.AddHandler(handlers.NewCommand("start", start))
}

func start(ctx *ext.Context, u *ext.Update) error {
	chatId := u.EffectiveChat().GetID()
	peerChatId := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatId.Type != int(storage.TypeUser) {
		return dispatcher.EndGroups
	}
	if len(config.ValueOf.AllowedUsers) != 0 && !utils.Contains(config.ValueOf.AllowedUsers, chatId) {
		ctx.Reply(u, "You are not allowed to use this bot.", nil)
		return dispatcher.EndGroups
	}

	// Check if force sub is enabled and user is subscribed
	if config.ValueOf.ForceSubChannelID != 0 {
		isSubscribed, err := utils.IsUserSubscribed(ctx, ctx.GetClient(), chatId)
		if err != nil {
			ctx.Reply(u, "Error checking subscription status. Please try again later.", nil)
			return dispatcher.EndGroups
		}
		if !isSubscribed {
			row := tg.KeyboardButtonRow{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonURL{
						Text: "Join Channel",
						URL:  fmt.Sprintf("https://t.me/%d", config.ValueOf.ForceSubChannelID),
					},
				},
			}
			markup := &tg.ReplyInlineMarkup{
				Rows: []tg.KeyboardButtonRow{row},
			}
			ctx.Reply(u, "Please join our channel to use this bot.", &ext.ReplyOpts{
				Markup: markup,
			})
			return dispatcher.EndGroups
		}
	}

	ctx.Reply(u, "Need a direct streamable link to a file? Send it my way! 🤓\n\nJoin my Update Channel @haris_garage 🗿 for more updates.\n\nLink validity: 24 hours ⏳\n\nPro Tip: Use 1DM Browser for lightning-fast downloads! 🔥", nil)
	return dispatcher.EndGroups
}
