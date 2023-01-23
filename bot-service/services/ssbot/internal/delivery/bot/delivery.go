package bot

import (
	"bot/services/bot/pkg/delivery/bot/commands"
	"bot/services/bot/pkg/delivery/bot/executor"
	"bot/services/bot/pkg/delivery/bot/router"
	"bot/services/bot/pkg/repository"
	routerIml "bot/services/ssbot/internal/delivery/bot/router"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"time"
)

type Delivery struct {
	router *routerIml.BotRouter
}

type Options struct {
	Token string
}

func New(options Options, repoProduct repository.ProductRepository, repoChat repository.ChatRepository, repoOrder repository.OrderRepository) (*Delivery, error) {
	b, err := tgbotapi.NewBotAPI(options.Token)
	if err != nil {
		return nil, err
	}

	b.Debug = false

	r := &routerIml.BotRouter{
		RepoProduct: repoProduct,
		RepoChat:    repoChat,
		RepoOrder:   repoOrder,
		Bot:         b,
	}

	return &Delivery{
		router: r,
	}, nil
}

func (d *Delivery) Start() error {
	log.Printf("Bot online %s", d.router.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := d.router.Bot.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {

		command, err := d.router.GetCommand(update)

		if err != nil {
			if router.IsCommandNotFoundError(err) {
				log.Printf("bot: chatId: %v, err: %s", update.FromChat().ID, err)
				continue
			}
			log.Printf("bot: chatId: %v, err: %s", update.FromChat().ID, err)
			continue
		}

		err = executor.Answer(command)
		if err != nil {
			if commands.IsRespondError(err) {
				log.Printf("bot: chatId: %v, err: %s", update.FromChat().ID, err)
				continue
			}

			log.Println("bot: process update %s", err)
			continue
		}
	}

	return nil
}
