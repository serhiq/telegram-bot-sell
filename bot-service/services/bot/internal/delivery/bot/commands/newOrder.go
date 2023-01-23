package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BUTTON_START_NEW_ORDER = "🍔  Заказать"

func (p Performer) newOrder(chatId int64) {
	state, err := p.RepoChat.GetOrCreateChat(chatId)

	if err != nil {
		fmt.Printf("performer: faillied get chat  %s", err)
		return
	}

	if !state.HaveUserName() {
		p.inputName(state)
		return
	}

	if !state.HaveUserPhone() {
		p.inputPhone(state)
		return
	}
	p.displayRootMenu(state.ChatId)

	msg := tgbotapi.NewMessage(chatId, "Добрый день")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUTTON_START_NEW_ORDER),
		))

	Respond(p.Bot, msg)
}
