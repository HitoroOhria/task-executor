package adapter

import (
	"fmt"

	"github.com/HitoroOhria/task-executer/domain/console"
	"github.com/HitoroOhria/task-executer/domain/value"
	"github.com/HitoroOhria/task-executer/io"
)

type RunnerImpl struct {
	readFile       func(path string) ([]byte, error)
	input          func(prompt string) string
	selectTaskName func(taskfile string) (string, error)
	runTask        func(taskfile string, name string, args ...string) error
}

func NewRunner() console.Runner {
	return &RunnerImpl{}
}

func (c *RunnerImpl) ReadFile(path string) ([]byte, error) {
	return io.ReadFile(path)
}

func (c *RunnerImpl) SelectTaskName(taskfile string) (value.FullTaskName, error) {
	name, err := io.SelectTaskName(taskfile)
	if err != nil {
		return "", fmt.Errorf("io.SelectTaskName: %w", err)
	}

	return value.NewFullTaskName(name), nil
}

func (c *RunnerImpl) Input(prompt string) string {
	return io.Input(prompt)
}

func (c *RunnerImpl) RunTask(taskfile string, fullName value.FullTaskName, args ...string) error {
	return io.RunTask(taskfile, string(fullName), args...)
}
