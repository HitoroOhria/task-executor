package main

import (
	"fmt"

	"github.com/go-task/task/v3/taskfile/ast"
	"gopkg.in/yaml.v3"
)

func NewTaskfile(file []byte) (*ast.Taskfile, error) {
	tf := &ast.Taskfile{}
	err := yaml.Unmarshal(file, tf)
	if err != nil {
		panic(fmt.Errorf("yaml.Unmarshal: %w", err))
	}

	return tf, nil
}

// NoSort
// TODO "github.com/go-task/task/v3@v3.44.0/internal/sort/sorter.go" の関数を参照する
func NoSort(items []string, namespaces []string) []string {
	return items
}
