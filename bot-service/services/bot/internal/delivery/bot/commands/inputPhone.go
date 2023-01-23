package commands

import (
	"bot/services/bot/pkg/repository/chat"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p Performer) inputPhone(state *chat.Chat) {
	state.ChatState = chat.INPUT_PHONE
	err := p.RepoChat.UpdateChat(state)
	if err != nil {
		fmt.Printf("performer: fallied update chat %s", err)
		return

	}

	msg := tgbotapi.NewMessage(state.ChatId, "Укажите контактный телефон")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	Respond(p.Bot, msg)
}
