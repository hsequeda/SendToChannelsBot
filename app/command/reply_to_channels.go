package command

import (
	"context"

	"github.com/stdevHsequeda/SendToChannelsBot/domain"
)

type ReplyToChannels struct {
	Text        string
	HashtagList []string
	UserName    string
	UserId      string
	MessageId   int
}

type ReplyToChannelsHandler struct {
	channelRepository domain.ChannelRepository
	messageSender     domain.MessageSender
}

func NewReplyToChannelsHandler(channelRepository domain.ChannelRepository, messageSender domain.MessageSender) ReplyToChannelsHandler {
	return ReplyToChannelsHandler{channelRepository, messageSender}
}

func (h ReplyToChannelsHandler) Handle(ctx context.Context, cmd ReplyToChannels) error {
	channels, err := h.channelRepository.GetChannelsByHashtags(ctx, cmd.HashtagList)
	if err != nil {
		return err
	}

	for index := range channels {
		if err := h.messageSender.SendMessageToTgChan(channels[index].Id, cmd.Text, cmd.UserId, cmd.UserName, cmd.MessageId); err != nil {
			return err
		}
	}

	return nil
}
