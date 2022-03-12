package command

import (
	"context"

	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/domain"
)

// AddInput TODO
type AddInput struct {
	UserID      int64
	ChatID      int64
	InputType   string
	Name        string
	Description string
}

// AddInputHandler is a handler for the AddInput command.
type AddInputHandler struct {
	inputFactory    domain.InputFactory
	inputRepository domain.InputRepository
}

// NewAddInputHandler insantiate an instance of AddInputHandler.
func NewAddInputHandler(inputFactory domain.InputFactory, inputRepository domain.InputRepository) AddInputHandler {
	return AddInputHandler{inputFactory: inputFactory}
}

// Handle execute the AddInput command.
func (h AddInputHandler) Handle(ctx context.Context, cmd AddInput) error {
	input, err := h.inputFactory.Build(ctx, cmd.ChatID, cmd.UserID, domain.NewInputTypeFromStr(cmd.InputType), cmd.Name, cmd.Description)
	if err != nil {
		return err
	}

	if err := h.inputRepository.Save(ctx, input); err != nil {
		return err
	}

	return nil
}
