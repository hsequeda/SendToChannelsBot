package adapter

import (
	"context"

	"github.com/stdevHsequeda/SendToChannelsBot/domain"
)

type MemoryChannelRepository struct {
	channels map[string]domain.Channel
}

func NewMemoryChannelRepository() *MemoryChannelRepository {
	return &MemoryChannelRepository{channels: map[string]domain.Channel{
		"-1001768605818": {
			Id:       "-1001768605818",
			Hashtags: []string{"#te"},
		},
	}}
}

func (m *MemoryChannelRepository) Save(ctx context.Context, channel domain.Channel) error {
	m.channels[channel.Id] = channel
	return nil
}

func (m *MemoryChannelRepository) GetChannelsByHashtags(ctx context.Context, hashtags []string) ([]domain.Channel, error) {
	for _, v := range m.channels {
		for _, hashtag := range v.Hashtags {
			for _, h := range hashtags {
				if h == hashtag {
					return []domain.Channel{v}, nil
				}
			}
		}
	}
	return nil, nil
}
