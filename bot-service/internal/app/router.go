package app

import (
	"bot/internal/entity"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func startsWith(prefix string, content string) bool {
	return (strings.Split(content, " ")[0] == prefix)
}

func CommandRouter(msg *tgbotapi.Message, a *An) {
	state, err := a.Db.GetOrCreateChat(msg.Chat.ID)

	if err != nil {
		log.Println("Failed to get chat  %s", err)
	}

	// обработка отдельных команд в отдельных методах выглядела бы красивей
	if startsWith("/start", msg.Text) {
		// почему здесь в горутине, а остальные обработчики - нет?
		go StartCommand(a.Bot, msg)
	} else if msg.Text == BUTTON_START_NEW_ORDER {
		if !state.HaveUserName() {
			InputName(a.Bot, a.Db, state)
			return
		}

		if !state.HaveUserPhone() {
			InputPhone(a.Bot, a.Db, state)
			return
		}

		DisplayMenuRootMenu(a, state)

	} else if msg.Text == DISPLAY_MENU_BUTTON {
		DisplayMenuRootMenu(a, state)
		// почему где-то есть `return`, а где-то - нет?
		// в перспективе создаст проблемы
	} else if startsWith("🛒", msg.Text) {
		DisplayOrder(a, state)
		return

	} else if msg.Text == SEND_ORDER_BUTTON {
		order := state.GetOrder()
		order.Contacts.Phone = state.PhoneUser
		order.Comment = "от " + state.NameUser
		result, err := PostOrder(a, order)
		if err != nil {
			log.Println("post order: error ", err)
			// если мы получили ошибку, то зачем идти дальше?
			// более того, иногда метод может вернуть и результат, и ошибку. Но при этом результат неактуален
			// понтяно, что тут метод Ваш и Ваша логика, но я бы так не делал
			// если вернулась ошибка, то обрабатывать результат не стоит
		}

		if result != nil {
			log.Print(result)

			state.OrderStr = ""
			state.ChatState = entity.STATE_PREPARE_ORDER
			a.Db.UpdateChat(state)

			msg := tgbotapi.NewMessage(state.ChatId, "Ожидайте, мы с Вами свяжемся, для оплаты и уточнения деталей заказа")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(BUTTON_START_NEW_ORDER),
				))
			Respond(a.Bot, msg)
			return
		}

	} else if msg.Text == CLEAR_ORDER_BUTTON {

		state.NewOrder()
		a.Db.UpdateChat(state)
		msg := tgbotapi.NewMessage(state.ChatId, "Заказ очищен")
		msg.ReplyMarkup = makeOrderKeyboard("0")
		Respond(a.Bot, msg)

		DisplayMenu(a, state, "")

	} else if msg.Text == BACK_ORDER_BUTTON {
		DisplayMenuRootMenu(a, state)
	} else {
		if state.ChatState == entity.INPUT_NAME {

			if !state.IsCorrectName(msg.Text) {
				msg := tgbotapi.NewMessage(state.ChatId, "Имя не может быть пустым")
				Respond(a.Bot, msg)
				return
			}

			state.NameUser = msg.Text
			InputPhone(a.Bot, a.Db, state)
			return
		}

		if state.ChatState == entity.INPUT_PHONE {

			if !state.IsCorrectPhone(msg.Text) {
				msg := tgbotapi.NewMessage(state.ChatId, "Введите телефон в федеральном формате")
				Respond(a.Bot, msg)
				return
			}

			state.PhoneUser = msg.Text
			DisplayMenu(a, state, "")
			return
		}
		// echo command
		//Respond(a.Bot, tgbotapi.NewMessage(state.ChatId, msg.Text))
	}
}

func ProcessKeyboardInput(data *UserCommand, chatId int64, a *An) {
	state, err := a.Db.GetOrCreateChat(chatId)
	if err != nil {
		// пользователю тоже стоит что-то сообщить
		log.Println("Failed to get chat  %s", err)
		return
	}

	switch data.Command {
	case CLICK_ON_POSITION:
		{
			menuItem, _ := a.Db.GetMenu(data.Uuid)
			if menuItem != nil && menuItem.Group {
				DisplayMenu(a, state, data.Uuid)
				return
			}
			if menuItem != nil && !menuItem.Group {
				DisplayMenuItem(a, state, menuItem)
				return
			}
		}

	case ADD_POSITION:
		{
			menuItem, _ := a.Db.GetMenu(data.Uuid)
			if menuItem != nil && !menuItem.Group {
				addPositionToOrder(a, state, menuItem)
				DisplayMenuItem(a, state, menuItem)
				return
			}
		}

	case DECREASE_POSITION:
		{
			menuItem, _ := a.Db.GetMenu(data.Uuid)
			if menuItem != nil && !menuItem.Group {
				deletePositionFromOrder(a, state, menuItem)
				DisplayMenuItem(a, state, menuItem)
				return
			}
		}
	}
}

func addPositionToOrder(a *An, state *entity.Chat, item *entity.MenuItemDatabase) {
	order := state.GetOrder()
	order.AddMenuItem(item)

	var msgText = "В заказ добавлена позиция " + item.Name + " " + item.PriceString()
	state.OrderStr = order.ToJson()
	a.Db.UpdateChat(state)

	msg := tgbotapi.NewMessage(state.ChatId, msgText)
	msg.ReplyMarkup = makeOrderKeyboard(order.CountPosition())
	Respond(a.Bot, msg)
}

func deletePositionFromOrder(a *An, state *entity.Chat, item *entity.MenuItemDatabase) {
	order := state.GetOrder()
	order.DecreaseMenuItem(item)

	var msgText = "удалена позиция " + item.Name + " " + item.PriceString()
	state.OrderStr = order.ToJson()
	a.Db.UpdateChat(state)

	msg := tgbotapi.NewMessage(state.ChatId, msgText)
	msg.ReplyMarkup = makeOrderKeyboard(order.CountPosition())
	Respond(a.Bot, msg)
}

// почему возвращается `interface{}`, а не конкретный тип?
// и я бы передавать не количество позиций, а весь заказ - тогда можно быть отобразить и сумму, что даже полезней количества позиций
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
