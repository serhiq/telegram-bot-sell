package commands

import (
	p "bot/pkg/type/product"
	"bot/services/bot/pkg/repository/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (p Performer) decreasePositionFromOrder(item *p.Product, state *chat.Chat) {
	order := state.GetOrder()
	order.DecreaseMenuItem(item)

	var msgText = "удалена позиция " + item.Name + " " + item.PriceString()

	strOrder, err := order.ToJson()
	if err != nil {
		log.Print("command: remove position fallied")
		return
		//	todo notify user
	}

	state.OrderStr = strOrder

	err = p.RepoChat.UpdateChat(state)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(state.ChatId, msgText)
	msg.ReplyMarkup = makeOrderKeyboard(order.CountPosition())
	Respond(p.Bot, msg)
}
