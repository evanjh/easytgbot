package easytgbot_test

import (
	"fmt"
	"testing"

	"github.com/mylukin/easytgbot"
)

const (
	TestToken = "951886466:AAEhTr7--GVVIkEhVuWUZZqGNC1nxMBVQ7o"
)

func getBot(t *testing.T) (*easytgbot.BotAPI, error) {
	bot, err := easytgbot.New(TestToken)
	bot.Debug = true
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	return bot, err
}

func TestNew(t *testing.T) {
	bot, err := getBot(t)
	self := bot.Self
	fmt.Printf("Self: %T %+[1]v\n", self)

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
func TestNewBotAPIWithNoToken(t *testing.T) {
	_, err := easytgbot.NewBotAPIWith("", easytgbot.APIEndpoint)

	if err == nil {
		t.Error(err)
		t.Fail()
	}
}

func TestGetUpdates(t *testing.T) {
	bot, _ := getBot(t)
	_, err := bot.GetUpdates(easytgbot.JSONBody{
		"offset": 0,
		"limit":  1,
	})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
