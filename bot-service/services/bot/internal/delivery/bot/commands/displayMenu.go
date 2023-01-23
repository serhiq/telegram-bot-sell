package commands

import (
	"bot/services/bot/pkg/delivery/bot/user_commands"
	"bot/services/bot/pkg/repository/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p Performer) displayMenu(state *chat.Chat) {

	msg := tgbotapi.NewMessage(state.ChatId, "Меню")
	// todo сделать подсчет позитций
	//msg.ReplyMarkup = makeOrderKeyboard(state.GetOrder().CountPosition())
	Respond(p.Bot, msg)
	p.displayMenuByUuid("", state.ChatId)
}

func (p Performer) displayMenuByUuid(parentUuid string, chatId int64) {

	//state.ChatState = repository.STATE_PREPARE_ORDER
	//a.Db.UpdateChat(state)
	msg := tgbotapi.NewMessage(chatId, "Выберите категорию")
	msg.ReplyMarkup = p.makeMenuKeyboard(parentUuid)
	Respond(p.Bot, msg)
}

func (p Performer) displayRootMenu(chatId int64) {
	p.displayMenuByUuid("", chatId)
}

func (p Performer) makeMenuKeyboard(parentUuid string) tgbotapi.InlineKeyboardMarkup {

	menuitems, err := p.RepoProduct.GetMenuItemByParent(parentUuid)
	if err != nil {

	}

	buttons := []tgbotapi.InlineKeyboardButton{}

	for _, menuitem := range menuitems {
		if menuitem.Group {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("🗀  "+menuitem.Name, user_commands.ClickOnFolder(menuitem.UUID).ToJson()))
		} else {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(menuitem.Name, user_commands.ClickOnProductItem(menuitem.UUID).ToJson()))
		}
	}

	rows := chunkSlice(buttons, 3)
	return tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)
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

func chunkSlice(items []tgbotapi.InlineKeyboardButton, chunkSize int) (chunks [][]tgbotapi.InlineKeyboardButton) {
	for chunkSize < len(items) {
		chunks = append(chunks, items[0:chunkSize])
		items = items[chunkSize:]
	}
	return append(chunks, items)
}
