package types

type Audio struct {
	FileId       string
	FileUniqueId string
	Duration     int
	Performer    string     
	Title        string     
	MimeType     string     
	FileSize     int        
	Thumb        *PhotoSize 
}
