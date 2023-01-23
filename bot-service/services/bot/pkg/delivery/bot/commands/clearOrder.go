package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ClearOrder struct {
	ChatID int64

	Bot         *tgbotapi.BotAPI
	RepoChat    repository.ChatRepository
	RepoProduct repository.ProductRepository
}

func (c ClearOrder) Execute() (bot.Command, error) {

	state, err := c.RepoChat.GetChat(c.ChatID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat  %s", err)
	}

	state.NewOrder()
	err = c.RepoChat.UpdateChat(state)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(state.ChatId, "Заказ очищен")
	msg.ReplyMarkup = makeOrderKeyboard("0")

	_, err = Respond(c.Bot, msg)
	if err != nil {
		return nil, err
	}

	return DisplayMenuByUuid{
		ChatID:      c.ChatID,
		FolderUuid:  "",
		RepoProduct: c.RepoProduct,
		Bot:         c.Bot,
	}, nil
}
