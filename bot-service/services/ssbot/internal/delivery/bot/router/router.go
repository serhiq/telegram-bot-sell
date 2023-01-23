package router

import (
	"bot/services/bot/pkg/delivery/bot"
	"bot/services/bot/pkg/delivery/bot/commands"
	"bot/services/bot/pkg/delivery/bot/router"
	userCommand "bot/services/bot/pkg/delivery/bot/user_commands"
	"bot/services/bot/pkg/repository"
	"bot/services/bot/pkg/repository/chat"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotRouter struct {
	RepoProduct repository.ProductRepository
	RepoChat    repository.ChatRepository
	RepoOrder   repository.OrderRepository

	Bot *tgbotapi.BotAPI
}

func (b BotRouter) GetCommand(update tgbotapi.Update) (bot.Command, error) {

	if update.Message == nil {
		inputData := update.CallbackData()
		fromChat := update.FromChat()

		c := userCommand.New(inputData)
		if c != nil || c.IsNotEmpty() {
			return b.AnswerUserCommand(c, fromChat.ID)
		}
	}

	commandFromButton, err := b.AnswerOnClickButton(update.Message.Text, update.FromChat().ID)
	if err != nil {
		// возможно это свободный ввод
		if router.IsCommandNotFoundError(err) {
			return b.answerFreeInput(update.FromChat().ID, update.Message.Text)
		}
		return nil, err
	}

	return commandFromButton, nil
}

func (b BotRouter) answerFreeInput(chatID int64, input string) (bot.Command, error) {
	state, err := b.RepoChat.GetChat(chatID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat  %s", err)
	}

	switch state.ChatState {
	case chat.INPUT_NAME:
		return commands.SaveInputName{
			Input:    input,
			State:    state,
			Bot:      b.Bot,
			RepoChat: b.RepoChat,
		}, nil

	case chat.INPUT_PHONE:
		return commands.SaveInputPhone{
			Input:       input,
			State:       state,
			Bot:         b.Bot,
			RepoChat:    b.RepoChat,
			RepoProduct: b.RepoProduct,
		}, nil

	}
	return nil, router.NewCommandNotFound(input)
}

func (b BotRouter) AnswerUserCommand(c *userCommand.UserCommand, chatID int64) (bot.Command, error) {
	switch c.Command {
	case userCommand.CLICK_ON_FOLDER:
		return commands.DisplayMenuByUuid{
			Bot:         b.Bot,
			RepoProduct: b.RepoProduct,
			ChatID:      chatID,
			FolderUuid:  c.Uuid,
		}, nil

	case userCommand.CLICK_ON_PRODUCT_ITEM:
		return commands.DisplayMenuItemByUuidCommand{
			Bot:         b.Bot,
			RepoProduct: b.RepoProduct,
			RepoChat:    b.RepoChat,
			ChatID:      chatID,
			ProductUuid: c.Uuid,
		}, nil

	case userCommand.ADD_POSITION:
		return commands.AddPositionToOrder{
			Bot:         b.Bot,
			RepoProduct: b.RepoProduct,
			RepoChat:    b.RepoChat,
			ChatID:      chatID,
			ProductUuid: c.Uuid,
		}, nil

	case userCommand.DECREASE_POSITION:
		return commands.DecreasePosition{
			Bot:         b.Bot,
			RepoProduct: b.RepoProduct,
			RepoChat:    b.RepoChat,
			ChatID:      chatID,
			ProductUuid: c.Uuid,
		}, nil
	default:
		return nil, router.NewCommandNotFound(c.ToJson())
	}

}

func (b BotRouter) AnswerOnClickButton(text string, chatID int64) (bot.Command, error) {

	if router.StartsWith("/start", text) {
		return commands.StartCommand{
			ChatID: chatID,
			Bot:    b.Bot,
		}, nil
	}

	if text == commands.BUTTON_START_NEW_ORDER {
		return commands.NewOrder{
			Bot:         b.Bot,
			RepoProduct: b.RepoProduct,
			RepoChat:    b.RepoChat,
			ChatID:      chatID,
		}, nil
	}

	if text == commands.DISPLAY_MENU_BUTTON {
		return commands.DisplayMenuByUuid{
			ChatID:      chatID,
			FolderUuid:  "",
			RepoProduct: b.RepoProduct,
			Bot:         b.Bot,
		}, nil
	}

	if router.StartsWith("🛒", text) {
		return commands.DisplayOrder{
			ChatID:   chatID,
			Bot:      b.Bot,
			RepoChat: b.RepoChat,
		}, nil
	}

	if text == commands.SEND_ORDER_BUTTON {
		return commands.SendOrder{
			ChatID:    chatID,
			RepoChat:  b.RepoChat,
			RepoOrder: b.RepoOrder,
			Bot:       b.Bot,
		}, nil
	}

	if text == commands.CLEAR_ORDER_BUTTON {
		return commands.ClearOrder{
			ChatID:      chatID,
			RepoChat:    b.RepoChat,
			RepoProduct: b.RepoProduct,
			Bot:         b.Bot,
		}, nil
	}

	if text == commands.BACK_ORDER_BUTTON {
		return commands.DisplayMenuByUuid{
			ChatID:      chatID,
			FolderUuid:  "",
			RepoProduct: b.RepoProduct,
			Bot:         b.Bot,
		}, nil
	}
	return nil, router.NewCommandNotFound(text)
}
