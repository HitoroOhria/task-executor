package model

import (
	"errors"
	"fmt"

	"github.com/HitoroOhria/task-executer/domain/console"
	"github.com/HitoroOhria/task-executer/domain/value"
	"github.com/go-task/task/v3/taskfile/ast"
	"gopkg.in/yaml.v3"
)

var ErrTaskNotFound = errors.New("task not found")

type Taskfile struct {
	tf   *ast.Taskfile
	deps *console.Deps

	FilePath string
	Tasks    Tasks
	Includes Includes
}

func NewTaskfile(filePath string, parentIncludeNames []string, deps *console.Deps) (*Taskfile, error) {
	file, err := deps.Command.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("io.ReadFile: %w", err)
	}

	tf := &ast.Taskfile{}
	err = yaml.Unmarshal(file, tf)
	if err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	ts := make(Tasks, 0)
	for _, task := range tf.Tasks.All(NoSort) {
		t := NewTask(task, parentIncludeNames, deps)
		ts = append(ts, t)
	}

	is, err := NewIncludes(filePath, tf.Includes, parentIncludeNames, deps)
	if err != nil {
		return nil, fmt.Errorf("NewIncludes: %w", err)
	}

	return &Taskfile{
		tf:       tf,
		deps:     deps,
		FilePath: filePath,
		Tasks:    ts,
		Includes: is,
	}, nil
}

func (tf *Taskfile) FindTaskFullByName(fullName value.FullTaskName) *Task {
	task := tf.Tasks.FindByFullName(fullName)
	if task != nil {
		return task
	}

	for _, i := range tf.Includes {
		task = i.Taskfile.FindTaskFullByName(fullName)
		if task != nil {
			return task
		}
	}

	return nil
}

func (tf *Taskfile) FindSelectedTask() *Task {
	found := tf.Tasks.FindSelected()
	if found != nil {
		return found
	}

	for _, i := range tf.Includes {
		found = i.Taskfile.FindSelectedTask()
		if found != nil {
			return found
		}
	}

	return nil
}

func (tf *Taskfile) SelectTask() (*Task, error) {
	fullName, err := tf.deps.Command.SelectTaskName(tf.FilePath)
	if err != nil {
		return nil, fmt.Errorf("cmd.SelectTaskName: %w", err)
	}

	task := tf.FindTaskFullByName(fullName)
	if task == nil {
		return nil, fmt.Errorf("%w: task = %s", ErrTaskNotFound, fullName)
	}

	task.Select()

	return task, nil
}

func (tf *Taskfile) RunSelectedTask() error {
	selected := tf.FindSelectedTask()
	if selected == nil {
		return fmt.Errorf("%w: selected task not found", ErrTaskNotFound)
	}

	return selected.Run(tf.FilePath)

}

// NoSort
// TODO "github.com/go-task/task/v3@v3.44.0/internal/sort/sorter.go" の関数を参照する
func NoSort(items []string, namespaces []string) []string {
	return items
}
