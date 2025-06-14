package value

import (
	"strings"
)

// FullTaskName は完全タスク名
// 包含された Taskfile のタスクの場合、ルート Taskfile からの完全なタスク名となる
// その場合は "include1:include2:task" のようにダブルコロン区切りで続く
type FullTaskName string

func NewFullTaskName(name string) FullTaskName {
	return FullTaskName(name)
}

func NewIncludedFullTaskName(includeNames []string, taskName string) FullTaskName {
	names := append(includeNames, taskName)
	name := strings.Join(names, ":")

	return NewFullTaskName(name)
}
