package domain

type MessageSender interface {
	SendMessageToTgChan(tgChanId string, text string, userId string, userName string, messageId int) error
}

// TgMessageSender is a Telegram message sender.
type TgMessageSender interface {
	SendMessageToTgChan(tgChanId, text, username, originMessageId string, tgFile TgFile) error
}
