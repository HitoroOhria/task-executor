package command

import (
	"github.com/HitoroOhria/task-executer/io"
	"github.com/HitoroOhria/task-executer/model"
)

type CommandImpl struct {
	readFile       func(path string) ([]byte, error)
	input          func(prompt string) string
	selectTaskName func(taskfile string) (string, error)
	runTask        func(taskfile string, name string, args ...string) error
}

func NewCommand() model.Command {
	return &CommandImpl{}
}

func (c *CommandImpl) ReadFile(path string) ([]byte, error) {
	return io.ReadFile(path)
}

func (c *CommandImpl) SelectTaskName(taskfile string) (string, error) {
	return io.SelectTaskName(taskfile)
}

func (c *CommandImpl) Input(prompt string) string {
	return io.Input(prompt)
}

func (c *CommandImpl) RunTask(taskfile string, name string, args ...string) error {
	return io.RunTask(taskfile, name, args...)
}
