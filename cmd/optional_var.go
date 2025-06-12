package main

import (
	"fmt"
	"regexp"

	"github.com/go-task/task/v3/taskfile/ast"
)

type OptionalVar struct {
	v        *ast.Var
	inputter *Inputter

	Name       string
	Value      string
	InputValue *string
}

func NewOptionalVar(name string, v *ast.Var) *OptionalVar {
	i := NewInputter()

	value := ""
	switch v.Value.(type) {
	case string:
		value = v.Value.(string)
	case *string:
		value = *v.Value.(*string)
	default:
		fmt.Printf("unknown var value type: type = %T. value = %+v\n", v.Value, v.Value)
	}

	return &OptionalVar{
		v:          v,
		inputter:   i,
		Name:       name,
		Value:      value,
		InputValue: nil,
	}
}

func (v *OptionalVar) MustInputValue() string {
	if v.InputValue == nil {
		panic("required var value is not set")
	}

	return *v.InputValue
}

func (v *OptionalVar) Arg() string {
	return fmt.Sprintf(`%s="%s"`, v.Name, v.MustInputValue())
}

func (v *OptionalVar) SelfValue() string {
	return fmt.Sprintf("{{.%s}}", v.Name)
}

func (v *OptionalVar) SelfWithDefaultRegex() *regexp.Regexp {
	re := fmt.Sprintf(`%s ?| default ".+"`, v.SelfValue())
	return regexp.MustCompile(re)
}

// IsInputtable は変数の値が入力可能かを判定する
func (v *OptionalVar) IsInputtable() bool {
	if v.Value == v.SelfValue() {
		return true
	}
	if v.SelfWithDefaultRegex().MatchString(v.Value) {
		return true
	}

	return false
}

func (v *OptionalVar) Input(maxNameLen int) {
	prompt := v.inputter.Prompt(false, maxNameLen, v.Name)
	value := v.inputter.Input(prompt)

	v.InputValue = &value
}
