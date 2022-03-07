package adapter

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/domain"
)

type TelegramServiceTgBotApi struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramServiceTgBotApi TODO
func NewTelegramServiceTgBotApi(bot *tgbotapi.BotAPI) *TelegramServiceTgBotApi {
	return &TelegramServiceTgBotApi{bot: bot}
}

// IsChatValid verify if chat below to the user ID and if the current bot
// have the required permissions on that chat.
func (t *TelegramServiceTgBotApi) IsChatValid(ctx context.Context, chatID int64, userID int64) (bool, error) {
	fmt.Printf("chatID = %#v\n", chatID)
	chat, err := t.bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: int64(chatID),
		},
	})
	if err != nil {
		if err.Error() == "Bad Request: chat not found" {
			return false, domain.ChatNotFoundError{ChatID: chatID}
		}

		panic(err) // Unexpected case
	}

	if chat.Type == "private" {
		if userID == chatID {
			return true, nil
		}

		return false, domain.ChatDoesNotBelowToUserError{ChatID: chatID, UserID: userID}
	}

	user, err := t.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: int(userID),
		},
	})
	if err != nil {
		panic(err) // Unexpected case
	}

	if user.Status != "creator" {
		return true, domain.ChatDoesNotBelowToUserError{ChatID: chatID, UserID: userID}
	}

	return true, nil
}

var (
	_ domain.TelegramService = &TelegramServiceTgBotApi{}
)

func panicOnUnexpected(err error) {
	panic(err)
}
