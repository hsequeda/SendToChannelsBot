package domain

type Input struct {
	id          int64
	name        string
	ownerID     int64
	inputType   InputType
	description string
	isValid     bool
	version     int64
}

func (i Input) ID() int64 {
	return i.id
}

func (i Input) Name() string {
	return i.name
}

func (i Input) OwnerID() int64 {
	return i.ownerID
}

func (i Input) InputType() InputType {
	return i.inputType
}

func (i Input) Description() string {
	return i.description
}

func (i Input) Version() int64 {
	return i.version
}
