package types

type VideoNote struct {
	FileId       string
	FileUniqueId string
	Length       int
	Duration     int
	Thumb        *PhotoSize
	FileSize     int
}
