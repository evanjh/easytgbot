package easytgbot

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	body := `{"update_id":818052699,"message":{"message_id":2100,"from":{"id":949939724,"is_bot":false,"first_name":"Hao1234Admin","username":"hao1234admin","language_code":"zh-hans"},"chat":{"id":949939724,"first_name":"Hao1234Admin","username":"hao1234admin","type":"private"},"date":1587018473,"text":"/ping","entities":[{"offset":0,"length":5,"type":"bot_command"}]}}`
	update := NewUpdate(string(body))
	fmt.Printf("%T %+[1]v", update.Command())

}
