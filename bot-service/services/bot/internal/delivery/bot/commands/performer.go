package commands

import (
	userCommand "bot/services/bot/pkg/delivery/bot/user_commands"
	"bot/services/bot/pkg/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Performer struct {
	RepoProduct repository.ProductRepository
	RepoChat    repository.ChatRepository
	RepoOrder   repository.OrderRepository
	Bot         *tgbotapi.BotAPI
}

func (p Performer) Answer(update tgbotapi.Update) error {
	if !p.answerUserCommand(update) {
		err := p.answerMessage(update.Message)
		if err != nil {
			log.Printf("bot: answer error %s", err)
		}
	}
	return nil
}

func (p Performer) answerMessage(msg *tgbotapi.Message) error {
	if startsWith("/start", msg.Text) {
		p.startCommand(msg)
	} else if msg.Text == BUTTON_START_NEW_ORDER {
		p.newOrder(msg.Chat.ID)

	} else if msg.Text == DISPLAY_MENU_BUTTON {
		p.displayRootMenu(msg.Chat.ID)
	} else if startsWith("🛒", msg.Text) {
		p.displayOrder(msg.Chat.ID)

	} else if msg.Text == SEND_ORDER_BUTTON {
		p.sendOrder(msg.Chat.ID)

	} else if msg.Text == CLEAR_ORDER_BUTTON {
		p.clearOrder(msg.Chat.ID)

	} else if msg.Text == BACK_ORDER_BUTTON {
		p.displayRootMenu(msg.Chat.ID)
	} else {
		// it free input

		p.processUserInput(msg.Text, msg.Chat.ID)

	}

	return nil
}
func (p Performer) processUserInput(input string, chatId int64) {

	state, err := p.RepoChat.GetOrCreateChat(chatId)

	if err != nil {
		fmt.Printf("performer: faillied get chat  %s", err)
		return
	}

	if !state.HaveUserName() {
		if !state.IsCorrectName(input) {
			msg := tgbotapi.NewMessage(state.ChatId, "Имя не может быть пустым")
			Respond(p.Bot, msg)
			return
		}

		state.NameUser = input
		p.inputName(state)
		return
	}

	if state.HaveUserPhone() {
		if !state.IsCorrectPhone(input) {
			msg := tgbotapi.NewMessage(state.ChatId, "Введите телефон в федеральном формате")
			Respond(p.Bot, msg)
			return
		}

		state.PhoneUser = input
		p.displayRootMenu(chatId)

		return
	}
	//	 todo комментарий к заказу,
}

func (p Performer) answerUserCommand(update tgbotapi.Update) bool {
	if update.Message == nil {
		inputData := update.CallbackData()
		fromChat := update.FromChat()

		userCommand := userCommand.New(inputData)
		if userCommand != nil || userCommand.IsNotEmpty() {
			p.PerformUserCommand(userCommand, fromChat.ID)
			return true
		}
	}
	return false
}

func (p Performer) PerformUserCommand(data *userCommand.UserCommand, chatId int64) {
	state, err := p.RepoChat.GetOrCreateChat(chatId)
	if err != nil {
		log.Println("Failed to get chat  %s", err)
		return
	}

	switch data.Command {
	case userCommand.CLICK_ON_FOLDER:
		p.displayMenuByUuid(data.Uuid, chatId)
	case userCommand.CLICK_ON_PRODUCT_ITEM:
		{
			menuItem, _ := p.RepoProduct.GetMenu(data.Uuid)
			p.displayMenuItem(state, menuItem)
		}

	case userCommand.ADD_POSITION:
		{
			menuItem, _ := p.RepoProduct.GetMenu(data.Uuid)
			if menuItem != nil && !menuItem.Group {
				p.addPositionToOrder(menuItem, state)
				p.displayMenuItem(state, menuItem)
				return
			}
		}

	case userCommand.DECREASE_POSITION:
		{
			menuItem, _ := p.RepoProduct.GetMenu(data.Uuid)
			if menuItem != nil && !menuItem.Group {
				p.decreasePositionFromOrder(menuItem, state)
				p.displayMenuItem(state, menuItem)
				return
			}
		}
	}
}

func startsWith(prefix string, content string) bool {
	return (strings.Split(content, " ")[0] == prefix)
}
