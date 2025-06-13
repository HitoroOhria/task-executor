package model

import (
	"github.com/HitoroOhria/task-executer/domain/console"
	"github.com/HitoroOhria/task-executer/domain/value"
	"github.com/go-task/task/v3/taskfile/ast"
)

type Task struct {
	t    *ast.Task
	deps *console.Deps

	Name     value.TaskName
	FullName value.FullTaskName
	Vars     *Vars
	Selected bool
}

func NewTask(t *ast.Task, includeNames []string, deps *console.Deps) *Task {
	name := value.NewTaskName(t.Name())
	fullName := value.NewFullTaskNameForIncluded(includeNames, t.Name())
	vs := NewVars(t, deps)

	return &Task{
		t:        t,
		deps:     deps,
		Name:     name,
		FullName: fullName,
		Vars:     vs,
		Selected: false,
	}
}

func (t *Task) Select() {
	t.Selected = true
}

func (t *Task) Input() error {
	return t.Vars.Input()
}

// CommandArgs はコマンドの引数を組み立てる
// e.g. { "NAME": "john", "age": "25" } => [NAME="john", age="25"]
func (t *Task) CommandArgs() []string {
	return t.Vars.CommandArgs()
}

func (t *Task) Run(taskfile string) error {
	t.deps.Printer.LineBreaks()
	t.deps.Printer.ExecutionTask(taskfile, t.FullName, t.CommandArgs()...)

	return t.deps.Command.RunTask(taskfile, t.FullName, t.CommandArgs()...)
}
