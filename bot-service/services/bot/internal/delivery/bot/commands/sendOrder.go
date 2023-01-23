package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func (p Performer) sendOrder(chatId int64) {
	state, err := p.RepoChat.GetOrCreateChat(chatId)

	if err != nil {
		fmt.Printf("performer: faillied get chat  %s", err)
		return
	}

	order := state.GetOrder()
	order.Contacts.Phone = state.PhoneUser
	order.Comment = "от " + state.NameUser

	err = p.RepoOrder.Send(order)

	//result, err := client.PostOrder(order)
	if err != nil {
		log.Println("post order: error ", err)
	}

	//if result != nil {
	//	log.Print(result)
	//
	//	state.OrderStr = ""
	//	state.ChatState = repository.STATE_PREPARE_ORDER
	//	p.RepoChat.UpdateChat(state)
	//
	//	msg := tgbotapi.NewMessage(state.ChatId, "Ожидайте, мы с Вами свяжемся, для оплаты и уточнения деталей заказа")
	//	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
	//		tgbotapi.NewKeyboardButtonRow(
	//			tgbotapi.NewKeyboardButton(BUTTON_START_NEW_ORDER),
	//		))
	//	Respond(p.Bot, msg)
	//}

	//
	//
	//state.NewOrder()
	//err = p.RepoChat.UpdateChat(state)
	//if err != nil {
	//	// todo
	//	//p.sayError()
	//	return
	//}
	//
	//msg := tgbotapi.NewMessage(state.ChatId, "Заказ очищен")
	//msg.ReplyMarkup = makeOrderKeyboard("0")
	//command.Respond(p.Bot, msg)
	//
	//p.displayRootMenu(chatId)
}
