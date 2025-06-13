package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executer/domain/console"
	"github.com/go-task/task/v3/taskfile/ast"
)

type RequiredVar struct {
	v    *ast.VarsWithValidation
	deps *console.Deps

	Name       string
	Prompt     *Prompt
	InputValue *string
}

func NewRequiredVar(v *ast.VarsWithValidation, deps *console.Deps) *RequiredVar {
	prompt := NewPrompt(v.Name)

	return &RequiredVar{
		v:          v,
		deps:       deps,
		Name:       v.Name,
		Prompt:     prompt,
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
	prompt := v.Prompt.Generate(maxNameLen)
	value := v.deps.Command.Input(prompt)

	if value == "" {
		return fmt.Errorf("variable %s is required", v.Name)
	}

	v.InputValue = &value
	return nil
}
