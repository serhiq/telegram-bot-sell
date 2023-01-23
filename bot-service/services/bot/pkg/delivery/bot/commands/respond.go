package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Respond(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) (int, error) {
	resultMsg, err := bot.Send(msg)

	if err != nil {
		return 0, NewErrorRespond(err)
	}
	return resultMsg.MessageID, nil
}

func NewErrorRespond(err error) *ErrRespond {
	return &ErrRespond{
		err: err.Error(),
	}

}

type ErrRespond struct {
	err string
}

func (e ErrRespond) Error() string {
	return fmt.Sprintf("Failed to respond  %s", e.err)
}

func IsRespondError(err error) bool {
	_, ok := err.(ErrRespond)
	return ok
}
