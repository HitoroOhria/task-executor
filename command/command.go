//go:generate mockgen -source=command.go -destination=mock/command.go -package=mock_command

package command

type Command interface {
	ReadFile(path string) ([]byte, error)
	SelectTaskName(taskfile string) (string, error)
	Prompt(maxNameLen int, varName string) string
	Input(prompt string) string
	RunTask(taskfile string, name string, args ...string) error
}
