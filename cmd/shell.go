package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const incrementalSearchTool = "peco"

var (
	searchTaskfiles = []string{"Taskfile.yml", "Taskfile.yaml"}

	selectTaskNameCommand = func(taskfile string) string {
		return fmt.Sprintf(`
	  task -t %s -l --sort none | \
	    tail -n +2 | \
	    sed 's/^\*//' | \
	    %s | \
	    sed -E 's/^ ([^ ]+):.*/\1/' | \
	    sed -E 's/:$//'
	`, taskfile, incrementalSearchTool)
	}

	ErrSelectedTaskfileNotFound = fmt.Errorf("taskfile not found")
)

// findTaskfileName は、カレントディレクトリの Taskfile を探索し、ファイル名を返却する
func findTaskfileName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	taskfileName := ""
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		for _, taskfile := range searchTaskfiles {
			if info.Name() == taskfile {
				taskfileName = info.Name()
				return nil
			}
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("filepath.Walk: %w", err)
	}

	if taskfileName == "" {
		return "", fmt.Errorf("taskfile not found")
	}

	return taskfileName, nil
}

func selectTaskName(taskfile string) (string, error) {
	cmd := exec.Command("sh", "-c", selectTaskNameCommand(taskfile))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.Output: %w", err)
	}

	o := strings.TrimSpace(string(output))
	if o == "" {
		return "", ErrSelectedTaskfileNotFound
	}

	return o, nil
}

func readFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return file, nil
}

func readOptionalInput(varName string, padding int) string {
	printInputPrompt(varName, padding, false)
	return readInput()
}

func readRequiredInput(varName string, padding int) string {
	printInputPrompt(varName, padding, true)
	return readInput()
}

func printInputPrompt(varName string, padding int, required bool) {
	necessity := "optional"
	if required {
		necessity = "required"
	}

	promptPadding := padding + 2
	promptVarName := fmt.Sprintf(`"%s"`, varName)

	fmt.Printf(`Enter %-*s (%s): `, promptPadding, promptVarName, necessity)
}

// readInput は値の入力を受け付ける
// Ctrl + C でキャンセルされた場合は、プログラムを正常終了する
func readInput() string {
	// Ctrl+C (SIGINT) を補足
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	// 入力がキャンセルされた場合、何もしない
	go func() {
		<-sigCh
		fmt.Println()
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		// 入力が完了したらシグナルハンドラーを停止
		signal.Reset(syscall.SIGINT)

		return scanner.Text()
	}

	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: scanner.Err():", err)
	}

	return ""
}

func runTask(taskfile string, name string, args ...string) error {
	fmt.Printf("run: task -t %s %s %s\n", taskfile, name, strings.Join(args, " "))

	cmdArgs := append([]string{"-t", taskfile, name}, args...)
	cmd := exec.Command("task", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run: %w", err)
	}

	return nil
}
