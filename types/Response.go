package types

import json "github.com/json-iterator/go"

type Response struct {
	Ok          bool
	Result      json.RawMessage
	ErrorCode   int
	Description string
	Parameters  *ResponseParameters
}
