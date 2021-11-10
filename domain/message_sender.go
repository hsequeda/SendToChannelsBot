package domain

type MessageSender interface {
	SendMessageToTgChan(tgChanId string, text string, userId string, userName string, messageId int) error
}
