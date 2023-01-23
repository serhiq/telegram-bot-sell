package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DecreasePosition struct {
	ChatID      int64
	ProductUuid string

	Bot         *tgbotapi.BotAPI
	RepoChat    repository.ChatRepository
	RepoProduct repository.ProductRepository
}

func (c DecreasePosition) Execute() (bot.Command, error) {

	state, err := c.RepoChat.GetOrCreateChat(c.ChatID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat  %s", err)
	}

	menuItem, err := c.RepoProduct.GetMenu(c.ProductUuid)

	if err != nil {
		return nil, fmt.Errorf("Failed to get product uuid =%s; err =  %s", c.ProductUuid, err)
	}

	if menuItem.Group {
		return nil, fmt.Errorf("incorrect command add to position for folder =%s", c.ProductUuid)
	}

	order := state.GetOrder()
	order.DecreaseMenuItem(menuItem)

	var msgText = "удалена позиция " + menuItem.Name + " " + menuItem.PriceString()
	strOrder, err := order.ToJson()
	if err != nil {
		return nil, fmt.Errorf("Decrease position command, json error for product  =%s", c.ProductUuid)
	}

	state.OrderStr = strOrder

	err = c.RepoChat.UpdateChat(state)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(state.ChatId, msgText)
	msg.ReplyMarkup = makeOrderKeyboard(order.CountPosition())

	_, err = Respond(c.Bot, msg)
	if err != nil {
		return nil, err
	}

	return DisplayMenuItemByUuidCommand{
		ProductUuid: c.ProductUuid,
		Bot:         c.Bot,
		RepoProduct: c.RepoProduct,
		RepoChat:    c.RepoChat,
		ChatID:      c.ChatID,
	}, nil
}
