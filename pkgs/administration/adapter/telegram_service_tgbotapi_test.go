package adapter_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramServiceTgBotApi_IsChatValid(t *testing.T) {
	// >> Initialize variables
	tgService := initializeTelegramServiceTgBotApi(t)
	validChatIDStr, ok := os.LookupEnv("VALID_CHAT_ID")
	if !ok {
		t.Skip("skip: VALID_CHAT_ID required")
	}

	validChatID, err := strconv.ParseInt(validChatIDStr, 10, 64)
	require.NoError(t, err)

	validUserIDStr, ok := os.LookupEnv("VALID_USER_ID")
	if !ok {
		t.Skip("skip: VALID_USER_ID  required ")
	}

	invalidChatIDStr, ok := os.LookupEnv("INVALID_CHAT_ID")
	if !ok {
		t.Skip("skip: VALID_CHAT_ID required")
	}

	invalidChatID, err := strconv.ParseInt(invalidChatIDStr, 10, 64)
	require.NoError(t, err)

	validUserID, err := strconv.ParseInt(validUserIDStr, 10, 64)
	require.NoError(t, err)

	t.Run("Ok case", func(t *testing.T) {
		t.Parallel()
		ok, err = tgService.IsChatValid(context.Background(), validChatID, validUserID)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("Error: Chat doesn't exist", func(t *testing.T) {
		t.Parallel()
		validUserID, err := strconv.ParseInt(validUserIDStr, 10, 64)
		require.NoError(t, err)

		ok, err = tgService.IsChatValid(context.Background(), 1000000000, validUserID)
		require.Error(t, err)
		assert.IsType(t, domain.ChatNotFoundError{}, err)
		assert.False(t, ok)
	})

	t.Run("Error: Chat doesn't below to user", func(t *testing.T) {
		validUserID, err := strconv.ParseInt(validUserIDStr, 10, 64)
		require.NoError(t, err)

		ok, err = tgService.IsChatValid(context.Background(), invalidChatID, validUserID)
		require.Error(t, err)
		assert.IsType(t, domain.ChatDoesNotBelowToUserError{}, err)
		assert.False(t, ok)
	})
}

func initializeTelegramServiceTgBotApi(t *testing.T) *adapter.TelegramServiceTgBotApi {
	t.Helper()
	testBotToken := os.Getenv("TEST_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(testBotToken)
	require.NoError(t, err)

	return adapter.NewTelegramServiceTgBotApi(bot)
}
