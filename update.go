package easytgbot

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// NewUpdate is create update instance
func NewUpdate(data string) *Update {
	return &Update{gjson.Parse(data)}
}

// GetType get message type
func (update *Update) GetType() string {
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
		"forward_date", // forward
	}
	message, err := update.Message()
	if err == nil {
		for _, key := range MessageSubTypes {
			if message.Get(key).Exists() {
				if key == "forward_date" {
					return "forward"
				}
				return key
			}
		}
	}
	return "unknown"
}

// Command get command
func (update *Update) Command() (string, string) {
	entities := update.Entities()
	for _, entity := range entities {
		etype := entity.Get("type").String()
		switch etype {
		case "bot_command":
			message, err := update.Message()
			if err == nil {
				text := message.Get("text").String()
				offset := entity.Get("offset").Int()
				length := offset + entity.Get("length").Int()
				command := text[offset:length]
				payload := strings.TrimSpace(text[length:])
				return command, payload
			}
		}
	}
	return "", ""
}

// Entities is Entities
func (update *Update) Entities() []Update {
	message, err := update.Message()
	if err == nil {
		if message.Get("entities").Exists() {
			return message.Get("entities").Array()
		}

		if message.Get("caption_entities").Exists() {
			return message.Get("caption_entities").Array()
		}
	}

	return []Update{}
}

// Message get message
func (update *Update) Message() (Update, error) {
	message := update.Get("message")
	if message.Exists() {
		return message, nil
	}

	editedMessage := update.Get("edited_message")
	if editedMessage.Exists() {
		return editedMessage, nil
	}

	channelPost := update.Get("channel_post")
	if channelPost.Exists() {
		return channelPost, nil
	}

	editedChannelPost := update.Get("edited_channel_post")
	if editedChannelPost.Exists() {
		return editedChannelPost, nil
	}

	callbackQuery := update.Get("callback_query")
	if callbackQuery.Exists() {
		message := callbackQuery.Get("message")
		if message.Exists() {
			return message, nil
		}
	}

	return Update{}, fmt.Errorf("chat is not found")
}

// Chat get update
func (update *Update) Chat() (Update, error) {
	message, err := update.Message()
	if err != nil {
		return Update{}, fmt.Errorf("chat is not found")
	}
	return message.Get("chat"), nil
}

// From get update
func (update *Update) From() (Update, error) {

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

	message, err := update.Message()
	if err == nil {
		return message.Get("from"), nil
	}

	return Update{}, fmt.Errorf("from is not found")
}

// Reply get telegram bot
func (update *Update) Reply(text string, extra JSONBody) JSONBody {
	message := update.Get("message")
	messageID := message.Get("message_id").Int()
	chat, _ := update.Chat()
	chatID := chat.Get("id").Int()
	return mergeJSON(JSONBody{
		"method":              "sendMessage",
		"chat_id":             chatID,
		"reply_to_message_id": messageID,
		"text":                text,
	}, extra)
}

// mergeJSON merge json body
func mergeJSON(map1 JSONBody, map2 JSONBody) JSONBody {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}
