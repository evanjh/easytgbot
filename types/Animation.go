package types

type Animation struct {
	FileId       string
	FileUniqueId string
	Width        int
	Height       int
	Duration     int
	Thumb        *PhotoSize
	FileName     string
	MimeType     string
	FileSize     int
}
