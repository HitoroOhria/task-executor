//go:generate mockgen -source=command.go -destination=mock/command.go -package=mock_console

package console

import "github.com/HitoroOhria/task-executer/domain/value"

type Command interface {
	ReadFile(path string) ([]byte, error)
	SelectTaskName(taskfile string) (value.FullTaskName, error)
	Input(prompt string) string
	RunTask(taskfile string, fullName value.FullTaskName, args ...string) error
}
