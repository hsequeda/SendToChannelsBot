package domain

import "context"

type MessageRepository interface {
	Save(ctx context.Context, msg Message) error
	GetByID(ctx context.Context, messageID string) (Message, error)
}
