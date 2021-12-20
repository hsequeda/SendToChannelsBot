package domain

type TgFile struct {
	id         string
	tgFileType TgFileType
}

var EmptyFile = TgFile{}

func NewTgFile(id string, tgFileType TgFileType) (TgFile, error) {
	return TgFile{id: id, tgFileType: tgFileType}, nil
}
