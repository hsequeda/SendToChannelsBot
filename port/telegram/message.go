package telegram

import (
	"errors"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	EntityTypeMention     = "mention"
	EntityTypeTextMention = "text_mention"
	EntityTypeHashtag     = "hashtag"
)

type Message struct {
	tgbotapi.Message
}

func MessageFromTgbotApiMessage(message *tgbotapi.Message) (Message, error) {
	if message == nil {
		return Message{}, errors.New("error creating Message: passed tg-message is empty")
	}

	return Message{*message}, nil
}

// Utf16Text returns the message text in Utf-16 (Format using by telegram)
func (m Message) Utf16Text() []uint16 {
	return utf16.Encode([]rune(m.Text))
}

func (m Message) HasHashtag() bool {
	for i := range m.Entities {
		if m.Entities[i].Type == EntityTypeHashtag {
			return true
		}
	}

	return false
}

// Hashtags returns the list of entities of 'hashtag' type. This solution was
// taked from https://github.com/go-telegram-bot-api/telegram-bot-api/issues/231
func (m Message) Hashtags() HashtagEntities {
	if !m.HasHashtag() {
		return nil
	}

	entities := m.Entities // TODO add conditionals to handle caption entities

	hashtags := make([]HashtagEntity, 0)
	for _, e := range entities {
		if e.Type == EntityTypeHashtag {
			hashtagEntity, _ := NewHashtagEntity(e, m.Utf16Text())
			hashtags = append(hashtags, hashtagEntity)
		}
	}

	return hashtags
}

func (m Message) Mentions() []MentionEntity {
	entities := m.Entities // TODO add conditionals to handle caption entities

	mentions := make([]MentionEntity, 0)
	for _, e := range entities {
		if e.Type == EntityTypeMention {
			mentionEntity, _ := NewMentionEntity(e, m.Utf16Text())
			mentions = append(mentions, mentionEntity)
		}
	}

	return mentions
}

// IsCaption returns if Message is a "Caption Message"
func (m Message) IsCaption() bool {
	return m.Caption != ""
}
