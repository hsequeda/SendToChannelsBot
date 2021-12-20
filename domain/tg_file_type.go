package domain

type TgFileType struct {
	s string
}

var (
	TgFileTypePhoto    = TgFileType{"photo"}
	TgFileTypeDocument = TgFileType{"document"}
)
