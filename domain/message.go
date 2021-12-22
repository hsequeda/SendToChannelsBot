package domain

import "errors"

type Message struct {
	ID              string
	Text            string
	Hashtags        []string
	ChannelMessages []ChannelMessage
}

func NewMessage(id string, text string, hashtags []string) (Message, error) {
	if id == "" {
		return Message{}, errors.New("message id is empty")
	}

	if text == "" {
		return Message{}, errors.New("message text is empty")
	}

	if len(hashtags) == 0 {
		return Message{}, errors.New("can't create a message without hashtags")
	}

	return Message{
		ID:       id,
		Text:     text,
		Hashtags: hashtags,
	}, nil
}

func (m Message) AddChannelMessage(channelMessage ChannelMessage) {
	m.ChannelMessages = append(m.ChannelMessages, channelMessage)
}
