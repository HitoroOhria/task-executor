package model

import (
	"fmt"

	"github.com/go-task/task/v3/taskfile/ast"
)

type RequiredVar struct {
	v        *ast.VarsWithValidation
	inputter *Inputter

	Name       string
	InputValue *string
}

func NewRequiredVar(v *ast.VarsWithValidation) *RequiredVar {
	return &RequiredVar{
		v:          v,
		inputter:   inputter,
		Name:       v.Name,
		InputValue: nil,
	}
}

func (v *RequiredVar) MustInputValue() string {
	if v.InputValue == nil {
		panic("required var value is not set")
	}

	return *v.InputValue
}

func (v *RequiredVar) Arg() string {
	return fmt.Sprintf(`%s="%s"`, v.Name, v.MustInputValue())
}

func (v *RequiredVar) Input(maxNameLen int) error {
	prompt := v.inputter.Prompt(true, maxNameLen, v.Name)
	value := v.inputter.Input(prompt)

	if value == "" {
		return fmt.Errorf("variable %s is required", v.Name)
	}

	v.InputValue = &value
	return nil
}
