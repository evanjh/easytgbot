package easytgbot

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/tidwall/gjson"
)

// MessageNodes is message node name
var MessageNodes = []string{"message", "edited_message", "channel_post", "edited_channel_post", "my_chat_member", "chat_member"}

// MessageQueryNodes is message query node name
var MessageQueryNodes = []string{"callback_query", "inline_query", "shipping_query", "pre_checkout_query", "chosen_inline_result"}

// NewUpdate is create update instance
func NewUpdate(data string) Update {
	return Update{gjson.Parse(data)}
}

// GetType get message type
func (update Update) GetType() string {
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
		"my_chat_member",
		"chat_member",
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
	}
	message, err := update.Message()
	if err == nil {
		for _, key := range MessageSubTypes {
			if message.Get(key).Exists() {
				return key
			}
		}
	}
	for _, key := range MessageSubTypes {
		if update.Get(key).Exists() {
			return key
		}
	}
	return "unknown"
}

// Command get command
func (update Update) Command() (string, string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("bot command: %s\n", err)
		}
	}()
	message, err := update.Message()
	if err != nil {
		return "", ""
	}
	text := message.Get("text").String()

	var res []string
	bText := []byte(text)
	for len(bText) > 0 {
		r, size := utf8.DecodeRune(bText)
		if size == 4 {
			res = append(res, fmt.Sprintf("%c", r))
			res = append(res, "")
		} else {
			res = append(res, fmt.Sprintf("%c", r))
		}
		bText = bText[size:]
	}

	entities := update.Entities()
	for _, entity := range entities {
		etype := entity.Get("type").String()
		switch etype {
		case "bot_command":
			offset := entity.Get("offset").Int()
			if offset == 0 {
				length := offset + entity.Get("length").Int()
				command := strings.Join(res[offset:length], "")
				payload := strings.TrimSpace(strings.Join(res[length:], ""))
				return command, payload
			}
		}
	}
	return "", ""
}

// Entities is Entities
func (update Update) Entities() []Update {
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
func (update Update) Message() (Update, error) {
	for _, t := range MessageNodes {
		message := update.Get(t)
		if message.Exists() {
			return message, nil
		}
	}

	for _, q := range MessageQueryNodes {
		callbackQuery := update.Get(q)
		if callbackQuery.Exists() {
			for _, t := range MessageNodes {
				message := callbackQuery.Get(t)
				if message.Exists() {
					return message, nil
				}
			}
		}
	}

	return Update{}, fmt.Errorf("message is not found")
}

// Chat get update
func (update Update) Chat() (Update, error) {
	message, err := update.Message()
	if err != nil {
		return Update{}, fmt.Errorf("chat is not found")
	}
	return message.Get("chat"), nil
}

// From get update
func (update Update) From() (Update, error) {

	for _, t := range MessageQueryNodes {
		callbackQuery := update.Get(t)
		if callbackQuery.Exists() {
			return callbackQuery.Get("from"), nil
		}
	}

	message, err := update.Message()
	if err == nil {
		return message.Get("from"), nil
	}

	return Update{}, fmt.Errorf("from is not found")
}

// SendMessage is send message
func (update Update) SendMessage(text string, extra JSONBody) JSONBody {
	chat, _ := update.Chat()
	chatID := chat.Get("id").Int()

	result := JSONBody{
		"method":  "sendMessage",
		"chat_id": chatID,
		"text":    text,
	}
	// reply
	if _, ok := extra["reply"]; ok {
		message, _ := update.Message()
		messageID := message.Get("message_id").Int()
		result["reply_to_message_id"] = messageID
	}
	result = mergeJSON(result, extra)
	return result
}

// Reply reply message
func (update Update) Reply(text string, extra JSONBody) JSONBody {
	message, _ := update.Message()
	messageID := message.Get("message_id").Int()
	chat, _ := update.Chat()
	chatID := chat.Get("id").Int()
	callbackQuery := update.Get("callback_query")

	result := JSONBody{
		"chat_id": chatID,
		"text":    text,
	}

	// callback
	if callbackQuery.Exists() {
		result["method"] = "editMessageText"
		result["message_id"] = messageID
	} else {
		result["method"] = "sendMessage"
		result["reply_to_message_id"] = messageID
	}
	result = mergeJSON(result, extra)
	return result
}

// EditMessageText edit message
func (update Update) EditMessageText(text string, extra JSONBody) JSONBody {
	message, _ := update.Message()
	messageID := message.Get("message_id").Int()
	chat, _ := update.Chat()
	chatID := chat.Get("id").Int()
	return mergeJSON(JSONBody{
		"method":     "editMessageText",
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	}, extra)
}

// EditMessageReplyMarkup edit message
func (update Update) EditMessageReplyMarkup(extra JSONBody) JSONBody {
	message, _ := update.Message()
	messageID := message.Get("message_id").Int()
	chat, _ := update.Chat()
	chatID := chat.Get("id").Int()
	return mergeJSON(JSONBody{
		"method":     "editMessageReplyMarkup",
		"chat_id":    chatID,
		"message_id": messageID,
	}, extra)
}

// AnswerCallbackQuery is AnswerCallbackQuery
func (update Update) AnswerCallbackQuery(text string, extra JSONBody) JSONBody {
	callbackQuery := update.Get("callback_query")
	queryID := callbackQuery.Get("id").String()
	return mergeJSON(JSONBody{
		"method":            "answerCallbackQuery",
		"callback_query_id": queryID,
		"show_alert":        true,
		"text":              text,
	}, extra)
}

// DeleteMessage see: https://core.telegram.org/bots/api#deletemessage
func (update Update) DeleteMessage() JSONBody {
	message, _ := update.Message()
	messageID := message.Get("message_id").Int()
	chat, _ := update.Chat()
	chatID := chat.Get("id").Int()
	return JSONBody{
		"method":     "deleteMessage",
		"chat_id":    chatID,
		"message_id": messageID,
	}
}

// SendMediaGroup is send media group
func (update Update) SendMediaGroup(media []JSONBody, extra JSONBody) JSONBody {
	chat, _ := update.Chat()
	chatID := chat.Get("id").Int()

	result := JSONBody{
		"method":  "sendMediaGroup",
		"chat_id": chatID,
		"media":   media,
	}
	// reply
	if _, ok := extra["reply"]; ok {
		message, _ := update.Message()
		messageID := message.Get("message_id").Int()
		result["reply_to_message_id"] = messageID
	}
	result = mergeJSON(result, extra)
	return result
}

// mergeJSON merge json body
func mergeJSON(map1 JSONBody, map2 JSONBody) JSONBody {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}
