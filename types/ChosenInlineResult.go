package types

type ChosenInlineResult struct {
	ResultId        string
	From            *User
	Location        *Location
	InlineMessageId string
	Query           string
}
