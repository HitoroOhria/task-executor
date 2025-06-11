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
type Vars []*Var

type Var struct {
	Required bool
	Name     string
	Value    string
}

func (vs *Vars) FindByName(name string) *Var {
	for _, v := range *vs {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func (vs *Vars) SetNameAsOptional(name string) {
	*vs = append(*vs, &Var{
		Required: false,
		Name:     name,
		Value:    "",
	})
}

func (vs *Vars) SetNameAsRequired(name string) {
	*vs = append(*vs, &Var{
		Required: true,
		Name:     name,
		Value:    "",
	})
}

func (vs *Vars) SetValue(name, value string) error {
	v := vs.FindByName(name)
	if v == nil {
		return fmt.Errorf("variable %s is not found. vs = %+v", name, *vs)
	}

	if v.Required {
		if value == "" {
			return fmt.Errorf("variable %s is required", name)
		}
	}

	v.Value = value
	return nil
}

func (vs *Vars) GetMaxNameLen() int {
	maxLen := 0
	for _, v := range *vs {
		if len(v.Name) > maxLen {
			maxLen = len(v.Name)
		}
	}

	return maxLen
}

// CommandArgs はコマンドの引数を組み立てる
// e.g. { "NAME": "john", "age": "25" } => [NAME="john", age="25"]
func (vs *Vars) CommandArgs() []string {
	args := make([]string, 0, len(*vs))
	for _, v := range *vs {
		arg := fmt.Sprintf(`%s="%s"`, v.Name, v.Value)
		args = append(args, arg)
	}

	return args
}
