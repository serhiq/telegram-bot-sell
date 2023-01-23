package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
	"bot/services/bot/pkg/repository/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SaveInputName struct {
	Input string
	State *chat.Chat

	Bot      *tgbotapi.BotAPI
	RepoChat repository.ChatRepository
}

func (c SaveInputName) Execute() (bot.Command, error) {
	if !c.State.IsCorrectName(c.Input) {
		msg := tgbotapi.NewMessage(c.State.ChatId, "Имя не может быть пустым")
		_, err := Respond(c.Bot, msg)
		if err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		c.State.NameUser = c.Input
		c.State.ChatState = chat.INPUT_PHONE

		err := c.RepoChat.UpdateChat(c.State)
		if err != nil {
			return nil, err
		}
	}

	return InputPhone{
		ChatID:   0,
		Bot:      c.Bot,
		RepoChat: c.RepoChat,
	}, nil
}
