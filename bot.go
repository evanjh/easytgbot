package easytgbot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
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
	// Endpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	Endpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

// JSONBody is send message
type JSONBody map[string]interface{}

func (body JSONBody) Error() string {
	data, err := json.Marshal(body)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// MiddlewareFunc defines a function to process middleware.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(interface{}, *Bot, *Update) JSONBody

// Bot allows you to interact with the Telegram Bot API.
type Bot struct {
	Debug   bool
	ID      int64
	Name    string
	Token   string
	Webhook string
	Buffer  int
	Timeout time.Duration
	Self    Update

	handlers        map[string]interface{}
	client          *req.Req
	shutdownChannel chan interface{}
	apiEndpoint     string
	middleware      []MiddlewareFunc
}

// Settings represents a utility struct for passing certain
// properties of a bot around and is required to make bots.
type Settings struct {
	// debug
	Debug bool // default: false
	// Telegram API Url
	Endpoint string

	// Webhook
	Webhook string

	// Telegram token
	Token string

	// Updates channel capacity
	Updates int // Default: 100

	// Timeout
	Timeout time.Duration // Default: 10s

	Proxy string

	GetMe bool
}

// Update is a response from the Telegram API with the result stored raw.
type Update struct {
	gjson.Result
}

// Get searches result for the specified path.
// The result should be a Update array or object.
func (t Update) Get(path string) Update {
	return Update{gjson.Get(t.Raw, path)}
}

// Array returns back an array of values.
// If the result represents a non-existent value, then an empty array will be
// returned. If the result is not a Update array, the return value will be an
// array containing one result.
func (t Update) Array() []Update {
	res := []Update{}
	if t.IsArray() {
		t.ForEach(func(key, value gjson.Result) bool {
			res = append(res, Update{value})
			return true // keep iterating
		})
	}
	return res
}

// Map returns back an map of values. The result should be a Update array.
func (t Update) Map() map[string]Update {
	res := map[string]Update{}
	if t.IsObject() {
		t.ForEach(func(key, value gjson.Result) bool {
			res[key.String()] = Update{value}
			return true // keep iterating
		})
	}
	return res
}

// Error is an error containing extra information returned by the Telegram API.
type Error struct {
	Code       int64
	Message    string
	Parameters Update
}

func (e Error) Error() string {
	return e.Message
}

// New bot instance
func New(token string, opts Settings) (*Bot, error) {
	if len(token) == 0 {
		return &Bot{}, fmt.Errorf("token is empty")
	}

	if opts.Updates == 0 {
		opts.Updates = 100
	}

	if opts.Timeout == 0 {
		opts.Timeout = 10 * time.Second
	}

	if opts.Endpoint == "" {
		opts.Endpoint = Endpoint
	}

	client := req.New()
	// set proxy
	if opts.Proxy != "" {
		client.SetProxyUrl(opts.Proxy)
	}

	bot := &Bot{
		Debug:   opts.Debug == true,
		Token:   token,
		Webhook: opts.Webhook,

		Buffer:  opts.Updates,
		Timeout: opts.Timeout,

		client:      client,
		apiEndpoint: opts.Endpoint,
		handlers:    make(map[string]interface{}),
	}

	if opts.GetMe {
		self, err := bot.GetMe()
		if err != nil {
			return nil, err
		}

		bot.Self = self
	}

	return bot, nil
}

// MakeRequest makes a request to a specific endpoint with our token.
func (bot *Bot) MakeRequest(endpoint string, params JSONBody) (Update, error) {
	method := fmt.Sprintf(bot.apiEndpoint, bot.Token, endpoint)
	var jsonBody JSONBody
	if params == nil {
		jsonBody = JSONBody{}
	} else {
		jsonBody = params
	}

	// set timeout
	bot.client.SetTimeout(bot.Timeout)
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
		resp, err = bot.client.Post(method, header, fromParams, fileUploads)
	} else {
		resp, err = bot.client.Post(method, header, req.BodyJSON(&jsonBody))
	}

	if err != nil {
		return Update{}, err
	}
	if bot.Debug {
		log.Printf("method: %s, resp: %-v", method, resp)
	}
	data, _ := resp.ToString()
	apiJSON := Update{gjson.Parse(data)}
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
// and so you may get this data from Bot.Self without the need for
// another request.
func (bot *Bot) GetMe() (Update, error) {
	return bot.MakeRequest("getMe", nil)
}

// SetBotID
func (bot *Bot) SetBotID(botID int64) {
	bot.ID = botID
}

// SetBotID
func (bot *Bot) SetBotName(name string) {
	bot.Name = name
}

// GetUpdates starts and returns a channel for getting updates.
func (bot *Bot) GetUpdates(params JSONBody) (chan Update, error) {
	if bot.Webhook != "" {
		return bot.listenUpdates()
	}

	// first delete webbook
	bot.DeleteWebhook()

	ch := make(chan Update, bot.Buffer)
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
func (bot *Bot) listenUpdates() (chan Update, error) {
	updates := make(chan Update)
	defer func() {
		updates <- Update{}
	}()
	_, err := bot.SetWebhook(JSONBody{
		"url":             bot.Webhook,
		"max_connections": bot.Buffer,
	})

	if err != nil {
		return updates, err
	}
	return updates, nil
}

// GetWebhookInfo allows you to fetch information about a webhook and if
// one currently is set, along with pending update count and error messages.
func (bot *Bot) GetWebhookInfo() (Update, error) {
	return bot.MakeRequest("getWebhookInfo", nil)
}

// SetWebhook sets a webhook.
//
// If this is set, GetUpdates will not get any data!
//
// If you do not have a legitimate TLS certificate, you need to include
// your self signed certificate with the config.
func (bot *Bot) SetWebhook(params JSONBody) (Update, error) {
	info, err := bot.GetWebhookInfo()
	if err != nil {
		return Update{}, err
	}
	url := info.Get("url").String()
	if url == params["url"] {
		return Update{}, err
	}
	return bot.MakeRequest("setWebhook", params)
}

// DeleteWebhook unsets the webhook.
func (bot *Bot) DeleteWebhook() (Update, error) {
	return bot.MakeRequest("deleteWebhook", nil)
}

// SendMessage send message
func (bot *Bot) SendMessage(chatID int64, text string, extra JSONBody) (Update, error) {
	return bot.MakeRequest("sendMessage", mergeJSON(JSONBody{
		"chat_id": chatID,
		"text":    text,
	}, extra))
}

// EditMessageText see https://core.telegram.org/bots/api#editmessagetext
func (bot *Bot) EditMessageText(chatID int64, messageID int64, text string, extra JSONBody) (Update, error) {
	return bot.MakeRequest("editMessageText", mergeJSON(JSONBody{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	}, extra))
}

// PinChatMessage see https://core.telegram.org/bots/api#pinchatmessage
func (bot *Bot) PinChatMessage(chatID int64, messageID int64, extra JSONBody) (Update, error) {
	return bot.MakeRequest("pinChatMessage", mergeJSON(JSONBody{
		"chat_id":    chatID,
		"message_id": messageID,
	}, extra))
}

// UnpinAllChatMessages see https://core.telegram.org/bots/api#unpinallchatmessages
func (bot *Bot) UnpinAllChatMessages(chatID int64) (Update, error) {
	return bot.MakeRequest("unpinAllChatMessages", mergeJSON(JSONBody{
		"chat_id": chatID,
	}, nil))
}

// SendPhoto send message
func (bot *Bot) SendPhoto(chatID int64, fileID string, extra JSONBody) (Update, error) {
	return bot.MakeRequest("sendPhoto", mergeJSON(JSONBody{
		"chat_id": chatID,
		"photo":   fileID,
	}, extra))
}

// SendVideo send message
func (bot *Bot) SendVideo(chatID int64, fileID string, extra JSONBody) (Update, error) {
	return bot.MakeRequest("sendVideo", mergeJSON(JSONBody{
		"chat_id": chatID,
		"video":   fileID,
	}, extra))
}

// ForwardMessage send message
func (bot *Bot) ForwardMessage(chatID int64, fromChatID int64, messageID int64, extra JSONBody) (Update, error) {
	return bot.MakeRequest("forwardMessage", mergeJSON(JSONBody{
		"chat_id":      chatID,
		"from_chat_id": fromChatID,
		"message_id":   messageID,
	}, extra))
}

// AnswerCallbackQuery see https://core.telegram.org/bots/api#answercallbackquery
func (bot *Bot) AnswerCallbackQuery(queryID string, text string, extra JSONBody) (Update, error) {
	return bot.MakeRequest("answerCallbackQuery", mergeJSON(JSONBody{
		"callback_query_id": queryID,
		"text":              text,
		"show_alert":        true,
	}, extra))
}

// SendMediaGroup see https://core.telegram.org/bots/api#sendmediagroup
func (bot *Bot) SendMediaGroup(chatID int64, media []JSONBody, extra JSONBody) (Update, error) {
	return bot.MakeRequest("sendMediaGroup", mergeJSON(JSONBody{
		"chat_id": chatID,
		"media":   media,
	}, extra))
}

// GetFile see https://core.telegram.org/bots/api#getfile
func (bot *Bot) GetFile(fileID string) (string, error) {
	res, err := bot.MakeRequest("getFile", JSONBody{
		"file_id": fileID,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(FileEndpoint, bot.Token, res.Get("file_path").String()), nil
}

// GetChat see https://core.telegram.org/bots/api#getchat
func (bot *Bot) GetChat(param interface{}) (Update, error) {
	params := JSONBody{}
	switch chatID := param.(type) {
	case string:
		params["chat_id"] = chatID
	default:
		params["chat_id"] = chatID.(int64)
	}
	return bot.MakeRequest("getChat", params)
}

// GetChatMember see https://core.telegram.org/bots/api#getchatmember
func (bot *Bot) GetChatMember(param interface{}, userID int64) (Update, error) {
	params := JSONBody{
		"user_id": userID,
	}
	switch chatID := param.(type) {
	case string:
		params["chat_id"] = chatID
	default:
		params["chat_id"] = chatID.(int64)
	}
	return bot.MakeRequest("getChatMember", params)
}

// GetChatAdministrators see https://core.telegram.org/bots/api#getchatadministrators
func (bot *Bot) GetChatAdministrators(param interface{}) (Update, error) {
	params := JSONBody{}
	switch chatID := param.(type) {
	case string:
		params["chat_id"] = chatID
	default:
		params["chat_id"] = chatID.(int64)
	}
	return bot.MakeRequest("getChatAdministrators", params)
}

// GetChatMembersCount see https://core.telegram.org/bots/api#getchatmemberscount
func (bot *Bot) GetChatMembersCount(param interface{}) (Update, error) {
	params := JSONBody{}
	switch chatID := param.(type) {
	case string:
		params["chat_id"] = chatID
	default:
		params["chat_id"] = chatID.(int64)
	}
	return bot.MakeRequest("getChatMembersCount", params)
}

// DeleteMessage see https://core.telegram.org/bots/api#deletemessage
func (bot *Bot) DeleteMessage(chatID int64, messageID int64) (Update, error) {
	return bot.MakeRequest("deleteMessage", JSONBody{
		"chat_id":    chatID,
		"message_id": messageID,
	})
}

// KickChatMember see https://core.telegram.org/bots/api#kickchatmember
func (bot *Bot) KickChatMember(chatID int64, userID int64, untilDate int64) (Update, error) {
	return bot.MakeRequest("kickChatMember", JSONBody{
		"chat_id":    chatID,
		"user_id":    userID,
		"until_date": untilDate,
	})
}

// UnbanChatMember see https://core.telegram.org/bots/api#unbanchatmember
func (bot *Bot) UnbanChatMember(chatID int64, userID int64) (Update, error) {
	return bot.MakeRequest("unbanChatMember", JSONBody{
		"chat_id": chatID,
		"user_id": userID,
	})
}

// RestrictChatMember see https://core.telegram.org/bots/api#restrictchatmember
func (bot *Bot) RestrictChatMember(chatID int64, userID int64, permissions map[string]bool, untilDate int64) (Update, error) {
	return bot.MakeRequest("restrictChatMember", JSONBody{
		"chat_id":     chatID,
		"user_id":     userID,
		"permissions": permissions,
		"until_date":  untilDate,
	})
}

// ExportChatInviteLink see https://core.telegram.org/bots/api#exportchatinvitelink
func (bot *Bot) ExportChatInviteLink(chatID int64) (Update, error) {
	return bot.MakeRequest("exportChatInviteLink", JSONBody{
		"chat_id": chatID,
	})
}

// Handle lets you set the handler for some command name or
// one of the supported endpoints.
func (bot *Bot) Handle(endpoint string, handler interface{}) {
	bot.handlers[endpoint] = handler
}

// Action lets you set the handler for some command name or
// one of the supported endpoints.
func (bot *Bot) Action(endpoint interface{}, handler interface{}) {
	switch end := endpoint.(type) {
	case string:
		bot.handlers["\f^"+end+"$"] = handler
	case *regexp.Regexp:
		bot.handlers["\f"+end.String()] = handler
	default:
		panic("easytgbot: unsupported endpoint")
	}
}

// Use
func (bot *Bot) Use(middleware ...MiddlewareFunc) {
	bot.middleware = append(bot.middleware, middleware...)
}

// applyMiddleware
func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// ApplyHandlers is apply handler
func (bot *Bot) ApplyHandlers(context interface{}, update *Update) (JSONBody, error) {
	updateType := update.GetType()
	// callback_query
	callbackQuery := update.Get("callback_query")
	if callbackQuery.Exists() {
		data := callbackQuery.Get("data").String()
		updateType = "\f" + data
		// for handlers
		for endpoint, handler := range bot.handlers {
			// skip command
			if endpoint[0:1] != "\f" {
				continue
			}
			endpoint = endpoint[1:]
			if regexp.MustCompile(endpoint).FindStringIndex(data) != nil {
				if handler, ok := handler.(func(interface{}, *Bot, *Update) JSONBody); ok {
					if len(bot.middleware) > 0 {
						handler = applyMiddleware(handler, bot.middleware...)
					}
					return handler(context, bot, update), nil
				}
			}
		}
	}

	// command first
	command, _ := update.Command()
	if len(command) > 0 {
		if pos := strings.Index(command, "@"); pos > -1 {
			botName := command[pos+1:]
			if bot.Name != "" && strings.ToLower(botName) == strings.ToLower(bot.Name) {
				command = command[0:pos]
			} else if bot.Name == "" {
				command = command[0:pos]
			}
		}

		// found handler
		if _, ok := bot.handlers[command]; ok {
			updateType = command
		}
	}
	// check handler has exists
	handler, ok := bot.handlers[updateType]
	if !ok {
		return JSONBody{}, fmt.Errorf("unsupported update type")
	}

	// execute
	if handler, ok := handler.(func(interface{}, *Bot, *Update) JSONBody); ok {
		if len(bot.middleware) > 0 {
			handler = applyMiddleware(handler, bot.middleware...)
		}
		return handler(context, bot, update), nil
	}

	return JSONBody{}, fmt.Errorf("unsupported update type")
}
