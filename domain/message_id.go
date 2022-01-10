package domain

import "strconv"

type MessageID struct {
	val int64
}

// NewMessageIDFromInt64 build a new MessageID instance from a int64 value.
func NewMessageIDFromInt64(value int64) (MessageID, error) {
	return MessageID{val: value}, nil
}

// NewMessageIDFromInt64 build a new MessageID instance from a string value.
func NewMessageIDFromStr(value string) (MessageID, error) {
	intVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return MessageID{}, err
	}

	return MessageID{val: intVal}, nil
}

func (m MessageID) Value() int64 {
	return m.val
}

func (m MessageID) String() string {
	return strconv.FormatInt(m.val, 10)
}
