package telegram

import (
	"errors"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HashtagEntities []HashtagEntity

func (hes HashtagEntities) StrHashtags() []string {
	var strHashtag = make([]string, len(hes), cap(hes))
	for i, he := range hes {
		strHashtag[i] = he.HashtagStr()
	}

	return strHashtag
}

type HashtagEntity struct {
	tgbotapi.MessageEntity
	hashtag []uint16
}

func NewHashtagEntity(entity tgbotapi.MessageEntity, text []uint16) (HashtagEntity, error) {
	if entity.Type != EntityTypeHashtag {
		return HashtagEntity{}, errors.New("hashtag entity need to be of type 'hashtag'")
	}

	hashtag := text[entity.Offset : entity.Offset+entity.Length]
	return HashtagEntity{
		MessageEntity: entity,
		hashtag:       hashtag,
	}, nil
}

func (me HashtagEntity) HashtagStr() string {
	return string(utf16.Decode(me.hashtag))
}

func (me HashtagEntity) Hashtag() []uint16 {
	return me.hashtag
}
