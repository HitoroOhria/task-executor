package model

import (
	"fmt"

	"github.com/go-task/task/v3/taskfile/ast"
)

// Vars は変数名と値のセット
type Vars struct {
	deps *Deps

	Requires  []*RequiredVar
	Optionals []*OptionalVar
}

func NewVars(t *ast.Task, deps *Deps) *Vars {
	rvs := make([]*RequiredVar, 0)
	if t.Requires != nil {
		for _, v := range t.Requires.Vars {
			rvs = append(rvs, NewRequiredVar(v, deps))
		}
	}

	ovs := make([]*OptionalVar, 0)
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

// InputtableOptVars は入力可能はオプショナル変数を返す
func (vs *Vars) InputtableOptVars() []*OptionalVar {
	ovs := make([]*OptionalVar, 0, len(vs.Optionals))
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

		for _, r := range vs.Requires {
			err := r.Input(vs.GetMaxNameLen())
			if err != nil {
				return fmt.Errorf("r.Input: %w", err)
			}
		}
	}

	if len(vs.InputtableOptVars()) != 0 {
		vs.deps.Printer.OptionalHeader()

		for _, o := range vs.InputtableOptVars() {
			o.Input(vs.GetMaxNameLen())
		}
	}

	vs.deps.Printer.EndLine()

	return nil
}

func (vs *Vars) GetMaxNameLen() int {
	maxLen := 0

	for _, r := range vs.Requires {
		if len(r.Name) > maxLen {
			maxLen = len(r.Name)
		}
	}
	for _, o := range vs.InputtableOptVars() {
		if len(o.Name) > maxLen {
			maxLen = len(o.Name)
		}
	}

	return maxLen
}

// CommandArgs はコマンドの引数を組み立てる
// e.g. { "NAME": "john", "age": "25" } => [NAME="john", age="25"]
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
