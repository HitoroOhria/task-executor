package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executer/domain/value"
	"github.com/go-task/task/v3/taskfile/ast"
)

type Cmd struct {
	Command      *string
	TaskName     *value.TaskName
	FullTaskName *value.FullTaskName
}

func NewCmd(cmd *ast.Cmd, includeNames []string) (*Cmd, error) {
	var command *string
	if cmd.Cmd != "" {
		command = &cmd.Cmd
	}

	var tname *value.TaskName
	var fname *value.FullTaskName
	if cmd.Task != "" {
		name, err := value.NewTaskName(cmd.Task)
		if err != nil {
			return nil, fmt.Errorf("value.NewTaskName: %w", err)
		}
		tname = &name

		fullName := value.NewFullTaskNameForIncluded(includeNames, cmd.Task)
		fname = &fullName
	}

	return &Cmd{
		Command:      command,
		TaskName:     tname,
		FullTaskName: fname,
	}, nil
}
