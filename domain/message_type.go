package domain

type MessageType struct {
	s string
}

var (
	MessageTypeDocument = MessageType{"message_type_document"}
	MessageTypePhoto    = MessageType{"message_type_photo"}
	MessageTypeTextOnly = MessageType{"message_type_text_only"}
)

func NewMessageTypeFromStr(val string) MessageType {
	switch val {
	case MessageTypeDocument.Value():
		return MessageTypeDocument
	case MessageTypePhoto.Value():
		return MessageTypePhoto
	default:
		return MessageTypeTextOnly
	}
}

func (mt MessageType) Value() string {
	if mt.IsZero() {
		return MessageTypeTextOnly.s
	}

	return mt.s
}

func (mt MessageType) IsZero() bool {
	return mt.s == ""
}
