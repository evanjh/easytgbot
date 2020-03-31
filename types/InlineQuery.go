package types

// InlineQuery is a Query from Telegram for an inline request.
type InlineQuery struct {
	Id       string
	From     *User
	Location *Location
	Query    string
	Offset   string
}

// InlineQueryResultArticle is an inline query response article.
type InlineQueryResultArticle struct {
	Type                string
	Id                  string
	Title               string
	InputMessageContent interface{}
	ReplyMarkup         *InlineKeyboardMarkup
	Url                 string
	HideUrl             bool
	Description         string
	ThumbUrl            string
	ThumbWidth          int
	ThumbHeight         int
}

// InlineQueryResultPhoto is an inline query response photo.
type InlineQueryResultPhoto struct {
	Type                string
	Id                  string
	PhotoUrl            string
	ThumbUrl            string
	PhotoWidth          int
	PhotoHeight         int
	Title               string
	Description         string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

// InlineQueryResultGIF is an inline query response GIF.
type InlineQueryResultGif struct {
	Type                string
	Id                  string
	GifUrl              string
	GifWidth            int
	GifHeight           int
	GifDuration         int
	ThumbUrl            string
	Title               string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

// InlineQueryResultMPEG4GIF is an inline query response MPEG4 GIF.
type InlineQueryResultMpeg4Gif struct {
	Type                string
	Id                  string
	Mpeg4Url            string
	Mpeg4Width          int
	Mpeg4Height         int
	Mpeg4Duration       int
	ThumbUrl            string
	Title               string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

// InlineQueryResultVideo is an inline query response video.
type InlineQueryResultVideo struct {
	Type                string
	Id                  string
	VideoUrl            string
	MimeType            string
	ThumbUrl            string
	Title               string
	Caption             string
	ParseMode           string
	VideoWidth          int
	VideoHeight         int
	VideoDuration       int
	Description         string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

// InlineQueryResultAudio is an inline query response audio.
type InlineQueryResultAudio struct {
	Type                string
	Id                  string
	AudioUrl            string
	Title               string
	Caption             string
	ParseMode           string
	Performer           string
	AudioDuration       int
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

// InlineQueryResultVoice is an inline query response voice.
type InlineQueryResultVoice struct {
	Type                string
	Id                  string
	VoiceUrl            string
	Title               string
	Caption             string
	ParseMode           string
	VoiceDuration       int
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

// InlineQueryResultDocument is an inline query response document.
type InlineQueryResultDocument struct {
	Type                string
	Id                  string
	Title               string
	Caption             string
	ParseMode           string
	DocumentUrl         string
	MimeType            string
	Description         string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
	ThumbUrl            string
	ThumbWidth          int
	ThumbHeight         int
}

// InlineQueryResultLocation is an inline query response location.
type InlineQueryResultLocation struct {
	Type                string
	Id                  string
	Latitude            float64
	Longitude           float64
	Title               string
	LivePeriod          int
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
	ThumbUrl            string
	ThumbWidth          int
	ThumbHeight         int
}

// InlineQueryResultVenue is an inline query response venue.
type InlineQueryResultVenue struct {
	Type                string
	Id                  string
	Latitude            float64
	Longitude           float64
	Title               string
	Address             string
	FoursquareId        string
	FoursquareType      string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
	ThumbUrl            string
	ThumbWidth          int
	ThumbHeight         int
}
type InlineQueryResultContact struct {
	Type                string
	Id                  string
	PhoneNumber         string
	FirstName           string
	LastName            string
	Vcard               string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
	ThumbUrl            string
	ThumbWidth          int
	ThumbHeight         int
}

// InlineQueryResultGame is an inline query response game.
type InlineQueryResultGame struct {
	Type          string
	Id            string
	GameShortName string
	ReplyMarkup   *InlineKeyboardMarkup
}

type InlineQueryResultCachedPhoto struct {
	Type                string
	Id                  string
	PhotoFileId         string
	Title               string
	Description         string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

type InlineQueryResultCachedGif struct {
	Type                string
	Id                  string
	GifFileId           string
	Title               string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

type InlineQueryResultCachedMpeg4Gif struct {
	Type                string
	Id                  string
	Mpeg4FileId         string
	Title               string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

type InlineQueryResultCachedSticker struct {
	Type                string
	Id                  string
	StickerFileId       string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

type InlineQueryResultCachedDocument struct {
	Type                string
	Id                  string
	Title               string
	DocumentFileId      string
	Description         string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

type InlineQueryResultCachedVideo struct {
	Type                string
	Id                  string
	VideoFileId         string
	Title               string
	Description         string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

type InlineQueryResultCachedVoice struct {
	Type                string
	Id                  string
	VoiceFileId         string
	Title               string
	Caption             string
	ParseMode           string
	VoiceDuration       int
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}

type InlineQueryResultCachedAudio struct {
	Type                string
	Id                  string
	AudioFileId         string
	Caption             string
	ParseMode           string
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent interface{}
}
