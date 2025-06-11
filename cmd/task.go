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

// VarIsSpecifiable は変数の値が指定可能かを判定する
func VarIsSpecifiable(name string, variable ast.Var) bool {
	specifiable := []string{
		"{{.}}",
		fmt.Sprintf("{{.%s}}", name),
	}

	for _, s := range specifiable {
		if variable.Value == s {
			return true
		}
	}

	return false
}

// Vars は変数名と値のセット
type Vars map[string]string

func (v Vars) SetOptional(name, value string) {
	v[name] = value
}

func (v Vars) SetRequired(name, value string) error {
	if value == "" {
		return fmt.Errorf("variable %s is required", name)
	}

	v[name] = value
	return nil
}

// CommandArgs はコマンドの引数を組み立てる
// e.g. { "NAME": "john", "age": "25" } => [NAME="john", age="25"]
func (v Vars) CommandArgs() []string {
	args := make([]string, 0, len(v))
	for name, value := range v {
		arg := fmt.Sprintf(`%s="%s"`, name, value)
		args = append(args, arg)
	}

	return args
}
