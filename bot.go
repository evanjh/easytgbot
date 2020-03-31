package easytgbot

import (
	json "github.com/json-iterator/go"
	"github.com/mylukin/easytgbot/types"
)

type Response struct {
	Ok          bool
	Result      json.RawMessage
	ErrorCode   int
	Description string
	Parameters  *types.ResponseParameters
}

func Hello() string {
	return "Hello, world."
}
