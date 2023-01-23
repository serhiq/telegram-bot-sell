package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
	"bot/services/bot/pkg/repository/chat"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InputPhone struct {
	ChatID int64

	Bot      *tgbotapi.BotAPI
	RepoChat repository.ChatRepository
}

func (c InputPhone) Execute() (bot.Command, error) {
	state, err := c.RepoChat.GetOrCreateChat(c.ChatID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat  %s", err)
	}

	state.ChatState = chat.INPUT_PHONE

	err = c.RepoChat.UpdateChat(state)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(state.ChatId, "Укажите контактный телефон")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err = Respond(c.Bot, msg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
