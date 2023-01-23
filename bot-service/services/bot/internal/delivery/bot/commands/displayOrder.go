package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	DISPLAY_ORDER_BUTTON = "🛒  Корзина"
	SEND_ORDER_BUTTON    = "✅  Подтвердить"
	CLEAR_ORDER_BUTTON   = "🗑  Очистить"
	BACK_ORDER_BUTTON    = "←  назад"
)

func (p Performer) displayOrder(chatId int64) error {
	state, err := p.RepoChat.GetOrCreateChat(chatId)

	if err != nil {
		//fmt.Printf("performer: faillied get chat  %s", err)
		return fmt.Errorf("performer: faillied get chat  %s", err)
	}

	headerBuilder := strings.Builder{}
	headerBuilder.WriteString("Имя: ")
	headerBuilder.WriteString(state.NameUser)
	headerBuilder.WriteString("\nТелефон: ")
	headerBuilder.WriteString(state.PhoneUser)

	msg := tgbotapi.NewMessage(state.ChatId, headerBuilder.String()+state.GetOrder().OrderDescription())

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(SEND_ORDER_BUTTON),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CLEAR_ORDER_BUTTON),
			tgbotapi.NewKeyboardButton(BACK_ORDER_BUTTON),
		),
	)
	Respond(p.Bot, msg)
	return nil
}
