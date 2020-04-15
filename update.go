package easytgbot

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// Update is message
type Update struct {
	JSON
}

// NewUpdate is create update instance
func NewUpdate(data string) *Update {
	return &Update{JSON{gjson.Parse(data)}}
}

// Chat get update
func (update *Update) Chat() (JSON, error) {
	message := update.Get("message")
	if message.Exists() {
		return message.Get("chat"), nil
	}

	editedMessage := update.Get("edited_message")
	if editedMessage.Exists() {
		return editedMessage.Get("chat"), nil
	}

	channelPost := update.Get("channel_post")
	if channelPost.Exists() {
		return channelPost.Get("chat"), nil
	}

	editedChannelPost := update.Get("edited_channel_post")
	if editedChannelPost.Exists() {
		return editedChannelPost.Get("chat"), nil
	}

	callbackQuery := update.Get("callback_query")
	if callbackQuery.Exists() {
		message := callbackQuery.Get("message")
		if message.Exists() {
			return message.Get("chat"), nil
		}
	}

	return JSON{}, fmt.Errorf("chat is not found")
}

// From get update
func (update *Update) From() (JSON, error) {
	message := update.Get("message")
	if message.Exists() {
		return message.Get("from"), nil
	}

	editedMessage := update.Get("edited_message")
	if editedMessage.Exists() {
		return editedMessage.Get("from"), nil
	}

	channelPost := update.Get("channel_post")
	if channelPost.Exists() {
		return channelPost.Get("from"), nil
	}

	editedChannelPost := update.Get("edited_channel_post")
	if editedChannelPost.Exists() {
		return editedChannelPost.Get("from"), nil
	}

	callbackQuery := update.Get("callback_query")
	if callbackQuery.Exists() {
		return callbackQuery.Get("from"), nil
	}

	inlineQuery := update.Get("inline_query")
	if inlineQuery.Exists() {
		return inlineQuery.Get("from"), nil
	}

	shippingQuery := update.Get("shipping_query")
	if shippingQuery.Exists() {
		return shippingQuery.Get("from"), nil
	}

	preCheckoutQuery := update.Get("pre_checkout_query")
	if preCheckoutQuery.Exists() {
		return preCheckoutQuery.Get("from"), nil
	}

	chosenInlineResult := update.Get("chosen_inline_result")
	if chosenInlineResult.Exists() {
		return chosenInlineResult.Get("from"), nil
	}

	return JSON{}, fmt.Errorf("from is not found")
}
