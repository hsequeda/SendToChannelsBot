package domain

import "errors"

type Message struct {
	ID              MessageID
	Text            string
	Hashtags        []string
	ChannelMessages []ChannelMessage
}

func NewMessage(id MessageID, text string, hashtags []string) (Message, error) {
	if len(hashtags) == 0 {
		return Message{}, errors.New("can't create a message without hashtags")
	}

	return Message{
		ID:       id,
		Text:     text,
		Hashtags: hashtags,
	}, nil
}

func (m *Message) AddChannelMessage(channelMessage ChannelMessage) {
	m.ChannelMessages = append(m.ChannelMessages, channelMessage)
}
