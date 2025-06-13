package cmdimpl

import "github.com/HitoroOhria/task-executer/command"

type CommandImpl struct {
	readFile       func(path string) ([]byte, error)
	input          func(prompt string) string
	selectTaskName func(taskfile string) (string, error)
	runTask        func(taskfile string, name string, args ...string) error
}

type NewCommandArgs struct {
	ReadFile       func(path string) ([]byte, error)
	Input          func(prompt string) string
	SelectTaskName func(taskfile string) (string, error)
	RunTask        func(taskfile string, name string, args ...string) error
}

func NewCommand(args *NewCommandArgs) command.Command {
	return &CommandImpl{
		readFile:       args.ReadFile,
		input:          args.Input,
		selectTaskName: args.SelectTaskName,
		runTask:        args.RunTask,
	}
}

func (c *CommandImpl) ReadFile(path string) ([]byte, error) {
	return c.readFile(path)
}

func (c *CommandImpl) SelectTaskName(taskfile string) (string, error) {
	return c.selectTaskName(taskfile)
}

func (c *CommandImpl) Input(prompt string) string {
	return c.input(prompt)
}

func (c *CommandImpl) RunTask(taskfile string, name string, args ...string) error {
	return c.runTask(taskfile, name, args...)
}
