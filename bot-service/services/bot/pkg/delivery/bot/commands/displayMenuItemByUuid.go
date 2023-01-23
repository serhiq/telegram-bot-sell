package commands

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/delivery/bot/user_commands"
	"bot/services/bot/pkg/repository"
	"bot/services/bot/pkg/repository/chat"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"strconv"
)

const INCREASE_POSITION_BUTTON = "➕ 1 шт."
const DECREASE_POSITION_BUTTON = "➖ 1 шт."

type DisplayMenuItemByUuidCommand struct {
	ChatID      int64
	ProductUuid string

	Bot         *tgbotapi.BotAPI
	RepoChat    repository.ChatRepository
	RepoProduct repository.ProductRepository
}

func (d DisplayMenuItemByUuidCommand) Execute() (bot.Command, error) {

	state, err := d.RepoChat.GetOrCreateChat(d.ChatID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat  %s", err)
	}

	menuItem, err := d.RepoProduct.GetMenu(d.ProductUuid)

	if err != nil {
		return nil, fmt.Errorf("Failed to get product uuid =%s; err =  %s", d.ProductUuid, err)
	}

	if menuItem.Group {
		return nil, fmt.Errorf("incorrect command display MenuItem for folder =%s", d.ProductUuid)
	}

	last := state.GetLastEditedMenuItem()

	if last.UuidMenuItem == menuItem.UUID {

		count := state.GetOrder().CountItemPosition(menuItem.UUID)

		text := menuItem.Name + "\n" + "Цена за " + menuItem.MeasureName + ":" + strconv.FormatInt(int64(menuItem.Price), 10) + " руб" + "\n" + "В корзине: " + count
		msg := tgbotapi.NewEditMessageText(state.ChatId, last.MessageId, text)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(INCREASE_POSITION_BUTTON, user_commands.AddPosition(menuItem.UUID).ToJson()),
				tgbotapi.NewInlineKeyboardButtonData(DECREASE_POSITION_BUTTON, user_commands.DecreasePosition(menuItem.UUID).ToJson()),
			))

		msg.ReplyMarkup = &keyboard

		_, err = d.Bot.Send(msg)

		if err != nil {
			return nil, err
		}

		return nil, nil

	} else {
		state.LastEditedMenuItemStr = ""
	}

	if menuItem.Image != "" {
		src := menuItem.Image

		file, err := ioutil.ReadFile(src)
		if err != nil {
			fmt.Printf("bot: error loading image %s", err)
		} else {
			photoFileBytes := tgbotapi.FileBytes{
				Name:  "picture",
				Bytes: file,
			}
			_, err = d.Bot.Send(tgbotapi.NewPhoto(state.ChatId, photoFileBytes))
			if err != nil {
				fmt.Printf("bot: error loading image %s", err)
			}
		}
	}

	count := state.GetOrder().CountItemPosition(menuItem.UUID)

	text := menuItem.Name + "\n" + "Цена за " + menuItem.MeasureName + ":" + strconv.FormatInt(int64(menuItem.Price), 10) + " руб" + "\n" + "В корзине: " + count
	msg := tgbotapi.NewMessage(state.ChatId, text)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(INCREASE_POSITION_BUTTON, user_commands.AddPosition(menuItem.UUID).ToJson()),
			tgbotapi.NewInlineKeyboardButtonData(DECREASE_POSITION_BUTTON, user_commands.DecreasePosition(menuItem.UUID).ToJson()),
		),
	)
	lastMenuItemMessage, err := Respond(d.Bot, msg)
	if err != nil {
		return nil, err
	}

	state.SaveLaseEdited(chat.LastEditedMenuItem{
		UuidMenuItem: menuItem.UUID,
		MessageId:    lastMenuItemMessage,
	})

	err = d.RepoChat.UpdateChat(state)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
