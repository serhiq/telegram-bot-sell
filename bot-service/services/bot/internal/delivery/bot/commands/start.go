package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (p Performer) startCommand(input *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(input.Chat.ID, "Добрый день")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUTTON_START_NEW_ORDER),
		))
	Respond(p.Bot, msg)
}
