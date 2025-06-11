package main

import (
	"fmt"
	"os"

	"github.com/go-task/task/v3/taskfile/ast"
	"gopkg.in/yaml.v3"
)

func main() {
	// Taskfile.yml を開く
	data, err := os.ReadFile("Taskfile.yml")
	if err != nil {
		panic(fmt.Errorf("failed to read Taskfile.yml: %w", err))
	}

	// パース用構造体に読み込む
	tf := &ast.Taskfile{}
	err = yaml.Unmarshal(data, tf)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal Taskfile.yml: %w", err))
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

// NoSort
// TODO "github.com/go-task/task/v3@v3.44.0/internal/sort/sorter.go" の関数を参照する
func NoSort(items []string, namespaces []string) []string {
	return items
}
