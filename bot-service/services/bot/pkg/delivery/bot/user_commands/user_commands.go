package user_commands

import (
	"encoding/json"
)

type UserCommand struct {
	Command string
	Uuid    string
}

func New(str string) *UserCommand {
	var u = &UserCommand{}
	err := json.Unmarshal([]byte(str), u)
	if err != nil {
		// ничего не делаем
	}
	return u
}

func (c *UserCommand) ToJson() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (c *UserCommand) IsNotEmpty() bool {
	return c.Command != ""
}

func AddPosition(uuid string) *UserCommand {
	return &UserCommand{
		Command: ADD_POSITION,
		Uuid:    uuid,
	}
}

func DecreasePosition(uuid string) *UserCommand {
	return &UserCommand{
		Command: DECREASE_POSITION,
		Uuid:    uuid,
	}
}

func ClickOnFolder(uuid string) *UserCommand {
	return &UserCommand{
		Command: CLICK_ON_FOLDER,
		Uuid:    uuid,
	}
}

func ClickOnProductItem(uuid string) *UserCommand {
	return &UserCommand{
		Command: CLICK_ON_PRODUCT_ITEM,
		Uuid:    uuid,
	}
}

const ADD_POSITION = "+"
const DECREASE_POSITION = "-"
const CLICK_ON_FOLDER = "@"
const CLICK_ON_PRODUCT_ITEM = "#"
