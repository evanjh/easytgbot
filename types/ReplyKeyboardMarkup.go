package types

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton
	ResizeKeyboard  bool
	OneTimeKeyboard bool
	Selective       bool
}
