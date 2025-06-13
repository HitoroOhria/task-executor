package model

import (
	"errors"
	"fmt"

	"github.com/HitoroOhria/task-executer/command"
	"github.com/go-task/task/v3/taskfile/ast"
	"gopkg.in/yaml.v3"
)

var ErrTaskNotFound = errors.New("task not found")

type Taskfile struct {
	tf  *ast.Taskfile
	cmd command.Command

	FilePath string
	Tasks    Tasks
	Includes Includes
}

func NewTaskfile(filePath string, cmd command.Command) (*Taskfile, error) {
	file, err := cmd.ReadFile(filePath)
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
		t := NewTask(task, cmd)
		ts = append(ts, t)
	}

	is, err := NewIncludes(filePath, tf.Includes, cmd)
	if err != nil {
		return nil, fmt.Errorf("NewIncludes: %w", err)
	}

	return &Taskfile{
		tf:       tf,
		cmd:      cmd,
		FilePath: filePath,
		Tasks:    ts,
		Includes: is,
	}, nil
}

func (tf *Taskfile) SelectTask() (*Task, error) {
	taskName, err := tf.cmd.SelectTaskName(tf.FilePath)
	if err != nil {
		return nil, fmt.Errorf("cmd.SelectTaskName: %w", err)
	}

	task := tf.Tasks.FindByName(taskName)
	if task == nil {
		return nil, fmt.Errorf("%w: task = %s", ErrTaskNotFound, taskName)
	}

	task.Select()

	return task, nil
}

func (tf *Taskfile) RunSelectedTask() error {
	selected := tf.Tasks.FindSelected()
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
