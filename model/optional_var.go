package model

import (
	"fmt"

	"github.com/go-task/task/v3/taskfile/ast"
)

type OptionalVar struct {
	v        *ast.Var
	inputter *Inputter

	Name       string
	Value      VarValue
	InputValue *string
}

func NewOptionalVar(name string, v *ast.Var) *OptionalVar {
	value := NewVarValue(v.Value)

	return &OptionalVar{
		v:          v,
		inputter:   inputter,
		Name:       name,
		Value:      value,
		InputValue: nil,
	}
}

func (v *OptionalVar) MustInputValue() string {
	if v.InputValue == nil {
		panic(fmt.Sprintf("input value is not set: %s", v.Name))
	}

	return *v.InputValue
}

func (v *OptionalVar) Arg() *string {
	if v.InputValue == nil || *v.InputValue == "" {
		return nil
	}

	arg := fmt.Sprintf(`%s="%s"`, v.Name, v.MustInputValue())
	return &arg
}

// IsInputtable は変数の値が入力可能かを判定する
func (v *OptionalVar) IsInputtable() bool {
	if v.Value.IsOptional(v.Name) {
		return true
	}
	if v.Value.IsOptionalWithDefault(v.Name) {
		return true
	}

	return false
}

func (v *OptionalVar) Input(maxNameLen int) {
	prompt := v.inputter.Prompt(maxNameLen, v.Name)
	value := v.inputter.Input(prompt)

	v.InputValue = &value
}
