package model

import (
	"fmt"

	"github.com/go-task/task/v3/taskfile/ast"
)

type Cmds []*Cmd

func NewCmds(cmds []*ast.Cmd, includeNames []string) (Cmds, error) {
	cs := make([]*Cmd, 0, len(cmds))
	for _, cmd := range cmds {
		c, err := NewCmd(cmd, includeNames)
		if err != nil {
			return nil, fmt.Errorf("NewCmd: %w", err)
		}
		cs = append(cs, c)
	}

	return cs, nil
}

func (cs Cmds) DependencyTasks() []*Cmd {
	cmds := make([]*Cmd, 0, len(cs))
	for _, c := range cs {
		if c.TaskName == nil {
			continue
		}

		cmds = append(cmds, c)
	}

	return cmds
}
