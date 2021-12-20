package telegram

import (
	"errors"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MentionEntity struct {
	tgbotapi.MessageEntity
	username []uint16
}

func NewMentionEntity(entity tgbotapi.MessageEntity, utf16Text []uint16) (MentionEntity, error) {
	if entity.Type != EntityTypeMention {
		return MentionEntity{}, errors.New("mention entity need to be of type 'mention'")
	}

	username := utf16Text[entity.Offset+1 : entity.Offset+entity.Length]
	return MentionEntity{
		MessageEntity: entity,
		username:      username,
	}, nil
}

func (me MentionEntity) UsernameStr() string {
	return string(utf16.Decode(me.username))
}

func (me MentionEntity) Username() []uint16 {
	return me.username
}
