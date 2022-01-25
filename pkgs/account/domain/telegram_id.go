package domain

import (
	"strconv"
)

// TelegramID represent a Telegram account ID.
type TelegramID struct {
	val int64
}

func NewTelegramID(val int64) (TelegramID, error) {
	return TelegramID{val}, nil
}

// Value returns the account value as int64.
func (t TelegramID) Value() int64 {
	return t.val
}

// String returns the account value as string.
func (t TelegramID) String() string {
	return strconv.FormatInt(t.val, 10)
}
