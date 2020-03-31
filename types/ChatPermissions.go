package types

type ChatPermissions struct {
	CanSendMessages       bool
	CanSendMediaMessages  bool
	CanSendPolls          bool
	CanSendOtherMessages  bool
	CanAddWebPagePreviews bool
	CanChangeInfo         bool
	CanInviteUsers        bool
	CanPinMessages        bool
}
