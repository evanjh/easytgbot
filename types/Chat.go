package types

type Chat struct {
	Id               int64  
	Type             string 
	Title            string 
	UserName         string 
	FirstName        string 
	LastName         string 
	Photo            *ChatPhoto 
	Description      string           
	InviteLink       string           
	PinnedMessage    *Message         
	Permissions      *ChatPermissions 
	SlowModeDelay    int              
	StickerSetName   string           
	CanSetStickerSet bool             
}

// IsPrivate returns if the Chat is a private conversation.
func (c Chat) IsPrivate() bool {
	return c.Type == "private"
}

// IsGroup returns if the Chat is a group.
func (c Chat) IsGroup() bool {
	return c.Type == "group"
}

// IsSuperGroup returns if the Chat is a supergroup.
func (c Chat) IsSuperGroup() bool {
	return c.Type == "supergroup"
}

// IsChannel returns if the Chat is a channel.
func (c Chat) IsChannel() bool {
	return c.Type == "channel"
}
