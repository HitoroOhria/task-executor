package model

import (
	"fmt"
	"path/filepath"

	"github.com/HitoroOhria/task-executer/domain/console"
	"github.com/go-task/task/v3/taskfile/ast"
)

type Include struct {
	Name     string
	Taskfile *Taskfile
}

func NewInclude(name string, taskfile *Taskfile) *Include {
	return &Include{
		Name:     name,
		Taskfile: taskfile,
	}
}

type Includes []*Include

func NewIncludes(parentTaskfilePath string, includes *ast.Includes, parentIncludeNames []string, deps *console.Deps) (Includes, error) {
	dir := filepath.Dir(parentTaskfilePath)

	is := make(Includes, 0)
	for name, i := range includes.All() {
		path := filepath.Join(dir, i.Taskfile)
		includeNames := append(parentIncludeNames, name)

		tf, err := NewTaskfile(path, includeNames, deps)
		if err != nil {
			return nil, fmt.Errorf("NewTaskfile: %w", err)
		}

		is = append(is, NewInclude(name, tf))
	}

	return is, nil
}
