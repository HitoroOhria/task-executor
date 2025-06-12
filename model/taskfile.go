package model

import (
	"fmt"

	"github.com/go-task/task/v3/taskfile/ast"
	"gopkg.in/yaml.v3"
)

type Taskfile struct {
	tf *ast.Taskfile

	Name  string
	Tasks Tasks
}

func NewTaskfile(name string, file []byte) (*Taskfile, error) {
	tf := &ast.Taskfile{}
	err := yaml.Unmarshal(file, tf)
	if err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	ts := make(Tasks, 0)
	for _, task := range tf.Tasks.All(NoSort) {
		t := NewTask(task)
		ts = append(ts, t)
	}

	return &Taskfile{
		tf:    tf,
		Name:  name,
		Tasks: ts,
	}, nil
}

// NoSort
// TODO "github.com/go-task/task/v3@v3.44.0/internal/sort/sorter.go" の関数を参照する
func NoSort(items []string, namespaces []string) []string {
	return items
}
