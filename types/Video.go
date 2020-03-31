package types

type Video struct {
	FileId       string
	FileUniqueId string
	Width        int
	Height       int
	Duration     int
	Thumb        *PhotoSize
	MimeType     string
	FileSize     int
}
