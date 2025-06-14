package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executor/domain/console"
	"github.com/HitoroOhria/task-executor/domain/value"
	"github.com/go-task/task/v3/taskfile/ast"
)

// Task はタスク
type Task struct {
	t *ast.Task

	Name     value.TaskName
	FullName value.FullTaskName
	Cmds     Cmds
	Vars     *Vars
	Selected bool
}

func NewTask(t *ast.Task, includeNames []string, deps *console.Deps) (*Task, error) {
	name, err := value.NewTaskName(t.Name())
	if err != nil {
		return nil, fmt.Errorf("value.NewTaskName: %w", err)
	}

	fullName := value.NewIncludedFullTaskName(includeNames, t.Name())
	cmds, err := NewCmds(t.Cmds, includeNames)
	if err != nil {
		return nil, fmt.Errorf("NewCmds: %w", err)
	}
	vs := NewVars(t, deps)

	return &Task{
		t:        t,
		Name:     name,
		FullName: fullName,
		Cmds:     cmds,
		Vars:     vs,
		Selected: false,
	}, nil
}

// Select はタスクを選択する
func (t *Task) Select() {
	t.Selected = true
}
