package domain

type TgFile struct {
	id         string
	tgFileType TgFileType
}

var EmptyFile = TgFile{}

func NewTgFile(id string, tgFileType TgFileType) (TgFile, error) {
	return TgFile{id: id, tgFileType: tgFileType}, nil
}

func (f TgFile) ID() string {
	return f.id
}

func (f TgFile) IsZero() bool {
	return f == EmptyFile
}
