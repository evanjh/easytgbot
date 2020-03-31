package types

// Document contains information about a document.
type Document struct {
	FileId       string
	FileUniqueId string
	Thumb        *PhotoSize
	FileName     string
	MimeType     string
	FileSize     int
}
