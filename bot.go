package easytgbot

import (
	"fmt"
	"log"
	"net/url"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

// BotAPI allows you to interact with the Telegram Bot API.
type BotAPI struct {
	Token string `json:"token"`
	Debug bool   `json:"debug"`

	Self   JSON     `json:"-"`
	Client *req.Req `json:"-"`

	apiEndpoint string
}

// JSON is a response from the Telegram API with the result stored raw.
type JSON struct {
	gjson.Result
}

// Get searches result for the specified path.
// The result should be a JSON array or object.
func (t JSON) Get(path string) JSON {
	return JSON{gjson.Get(t.Raw, path)}
}

// Error is an error containing extra information returned by the Telegram API.
type Error struct {
	Code       int64
	Message    string
	Parameters JSON
}

func (e Error) Error() string {
	return e.Message
}

// New bot instance
func New(token string) (*BotAPI, error) {
	return NewBotAPIWith(token, APIEndpoint)
}

// NewBotAPIWith creates a new BotAPI instance and allows you to pass API endpoint.
func NewBotAPIWith(token string, apiEndpoint string) (*BotAPI, error) {
	client := req.New()
	bot := &BotAPI{
		Token:       token,
		Client:      client,
		apiEndpoint: apiEndpoint,
	}

	self, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Self = self

	return bot, nil
}

// MakeRequest makes a request to a specific endpoint with our token.
func (bot *BotAPI) MakeRequest(endpoint string, params url.Values) (JSON, error) {
	method := fmt.Sprintf(bot.apiEndpoint, bot.Token, endpoint)
	resp, err := bot.Client.Get(method)
	if err != nil {
		return JSON{}, err
	}
	data, err := resp.ToString()
	fmt.Printf("data: %+v\n", data)
	if err != nil {
		log.Fatal(err)
	}

	apiResp := JSON{gjson.Parse(data)}
	ok := apiResp.Get("ok").Bool()
	if !ok {
		// error
		return apiResp, &Error{
			Code:       apiResp.Get("error_code").Int(),
			Message:    apiResp.Get("description").String(),
			Parameters: apiResp.Get("parameters"),
		}
	}

	result := apiResp.Get("result")

	return result, nil

}

// GetMe fetches the currently authenticated bot.
//
// This method is called upon creation to validate the token,
// and so you may get this data from BotAPI.Self without the need for
// another request.
func (bot *BotAPI) GetMe() (JSON, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return JSON{}, err
	}
	return resp, nil
}
