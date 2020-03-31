package types

type InlineKeyboardButton struct {
	Text                         string
	Url                          *string
	LoginUrl                     *LoginUrl
	CallbackData                 *string
	SwitchInlineQuery            *string
	SwitchInlineQueryCurrentChat *string
	CallbackGame                 *CallbackGame
	Pay                          bool
}
