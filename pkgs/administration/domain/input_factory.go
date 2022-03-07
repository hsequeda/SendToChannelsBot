package domain

import (
	"context"
	"errors"
)

var InvalidChatErr = errors.New("chat no valid")

type InputFactoryImpl struct {
	telegramClient TelegramService
}

func (f InputFactoryImpl) Build(ctx context.Context, id int64, ownerID int64, inputType InputType, name string, description string) (Input, error) {
	if name == "" {
		return Input{}, errors.New("name is empty")
	}

	ok, err := f.telegramClient.IsChatValid(ctx, id, ownerID)
	if err != nil {
		return Input{}, err
	}

	if !ok {
		return Input{}, InvalidChatErr
	}

	return Input{
		id:          id,
		name:        name,
		ownerID:     ownerID,
		inputType:   inputType,
		description: description,
		version:     0,
		isValid:     true,
	}, nil
}
