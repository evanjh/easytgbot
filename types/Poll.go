package types

type Poll struct {
	Id                    string
	question              string
	Options               []PollOption
	total_voter_count     int
	IsClosed              bool
	IsAnonymous           bool
	Type                  string
	AllowsMultipleAnswers bool
	CorrectOptionId       int
}
