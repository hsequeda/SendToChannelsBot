package adapter_test

import (
	"fmt"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/domain"
	"github.com/stretchr/testify/require"
)

func TestTgMessageSender_SendMessageToTgChan(t *testing.T) {
	tgMessageSender := newTgMessageSender(t)
	t.Run("Ok", func(t *testing.T) {
		chatID := os.Getenv("TEST_CHAT_ID")
		if chatID == "" {
			t.Error("test chat id is empty")
			return
		}

		message, err := tgMessageSender.SendMessageToTgChan(chatID, "a message to test the message sender", "test_username", "origin message id", domain.EmptyFile)
		require.NoError(t, err)
		fmt.Printf("message = %#v\n", message)
	})
}

func newTgMessageSender(t *testing.T) *adapter.TgMessageSender {
	t.Helper()
	testBotToken := os.Getenv("TEST_BOT_TOKEN")
	if testBotToken == "" {
		t.Error("test bot token is empty")
	}

	bot, err := tgbotapi.NewBotAPI(testBotToken)
	require.NoError(t, err)
	return adapter.NewMessageSender(bot)
}
