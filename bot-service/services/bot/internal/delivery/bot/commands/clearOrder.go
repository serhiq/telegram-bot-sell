package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p Performer) clearOrder(chatId int64) {
	state, err := p.RepoChat.GetOrCreateChat(chatId)

	if err != nil {
		fmt.Printf("performer: faillied get chat  %s", err)
		return
	}

	state.NewOrder()
	err = p.RepoChat.UpdateChat(state)
	if err != nil {
		// todo
		//p.sayError()
		return
	}

	msg := tgbotapi.NewMessage(state.ChatId, "Заказ очищен")
	msg.ReplyMarkup = makeOrderKeyboard("0")
	Respond(p.Bot, msg)

	p.displayRootMenu(chatId)
}
