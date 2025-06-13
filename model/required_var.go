package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executer/command"
	"github.com/go-task/task/v3/taskfile/ast"
)

type RequiredVar struct {
	v   *ast.VarsWithValidation
	cmd command.Command

	Name       string
	InputValue *string
}

func NewRequiredVar(v *ast.VarsWithValidation, cmd command.Command) *RequiredVar {
	return &RequiredVar{
		v:          v,
		cmd:        cmd,
		Name:       v.Name,
		InputValue: nil,
	}
}

func (v *RequiredVar) MustInputValue() string {
	if v.InputValue == nil {
		panic(fmt.Sprintf("input value is not set: %s", v.Name))
	}

	return *v.InputValue
}

func (v *RequiredVar) Arg() string {
	return fmt.Sprintf(`%s="%s"`, v.Name, v.MustInputValue())
}

func (v *RequiredVar) Input(maxNameLen int) error {
	prompt := v.cmd.Prompt(maxNameLen, v.Name)
	value := v.cmd.Input(prompt)

	if value == "" {
		return fmt.Errorf("variable %s is required", v.Name)
	}

	v.InputValue = &value
	return nil
}
