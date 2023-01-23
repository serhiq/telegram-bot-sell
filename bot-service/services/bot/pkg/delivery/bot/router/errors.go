package router

import (
	"fmt"
	"strings"
)

func NewCommandNotFound(command string) *ErrCommandNotFound {
	return &ErrCommandNotFound{
		command: command,
	}
}

type ErrCommandNotFound struct {
	command string
}

func (e ErrCommandNotFound) Error() string {
	return fmt.Sprintf("command not found %s", e.command)
}

func IsCommandNotFoundError(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "command not found")
}

func StartsWith(prefix string, content string) bool {
	return (strings.Split(content, " ")[0] == prefix)
}
