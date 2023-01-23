package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/delivery/bot/user_commands"
	"bot/services/bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DisplayMenuByUuid struct {
	ChatID     int64
	FolderUuid string

	Bot         *tgbotapi.BotAPI
	RepoProduct repository.ProductRepository
}

func (d DisplayMenuByUuid) Execute() (bot.Command, error) {
	msg := tgbotapi.NewMessage(d.ChatID, "Выберите категорию")
	msg.ReplyMarkup = makeMenuKeyboard(d.FolderUuid, d.RepoProduct)
	_, err := Respond(d.Bot, msg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func makeMenuKeyboard(parentUuid string, repoProduct repository.ProductRepository) tgbotapi.InlineKeyboardMarkup {

	menuitems, err := repoProduct.GetMenuItemByParent(parentUuid)
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

func chunkSlice(items []tgbotapi.InlineKeyboardButton, chunkSize int) (chunks [][]tgbotapi.InlineKeyboardButton) {
	for chunkSize < len(items) {
		chunks = append(chunks, items[0:chunkSize])
		items = items[chunkSize:]
	}
	return append(chunks, items)
}
