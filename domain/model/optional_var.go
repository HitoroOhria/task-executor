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
	InputValue *string
}

func NewOptionalVar(name string, v *ast.Var, deps *console.Deps) *OptionalVar {
	value := NewVarValue(v.Value)

	return &OptionalVar{
		v:          v,
		deps:       deps,
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

// Arg は引数を返す
// 引数は "<var_name>=<var_value>" の形式である
// 変数が未入力の場合、nil を返す
// 変数が見入力であり、デフォルト値がある場合、デフォルト値の引数を返す
func (v *OptionalVar) Arg() *string {
	var value string
	if v.InputValue != nil && *v.InputValue != "" {
		value = *v.InputValue
	}
	if value == "" && v.Value.Default() != "" {
		value = v.Value.Default()
	}

	if value == "" {
		return nil
	}

	arg := makeArg(v.Name, value)
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
