package domain

import (
	"context"
	"errors"
)

var ErrInputConflictOnSave = errors.New("conflic when saving an input: the input version is different from the persisted input")

type InputRepository interface {
	Save(ctx context.Context, input Input) error
}
