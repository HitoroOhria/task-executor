//go:generate mockgen -source=command.go -destination=mock/command.go -package=mock_model

package model

type Command interface {
	ReadFile(path string) ([]byte, error)
	SelectTaskName(taskfile string) (FullTaskName, error)
	Input(prompt string) string
	RunTask(taskfile string, fullName FullTaskName, args ...string) error
}
