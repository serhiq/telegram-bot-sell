package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SendOrder struct {
	ChatID int64

	Bot       *tgbotapi.BotAPI
	RepoChat  repository.ChatRepository
	RepoOrder repository.OrderRepository
}

func (c SendOrder) Execute() (bot.Command, error) {

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

	order := state.GetOrder()
	order.Contacts.Phone = state.PhoneUser
	order.Comment = "от " + state.NameUser

	err = c.RepoOrder.Send(order)

	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(state.ChatId, "Ваш заказ отправлен")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUTTON_START_NEW_ORDER),
		))

	_, err = Respond(c.Bot, msg)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
