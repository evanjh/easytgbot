package configs

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

// Constant values for ChatActions
const (
	ChatTyping         = "typing"
	ChatUploadPhoto    = "upload_photo"
	ChatRecordVideo    = "record_video"
	ChatUploadVideo    = "upload_video"
	ChatRecordAudio    = "record_audio"
	ChatUploadAudio    = "upload_audio"
	ChatUploadDocument = "upload_document"
	ChatFindLocation   = "find_location"
)

// API errors
const (
	// ErrAPIForbidden happens when a token is bad
	ErrAPIForbidden = "forbidden"
)

// Constant values for ParseMode in MessageConfig
const (
	ModeMarkdown = "Markdown"
	ModeHTML     = "HTML"
)

// Library errors
const (
	// ErrBadFileType happens when you pass an unknown type
	ErrBadFileType = "bad file type"
	ErrBadURL      = "bad or empty url"
)
