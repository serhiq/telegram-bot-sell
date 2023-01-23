package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddPositionToOrder struct {
	ChatID      int64
	ProductUuid string

	Bot         *tgbotapi.BotAPI
	RepoChat    repository.ChatRepository
	RepoProduct repository.ProductRepository
}

func (d AddPositionToOrder) Execute() (bot.Command, error) {

	state, err := d.RepoChat.GetOrCreateChat(d.ChatID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat  %s", err)
	}

	menuItem, err := d.RepoProduct.GetMenu(d.ProductUuid)

	if err != nil {
		return nil, fmt.Errorf("Failed to get product uuid =%s; err =  %s", d.ProductUuid, err)
	}

	if menuItem.Group {
		return nil, fmt.Errorf("incorrect command add to position for folder =%s", d.ProductUuid)
	}

	order := state.GetOrder()
	order.AddItem(menuItem)

	var msgText = "В заказ добавлена позиция " + menuItem.Name + " " + menuItem.PriceString()
	strOrder, err := order.ToJson()
	if err != nil {
		return nil, fmt.Errorf("Add poistion command, json error for product  =%s", d.ProductUuid)
	}

	state.OrderStr = strOrder

	err = d.RepoChat.UpdateChat(state)
	if err != nil {
		return nil, err

	}

	msg := tgbotapi.NewMessage(state.ChatId, msgText)
	msg.ReplyMarkup = makeOrderKeyboard(order.CountPosition())

	_, err = Respond(d.Bot, msg)
	if err != nil {
		return nil, err
	}

	return DisplayMenuItemByUuidCommand{
		Bot:         d.Bot,
		RepoProduct: d.RepoProduct,
		RepoChat:    d.RepoChat,
		ChatID:      d.ChatID,
		ProductUuid: d.ProductUuid,
	}, nil //return nil, nil
}

func makeOrderKeyboard(count string) interface{} {
	var textBucket = "🛒  Корзина (" + count + ")"

	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DISPLAY_MENU_BUTTON),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(textBucket),
		),
	)
}

const DISPLAY_MENU_BUTTON = "🍽  Меню"
