package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executor/domain/console"
	"github.com/go-task/task/v3/taskfile/ast"
)

// OptionalVar はオプショナル変数
// required ではない変数のこと
type OptionalVar struct {
	v    *ast.Var
	deps *console.Deps

	Name       string
	Value      VarValue
	Prompt     *Prompt
	InputValue *string
}

func NewOptionalVar(name string, v *ast.Var, deps *console.Deps) *OptionalVar {
	value := NewVarValue(v.Value)
	prompt := NewPrompt(name)

	return &OptionalVar{
		v:          v,
		deps:       deps,
		Name:       name,
		Value:      value,
		Prompt:     prompt,
		InputValue: nil,
	}
}

func (v *OptionalVar) MustInputValue() string {
	if v.InputValue == nil {
		panic(fmt.Sprintf("input value is not set: %s", v.Name))
	}

	return *v.InputValue
}

// Arg は引数を返す
// 引数は "<var_name>=<var_value>" の形式である
// 変数が未入力の場合、nil を返す
func (v *OptionalVar) Arg() *string {
	if v.InputValue == nil || *v.InputValue == "" {
		return nil
	}

	arg := fmt.Sprintf(`%s=%s`, v.Name, v.MustInputValue())
	return &arg
}

// IsInputtable は変数が入力可能かを判定する
// 入力可能な形式は "{{.VARIABLE}}" のような値である場合である
// 固定値がセットされている場合や、コマンド結果などの場合は入力可能ではない
func (v *OptionalVar) IsInputtable() bool {
	if v.Value.IsOptional(v.Name) {
		return true
	}
	if v.Value.IsOptionalWithDefault(v.Name) {
		return true
	}

	return false
}

// Input は変数の値を入力する
func (v *OptionalVar) Input(maxNameLen int) {
	prompt := v.Prompt.Generate(maxNameLen)
	value := v.deps.Runner.Input(prompt)

	v.InputValue = &value
}
