package types

// Update is an update response, from GetUpdates.
type Update struct {
	UpdateId           int // required
	Message            *Message
	EditedMessage      *Message
	ChannelPost        *Message
	EditedChannelPost  *Message
	InlineQuery        *InlineQuery
	ChosenInlineResult *ChosenInlineResult
	CallbackQuery      *CallbackQuery
	ShippingQuery      *ShippingQuery
	PreCheckoutQuery   *PreCheckoutQuery
	Poll               *Poll
	PollAnswer         *PollAnswer
}

// UpdatesChannel is the channel for getting updates.
type UpdatesChannel <-chan Update

// Clear discards all unprocessed incoming updates.
func (ch UpdatesChannel) Clear() {
	for len(ch) != 0 {
		<-ch
	}
}
