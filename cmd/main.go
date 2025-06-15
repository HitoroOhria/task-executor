package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/HitoroOhria/task-executor/adapter"
	"github.com/HitoroOhria/task-executor/domain/console"
	"github.com/HitoroOhria/task-executor/domain/model"
	"github.com/HitoroOhria/task-executor/io"
	"github.com/go-task/task/v3/errors"
)

var deps *console.Deps

func init() {
	deps = console.NewDeps(
		adapter.NewRunner(),
		adapter.NewPrinter(),
		adapter.NewVariableInputter(),
	)
}

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

	tf, err := model.NewTaskfile(taskfilePath, nil, deps)
	if err != nil {
		handleError(err, "failed to new Taskfile")
		return
	}

	task, err := tf.SelectTask()
	if err != nil {
		// インクリメンタルサーチ中にキャンセルされた場合、何もしない
		if errors.Is(err, io.ErrCanceledIncrementalSearch) {
			os.Exit(0)
			return
		}

		handleError(err, "failed to select task")
		return
	}

	err = tf.InputVars(task.FullName)
	if err != nil {
		handleError(err, "failed to input task vars")
		return
	}

	err = tf.RunSelectedTask()
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
