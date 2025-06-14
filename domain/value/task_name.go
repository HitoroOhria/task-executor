package value

import "errors"

var ErrTaskNameIsEmpty = errors.New("task name is empty")

// TaskName はタスク名
// 個別の Taskfile の単一タスク名である
type TaskName string

func NewTaskName(name string) (TaskName, error) {
	if name == "" {
		return "", ErrTaskNameIsEmpty
	}

	return TaskName(name), nil
}
