package easytgbot

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// NewUpdate is create update instance
func NewUpdate(data string) *Update {
	return &Update{gjson.Parse(data)}
}

// GetType get message type 
func (update *Update) GetType() string  {
	MessageSubTypes := []string{
		"voice",
		"video_note",
		"video",
		"animation",
		"venue",
		"text",
		"supergroup_chat_created",
		"successful_payment",
		"sticker",
		"pinned_message",
		"photo",
		"new_chat_title",
		"new_chat_photo",
		"new_chat_members",
		"migrate_to_chat_id",
		"migrate_from_chat_id",
		"location",
		"left_chat_member",
		"invoice",
		"group_chat_created",
		"game",
		"document",
		"delete_chat_photo",
		"contact",
		"channel_chat_created",
		"audio",
		"connected_website",
		"passport_data",
		"poll",
		"forward_date",
	}
	if update.Get("message").Exists() || update.Get("channel_post").Exists() {
		for _, key := range MessageSubTypes {
			if update.Get(key).Exists()	{
				if key == "forward_date" {
					return "forward"
				} else {
					return key
				}
				
			}
		}
	}
	return "unknown"
}

// Chat get update
func (update *Update) Chat() (Update, error) {
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

	return Update{}, fmt.Errorf("chat is not found")
}

// From get update
func (update *Update) From() (Update, error) {
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

	return Update{}, fmt.Errorf("from is not found")
}
