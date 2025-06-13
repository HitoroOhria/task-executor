package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	cmdimpl "github.com/HitoroOhria/task-executer/command/impl"
	"github.com/HitoroOhria/task-executer/io"
	"github.com/HitoroOhria/task-executer/model"
	"github.com/go-task/task/v3/errors"
)

func getArgs() string {
	taskfilePath := flag.String("taskfile", "", "Taskfile path.")
	flag.Parse()

	return *taskfilePath
}

func main() {
	var err error
	taskfilePath := getArgs()

	// 引数で Taskfile の指定がなければ、カレントディレクトリから探索する
	if taskfilePath == "" {
		taskfilePath, err = io.FindTaskfileName()
		if err != nil {
			handleError(err, "failed to get taskfile name")
			return
		}
	}

	cmd := cmdimpl.NewCommand(&cmdimpl.NewCommandArgs{
		ReadFile:       io.ReadFile,
		Prompt:         io.Prompt,
		Input:          io.Input,
		SelectTaskName: io.SelectTaskName,
	})

	tf, err := model.NewTaskfile(taskfilePath, cmd)
	if err != nil {
		handleError(err, "failed to new Taskfile")
		return
	}

	taskName, err := tf.SelectTask()
	if err != nil {
		if errors.Is(err, io.ErrTaskfileNotFound) {
			handleError(err, fmt.Sprintf("taskfile not found: %s", taskfilePath))
			return
		}
		// インクリメンタルサーチ中にキャンセルされた場合、何もしない
		if errors.Is(err, io.ErrCanceledIncrementalSearch) {
			os.Exit(0)
			return
		}

		handleError(err, "failed to select task name")
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
	err = io.RunTask(tf.FilePath, task.Name, task.CommandArgs()...)
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
