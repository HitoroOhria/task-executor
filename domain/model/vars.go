package model

import (
	"errors"
	"fmt"

	"github.com/HitoroOhria/task-executor/domain/console"
	"github.com/go-task/task/v3/taskfile/ast"
)

var (
	ErrRequiredVarIsEmpty = errors.New("required variable is empty")
)

// Vars は変数の集合
type Vars struct {
	deps *console.Deps

	Requires  RequiredVars
	Optionals OptionalVars
}

func NewVars(t *ast.Task, deps *console.Deps) *Vars {
	rvs := make(RequiredVars, 0)
	if t.Requires != nil {
		for _, v := range t.Requires.Vars {
			rvs = append(rvs, NewRequiredVar(v, deps))
		}
	}

	ovs := make(OptionalVars, 0)
	if t.Vars != nil {
		for name, v := range t.Vars.All() {
			ovs = append(ovs, NewOptionalVar(name, &v, deps))
		}
	}

	return &Vars{
		deps:      deps,
		Requires:  rvs,
		Optionals: ovs,
	}
}

func (vs *Vars) inputtableLen() int {
	return len(vs.Requires) + len(vs.InputtableOptVars())
}

func (vs *Vars) Duplicate() *Vars {
	return &Vars{
		deps:      vs.deps,
		Requires:  vs.Requires,
		Optionals: vs.Optionals,
	}
}

func (vs *Vars) Merge(target *Vars) {
	vs.Requires = append(vs.Requires, target.Requires...)
	vs.Optionals = append(vs.Optionals, target.Optionals...)
}

// InputtableOptVars は入力可能はオプショナル変数を返す
func (vs *Vars) InputtableOptVars() OptionalVars {
	ovs := make(OptionalVars, 0, len(vs.Optionals))
	for _, ov := range vs.Optionals {
		if ov.IsInputtable() {
			ovs = append(ovs, ov)
		}
	}

	return ovs
}

func (vs *Vars) Input() error {
	if len(vs.Requires) == 0 && len(vs.InputtableOptVars()) == 0 {
		return nil
	}

	vars := make([]*console.Variable, 0, vs.inputtableLen())
	if len(vs.Requires) != 0 {
		for _, r := range vs.Requires.Distinct() {
			vars = append(vars, console.NewVariable(r.Name, true, ""))
		}
	}
	if len(vs.InputtableOptVars()) != 0 {
		for _, o := range vs.InputtableOptVars().Distinct() {
			vars = append(vars, console.NewVariable(o.Name, false, o.Value.Default()))
		}
	}

	results, err := vs.deps.VariableInputter.Input(vars)
	if err != nil {
		return fmt.Errorf("vs.deps.VariableInputter.Input: %w", err)
	}

	for _, result := range results {
		if result.Required {
			for _, r := range vs.Requires {
				if r.Name != result.Name {
					continue
				}

				if result.InputValue == "" {
					return ErrRequiredVarIsEmpty
				}
				r.InputValue = &result.InputValue
			}
		}

		for _, o := range vs.InputtableOptVars() {
			if o.Name != result.Name {
				continue
			}

			o.InputValue = &result.InputValue
		}
	}

	return nil
}

// CommandArgs はコマンドの引数を組み立てる
// e.g. { "NAME": "john", "age": "25" } => [NAME=john, age=25]
func (vs *Vars) CommandArgs() []string {
	args := make([]string, 0)

	for _, r := range vs.Requires {
		args = append(args, r.Arg())
	}
	for _, o := range vs.InputtableOptVars() {
		arg := o.Arg()
		if arg == nil {
			continue
		}

		args = append(args, *arg)
	}

	return args
}
