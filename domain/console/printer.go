//go:generate mockgen -source=printer.go -destination=mock/printer.go -package=mock_console

package console

import "github.com/HitoroOhria/task-executer/domain/value"

type Printer interface {
	RequiredHeader()
	OptionalHeader()
	EndLine()
	LineBreaks()
	ExecutionTask(taskfile string, fullName value.FullTaskName, args ...string)
}
