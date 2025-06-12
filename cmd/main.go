package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-task/task/v3/errors"
)

func getArgs() string {
	taskfileName := flag.String("taskfile", "", "Taskfile name.")
	flag.Parse()

	return *taskfileName
}

func main() {
	var err error
	taskfileName := getArgs()

	// 引数で Taskfile の指定がなければ、カレントディレクトリから探索する
	if taskfileName == "" {
		taskfileName, err = findTaskfileName()
		if err != nil {
			handleError(err, "failed to get taskfile name")
			return
		}
	}

	taskName, err := selectTaskName(taskfileName)
	if err != nil {
		if errors.Is(err, ErrSpecifiedTaskfileNotFound) {
			handleError(err, fmt.Sprintf("taskfile not found: %s", taskfileName))
			return
		}
		// インクリメンタルサーチ中にキャンセルされた場合、何もしない
		if errors.Is(err, ErrCanceledIncrementalSearch) {
			os.Exit(0)
			return
		}

		handleError(err, "failed to select task name")
		return
	}

	file, err := readFile(taskfileName)
	if err != nil {
		handleError(err, "failed to read file")
		return
	}

	tf, err := NewTaskfile(taskfileName, file)
	if err != nil {
		handleError(err, "failed to new Taskfile")
		return
	}

	task := tf.Tasks.FindByName(taskName)
	if task == nil {
		handleError(fmt.Errorf("task '%s' not found", taskName), "failed to find task")
		return
	}

	err = task.Vars.Input()
	if err != nil {
		handleError(err, "failed to input vars")
		return
	}

	// タスクを実行
	err = runTask(tf.Name, task.Name, task.CommandArgs()...)
	if err != nil {
		handleError(err, "failed to run task")
		return
	}
}

func handleError(err error, msg string) {
	_, printErr := fmt.Fprintf(os.Stderr, "%s.\n%v.\n", msg, err)
	if printErr != nil {
		log.Fatalf("fmt.Fprintf: %v.\n%s.\n%v.\n", printErr, msg, err)
	}

	os.Exit(1)
}
