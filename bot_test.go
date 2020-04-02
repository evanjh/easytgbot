package easytgbot

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	bot, err := New("951886466:AAEhTr7--GVVIkEhVuWUZZqGNC1nxMBVQ7o")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("bot: %+v\n", bot)
}
