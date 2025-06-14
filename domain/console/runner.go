//go:generate mockgen -source=runner.go -destination=mock/runner.go -package=mock_console

package console

import "github.com/HitoroOhria/task-executor/domain/value"

// Runner はコマンドを走らせるもの
type Runner interface {
	ReadFile(path string) ([]byte, error)
	SelectTaskName(taskfile string) (value.FullTaskName, error)
	Input(prompt string) string
	RunTask(taskfile string, fullName value.FullTaskName, args ...string) error
}
