package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executor/domain/console"
	"github.com/go-task/task/v3/taskfile/ast"
)

// RequiredVar は必須変数
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

// Arg は引数を返す
// 引数は "<var_name>=<var_value>" の形式である
func (v *RequiredVar) Arg() string {
	return makeArg(v.Name, v.MustInputValue())
}

// Deprecated: Input は変数の値を入力する
func (v *RequiredVar) Input(maxDisplayLen int) error {
	prompt := v.Prompt.Generate(maxDisplayLen, "")
	value := v.deps.Runner.Input(prompt)

	if value == "" {
		return fmt.Errorf("variable %s is required", v.Name)
	}

	v.InputValue = &value
	return nil
}
