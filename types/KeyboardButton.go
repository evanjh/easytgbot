package types

type KeyboardButton struct {
	Text            string
	RequestContact  bool
	RequestLocation bool
	RequestPoll     *KeyboardButtonPollType
}
