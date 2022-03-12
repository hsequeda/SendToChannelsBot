package domain

import (
	"context"
	"errors"
)

type InputFactory struct {
	telegramClient TelegramService
}

func NewInputFactory(telegramClient TelegramService) InputFactory {
	return InputFactory{telegramClient}
}

func (f InputFactory) Build(ctx context.Context, id int64, ownerID int64, inputType InputType, name string, description string) (Input, error) {
	if name == "" {
		return Input{}, errors.New("name is empty")
	}

	ok, err := f.telegramClient.IsChatValid(ctx, id, ownerID)
	if err != nil {
		return Input{}, err
	}

	if !ok {
		return Input{}, errors.New("chat no valid")
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
