package domain

import (
	"context"
	"errors"
)

var ErrAccountConflictOnSave = errors.New("conflic when saving an account: the account version is different from the persisted account")

// AccountRepository is an implementation of the repository pattern for the Account entity.
type AccountRepository interface {
	// Save create or update an account in the repository. Verify optimistic locking.
	Save(ctx context.Context, account *Account) error
	// GetOneByTgID returns an account by a TelegramID.
	GetOneByTgID(ctx context.Context, tgID TelegramID) (*Account, error)
	// GetOneByID returns an account by an ID.
	GetOneByID(ctx context.Context, id string) (*Account, error)
}
