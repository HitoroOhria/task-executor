package model

import (
	"strings"
)

// FullTaskName は完全タスク名
// includes された Taskfile のタスクの場合、ルート Taskfile からの完全なタスク名となる
// "include1:include2:task" のようにダブルコロン区切りで続く
type FullTaskName string

func NewFullTaskName(name string) FullTaskName {
	return FullTaskName(name)
}

func NewFullTaskNameForIncluded(includeNames []string, taskName string) FullTaskName {
	names := append(includeNames, taskName)
	name := strings.Join(names, ":")

	return NewFullTaskName(name)
}
