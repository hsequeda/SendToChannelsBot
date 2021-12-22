package domain

// TgMessageSender is a Telegram message sender.
type TgMessageSender interface {
	SendMessageToTgChan(tgChanId, text, username, originMessageId string, tgFile TgFile) (ChannelMessage, error)
}
