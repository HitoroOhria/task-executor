package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executer/domain/value"
	"github.com/go-task/task/v3/taskfile/ast"
)

// Cmd はコマンド
type Cmd struct {
	Command        *string
	DependencyTask *DependencyTask
}

type DependencyTask struct {
	Name     value.TaskName
	FullName value.FullTaskName
}

func NewCmd(cmd *ast.Cmd, includeNames []string) (*Cmd, error) {
	var command *string
	if cmd.Cmd != "" {
		command = &cmd.Cmd
	}

	var dependencyTask *DependencyTask
	if cmd.Task != "" {
		name, err := value.NewTaskName(cmd.Task)
		if err != nil {
			return nil, fmt.Errorf("value.NewTaskName: %w", err)
		}
		fullName := value.NewFullTaskNameForIncluded(includeNames, cmd.Task)

		dependencyTask = &DependencyTask{
			Name:     name,
			FullName: fullName,
		}
	}

	return &Cmd{
		Command:        command,
		DependencyTask: dependencyTask,
	}, nil
}
