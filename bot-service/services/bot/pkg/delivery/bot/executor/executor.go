package executor

import (
	"bot/services/bot/pkg/delivery/bot"
)

func Answer(c bot.Command) error {

	answer, err := c.Execute()
	if err != nil {
		return err
	}
	if answer != nil {
		return Answer(answer)
	}

	return nil
}
