package model

import (
	"fmt"

	"github.com/HitoroOhria/task-executor/domain/console"
	"github.com/go-task/task/v3/taskfile/ast"
)

// Vars は変数の集合
type Vars struct {
	deps *console.Deps

	Requires  RequiredVars
	Optionals OptionalVars
}

type RequiredVars []*RequiredVar
type OptionalVars []*OptionalVar

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

	if len(vs.Requires) != 0 {
		vs.deps.Printer.RequiredHeader()

		for _, r := range vs.Requires.Distinct() {
			err := r.Input(vs.GetMaxRequireVarsDisplayLen())
			if err != nil {
				return fmt.Errorf("r.Input: %w", err)
			}
		}
	}

	if len(vs.InputtableOptVars()) != 0 {
		vs.deps.Printer.OptionalHeader()

		for _, o := range vs.InputtableOptVars().Distinct() {
			o.Input(vs.GetMaxOptionalVarsDisplayLen())
		}
	}

	vs.deps.Printer.EndLine()

	return nil
}

func (vs *Vars) GetMaxRequireVarsDisplayLen() int {
	varDisplays := make([]string, 0)
	for _, r := range vs.Requires {
		disp := generateVarDisplay(r.Name, "")
		varDisplays = append(varDisplays, disp)
	}

	maxLen := 0
	for _, disp := range varDisplays {
		if len(disp) > maxLen {
			maxLen = len(disp)
		}
	}

	return maxLen
}

func (vs *Vars) GetMaxOptionalVarsDisplayLen() int {
	varDisplays := make([]string, 0)
	for _, i := range vs.InputtableOptVars() {
		disp := generateVarDisplay(i.Name, i.Value.Default())
		varDisplays = append(varDisplays, disp)
	}

	maxLen := 0
	for _, disp := range varDisplays {
		if len(disp) > maxLen {
			maxLen = len(disp)
		}
	}

	return maxLen
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

func (vs RequiredVars) existByName(name string) bool {
	for _, v := range vs {
		if v.Name == name {
			return true
		}
	}

	return false
}

func (vs RequiredVars) Distinct() RequiredVars {
	vars := make(RequiredVars, 0, len(vs))
	for _, v := range vs {
		if vars.existByName(v.Name) {
			continue
		}

		vars = append(vars, v)
	}

	return vars
}

func (vs OptionalVars) existByName(name string) bool {
	for _, v := range vs {
		if v.Name == name {
			return true
		}
	}

	return false
}

func (vs OptionalVars) Distinct() OptionalVars {
	vars := make(OptionalVars, 0, len(vs))
	for _, v := range vs {
		if vars.existByName(v.Name) {
			continue
		}

		vars = append(vars, v)
	}

	return vars
}
