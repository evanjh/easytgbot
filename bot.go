package easytgbot

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
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
)

// JSONBody is send message
type JSONBody map[string]interface{}

// BotOption type for additional Server options
type BotOption func(*BotAPI)

// BotAPI allows you to interact with the Telegram Bot API.
type BotAPI struct {
	Debug   bool
	Token   string
	Webhook string
	Buffer  int

	Self            JSON
	Client          *req.Req
	shutdownChannel chan interface{}
	Timeout         time.Duration

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

// Value returns one of these types:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	Number, for JSON numbers
//	string, for JSON string literals
//	nil, for JSON null
//	map[string]interface{}, for JSON objects
//	[]interface{}, for JSON arrays
//
func (t JSON) Value() interface{} {
	return t.Value()
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
func New(token string, options ...BotOption) (*BotAPI, error) {
	if len(token) == 0 {
		return &BotAPI{}, fmt.Errorf("token is empty")
	}
	client := req.New()
	bot := &BotAPI{
		Token:       token,
		Client:      client,
		Buffer:      100,
		Timeout:     10,
		apiEndpoint: APIEndpoint,
	}
	for _, optFunc := range options {
		optFunc(bot)
	}

	self, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Self = self

	return bot, nil
}

// WithDebug set debug mode
func WithDebug(isDebug bool) BotOption {
	return func(bot *BotAPI) {
		bot.Debug = isDebug
	}
}

// WithEndpoint set base api
func WithEndpoint(endpoint string, timeout time.Duration) BotOption {
	return func(bot *BotAPI) {
		bot.apiEndpoint = endpoint
		bot.Timeout = timeout
	}
}

// WithWebhook returns BotOption for given Webhook URL and Server address to listen.
func WithWebhook(url string) BotOption {
	return func(bot *BotAPI) {
		bot.Webhook = url
	}
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
	} else {
		log.Printf("%-v", resp)
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

// GetUpdates starts and returns a channel for getting updates.
func (bot *BotAPI) GetUpdates(params JSONBody) (chan JSON, error) {
	if bot.Webhook != "" {
		return bot.listenUpdates()
	}

	// first delete webbook
	bot.DeleteWebhook()

	ch := make(chan JSON, bot.Buffer)
	offset, _ := strconv.ParseInt(strconv.Itoa(params["offset"].(int)), 10, 64)

	go func() {
		for {
			select {
			case <-bot.shutdownChannel:
				close(ch)
				return
			default:
			}

			resp, err := bot.MakeRequest("getUpdates", params)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			updates := resp.Array()
			for _, update := range updates {
				if update.Get("update_id").Int() >= offset {
					params["offset"] = update.Get("update_id").Int() + 1
					ch <- update
				}
			}
		}
	}()

	return ch, nil
}

// listenUpdates
func (bot *BotAPI) listenUpdates() (chan JSON, error) {
	updates := make(chan JSON)
	defer func() {
		updates <- JSON{}
	}()
	_, err := bot.SetWebhook(JSONBody{
		"url":             bot.Webhook,
		"max_connections": 100,
	})

	if err != nil {
		return updates, err
	}
	return updates, nil
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
