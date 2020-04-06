package easytgbot_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/imroc/req"
	"github.com/mylukin/easytgbot"
)

const (
	TestToken = "951886466:AAEhTr7--GVVIkEhVuWUZZqGNC1nxMBVQ7o"
)

func getBot(t *testing.T) (*easytgbot.BotAPI, error) {
	bot, err := easytgbot.New(
		TestToken,
		easytgbot.WithDebug(true),
	)
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
	_, err := easytgbot.New(
		"",
		easytgbot.WithDebug(true),
	)

	if err == nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDeleteWebhook(t *testing.T) {
	bot, _ := getBot(t)
	_, err := bot.DeleteWebhook()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSetWebhookWithCert(t *testing.T) {
	bot, _ := getBot(t)
	file, _ := os.Open("./ca.cert")
	_, err := bot.SetWebhook(easytgbot.JSONBody{
		"url":             "https://test01.tg.atmy.work/",
		"max_connections": 10,
		"allowed_updates": []string{
			"message",
			"edited_channel_post",
			"callback_query",
		},
		"certificate": req.FileUpload{
			File:      file,
			FieldName: "certificate",
			FileName:  "ca.cert",
		},
	})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSetWebhooks(t *testing.T) {
	bot, _ := getBot(t)
	_, err := bot.SetWebhook(easytgbot.JSONBody{
		"url":             "https://test01.tg.atmy.work/",
		"max_connections": 100,
		"allowed_updates": []string{
			"message",
			"callback_query",
		},
	})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestGetWebhookInfo(t *testing.T) {
	bot, _ := getBot(t)
	_, err := bot.GetWebhookInfo()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestUpdates(t *testing.T) {
	bot, _ := getBot(t)
	updates, err := bot.GetUpdates(easytgbot.JSONBody{
		"offset": 0,
		"limit":  1,
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	for update := range updates {
		fmt.Printf("update: %T %+[1]v\n", update)
	}
}
