package app

import (
	"bot/internal/entity"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func Respond(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) int {
	resultMsg, err := bot.Send(msg)

	if err != nil {
		log.Println("Failed to respond  %s", err)
	}
	return resultMsg.MessageID
}

///////////////////////////////////////////////////////////////////////////////
const BUTTON_START_NEW_ORDER = "🍔  Заказать"

func StartCommand(bot *tgbotapi.BotAPI, input *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(input.Chat.ID, "Добрый день")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BUTTON_START_NEW_ORDER),
		))
	Respond(bot, msg)
}

////////////////////////////////////////////////////////////////////

func InputName(bot *tgbotapi.BotAPI, g *entity.GormDatabase, state *entity.Chat) {
	state.ChatState = entity.INPUT_NAME
	g.UpdateChat(state)
	msg := tgbotapi.NewMessage(state.ChatId, "Введите пожалуйста, Ваше имя")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	Respond(bot, msg)
}

func InputPhone(bot *tgbotapi.BotAPI, g *entity.GormDatabase, state *entity.Chat) {
	state.ChatState = entity.INPUT_PHONE
	g.UpdateChat(state)

	msg := tgbotapi.NewMessage(state.ChatId, "Укажите контактный телефон")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	Respond(bot, msg)
}

///////////////////////////////////////////////////////////////////////////

const DISPLAY_MENU_BUTTON = "🍽  Меню"

func DisplayMenu(a *An, state *entity.Chat, parentUuid string) {
	state.ChatState = entity.STATE_PREPARE_ORDER
	a.Db.UpdateChat(state)
	msg := tgbotapi.NewMessage(state.ChatId, "Выберите категорию")
	msg.ReplyMarkup = makeMenuKeyboard(a, parentUuid)
	Respond(a.Bot, msg)
}

func makeMenuKeyboard(a *An, parentUuid string) interface{} {
	menuitems, err := a.Db.GetMenuItemByParent(parentUuid)
	if err != nil {
		// ай-яй-яй
	}

	buttons := []tgbotapi.InlineKeyboardButton{}

	for _, menuitem := range menuitems {
		if menuitem.Group {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("🗀  "+menuitem.Name, ClickOnPosition(menuitem.UUID).ToJson()))
		} else {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(menuitem.Name, ClickOnPosition(menuitem.UUID).ToJson()))
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

func DisplayMenuRootMenu(a *An, state *entity.Chat) {
	msg := tgbotapi.NewMessage(state.ChatId, "Меню")
	msg.ReplyMarkup = makeOrderKeyboard(state.GetOrder().CountPosition())
	Respond(a.Bot, msg)
	DisplayMenu(a, state, "")

}

//////////////////////////////////////////////////////////////////
const INCREASE_POSITION_BUTTON = "➕ 1 шт."
const DECREASE_POSITION_BUTTON = "➖ 1 шт."

func DisplayMenuItem(a *An, state *entity.Chat, item *entity.MenuItemDatabase) {
	last := state.GetLastEditedMenuItem()
	log.Println(last)

	if last.UuidMenuItem == item.UUID {

		count := state.GetOrder().CountItemPosition(item.UUID)

		text := item.Name + "\n" + "Цена за " + item.MeasureName + ":" + strconv.FormatInt(int64(item.Price), 10) + " руб" + "\n" + "В корзине: " + count
		msg := tgbotapi.NewEditMessageText(state.ChatId, last.MessageId, text)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(INCREASE_POSITION_BUTTON, AddPosition(item.UUID).ToJson()),
				tgbotapi.NewInlineKeyboardButtonData(DECREASE_POSITION_BUTTON, DecreasePosition(item.UUID).ToJson()),
			))

		msg.ReplyMarkup = &keyboard

		_, err := a.Bot.Send(msg)

		if err != nil {
			log.Println("Failed to respond  %s", err)
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
	_, err = a.Bot.Send(tgbotapi.NewPhoto(state.ChatId, photoFileBytes))

	count := state.GetOrder().CountItemPosition(item.UUID)

	text := item.Name + "\n" + "Цена за " + item.MeasureName + ":" + strconv.FormatInt(int64(item.Price), 10) + " руб" + "\n" + "В корзине: " + count
	msg := tgbotapi.NewMessage(state.ChatId, text)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(INCREASE_POSITION_BUTTON, AddPosition(item.UUID).ToJson()),
			tgbotapi.NewInlineKeyboardButtonData(DECREASE_POSITION_BUTTON, DecreasePosition(item.UUID).ToJson()),
		),
	)
	lastMenuItemMessage := Respond(a.Bot, msg)

	state.SaveLaseEdited(entity.LastEditedMenuItem{
		UuidMenuItem: item.UUID,
		MessageId:    lastMenuItemMessage,
	})
	a.Db.UpdateChat(state)
}

////////////////////////////////////////////////////////////
const (
	DISPLAY_ORDER_BUTTON = "🛒  Корзина"
	SEND_ORDER_BUTTON    = "✅  Подтвердить"
	CLEAR_ORDER_BUTTON   = "🗑  Очистить"
	BACK_ORDER_BUTTON    = "←  назад"
)

func DisplayOrder(a *An, state *entity.Chat) {
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
	Respond(a.Bot, msg)
}
