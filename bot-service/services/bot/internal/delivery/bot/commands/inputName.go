package commands

import (
	"bot/services/bot/pkg/repository/chat"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p Performer) inputName(state *chat.Chat) {
	state.ChatState = chat.INPUT_NAME
	err := p.RepoChat.UpdateChat(state)
	if err != nil {
		fmt.Printf("perpormer: fallied update chat %s", err)
		return

	}

	msg := tgbotapi.NewMessage(state.ChatId, "Введите пожалуйста, Ваше имя")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	Respond(p.Bot, msg)
}
