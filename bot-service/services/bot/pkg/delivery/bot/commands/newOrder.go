package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewOrder struct {
	ChatID int64

	Bot         *tgbotapi.BotAPI
	RepoChat    repository.ChatRepository
	RepoProduct repository.ProductRepository
}

func (c NewOrder) Execute() (bot.Command, error) {

	state, err := c.RepoChat.GetOrCreateChat(c.ChatID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat  %s", err)
	}

	if !state.HaveUserName() {
		return InputName{
			Bot:      c.Bot,
			RepoChat: c.RepoChat,
			ChatID:   c.ChatID,
		}, nil

	}

	if !state.HaveUserPhone() {
		return InputPhone{
			Bot:      c.Bot,
			RepoChat: c.RepoChat,
			ChatID:   c.ChatID,
		}, nil

	}

	msg := tgbotapi.NewMessage(state.ChatId, "Добрый день")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUTTON_START_NEW_ORDER),
		))

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
