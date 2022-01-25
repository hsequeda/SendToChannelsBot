package domain

import "errors"

// Account is a user account
type Account struct {
	id         string
	telegramID TelegramID
	version    uint
	isValid    bool
}

// NewAccount instantiate a new user account in the application.
func NewAccount(id string, telegramID TelegramID) (*Account, error) {
	if id == "" {
		return nil, errors.New("account ID is empty")
	}

	return &Account{
		id:         id,
		telegramID: telegramID,
		version:    1,
		isValid:    true,
	}, nil
}

// TODO: implement me!
// func MapAccountFromRepository(id string, telegramID TelegramID, version uint) *Account {

// 	return &Account{
// 		id:         id,
// 		telegramID: telegramID,
// 		version:    1,
// 		isValid:    true,
// 	}
// }

// ID is the unique identifier of an account
func (a Account) ID() string {
	a.panicIfNotValid()
	return a.id
}

// TelegramID is the telegram account ID related to the current account.
func (a Account) TelegramID() TelegramID {
	a.panicIfNotValid()
	return a.telegramID
}

// Version is a util property to ensure the optimistic locking of the account.
func (a Account) Version() uint {
	a.panicIfNotValid()
	return a.version
}

func (a Account) panicIfNotValid() {
	if !a.isValid {
		panic("an invalid Account entity has been instantiate")
	}
}
