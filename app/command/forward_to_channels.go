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
	messageRepository domain.MessageRepository
	messageSender     domain.TgMessageSender
}

func NewForwardToChannelsHandler(
	channelRepository domain.ChannelRepository,
	messageRepository domain.MessageRepository,
	messageSender domain.TgMessageSender,
) ForwardToChannelsHandler {
	return ForwardToChannelsHandler{channelRepository, messageRepository, messageSender}
}

func (h ForwardToChannelsHandler) Handle(ctx context.Context, cmd ForwardToChannels) error {
	channels, err := h.channelRepository.GetChannelsByHashtags(ctx, cmd.HashtagList)
	if err != nil {
		return err
	}

	message, err := domain.NewMessage(fmt.Sprint(cmd.MessageID), cmd.Text, cmd.HashtagList)
	if err != nil {
		return err
	}

	for index := range channels {
		channelMessage, err := h.messageSender.SendMessageToTgChan(channels[index].Id, cmd.Text, cmd.UserName, fmt.Sprint(cmd.MessageID), cmd.File)
		if err != nil {
			return err
		}

		message.AddChannelMessage(channelMessage)
	}

	return h.messageRepository.Save(ctx, message)
}
