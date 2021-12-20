package command

import (
	"context"
	"fmt"

	"github.com/stdevHsequeda/SendToChannelsBot/domain"
)

type ForwardToChannels struct {
	Text        string
	HashtagList []string
	UserName    string
	UserID      string
	MessageID   int
	MessageType domain.MessageType
	File        domain.TgFile
}

type ForwardToChannelsHandler struct {
	channelRepository domain.ChannelRepository
	messageSender     domain.TgMessageSender
}

func NewForwardToChannelsHandler(channelRepository domain.ChannelRepository, messageSender domain.TgMessageSender) ForwardToChannelsHandler {
	return ForwardToChannelsHandler{channelRepository, messageSender}
}

func (h ForwardToChannelsHandler) Handle(ctx context.Context, cmd ForwardToChannels) error {
	channels, err := h.channelRepository.GetChannelsByHashtags(ctx, cmd.HashtagList)
	if err != nil {
		return err
	}

	for index := range channels {
		if err := h.messageSender.SendMessageToTgChan(channels[index].Id, cmd.Text, cmd.UserName, fmt.Sprint(cmd.MessageID), cmd.File); err != nil {
			return err
		}
	}

	return nil
}
