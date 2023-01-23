package commands

import (
	p "bot/pkg/type/product"
	"bot/services/bot/pkg/delivery/bot/user_commands"
	"bot/services/bot/pkg/repository/chat"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"strconv"
)

const INCREASE_POSITION_BUTTON = "➕ 1 шт."
const DECREASE_POSITION_BUTTON = "➖ 1 шт."

func (p Performer) displayMenuItem(state *chat.Chat, item *p.Product) {

	last := state.GetLastEditedMenuItem()
	log.Println(last)

	if last.UuidMenuItem == item.UUID {

		count := state.GetOrder().CountItemPosition(item.UUID)

		text := item.Name + "\n" + "Цена за " + item.MeasureName + ":" + strconv.FormatInt(int64(item.Price), 10) + " руб" + "\n" + "В корзине: " + count
		msg := tgbotapi.NewEditMessageText(state.ChatId, last.MessageId, text)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(INCREASE_POSITION_BUTTON, user_commands.AddPosition(item.UUID).ToJson()),
				tgbotapi.NewInlineKeyboardButtonData(DECREASE_POSITION_BUTTON, user_commands.DecreasePosition(item.UUID).ToJson()),
			))

		msg.ReplyMarkup = &keyboard

		_, err := p.Bot.Send(msg)

		if err != nil {
			log.Printf("Failed to respond  %s", err)
		}
		return

	} else {
		state.LastEditedMenuItemStr = ""
	}

	src := item.Image

	file, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
	}

	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: file,
	}
	_, err = p.Bot.Send(tgbotapi.NewPhoto(state.ChatId, photoFileBytes))

	count := state.GetOrder().CountItemPosition(item.UUID)

	text := item.Name + "\n" + "Цена за " + item.MeasureName + ":" + strconv.FormatInt(int64(item.Price), 10) + " руб" + "\n" + "В корзине: " + count
	msg := tgbotapi.NewMessage(state.ChatId, text)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(INCREASE_POSITION_BUTTON, user_commands.AddPosition(item.UUID).ToJson()),
			tgbotapi.NewInlineKeyboardButtonData(DECREASE_POSITION_BUTTON, user_commands.DecreasePosition(item.UUID).ToJson()),
		),
	)
	lastMenuItemMessage := Respond(p.Bot, msg)

	state.SaveLaseEdited(chat.LastEditedMenuItem{
		UuidMenuItem: item.UUID,
		MessageId:    lastMenuItemMessage,
	})

	p.RepoChat.UpdateChat(state)
}
