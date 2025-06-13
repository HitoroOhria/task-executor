package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executer/io"
	"github.com/go-task/task/v3/taskfile/ast"
	"gopkg.in/yaml.v3"
)

type Taskfile struct {
	tf *ast.Taskfile

	FilePath string
	Tasks    Tasks
	Includes Includes
}

func NewTaskfile(filePath string) (*Taskfile, error) {
	// FIXME DI する
	file, err := io.ReadFile(filePath)
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
		t := NewTask(task)
		ts = append(ts, t)
	}

	is, err := NewIncludes(filePath, tf.Includes)
	if err != nil {
		return nil, fmt.Errorf("NewIncludes: %w", err)
	}

	return &Taskfile{
		tf:       tf,
		FilePath: filePath,
		Tasks:    ts,
		Includes: is,
	}, nil
}

// NoSort
// TODO "github.com/go-task/task/v3@v3.44.0/internal/sort/sorter.go" の関数を参照する
func NoSort(items []string, namespaces []string) []string {
	return items
}
