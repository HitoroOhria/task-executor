package cmdimpl

import "github.com/HitoroOhria/task-executer/command"

type CommandImpl struct {
	readFile       func(path string) ([]byte, error)
	prompt         func(maxNameLen int, varName string) string
	input          func(prompt string) string
	selectTaskName func(taskfile string) (string, error)
}

type NewCommandArgs struct {
	ReadFile       func(path string) ([]byte, error)
	Prompt         func(maxNameLen int, varName string) string
	Input          func(prompt string) string
	SelectTaskName func(taskfile string) (string, error)
}

func NewCommand(args *NewCommandArgs) command.Command {
	return &CommandImpl{
		readFile:       args.ReadFile,
		prompt:         args.Prompt,
		input:          args.Input,
		selectTaskName: args.SelectTaskName,
	}
}

func (c *CommandImpl) ReadFile(path string) ([]byte, error) {
	return c.readFile(path)
}

func (c *CommandImpl) SelectTaskName(taskfile string) (string, error) {
	return c.selectTaskName(taskfile)
}

func (c *CommandImpl) Prompt(maxNameLen int, varName string) string {
	return c.prompt(maxNameLen, varName)
}

func (c *CommandImpl) Input(prompt string) string {
	return c.input(prompt)
}
