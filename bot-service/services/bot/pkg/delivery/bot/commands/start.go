package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BUTTON_START_NEW_ORDER = "🍔  Заказать"

type StartCommand struct {
	ChatID int64

	Bot *tgbotapi.BotAPI
}

func (s StartCommand) Execute() (bot.Command, error) {
	msg := tgbotapi.NewMessage(s.ChatID, "Добрый день")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUTTON_START_NEW_ORDER),
		))

	_, err := Respond(s.Bot, msg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//
//msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
//tgbotapi.NewKeyboardButtonContact("Предоставить номер телефона!"),
//))
