package command

import (
	"context"

	"github.com/stdevHsequeda/SendToChannelsBot/account/domain"
)

// RegisterAccount TODO
type RegisterAccount struct {
	ID         string
	TelegramID int64
}

// RegisterAccountHandler is a handler for the RegisterAccount command.
type RegisterAccountHandler struct {
	accountRepository domain.AccountRepository
}

// NewRegisterAccountHandler insantiate an instance of NewRegisterAccount.
func NewRegisterAccountHandler(accountRepository domain.AccountRepository) RegisterAccountHandler {
	return RegisterAccountHandler{accountRepository}
}

// Handle execute the RegisterAccount command.
func (h RegisterAccountHandler) Handle(ctx context.Context, cmd RegisterAccount) error {
	telegramID, err := domain.NewTelegramID(cmd.TelegramID)
	if err != nil {
		return err
	}

	account, err := domain.NewAccount(cmd.ID, telegramID)
	if err != nil {
		return err
	}

	if err := h.accountRepository.Save(ctx, account); err != nil {
		return err
	}

	return nil
}
