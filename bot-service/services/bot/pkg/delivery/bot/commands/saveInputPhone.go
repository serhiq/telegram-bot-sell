package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
	"bot/services/bot/pkg/repository/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SaveInputPhone struct {
	Input string
	State *chat.Chat

	Bot         *tgbotapi.BotAPI
	RepoChat    repository.ChatRepository
	RepoProduct repository.ProductRepository
}

func (c SaveInputPhone) Execute() (bot.Command, error) {
	if !c.State.IsCorrectPhone(c.Input) {
		msg := tgbotapi.NewMessage(c.State.ChatId, "Введите телефон в федеральном формате")
		_, err := Respond(c.Bot, msg)
		if err != nil {
			return nil, err
		}
		return nil, nil

	} else {
		c.State.PhoneUser = c.Input
		c.State.ChatState = chat.STATE_PREPARE_ORDER

		err := c.RepoChat.UpdateChat(c.State)
		if err != nil {
			return nil, err
		}
	}
	return DisplayMenuByUuid{
		ChatID:      c.State.ChatId,
		FolderUuid:  "",
		RepoProduct: c.RepoProduct,
		Bot:         c.Bot,
	}, nil
}
