package io

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-task/task/v3/errors"
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

	ErrSpecifiedTaskfileNotFound = errors.New("specified taskfile not found")
	ErrCanceledIncrementalSearch = errors.New("canceled incremental search")
)

func findFileByName(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	foundName := ""
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if info.Name() == name {
			foundName = info.Name()
			return nil
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("filepath.Walk: %w", err)
	}

	return foundName, nil
}

// FindTaskfileName は、カレントディレクトリの Taskfile を探索し、ファイル名を返却する
func FindTaskfileName() (string, error) {
	taskfileName := ""
	for _, taskfile := range searchTaskfiles {
		found, err := findFileByName(taskfile)
		if err != nil {
			// 期待するファイル名をループ探索しているので、ファイルが見つからなくてもスルーする
			if errors.Is(err, ErrSpecifiedTaskfileNotFound) {
				continue
			}

			return "", fmt.Errorf("findFileByName: %w", err)
		}
		if found == "" {
			continue
		}

		taskfileName = found
		break
	}
	if taskfileName == "" {
		return "", fmt.Errorf("taskfile not found")
	}

	return taskfileName, nil
}

func SelectTaskName(taskfile string) (string, error) {
	found, err := findFileByName(taskfile)
	if err != nil {
		return "", fmt.Errorf("findFileByName: %w", err)
	}
	if found == "" {
		return "", ErrSpecifiedTaskfileNotFound
	}

	cmd := exec.Command("sh", "-c", selectTaskNameCommand(taskfile))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.Output: %w", err)
	}

	o := strings.TrimSpace(string(output))
	// インクリメンタルサーチ中に Ctrl + C で中断された場合は、特定のエラーを返す
	if o == "" {
		return "", ErrCanceledIncrementalSearch
	}

	return o, nil
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return file, nil
}

// ReadInput は値の入力を受け付ける
// Ctrl + C でキャンセルされた場合は、プログラムを正常終了する
// FIXME context & goroutine を使用した方法もあるので、検討する
func ReadInput() string {
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

func RunTask(taskfile string, name string, args ...string) error {
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
