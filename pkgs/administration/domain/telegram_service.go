package domain

import (
	"context"
	"fmt"
)

// ChatNotFoundError is returned when the passed chatID doesn't match with
// any chat.
type ChatNotFoundError struct {
	ChatID int64
}

func (e ChatNotFoundError) Error() string {
	return fmt.Sprintf("chat with [id: '%d'] not found", e.ChatID)
}

// ChatDoesNotBelowToUserError is returned when the user isn't the owner
// of the chat.
type ChatDoesNotBelowToUserError struct {
	UserID, ChatID int64
}

func (e ChatDoesNotBelowToUserError) Error() string {
	return fmt.Sprintf("chat '%d' doesn't below to  user '%d'", e.ChatID, e.UserID)
}

// TelegramService represent a service to access to TelegramAPI
type TelegramService interface {
	// IsChatValid verify if chat below to the user ID and if the current bot
	// have the required permissions on that chat.
	IsChatValid(ctx context.Context, chatId, userId int64) (bool, error)
}
