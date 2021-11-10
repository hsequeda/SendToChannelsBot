package adapter

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageSender struct {
	tgBot *tgbotapi.BotAPI
}

func NewMessageSender(tgBot *tgbotapi.BotAPI) *MessageSender {
	return &MessageSender{tgBot}
}

func (ms *MessageSender) SendMessageToTgChan(tgChanId string, text string, userId string, userName string, messageId int) error {
	chatID, err := strconv.ParseInt(tgChanId, 10, 64)
	if err != nil {
		return err
	}

	toSend := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n  \n <a href=\"http://t.me/%s\">Escribir al autor👤(%s)</a>", text, userId, userName))
	toSend.ParseMode = "html"
	toSend.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/DondeHayEnLaHabana/%d", messageId)),
		),
	)
	_, err = ms.tgBot.Send(toSend)
	if err != nil {
		return err
	}

	return nil
}
