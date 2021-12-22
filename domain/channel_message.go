package domain

type ChannelMessage struct {
	ID        string
	ChannelID string
}

func NewChannelMessage(messageID string, channelID string) (ChannelMessage, error) {
	return ChannelMessage{
		ID:        messageID,
		ChannelID: channelID,
	}, nil
}
