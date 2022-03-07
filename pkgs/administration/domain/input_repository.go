package domain

import (
	"context"
)

type InputRepository interface {
	Save(ctx context.Context, input Input) error
}
