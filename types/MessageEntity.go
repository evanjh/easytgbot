package types

import (
	"errors"
	"net/url"

	"github.com/mylukin/easytgbot/configs"
)

// MessageEntity contains information about data in a Message.
type MessageEntity struct {
	Type     string
	Offset   int
	Length   int
	Url      string
	User     *User
	Language string
}

// ParseUrl attempts to parse a Url contained within a MessageEntity.
func (e MessageEntity) ParseUrl() (*url.URL, error) {
	if e.Url == "" {
		return nil, errors.New(configs.ErrBadURL)
	}

	return url.Parse(e.Url)
}

// IsMention returns true if the type of the message entity is "mention" (@username).
func (e MessageEntity) IsMention() bool {
	return e.Type == "mention"
}

// IsHashtag returns true if the type of the message entity is "hashtag".
func (e MessageEntity) IsHashtag() bool {
	return e.Type == "hashtag"
}

// IsCommand returns true if the type of the message entity is "bot_command".
func (e MessageEntity) IsCommand() bool {
	return e.Type == "bot_command"
}

// IsUrl returns true if the type of the message entity is "url".
func (e MessageEntity) IsUrl() bool {
	return e.Type == "url"
}

// IsEmail returns true if the type of the message entity is "email".
func (e MessageEntity) IsEmail() bool {
	return e.Type == "email"
}

// IsBold returns true if the type of the message entity is "bold" (bold text).
func (e MessageEntity) IsBold() bool {
	return e.Type == "bold"
}

// IsItalic returns true if the type of the message entity is "italic" (italic text).
func (e MessageEntity) IsItalic() bool {
	return e.Type == "italic"
}

// IsCode returns true if the type of the message entity is "code" (monowidth string).
func (e MessageEntity) IsCode() bool {
	return e.Type == "code"
}

// IsPre returns true if the type of the message entity is "pre" (monowidth block).
func (e MessageEntity) IsPre() bool {
	return e.Type == "pre"
}

// IsTextLink returns true if the type of the message entity is "text_link" (clickable text Url).
func (e MessageEntity) IsTextLink() bool {
	return e.Type == "text_link"
}
