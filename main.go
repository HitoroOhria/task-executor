package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const defaultTaskfileName = "Taskfile.yml"

func getArgs() string {
	taskfileName := flag.String("taskfile", defaultTaskfileName, "Taskfile name.")
	flag.Parse()

	return *taskfileName
}

func main() {
	taskfileName := getArgs()

	taskName, err := selectTaskName()
	if err != nil {
		handleError(err, "failed to select task name")
		return
	}
	fmt.Printf("taskName = %+v\n", taskName)

	// Taskfile.yml を開く
	file, err := readFile(taskfileName)
	if err != nil {
		handleError(err, "failed to read file")
		return
	}

	// パース用構造体に読み込む
	tf, err := NewTaskfile(file)
	if err != nil {
		handleError(err, "failed to new Taskfile")
		return
	}

	// タスク一覧を表示
	for name, task := range tf.Tasks.All(NoSort) {
		fmt.Printf("Task: %s\n", name)
		fmt.Println("  Vars:")
		for varName := range task.Vars.All() {
			fmt.Printf("    - %s\n", varName)
		}
	}
}

func handleError(err error, msg string) {
	_, printErr := fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	if printErr != nil {
		log.Fatalf("fmt.Fprintf: %v. and %s: %v\n", printErr, msg, err)
	}

	os.Exit(1)
}

const selectTaskNameCommand = `
  task -l --sort none \
    | tail -n +2 \
    | peco \
    | sed -E 's/^\* ([^ ]+):.*/\1/' \
    | sed -E 's/:$//'
`

func selectTaskName() (string, error) {
	cmd := exec.Command("sh", "-c", selectTaskNameCommand)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.Output: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func readFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return file, nil
}
