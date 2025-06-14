package model

import (
	"errors"
	"fmt"

	"github.com/HitoroOhria/task-executer/domain/console"
	"github.com/HitoroOhria/task-executer/domain/value"
	"github.com/go-task/task/v3/taskfile/ast"
	"gopkg.in/yaml.v3"
)

var ErrTaskNotFound = errors.New("task not found")

// Taskfile はタスクファイル
type Taskfile struct {
	tf   *ast.Taskfile
	deps *console.Deps

	FilePath string
	Tasks    Tasks
	Includes Includes
}

func NewTaskfile(filePath string, parentIncludeNames []string, deps *console.Deps) (*Taskfile, error) {
	file, err := deps.Runner.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("io.ReadFile: %w", err)
	}

	tf := &ast.Taskfile{}
	err = yaml.Unmarshal(file, tf)
	if err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	ts := make(Tasks, 0)
	for _, task := range tf.Tasks.All(NoSort) {
		t, err := NewTask(task, parentIncludeNames, deps)
		if err != nil {
			return nil, fmt.Errorf("NewTask: %w", err)
		}
		ts = append(ts, t)
	}

	is, err := NewIncludes(filePath, tf.Includes, parentIncludeNames, deps)
	if err != nil {
		return nil, fmt.Errorf("NewIncludes: %w", err)
	}

	return &Taskfile{
		tf:       tf,
		deps:     deps,
		FilePath: filePath,
		Tasks:    ts,
		Includes: is,
	}, nil
}

// FindTaskByFullName は完全タスク名に一致するタスクを探す
func (tf *Taskfile) FindTaskByFullName(fullName value.FullTaskName) *Task {
	task := tf.Tasks.FindByFullName(fullName)
	if task != nil {
		return task
	}

	for _, i := range tf.Includes {
		task = i.Taskfile.FindTaskByFullName(fullName)
		if task != nil {
			return task
		}
	}

	return nil
}

// FindSelectedTask は選択されたタスクを探す
func (tf *Taskfile) FindSelectedTask() *Task {
	found := tf.Tasks.FindSelected()
	if found != nil {
		return found
	}

	for _, i := range tf.Includes {
		found = i.Taskfile.FindSelectedTask()
		if found != nil {
			return found
		}
	}

	return nil
}

// SelectTask はタスクを選択する
func (tf *Taskfile) SelectTask() (*Task, error) {
	fullName, err := tf.deps.Runner.SelectTaskName(tf.FilePath)
	if err != nil {
		return nil, fmt.Errorf("cmd.SelectTaskName: %w", err)
	}

	task := tf.FindTaskByFullName(fullName)
	if task == nil {
		return nil, fmt.Errorf("%w: task = %s", ErrTaskNotFound, fullName)
	}

	task.Select()

	return task, nil
}

// CollectAllVars は、タスクの変数を依存タスクまで再起的に探索して集める
func (tf *Taskfile) CollectAllVars(task *Task) (*Vars, error) {
	vars := task.Vars.Duplicate()
	for _, cmd := range task.Cmds.FilterByDependencyTask() {
		depsTask := tf.FindTaskByFullName(cmd.DependencyTask.FullName)
		if depsTask == nil {
			return nil, fmt.Errorf("%w: dependency task = %s", ErrTaskNotFound, cmd.DependencyTask.FullName)
		}

		vs, err := tf.CollectAllVars(depsTask)
		if err != nil {
			return nil, fmt.Errorf("tf.CollectAllVars: %w", err)
		}
		vars.Merge(vs)
	}

	return vars, nil
}

// InputVars はタスクの変数を入力する
func (tf *Taskfile) InputVars(fullName value.FullTaskName) error {
	task := tf.FindTaskByFullName(fullName)
	if task == nil {
		return fmt.Errorf("%w: task = %s", ErrTaskNotFound, fullName)
	}

	vars, err := tf.CollectAllVars(task)
	if err != nil {
		return fmt.Errorf("tf.CollectAllVars: %w", err)
	}

	err = vars.Input()
	if err != nil {
		return fmt.Errorf("vars.Input: %w", err)
	}

	return nil
}

// RunSelectedTask は選択されたタスクを実行する
func (tf *Taskfile) RunSelectedTask() error {
	selected := tf.FindSelectedTask()
	if selected == nil {
		return fmt.Errorf("%w: selected task not found", ErrTaskNotFound)
	}

	vars, err := tf.CollectAllVars(selected)
	if err != nil {
		return fmt.Errorf("tf.CollectAllVars: %w", err)
	}

	tf.deps.Printer.LineBreaks()
	tf.deps.Printer.ExecutionTask(tf.FilePath, selected.FullName, vars.CommandArgs()...)

	return tf.deps.Runner.RunTask(tf.FilePath, selected.FullName, vars.CommandArgs()...)
}

// NoSort
// TODO "github.com/go-task/task/v3@v3.44.0/internal/sort/sorter.go" の関数を参照する
func NoSort(items []string, namespaces []string) []string {
	return items
}
