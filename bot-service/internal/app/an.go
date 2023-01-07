package app

import (
	"bot/config"
	"bot/database"
	"bot/internal/entity"
	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"time"
)

type An struct {
	Cfg *config.Config

	Bot    *tgbotapi.BotAPI
	Client *resty.Client
	Db     *entity.GormDatabase
}

func New(cfg *config.Config) (*An, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false

	client := resty.New()

	db, err := database.Init()
	if err != nil {
		log.Panic("Can't connect to Mysql", err)
	}

	an := An{
		Cfg:    cfg,
		Bot:    bot,
		Client: client,
		Db:     entity.CreateGorm(db.Db),
	}
	return &an, nil
}

func (a *An) Start() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(int(10)).Minutes().Do(UpdateMenu, a)
	s.StartAsync()

	log.Printf("Bot online %s", a.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.Bot.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {

		if update.Message == nil {
			inputData := update.CallbackData()
			fromChat := update.FromChat()

			userCommand := FromJsonCommand(inputData)
			if userCommand != nil || userCommand.IsNotEmpty() {
				ProcessKeyboardInput(userCommand, fromChat.ID, a)
			}

			continue
		}
		CommandRouter(update.Message, a)
	}
}
func (a *An) Shutdown() {
	db, err := a.Db.Db.DB()
	if err != nil {
		log.Printf("database: error close database, %s", err)
	}
	err = db.Close()
	if err != nil {
		log.Printf("database: error close database, %s", err)
		return
	}
	log.Print("database: close")
}

//////////////////////////////////////////////////////

var UpdateMenu = func(a *An) {
	err := SyncMenu(a)
	if err != nil {
		log.Printf("Ошибка синхронизации, %s", err)
	}
}
