package bot

import (
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
