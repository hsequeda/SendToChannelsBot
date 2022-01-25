package command

import "context"

type AddBotToAccount struct{}

type AddBotToAccountHandler struct{}

func NewAddBotToAccountHandler() AddBotToAccountHandler {
	return AddBotToAccountHandler{}
}

func (h AddBotToAccountHandler) Handle(ctx context.Context, cmd AddBotToAccount) error {
	panic("not implemented") // TODO: Implement
}
