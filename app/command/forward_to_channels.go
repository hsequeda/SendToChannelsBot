package command

import (
	"context"
	"fmt"
	"log"

	"github.com/stdevHsequeda/SendToChannelsBot/domain"
)

// The ForwardToChannels command takes a listened message from the groups and
// gets the channels that match the hashtags in the message. It then sends a
// copy of the message to each channel with a reference to the owner and the
// original message.
type ForwardToChannels struct {
	Text        string
	HashtagList []string
	UserName    string
	UserID      string
	MessageID   int64
	MessageType domain.MessageType
	File        domain.TgFile
}

// ForwardToChannelsHandler handles the ForwardToChannels command.
type ForwardToChannelsHandler struct {
	channelRepository domain.ChannelRepository
	messageRepository domain.MessageRepository
	messageSender     domain.TgMessageSender
}

// NewForwardToChannelsHandler build a new instance of ForwardToChannelsHandler.
func NewForwardToChannelsHandler(
	channelRepository domain.ChannelRepository,
	messageRepository domain.MessageRepository,
	messageSender domain.TgMessageSender,
) ForwardToChannelsHandler {
	return ForwardToChannelsHandler{channelRepository, messageRepository, messageSender}
}

// Handle execute the command logic.
func (h ForwardToChannelsHandler) Handle(ctx context.Context, cmd ForwardToChannels) error {
	log.Println("[ForwardToChannels] info: start command")
	channels, err := h.channelRepository.GetChannelsByHashtags(ctx, cmd.HashtagList)
	if err != nil {
		log.Println("[ForwardToChannels] error: ", err.Error())
		return err
	}

	messageID, err := domain.NewMessageIDFromInt64(cmd.MessageID)
	if err != nil {
		return err
	}

	message, err := domain.NewMessage(messageID, cmd.Text, cmd.HashtagList)
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
