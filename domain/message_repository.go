package domain

type MessageRepository interface {
	Save() error
}
