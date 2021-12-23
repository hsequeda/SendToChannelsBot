package adapter

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/domain"
)

type TgMessageSender struct {
	tgBot *tgbotapi.BotAPI
}

func NewMessageSender(tgBot *tgbotapi.BotAPI) *TgMessageSender {
	return &TgMessageSender{tgBot}
}

func (t *TgMessageSender) SendMessageToTgChan(tgChanId string, text string, username string, originMessageId string, tgFile domain.TgFile) (domain.ChannelMessage, error) {
	if tgFile == domain.EmptyFile {
		return t.sendMessageToTgChan(tgChanId, text, username, originMessageId, tgFile)
	}

	return t.sendPhotoToTgChan(tgChanId, text, username, originMessageId, tgFile)
}

func (t *TgMessageSender) sendMessageToTgChan(tgChanId string, text string, username string, originMessageId string, tgFile domain.TgFile) (domain.ChannelMessage, error) {
	chatID, err := strconv.ParseInt(tgChanId, 10, 64)
	if err != nil {
		return domain.ChannelMessage{}, err
	}

	toSend := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n  \n <a href=\"http://t.me/%s\">Escribir al autor👤(@%s)</a>", text, username, username))
	toSend.ParseMode = "html"
	toSend.DisableWebPagePreview = true
	toSend.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/DondeHayEnLaHabana/%s", originMessageId)),
		),
	)
	sendedMessage, err := t.tgBot.Send(toSend)
	if err != nil {
		return domain.ChannelMessage{}, err
	}

	return domain.NewChannelMessage(fmt.Sprint(sendedMessage.MessageID), tgChanId)
}

func (t *TgMessageSender) sendPhotoToTgChan(tgChanId string, text string, username string, originMessageId string, tgFile domain.TgFile) (domain.ChannelMessage, error) {
	chatID, err := strconv.ParseInt(tgChanId, 10, 64)
	if err != nil {
		return domain.ChannelMessage{}, err
	}

	toSend := tgbotapi.NewPhotoShare(chatID, tgFile.ID())
	toSend.Caption = fmt.Sprintf("%s\n  \n <a href=\"http://t.me/%s\">Escribir al autor👤(@%s)</a>", text, username, username)
	toSend.ParseMode = "html"
	toSend.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Ir al mensaje", fmt.Sprintf("https://t.me/DondeHayEnLaHabana/%s", originMessageId)),
		),
	)
	sendedMessage, err := t.tgBot.Send(toSend)
	if err != nil {
		return domain.ChannelMessage{}, err
	}

	return domain.NewChannelMessage(fmt.Sprint(sendedMessage.MessageID), tgChanId)
}
