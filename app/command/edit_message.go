package command

import (
	"context"
	"fmt"

	"github.com/stdevHsequeda/SendToChannelsBot/domain"
)

type EditMessage struct {
	MessageID   string
	NewText     string
	NewHashtags []string
}

type EditMessageHandler struct {
	channelRepository domain.ChannelRepository
	messageRepository domain.MessageRepository
}

func (h EditMessageHandler) Hanlde(ctx context.Context, cmd EditMessage) error {
	message, err := h.messageRepository.GetByID(ctx, cmd.MessageID)
	if err != nil {
		return err
	}

	channels, err := h.channelRepository.GetChannelsByHashtags(ctx, message.Hashtags)
	if err != nil {
		return err
	}
	fmt.Printf("channels = %#v\n", channels)

	// filter channel with new hashtag list
	// remove message from channels without hashtags -> needs messageID in each channel
	// edit message in the rest of channels -> needs messageID in each channel

	return nil
}
