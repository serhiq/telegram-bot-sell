package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
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

type DisplayOrder struct {
	ChatID int64

	Bot      *tgbotapi.BotAPI
	RepoChat repository.ChatRepository
}

func (c DisplayOrder) Execute() (bot.Command, error) {

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

	_, err = Respond(c.Bot, msg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
