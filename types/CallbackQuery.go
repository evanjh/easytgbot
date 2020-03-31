package types

type CallbackQuery struct {
	Id              string
	From            *User
	Message         *Message
	InlineMessageId string
	ChatInstance    string
	Data            string
	GameShortName   string
}
