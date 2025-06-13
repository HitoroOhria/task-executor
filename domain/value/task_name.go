package value

// TaskName はタスク名
// 個別の Taskfile の単一タスク名である
type TaskName string

func NewTaskName(name string) TaskName {
	return TaskName(name)
}
