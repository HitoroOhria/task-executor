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
		// インクリメンタルサーチ中にキャンセルされた場合、何もしない
		if errors.Is(err, ErrSelectedTaskfileNotFound) {
			os.Exit(0)
		}

		handleError(err, "failed to select task name")
		return
	}

	file, err := readFile(taskfileName)
	if err != nil {
		handleError(err, "failed to read file")
		return
	}

	tf, err := NewTaskfile(file)
	if err != nil {
		handleError(err, "failed to new Taskfile")
		return
	}

	// タスクの変数を収集
	vars := make(Vars, 0)
	for name, task := range tf.Tasks.All(NoSort) {
		if name != taskName {
			continue
		}

		if task.Requires != nil {
			for _, variable := range task.Requires.Vars {
				vars.SetNameAsRequired(variable.Name)
			}
		}
		if task.Vars != nil {
			for varName, variable := range task.Vars.All() {
				if VarIsSpecifiable(varName, variable) {
					vars.SetNameAsOptional(varName)
				}
			}
		}
	}

	// タスクの変数の値を入力
	for _, v := range vars {
		padding := vars.GetMaxNameLen()

		if v.Required {
			value := readRequiredInput(v.Name, padding)
			err = vars.SetValue(v.Name, value)
			if err != nil {
				handleError(err, "failed to set value")
				return
			}

			continue
		}

		value := readOptionalInput(v.Name, padding)
		err = vars.SetValue(v.Name, value)
		if err != nil {
			handleError(err, "failed to set value")
			return
		}
	}

	// タスクを実行
	err = runTask(taskfileName, taskName, vars.CommandArgs()...)
	if err != nil {
		handleError(err, "failed to run task")
		return
	}
}

func handleError(err error, msg string) {
	_, printErr := fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	if printErr != nil {
		log.Fatalf("fmt.Fprintf: %v. and %s: %v\n", printErr, msg, err)
	}

	os.Exit(1)
}
