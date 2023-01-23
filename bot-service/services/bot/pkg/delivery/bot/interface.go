package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Sbot interface {
	Start() error
}

type Router interface {
	GetCommand(tgbotapi.Update) (Command, error)
}

type Command interface {
	Execute() (Command, error)
}
