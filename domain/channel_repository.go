package domain

import "context"

type ChannelRepository interface {
	Save(ctx context.Context, channel Channel) error
	GetChannelsByHashtags(ctx context.Context, hashtags []string) ([]Channel, error)
}
