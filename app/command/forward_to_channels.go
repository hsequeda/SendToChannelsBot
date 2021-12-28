package command

import (
	"context"
	"fmt"
	"log"

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
	log.Println("[ForwardToChannels] info: start command")
	channels, err := h.channelRepository.GetChannelsByHashtags(ctx, cmd.HashtagList)
	if err != nil {
		log.Println("[ForwardToChannels] error: ", err.Error())
		return err
	}

	message, err := domain.NewMessage(fmt.Sprint(cmd.MessageID), cmd.Text, cmd.HashtagList)
	if err != nil {
		log.Println("[ForwardToChannels] error: ", err.Error())
		return err
	}

	for index := range channels {
		channelMessage, err := h.messageSender.SendMessageToTgChan(channels[index].Id, cmd.Text, cmd.UserName, fmt.Sprint(cmd.MessageID), cmd.File)
		if err != nil {
			log.Println("[ForwardToChannels] error: ", err.Error())
			return err
		}

		message.AddChannelMessage(channelMessage)
	}

	if err := h.messageRepository.Save(ctx, message); err != nil {
		log.Println("[ForwardToChannels] error: ", err.Error())
	}

	log.Println("[ForwardToChannels] info: end OK")
	return nil
}
