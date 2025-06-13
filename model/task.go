package model

import (
	"github.com/HitoroOhria/task-executer/command"
	"github.com/go-task/task/v3/taskfile/ast"
)

type Task struct {
	t *ast.Task

	Name string
	Vars *Vars
}

func NewTask(t *ast.Task, cmd command.Command) *Task {
	vs := NewVars(t, cmd)

	return &Task{
		t:    t,
		Name: t.Name(),
		Vars: vs,
	}
}

// CommandArgs はコマンドの引数を組み立てる
// e.g. { "NAME": "john", "age": "25" } => [NAME="john", age="25"]
func (t *Task) CommandArgs() []string {
	return t.Vars.CommandArgs()
}
