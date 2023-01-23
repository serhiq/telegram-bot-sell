package bot

import (
	"bot/services/bot/internal/delivery/bot/commands"
	"bot/services/bot/pkg/repository"
	"bot/services/bot/pkg/repository/chat"
	"bot/services/bot/pkg/repository/product"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"time"
)

type Delivery struct {
	performer *commands.Performer
}

type Options struct {
	Token string
}

func New(options Options, repoProduct *product.Repository, repoChat *chat.Repository, repoOrder repository.OrderRepository) (*Delivery, error) {
	bot, err := tgbotapi.NewBotAPI(options.Token)
	if err != nil {
		return nil, err
	}

	bot.Debug = false

	performer := &commands.Performer{
		RepoProduct: repoProduct,
		RepoChat:    repoChat,
		RepoOrder:   repoOrder,
		Bot:         bot,
	}

	return &Delivery{
		performer: performer,
	}, nil
}

func (d *Delivery) Start() error {
	log.Printf("Bot online %s", d.performer.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := d.performer.Bot.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {

		err := d.performer.Answer(update)
		if err != nil {
			log.Println("bot: process update %s", err)
		}
	}

	return nil
}
