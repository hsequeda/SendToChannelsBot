package telegram

import (
	"context"
	"log"
	"strconv"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/app/command"
	"github.com/stdevHsequeda/SendToChannelsBot/domain"
)

const (
	CommandReported = "informo"
)

type TelegramBotUpdateHandler struct {
	updateChan        tgbotapi.UpdatesChannel
	forwardToChannels command.ForwardToChannelsHandler
}

func NewTelegramBotUpdateHandler(
	updateChan tgbotapi.UpdatesChannel,
	forwardToChannels command.ForwardToChannelsHandler,
) TelegramBotUpdateHandler {
	return TelegramBotUpdateHandler{
		updateChan:        updateChan,
		forwardToChannels: forwardToChannels,
	}
}

func (t *TelegramBotUpdateHandler) Run() {
	for telegramUpdate := range t.updateChan {
		if telegramUpdate.Message != nil {
			message, err := MessageFromTgbotApiMessage(telegramUpdate.Message)
			if err != nil {
				log.Println(err)
				continue
			}

			if message.IsCommand() {
				switch message.Command() {
				case CommandReported:
					t.reportCommand(message)
				default:
					log.Println("error: unhandle command")
				}

				continue
			}

			t.forwardToChannel(message)
		}
	}
}

func (t *TelegramBotUpdateHandler) reportCommand(message Message) {
	if len(message.Mentions()) == 0 {
		log.Println("error: missed mention in report command")
	}

	realMsgAuthorEntity := message.Mentions()[0]
	shortedText := message.Utf16Text()[realMsgAuthorEntity.Offset+realMsgAuthorEntity.Length:]
	if message.HasHashtag() {
		messageType := domain.MessageTypeTextOnly
		file := domain.EmptyFile

		if message.PhotoID() != "" {
			messageType = domain.MessageTypePhoto
			file, _ = domain.NewTgFile(message.PhotoID(), domain.TgFileTypePhoto)
		}

		if err := t.forwardToChannels.Handle(context.TODO(), command.ForwardToChannels{
			Text:        string(utf16.Decode(shortedText)),
			HashtagList: message.Hashtags().StrHashtags(),
			UserName:    realMsgAuthorEntity.UsernameStr(),
			UserID:      strconv.Itoa(message.From.ID),
			MessageID:   message.MessageID,
			MessageType: messageType,
			File:        file,
		}); err != nil {
			log.Println(err)
		}
	}
}

func (t *TelegramBotUpdateHandler) forwardToChannel(message Message) {
	if message.HasHashtag() {
		messageType := domain.MessageTypeTextOnly
		file := domain.EmptyFile

		if message.PhotoID() != "" {
			messageType = domain.MessageTypePhoto
			file, _ = domain.NewTgFile(message.PhotoID(), domain.TgFileTypePhoto)
		}

		if err := t.forwardToChannels.Handle(context.TODO(), command.ForwardToChannels{
			Text:        string(utf16.Decode(message.Utf16Text())),
			HashtagList: message.Hashtags().StrHashtags(),
			UserName:    message.From.UserName,
			UserID:      strconv.Itoa(message.From.ID),
			MessageID:   message.MessageID,
			MessageType: messageType,
			File:        file,
		}); err != nil {
			log.Println(err)
		}
	}
}
