package telegram

import (
	"context"
	"log"
	"strconv"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/app/command"
)

type TelegramBotUpdateHandler struct {
	updateChan      tgbotapi.UpdatesChannel
	replyToChannels command.ReplyToChannelsHandler
}

const (
	CommandReported = "informo"
)

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

			t.ReplyToChannels(message)
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
		if err := t.replyToChannels.Handle(context.TODO(), command.ReplyToChannels{
			Text:        string(utf16.Decode(shortedText)),
			HashtagList: message.Hashtags().StrHashtags(),
			UserName:    realMsgAuthorEntity.UsernameStr(),
			UserId:      strconv.Itoa(message.From.ID),
			MessageId:   message.MessageID,
		}); err != nil {
			log.Println(err)
		}
	}
}

func (t *TelegramBotUpdateHandler) ReplyToChannels(message Message) {
	if message.HasHashtag() { //  TODO add  ( or len(message.CaptionEntities) != 0)
		if err := t.replyToChannels.Handle(context.TODO(), command.ReplyToChannels{
			Text:        message.Text,
			HashtagList: message.Hashtags().StrHashtags(),
			UserName:    message.From.UserName,
			UserId:      strconv.Itoa(message.From.ID),
			MessageId:   message.MessageID,
		}); err != nil {
			log.Println(err)
		}
	}
}
