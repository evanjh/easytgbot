package types

type Game struct {
	Title        string
	Description  string
	Photo        []PhotoSize
	Text         string
	TextEntities []MessageEntity
	Animation    Animation
}

type CallbackGame struct{}

type GameHighScore struct {
	Position int
	User     User
	Score    int
}
