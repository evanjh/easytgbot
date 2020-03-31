package types

// ChatMember is information about a member in a chat.
type ChatMember struct {
	User                  *User
	Status                string
	CustomTitle           string
	UntilDate             int64
	CanBeEdited           bool
	CanPostMessages       bool
	CanEditMessages       bool
	CanDeleteMessages     bool
	CanRestrictMembers    bool
	CanPromoteMembers     bool
	CanChangeInfo         bool
	CanInviteUsers        bool
	CanPinMessages        bool
	CanSendMessages       bool
	CanSendMediaMessages  bool
	CanSendPolls          bool
	CanSendOtherMessages  bool
	CanAddWebPagePreviews bool
}

// IsCreator returns if the ChatMember was the creator of the chat.
func (user ChatMember) IsCreator() bool { return user.Status == "creator" }

// IsAdministrator returns if the ChatMember is a chat administrator.
func (user ChatMember) IsAdministrator() bool { return user.Status == "administrator" }

// IsMember returns if the ChatMember is a current member of the chat.
func (user ChatMember) IsMember() bool { return user.Status == "member" }

// HasLeft returns if the ChatMember left the chat.
func (user ChatMember) HasLeft() bool { return user.Status == "left" }

// WasKicked returns if the ChatMember was kicked from the chat.
func (user ChatMember) WasKicked() bool { return user.Status == "kicked" }
