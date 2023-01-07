package app

import (
	"encoding/json"
)

type UserCommand struct {
	Command string
	Uuid    string
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

func ClickOnPosition(uuid string) *UserCommand {
	return &UserCommand{
		Command: CLICK_ON_POSITION,
		Uuid:    uuid,
	}
}

func FromJsonCommand(str string) *UserCommand {
	var stu = &UserCommand{}
	json.Unmarshal([]byte(str), stu)
	return stu
}

const ADD_POSITION = "+"
const DECREASE_POSITION = "-"
const CLICK_ON_POSITION = "."

//////////////////////////////////////////////
