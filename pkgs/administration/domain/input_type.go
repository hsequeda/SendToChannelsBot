package domain

var (
	TelegramPrivate    InputType = InputType{s: "tg_private"}
	TelegramChannel    InputType = InputType{s: "tg_channel"}
	TelegramGroup      InputType = InputType{s: "tg_group"}
	TelegramSuperGroup InputType = InputType{s: "tg_supergroup"}
)

type InputType struct {
	s string
}

func NewInputTypeFromStr(value string) InputType {
	switch value {
	case "tg_private":
		return TelegramPrivate
	case "tg_channel":
		return TelegramChannel
	case "tg_group":
		return TelegramGroup
	case "tg_supergroup":
		return TelegramSuperGroup
	default:
		panic("")
	}
}

func (i InputType) Value() string {
	return i.s
}
