package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const defaultTaskfileName = "Taskfile.yml"

func getArgs() string {
	taskfileName := flag.String("taskfile", defaultTaskfileName, "Taskfile name.")
	flag.Parse()

	return *taskfileName
}

func main() {
	taskfileName := getArgs()

	taskName, err := selectTaskName(taskfileName)
	if err != nil {
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

	// タスクの変数の値を受け付け
	vars := make(Vars)
	for name, task := range tf.Tasks.All(NoSort) {
		if name != taskName {
			continue
		}

		if task.Vars != nil {
			for varName, variable := range task.Vars.All() {
				if VarIsSpecifiable(varName, variable) {
					value := readInput(varName)
					vars.SetOptional(varName, value)
				}
			}
		}

		if task.Requires != nil {
			for _, variable := range task.Requires.Vars {
				value := readInput(variable.Name)
				err = vars.SetRequired(variable.Name, value)
				if err != nil {
					handleError(err, "failed to set required variable")
					return
				}
			}
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
