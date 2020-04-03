package easytgbot

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// Telegram constants
const (
	// UserAgent is http user-agent header
	UserAgent = "EasyTGBot/1.0.0"
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

// JSONBody is send message
type JSONBody map[string]interface{}

// BotAPI allows you to interact with the Telegram Bot API.
type BotAPI struct {
	Token string
	Debug bool

	Self    JSON
	Client  *req.Req
	Timeout time.Duration

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

// Array returns back an array of values.
// If the result represents a non-existent value, then an empty array will be
// returned. If the result is not a JSON array, the return value will be an
// array containing one result.
func (t JSON) Array() []JSON {
	res := []JSON{}
	if t.IsArray() {
		t.ForEach(func(key, value gjson.Result) bool {
			res = append(res, JSON{value})
			return true // keep iterating
		})
	}
	return res
}

// Map returns back an map of values. The result should be a JSON array.
func (t JSON) Map() map[string]JSON {
	res := map[string]JSON{}
	if t.IsObject() {
		t.ForEach(func(key, value gjson.Result) bool {
			res[key.String()] = JSON{value}
			return true // keep iterating
		})
	}
	return res
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
		Timeout:     10,
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
func (bot *BotAPI) MakeRequest(endpoint string, params JSONBody) (JSON, error) {
	method := fmt.Sprintf(bot.apiEndpoint, bot.Token, endpoint)
	var jsonBody JSONBody
	if params == nil {
		jsonBody = JSONBody{}
	} else {
		jsonBody = params
	}

	// set timeout
	bot.Client.SetTimeout(bot.Timeout * time.Second)
	// post data
	var (
		resp *req.Resp
		err  error
	)

	header := req.Header{
		"User-Agent": UserAgent,
	}

	if endpoint == "setWebhook" {
		var fileUploads []req.FileUpload
		fromParams := url.Values{}
		for key, value := range params {
			switch value.(type) {
			case string, int:
				fromParams.Add(key, fmt.Sprintf("%v", value))
			case []string:
				fromParams.Add(key, fmt.Sprintf("[\"%v\"]", strings.Join(value.([]string), "\",\"")))
			case req.FileUpload:
				fileUploads = append(fileUploads, value.(req.FileUpload))
			}
		}
		resp, err = bot.Client.Post(method, header, fromParams, fileUploads)
	} else {
		resp, err = bot.Client.Post(method, header, req.BodyJSON(&jsonBody))
	}

	if err != nil {
		return JSON{}, err
	}
	if bot.Debug {
		log.Printf("%+v", resp)
	}
	data, _ := resp.ToString()
	apiJSON := JSON{gjson.Parse(data)}
	ok := apiJSON.Get("ok").Bool()
	if !ok {
		// error
		return apiJSON, &Error{
			Code:       apiJSON.Get("error_code").Int(),
			Message:    apiJSON.Get("description").String(),
			Parameters: apiJSON.Get("parameters"),
		}
	}

	result := apiJSON.Get("result")

	return result, nil

}

// GetMe fetches the currently authenticated bot.
//
// This method is called upon creation to validate the token,
// and so you may get this data from BotAPI.Self without the need for
// another request.
func (bot *BotAPI) GetMe() (JSON, error) {
	return bot.MakeRequest("getMe", nil)
}

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
func (bot *BotAPI) GetUpdates(params JSONBody) ([]JSON, error) {
	resp, err := bot.MakeRequest("getUpdates", params)
	if err != nil {
		return []JSON{}, err
	}
	return resp.Array(), nil
}

// GetWebhookInfo allows you to fetch information about a webhook and if
// one currently is set, along with pending update count and error messages.
func (bot *BotAPI) GetWebhookInfo() (JSON, error) {
	return bot.MakeRequest("getWebhookInfo", nil)
}

// SetWebhook sets a webhook.
//
// If this is set, GetUpdates will not get any data!
//
// If you do not have a legitimate TLS certificate, you need to include
// your self signed certificate with the config.
func (bot *BotAPI) SetWebhook(params JSONBody) (JSON, error) {
	return bot.MakeRequest("setWebhook", params)
}

// DeleteWebhook unsets the webhook.
func (bot *BotAPI) DeleteWebhook() (JSON, error) {
	return bot.MakeRequest("deleteWebhook", nil)
}
