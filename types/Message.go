package types

import (
	"strings"
	"time"
)

// Message is returned by almost every request, and contains data about
// almost anything.
type Message struct {
	MessageId             int
	From                  *User
	Date                  int
	Chat                  *Chat
	ForwardFrom           *User              
	ForwardFromChat       *Chat              
	ForwardFromMessageId  int                
	ForwardSignature      string             
	ForwardSenderName     string             
	ForwardDate           int                
	ReplyToMessage        *Message           
	EditDate              int                
	MediaGroupId          string             
	AuthorSignature       string             
	Text                  string             
	Entities              *[]MessageEntity   
	CaptionEntities       *[]MessageEntity   
	Audio                 *Audio             
	Document              *Document          
	Animation             *Animation         
	Game                  *Game              
	Photo                 *[]PhotoSize       
	Sticker               *Sticker           
	Video                 *Video             
	Voice                 *Voice             
	VideoNote             *VideoNote         
	Caption               string             
	Contact               *Contact           
	Location              *Location          
	Venue                 *Venue             
	Poll                  *Poll              
	NewChatMembers        *[]User            
	LeftChatMember        *User              
	NewChatTitle          string             
	NewChatPhoto          *[]PhotoSize       
	DeleteChatPhoto       bool               
	GroupChatCreated      bool               
	SuperGroupChatCreated bool               
	ChannelChatCreated    bool               
	MigrateToChatId       int64              
	MigrateFromChatId     int64              
	PinnedMessage         *Message           
	Invoice               *Invoice           
	SuccessfulPayment     *SuccessfulPayment 
	ConnectedWebsite      string             
	PassportData          *PassportData      
	ReplyMarkup           *InlineKeyboardMarkup
}

// Time converts the message timestamp into a Time.
func (m *Message) Time() time.Time {
	return time.Unix(int64(m.Date), 0)
}

// IsCommand returns true if message starts with a "bot_command" entity.
func (m *Message) IsCommand() bool {
	if m.Entities == nil || len(*m.Entities) == 0 {
		return false
	}

	entity := (*m.Entities)[0]
	return entity.Offset == 0 && entity.IsCommand()
}

// Command checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is removed. Use
// CommandWithAt() if you do not want that.
func (m *Message) Command() string {
	command := m.CommandWithAt()

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

// CommandWithAt checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is not removed. Use Command()
// if you want that.
func (m *Message) CommandWithAt() string {
	if !m.IsCommand() {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := (*m.Entities)[0]
	return m.Text[1:entity.Length]
}

// CommandArguments checks if the message was a command and if it was,
// returns all text after the command name. If the Message was not a
// command, it returns an empty string.
//
// Note: The first character after the command name is omitted:
// - "/foo bar baz" yields "bar baz", not " bar baz"
// - "/foo-bar baz" yields "bar baz", too
// Even though the latter is not a command conforming to the spec, the API
// marks "/foo" as command entity.
func (m *Message) CommandArguments() string {
	if !m.IsCommand() {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := (*m.Entities)[0]
	if len(m.Text) == entity.Length {
		return "" // The command makes up the whole message
	}

	return m.Text[entity.Length+1:]
}
