package types

type InputMediaPhoto struct {
	Type              string
	Media             string
	Thumb             string
	Caption           string
	ParseMode         string
	Width             int
	Height            int
	Duration          int
	SupportsStreaming bool
}

type InputMediaVideo struct {
	Type      string
	Media     string
	Thumb     string
	Caption   string
	ParseMode string
	Width     int
	Height    int
	Duration  int
}

type InputMediaAudio struct {
	Type      string
	Media     string
	Thumb     string
	Caption   string
	ParseMode string
	Duration  int
	Performer string
	Title     string
}

type InputMediaDocument struct {
	Type      string
	Media     string
	Thumb     string
	Caption   string
	ParseMode string
}
