package types

type InputTextMessageContent struct {
	MessageText           string
	ParseMode             string
	DisableWebPagePreview bool
}

type InputLocationMessageContent struct {
	Latitude   float64
	Longitude  float64
	LivePeriod int
}

type InputVenueMessageContent struct {
	Latitude       float64
	Longitude      float64
	Title          string
	Address        string
	FoursquareId   string
	FoursquareType string
}
type InputContactMessageContent struct {
	PhoneNumber string
	FirstName   string
	LastName    string
	Vcard       string
}
