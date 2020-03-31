package types

type PollAnswer struct {
	PollId    string
	User      *User
	OptionIds []int
}
